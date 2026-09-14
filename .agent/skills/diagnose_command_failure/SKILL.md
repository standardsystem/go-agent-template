---
name: diagnose_command_failure
description: エージェントがコマンドを実行した際に先頭 1 バイトが欠落する・日本語が文字化けする・引数エスケープに失敗するケースを診断する。問題発生時に呼び出す。
---

# コマンド先頭バイト欠落 診断スキル

エージェントがコマンドを実行した際に、**コマンド名の先頭 1 バイト（または数バイト）が
切れてコマンドが失敗する**現象を調査するためのスキル。付属スクリプトは
`.agent/skills/diagnose_command_failure/scripts/` にあり、リポジトリルートから実行する。

## 症状の例

- `git status` → `it status: command not found`
- `go build ./...` → `o build: command not found`
- `gh issue create ...` → `h issue: command not found`

## 診断手順

### STEP 1: ヘルパースクリプトを実行して環境情報を収集

```powershell
& ".\.agent\skills\diagnose_command_failure\scripts\collect_env.ps1"
```

### STEP 2: 失敗したコマンドの再現テスト

1. ユーザーから失敗したコマンドの**正確な文字列**を取得する
2. 以下のスクリプトでバイト列を解析する

```powershell
# コマンド文字列のバイト解析（変数 $cmd にコマンド文字列をセット）
$cmd = "<ここに失敗したコマンド文字列>"
& ".\.agent\skills\diagnose_command_failure\scripts\analyze_bytes.ps1" -InputString $cmd
```

### STEP 3: 一時ファイルの BOM 確認

エージェントが一時スクリプトファイル経由でコマンドを実行している場合、BOM が原因の
可能性がある。

```powershell
& ".\.agent\skills\diagnose_command_failure\scripts\check_temp_files.ps1"
```

### STEP 4: コードページとエンコーディングの確認

```powershell
& ".\.agent\skills\diagnose_command_failure\scripts\collect_env.ps1" -Verbose
```

### STEP 5: 診断結果をまとめてログに保存

保存先は `temp/`（Git 管理外。拡張子は `.log`）。

```powershell
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$logPath = "temp\cmd_diagnosis_${timestamp}.log"
& ".\.agent\skills\diagnose_command_failure\scripts\collect_env.ps1" | Tee-Object -FilePath $logPath
Write-Host "診断ログ保存先: $logPath"
```

## チェックリスト（調査時に確認する項目）

| 項目 | 確認方法 | 備考 |
|------|----------|------|
| シェルのバージョン | `$PSVersionTable.PSVersion` | **`7.x` であること**。5.1 は文字化け・BOM の原因 |
| 一時ファイルの BOM | `analyze_bytes.ps1` で先頭バイト確認 | EF BB BF = UTF-8 BOM |
| PATH 環境変数 | `$env:PATH` | コマンドが見つかるか |
| コマンド経路 | `Get-Command <cmd>` | which ではなく Get-Command |
| エラーの詳細 | `$Error[0] \| Format-List *` | 直近のエラー詳細 |

## 既知の原因パターン

### パターン A: BOM による先頭バイト欠落

**症状**: 一時ファイル経由でスクリプトを実行すると最初のコマンド文字が欠落

**原因**: PowerShell が一時 `.ps1` ファイルを **UTF-8 with BOM (EF BB BF)** で保存した
場合、一部のシェルやツールが BOM をコマンド文字として解釈してしまう

**対策**:

```powershell
# ファイル保存時に BOM なし UTF-8 を明示
$content | Out-File -FilePath $path -Encoding utf8NoBOM
```

### パターン B: 二重シェル (pwsh の上に powershell)

**症状**: `powershell -Command "..."` が `The term 'powershel...` で失敗

**原因**: すでに `pwsh` (PowerShell 7) 上で動いているのに、`powershell.exe`
(Windows PowerShell 5) を呼び出そうとして失敗（またはコマンドが切り詰められる）

**対策**: `powershell -Command` は使わない。直接 `pwsh` コマンドを書く

### パターン C: マルチバイト文字の途中切断・文字化け

**症状**: 日本語を含むコマンド引数や出力が途中で切れる・化ける

**原因**: **Windows PowerShell 5.1 で実行している**。5.1 は入出力の既定が UTF-8 ではなく、
`Set-Content` の既定も ANSI (日本語 Windows では CP932) なので、UTF-8 前提のツールとの
やり取りで化ける。`[Console]::OutputEncoding` を書き換えてもファイル出力側は直らない。

**対策**: シェルを PowerShell 7 (pwsh) にする。7 は入出力とも既定が UTF-8 で、
コードページの操作は不要。

```powershell
$PSVersionTable.PSVersion   # 7.x であることを確認する
```

導入手順は
[PowerShell 7 セットアップ手順書](../../../docs/manuals/POWERSHELL7_SETUP.md)
を参照。

### パターン D: エージェントのコマンド実行ツールの引数エスケープ

**症状**: エージェントが直接コマンド実行ツールで渡したコマンドの先頭が切れる

**原因**: ツール側での文字列処理（エスケープ、クォーティング）の問題

**対策**: コマンドを一時 `.ps1` ファイル（`temp/` 配下、BOM なし UTF-8）に書き出して
`& file.ps1` で実行する

## 情報収集後の対応

診断ログと収集した情報をもとに、以下を記録する:

1. **症状**: 失敗したコマンドとエラーメッセージ
2. **環境情報**: `collect_env.ps1` の出力
3. **再現条件**: どのスキル・ワークフローを使った時に発生したか
4. **発生頻度**: 毎回か特定条件か

これらを `temp/cmd_diagnosis_<timestamp>.log` に保存してユーザーに報告する。

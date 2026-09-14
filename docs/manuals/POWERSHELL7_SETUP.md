# PowerShell 7 (pwsh) セットアップ手順書

対象読者: 本プロジェクトの開発メンバー（Windows 環境）
所要時間: 5 分程度
関連: [AGENTS.md](../../AGENTS.md) BASE RULES

## 1. 目的

本リポジトリのターミナルコマンドは **PowerShell 7 (pwsh)** で実行する。Windows に
標準で入っている Windows PowerShell 5.1 は使わない（正典は
[AGENTS.md](../../AGENTS.md) BASE RULES）。

5.1 のままだと、書き出したファイルに BOM が付いて外部ツールが構文エラーを出す、
ネイティブコマンドに渡した JSON 引数が途中で割れる、日本語の出力が文字化けする、
といった問題が出る。7 は入出力の既定が UTF-8 なので、これらは起きない。

本手順を完了すると `$PSVersionTable.PSVersion` が `7.x` を返し、VS Code の統合
ターミナルもそのシェルで開くようになる。

## 2. 現在のシェルを確認する

VS Code の統合ターミナルで実行する。

```powershell
$PSVersionTable.PSVersion
```

| 結果 | 対応 |
|---|---|
| `Major` が `7` | 導入済み。§5 の確認だけ行えばよい |
| `Major` が `5`（`5.1.x`） | §3 へ進む（導入済みでも既定シェルが 5.1 になっている場合がある） |
| コマンドが見つからない | §3 へ進む |

## 3. PowerShell 7 を導入する

```powershell
winget install --id Microsoft.PowerShell --source winget
```

インストール後はターミナルを開き直す（`pwsh` が PATH に載る）。`winget` が使えない
場合は [PowerShell の GitHub リリース](https://github.com/PowerShell/PowerShell/releases)
から MSI を取得してもよい。

> **インストール先に注意**: Microsoft Store 版やユーザー単位インストールは
> `C:\Program Files\PowerShell\7\pwsh.exe` に入らない。sshd などサービスから呼ぶ
> 用途がある場合はマシン全体へのインストールが必要になる。上記 `winget` コマンドは
> マシン全体にインストールされる。

## 4. VS Code の既定シェルを切り替える

導入しただけでは統合ターミナルは 5.1 のままなので、明示的に切り替える。

1. `Ctrl+Shift+P` でコマンドパレットを開く
2. `Terminal: Select Default Profile` を選ぶ
3. **PowerShell**（`C:\Program Files\PowerShell\7\pwsh.exe`）を選ぶ

`Windows PowerShell` という名前のプロファイルは 5.1 なので選ばない。紛らわしいが、
**`PowerShell` が 7 で、`Windows PowerShell` が 5.1** である。本リポジトリの
[.vscode/settings.json](../../.vscode/settings.json) は既定プロファイルを `PowerShell`
にしてある。

## 5. 確認する

新しいターミナルを開いて、次が `7.x` を返せば完了。

```powershell
$PSVersionTable.PSVersion
```

## 6. うまくいかないとき

| 症状 | 原因と対処 |
|---|---|
| 切り替えたのに `5.1.x` が返る | 既存のターミナルタブが残っている。タブを閉じて開き直す |
| `pwsh` が見つからない | ターミナルを開き直していない（PATH が未反映）。VS Code 自体の再起動も試す |
| `winget` が見つからない | 「アプリ インストーラー」が未導入。Microsoft Store から入れるか、上記 MSI で導入する |
| 日本語が文字化けする | まだ 5.1 で動いている可能性が高い。§5 を再確認する。`[Console]::OutputEncoding` を書き換えるハックは使わない（7 では不要、5.1 では解決しない） |

## 7. 例外: 5.1 のまま書くべきスクリプト

pwsh がまだ入っていないホストを構築するためのスクリプト（存在する場合）だけは例外で、
`#Requires -Version 5.1` を維持し、5.1 でも動く書き方を守る（`utf8NoBOM` や
PowerShell 7 専用の演算子は使わない）。該当するスクリプトがあれば
[AGENTS.md](../../AGENTS.md) BASE RULES に例外として明記すること。

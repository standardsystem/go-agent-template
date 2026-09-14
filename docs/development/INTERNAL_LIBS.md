# 内部ライブラリ索引 (internal/)

`internal/` 配下の公開 API の索引。エージェントはソースを全部読む前にここで API と構造を
確認する（[AGENTS.md](../../AGENTS.md)「内部ライブラリの利用ルール」）。公開 API を
追加・変更・削除したら本書も更新すること。パスはリポジトリルート基準。

## [internal/scaffold](../../internal/scaffold)

テンプレート由来の名前（Go モジュールパス・プロジェクト名）を新しい名前へ一括置換する。
[cmd/cli/scaffold-init](../../cmd/cli/scaffold-init/main.go) の実体。

### 型

- **Options**: 置換の対象と内容
  - `Root string`: 走査するディレクトリ（通常はリポジトリルート）
  - `OldModule` / `NewModule string`: go.mod の module パス
  - `OldName` / `NewName string`: プロジェクト名
  - `DryRun bool`: true なら書き換えず報告だけ
- **Report**: 結果
  - `Scanned int`: 内容を検査したファイル数
  - `Changed []string`: 置換された（DryRun なら予定の）ファイルの相対パス。`/` 区切り・昇順

### 関数

- `func Run(opts Options) (Report, error)`: 置換を実行する。モジュールパスと名前は
  1 回の走査で同時に置換するため、`NewModule` に旧名が含まれても二重置換しない。
  走査から除くディレクトリ（`.git` `temp` `output` `data` `node_modules` 等）と対象
  拡張子（`.go` `.mod` `.md` `.toml` `.yaml` `.yml` `.json` `.sample` `.txt` `.ps1` 等）は
  パッケージ内の `skipDirs` / `targetExts` で定義

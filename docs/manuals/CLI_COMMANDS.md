# CLI コマンド一覧

`cmd/cli/<ツール名>/main.go` に 1 ツール 1 ディレクトリで置く。追加時は
[scaffold_new_tool](../../.agent/skills/scaffold_new_tool/SKILL.md) の手順に従い、本書に
節を追加する。ビルドは `mise run build`（`bin/` に出力。Git 管理外）。

## scaffold-init

- **パス**: `cmd/cli/scaffold-init`
- **概要**: テンプレート由来の Go モジュールパスとプロジェクト名を新しい名前へ一括置換
  する。テンプレートから複製した直後に 1 回だけ使う
- **使い方**:

```powershell
go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo> -dry-run   # 変更予定の表示
go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo>            # 実行 (名前はパス末尾)
go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo> -name <名前> -root .
```

- **環境変数**: `LOG_LEVEL`（DEBUG / INFO / WARN / ERROR。`.env` は mise が読み込む）
- **終了コード**: 0 成功 / 1 置換失敗 / 2 引数不正
- **備考**: 自分自身も置換対象なので、実行後は no-op になる。不要なら
  `cmd/cli/scaffold-init` と `internal/scaffold` を削除し、本書と
  [INTERNAL_LIBS.md](../development/INTERNAL_LIBS.md) から節を消す

---
name: scaffold_new_tool
description: cmd/cli/ に Go CLI ツールを 1 本追加する。雛形コード・テスト・文書更新の手順。
---

# Go CLI ツールの追加

新規ツールは Go で作る（[AGENTS.md](../../../AGENTS.md) BASE RULES）。
1 ツール 1 ディレクトリ、`main` は薄く、ロジックは `internal/` に置いてテストする。
実例は [cmd/cli/scaffold-init](../../../cmd/cli/scaffold-init/main.go) と
[internal/scaffold](../../../internal/scaffold/rename.go)。

## 手順 1: ディレクトリ作成

`cmd/cli/<tool-name>/`（小文字ケバブケース）を作る。

```powershell
New-Item -ItemType Directory -Force cmd/cli/<tool-name>
```

## 手順 2: main.go を雛形から作る

設定はフラグと環境変数で受ける（`.env` は mise が読み込むので、dotenv ライブラリは
不要）。ログは標準の `log/slog`。終了コードは 0 成功 / 1 処理失敗 / 2 引数不正。

```go
// <tool-name> は <一行説明>。
//
// 使い方:
//
//  go run ./cmd/cli/<tool-name> -input <path>
package main

import (
    "flag"
    "log/slog"
    "os"
)

func main() {
    var (
        input = flag.String("input", "", "入力ファイル (必須)")
    )
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
    if *input == "" {
        flag.Usage()
        logger.Error("-input は必須です")
        os.Exit(2)
    }
    if err := run(logger, *input); err != nil {
        logger.Error("処理に失敗しました", "error", err)
        os.Exit(1)
    }
}

// run は本体処理。テストしやすいよう main から分離し、ロジックは internal/ に委譲する。
func run(logger *slog.Logger, input string) error {
    logger.Info("開始", "input", input)
    return nil
}
```

## 手順 3: ロジックとテストを internal/ に置く

`internal/<pkg>/` にロジックを置き、`<pkg>_test.go` を
[テスト設計指針](../../../docs/development/TESTING_STANDARDS.md) に従って書く
（テスト対象情報の注釈・AAA 形式・`t.Logf` の日本語ログ）。

## 手順 4: 文書を更新する

- [CLI コマンド一覧](../../../docs/manuals/CLI_COMMANDS.md) に節を追加する
- `internal/` を追加・変更したら
  [内部ライブラリ索引](../../../docs/development/INTERNAL_LIBS.md) も更新する

## 手順 5: 検証

```powershell
mise run check
```

## 注意

- 一時出力は `temp/` だけ。デバッグ出力の拡張子は `.log`
- 外部依存を足すときは `go get` の後に `go mod tidy` し、追加理由をコミットメッセージに書く
- ビルド成果物（`*.exe`、`bin/`）はコミットしない（CI の tracked-ignored 検査で弾かれる）

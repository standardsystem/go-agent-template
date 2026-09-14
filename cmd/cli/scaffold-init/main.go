// scaffold-init は、テンプレートから複製した直後のリポジトリで、テンプレート由来の
// Go モジュールパスとプロジェクト名を新しい名前へ一括置換する。
//
// 使い方:
//
//	go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo> -dry-run   # 変更予定の表示だけ
//	go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo>            # 名前はパス末尾
//	go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo> -name <名前>
//
// 置換後は `go build ./...` と `mise run check` で緑を確認してから最初のコミットを作る。
// scaffold-init 自身も置換対象なので、実行後は下の定数が新しい名前になり、再実行しても
// 何も変わらない (no-op)。不要になったら cmd/cli/scaffold-init と internal/scaffold ごと
// 削除してよい。
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/standardsystem/go-agent-template/internal/scaffold"
)

// テンプレートの現在の名前。scaffold-init の実行で書き換わる。
const (
	currentModule = "github.com/standardsystem/go-agent-template"
	currentName   = "go-agent-template"
)

func main() {
	var (
		module = flag.String("module", "", "新しい Go モジュールパス (例: github.com/org/repo)。必須")
		name   = flag.String("name", "", "新しいプロジェクト名。省略時はモジュールパスの末尾要素")
		root   = flag.String("root", ".", "走査するリポジトリルート")
		dryRun = flag.Bool("dry-run", false, "書き換えずに変更予定のファイルだけ表示する")
	)
	flag.Parse()

	logger := newLogger()
	if *module == "" {
		flag.Usage()
		logger.Error("-module は必須です")
		os.Exit(2)
	}
	if *name == "" {
		*name = (*module)[strings.LastIndex(*module, "/")+1:]
	}

	report, err := scaffold.Run(scaffold.Options{
		Root:      *root,
		OldModule: currentModule,
		NewModule: *module,
		OldName:   currentName,
		NewName:   *name,
		DryRun:    *dryRun,
	})
	if err != nil {
		logger.Error("置換に失敗しました", "error", err)
		os.Exit(1)
	}
	for _, p := range report.Changed {
		fmt.Println(p)
	}
	logger.Info("scaffold-init 完了",
		"module", *module,
		"name", *name,
		"scanned", report.Scanned,
		"changed", len(report.Changed),
		"dry_run", *dryRun,
	)
}

// newLogger は環境変数 LOG_LEVEL (DEBUG / INFO / WARN / ERROR。.env は mise が読み込む) を
// 反映した slog ロガーを返す。未設定または不正な値なら INFO。
func newLogger() *slog.Logger {
	level := slog.LevelInfo
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if err := level.UnmarshalText([]byte(v)); err != nil {
			fmt.Fprintf(os.Stderr, "LOG_LEVEL=%q は不正です。INFO として扱います\n", v)
			level = slog.LevelInfo
		}
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

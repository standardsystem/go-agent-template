// Package scaffold は、テンプレートから複製した直後のリポジトリに残る
// テンプレート由来の名前 (Go モジュールパスとプロジェクト名) を新しい名前へ
// 一括置換する。
//
// cmd/cli/scaffold-init の実体であると同時に、このリポジトリにおける
// 「internal パッケージ + AAA 形式のテスト」の実例でもある。
package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options は置換の対象と内容を表す。
type Options struct {
	// Root は走査するディレクトリ。通常はリポジトリルート。
	Root string
	// OldModule と NewModule は go.mod の module パス (例: github.com/org/repo)。
	OldModule string
	NewModule string
	// OldName と NewName はプロジェクト名 (例: go-agent-template)。
	OldName string
	NewName string
	// DryRun が true のとき、ファイルは書き換えずに変更予定だけを報告する。
	DryRun bool
}

// Report は置換の結果。
type Report struct {
	// Scanned は内容を検査したファイル数。
	Scanned int
	// Changed は置換が発生した (DryRun のときは発生する予定の) ファイルの
	// Root からの相対パス。区切りは "/" に正規化し、昇順に並べる。
	Changed []string
}

// skipDirs は走査しないディレクトリ名。Git 管理外の置き場と依存物。
var skipDirs = map[string]bool{
	".git":         true,
	"temp":         true,
	"tmp":          true,
	"output":       true,
	"data":         true,
	"node_modules": true,
	".venv":        true,
	"bin":          true,
	"dist":         true,
	"vendor":       true,
}

// targetExts は置換対象とする拡張子。バイナリや生成物は対象にしない。
var targetExts = map[string]bool{
	".go":     true,
	".mod":    true,
	".md":     true,
	".toml":   true,
	".yaml":   true,
	".yml":    true,
	".json":   true,
	".jsonc":  true,
	".sample": true,
	".txt":    true,
	".ps1":    true,
}

// Run は opts に従って置換を行い、結果を返す。
//
// モジュールパスとプロジェクト名は 1 回の走査で同時に置換する (strings.Replacer)。
// そのため NewModule にテンプレート名が含まれていても (例: .../go-agent-template-v2)
// 名前の置換が二重に適用されることはない。
func Run(opts Options) (Report, error) {
	if err := validate(opts); err != nil {
		return Report{}, err
	}
	replacer := strings.NewReplacer(
		opts.OldModule, opts.NewModule,
		opts.OldName, opts.NewName,
	)
	oldModule := []byte(opts.OldModule)
	oldName := []byte(opts.OldName)

	var report Report
	err := filepath.WalkDir(opts.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != opts.Root && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !targetExts[filepath.Ext(d.Name())] {
			return nil
		}
		report.Scanned++

		before, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s の読み込みに失敗: %w", path, err)
		}
		if !bytes.Contains(before, oldModule) && !bytes.Contains(before, oldName) {
			return nil
		}
		after := replacer.Replace(string(before))

		rel, err := filepath.Rel(opts.Root, path)
		if err != nil {
			return err
		}
		report.Changed = append(report.Changed, filepath.ToSlash(rel))
		if opts.DryRun {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(after), info.Mode().Perm()); err != nil {
			return fmt.Errorf("%s の書き込みに失敗: %w", path, err)
		}
		return nil
	})
	if err != nil {
		return Report{}, err
	}
	sort.Strings(report.Changed)
	return report, nil
}

// validate は Options の整合性を検査する。ファイルには触らない。
func validate(opts Options) error {
	// エラー文字列は staticcheck ST1005 (先頭を大文字にしない) に合わせて日本語で始める。
	switch {
	case opts.Root == "":
		return errors.New("走査ルート (Root) が空です")
	case opts.OldModule == "" || opts.NewModule == "":
		return errors.New("モジュールパス (OldModule / NewModule) は必須です")
	case opts.OldName == "" || opts.NewName == "":
		return errors.New("プロジェクト名 (OldName / NewName) は必須です")
	case opts.OldModule == opts.NewModule && opts.OldName == opts.NewName:
		return errors.New("置換前後の名前が同じです (置換するものがありません)")
	}
	info, err := os.Stat(opts.Root)
	if err != nil {
		return fmt.Errorf("走査ルート (Root) を参照できません: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("走査ルート (Root) がディレクトリではありません: %s", opts.Root)
	}
	return nil
}

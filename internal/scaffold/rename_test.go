package scaffold

// [テスト対象情報]
// 対象ファイル: internal/scaffold/rename.go
// 対象要素: Run 関数 (Options による一括置換)、走査対象の判定 (skipDirs / targetExts)、validate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// フィクスチャ用の名前。テンプレート自身の名前とは無関係な値にして、
// scaffold-init で本リポジトリを改名してもこのテストが影響を受けないようにする。
const (
	fxOldModule = "example.com/old-org/proj-old"
	fxNewModule = "example.com/new-org/proj-new"
	fxOldName   = "proj-old"
	fxNewName   = "proj-new"
)

// TestRun_ReplacesOnlyTargetFiles は、対象拡張子のファイル内のモジュールパスと
// プロジェクト名が置換され、対象外 (temp/ 配下・対象外拡張子) は変更されないことを
// 確認する。
func TestRun_ReplacesOnlyTargetFiles(t *testing.T) {
	// Arrange: 一時ディレクトリに対象 3 件・対象外 2 件を用意する
	root := t.TempDir()
	writeFixture(t, root, "go.mod", "module "+fxOldModule+"\n\ngo 1.26\n")
	writeFixture(t, root, "cmd/app/main.go", "package main\n\nimport _ \""+fxOldModule+"/internal/x\"\n")
	writeFixture(t, root, "README.md", "# "+fxOldName+"\n")
	writeFixture(t, root, "temp/note.md", fxOldName+" は temp/ 配下なので置換しない\n")
	writeFixture(t, root, "asset.bin", fxOldName)

	// Act
	report, err := Run(Options{
		Root:      root,
		OldModule: fxOldModule,
		NewModule: fxNewModule,
		OldName:   fxOldName,
		NewName:   fxNewName,
	})

	// Assert
	if err != nil {
		t.Fatalf("Run がエラーを返した: %v", err)
	}
	wantChanged := []string{"README.md", "cmd/app/main.go", "go.mod"}
	if !reflect.DeepEqual(report.Changed, wantChanged) {
		t.Fatalf("置換されたファイル一覧が期待と異なる: want=%v got=%v", wantChanged, report.Changed)
	}
	if report.Scanned != 3 {
		t.Fatalf("検査したファイル数が期待と異なる: want=3 got=%d", report.Scanned)
	}
	assertContent(t, root, "go.mod", "module "+fxNewModule+"\n\ngo 1.26\n")
	assertContent(t, root, "cmd/app/main.go", "package main\n\nimport _ \""+fxNewModule+"/internal/x\"\n")
	assertContent(t, root, "README.md", "# "+fxNewName+"\n")
	assertContent(t, root, "temp/note.md", fxOldName+" は temp/ 配下なので置換しない\n")
	assertContent(t, root, "asset.bin", fxOldName)
	t.Logf("対象 3 ファイル %v が置換され、temp/ 配下と対象外拡張子は変更されないことを確認した (scanned=%d)",
		report.Changed, report.Scanned)
}

// TestRun_DryRunReportsWithoutWriting は、DryRun のとき変更予定のファイルは報告されるが
// 内容は書き換えられないことを確認する。
func TestRun_DryRunReportsWithoutWriting(t *testing.T) {
	// Arrange
	root := t.TempDir()
	original := "module " + fxOldModule + "\n"
	writeFixture(t, root, "go.mod", original)

	// Act
	report, err := Run(Options{
		Root:      root,
		OldModule: fxOldModule,
		NewModule: fxNewModule,
		OldName:   fxOldName,
		NewName:   fxNewName,
		DryRun:    true,
	})

	// Assert
	if err != nil {
		t.Fatalf("Run がエラーを返した: %v", err)
	}
	if !reflect.DeepEqual(report.Changed, []string{"go.mod"}) {
		t.Fatalf("変更予定の一覧が期待と異なる: want=[go.mod] got=%v", report.Changed)
	}
	assertContent(t, root, "go.mod", original)
	t.Logf("DryRun では go.mod が変更予定として報告され、内容は %q のまま保持されることを確認した", original)
}

// TestRun_DoesNotDoubleReplaceWhenNewModuleContainsOldName は、新しいモジュールパスに
// 旧プロジェクト名が含まれていても (例: proj-old-v2)、名前の置換が二重に適用されない
// ことを確認する。
func TestRun_DoesNotDoubleReplaceWhenNewModuleContainsOldName(t *testing.T) {
	// Arrange
	root := t.TempDir()
	newModule := "example.com/new-org/" + fxOldName + "-v2"
	writeFixture(t, root, "go.mod", "module "+fxOldModule+"\n")
	writeFixture(t, root, "README.md", "# "+fxOldName+"\n")

	// Act
	_, err := Run(Options{
		Root:      root,
		OldModule: fxOldModule,
		NewModule: newModule,
		OldName:   fxOldName,
		NewName:   fxNewName,
	})

	// Assert
	if err != nil {
		t.Fatalf("Run がエラーを返した: %v", err)
	}
	assertContent(t, root, "go.mod", "module "+newModule+"\n")
	assertContent(t, root, "README.md", "# "+fxNewName+"\n")
	t.Logf("go.mod は %q に置換され、その中の %q が %q へ再置換されないことを確認した",
		newModule, fxOldName, fxNewName)
}

// TestRun_RejectsInvalidOptions は、不正な Options に対して Run がファイルに触らずに
// エラーを返すことを確認する。
func TestRun_RejectsInvalidOptions(t *testing.T) {
	root := t.TempDir()
	valid := Options{
		Root:      root,
		OldModule: fxOldModule,
		NewModule: fxNewModule,
		OldName:   fxOldName,
		NewName:   fxNewName,
	}
	cases := []struct {
		name    string
		mutate  func(o *Options)
		wantErr string
	}{
		{
			name:    "Root が存在しない",
			mutate:  func(o *Options) { o.Root = filepath.Join(root, "missing") },
			wantErr: "走査ルート (Root) を参照できません",
		},
		{
			name:    "OldModule が空",
			mutate:  func(o *Options) { o.OldModule = "" },
			wantErr: "モジュールパス (OldModule / NewModule) は必須です",
		},
		{
			name:    "NewName が空",
			mutate:  func(o *Options) { o.NewName = "" },
			wantErr: "プロジェクト名 (OldName / NewName) は必須です",
		},
		{
			name: "置換前後が同じ",
			mutate: func(o *Options) {
				o.NewModule = o.OldModule
				o.NewName = o.OldName
			},
			wantErr: "置換前後の名前が同じです",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			opts := valid
			tc.mutate(&opts)

			// Act
			_, err := Run(opts)

			// Assert
			if err == nil {
				t.Fatalf("エラーを期待したが nil だった")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("エラーメッセージが期待と異なる: want=%q got=%q", tc.wantErr, err.Error())
			}
			t.Logf("%s のとき %q を含むエラーが返ることを確認した", tc.name, tc.wantErr)
		})
	}
}

// writeFixture は root からの相対パス rel にファイルを作る。親ディレクトリも作る。
func writeFixture(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("%s の親ディレクトリを作れない: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("%s を書けない: %v", rel, err)
	}
}

// assertContent は rel の内容が want と一致することを検証する。
func assertContent(t *testing.T, root, rel, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("%s を読めない: %v", rel, err)
	}
	if string(got) != want {
		t.Fatalf("%s の内容が期待と異なる:\nwant: %q\ngot:  %q", rel, want, string(got))
	}
}

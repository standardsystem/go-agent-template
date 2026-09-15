# go-agent-template

Go 製 CLI / 解析ツール群を **AI エージェント（Claude Code / Codex CLI 等）と協働して
開発する**ためのリポジトリテンプレート。`mise` によるツールチェーン固定、lint / test /
CI のハーネス、エージェント運用規約（AGENTS.md）、Git 管理外ディレクトリの命名規約を
最初から揃えた状態で新規プロジェクトを始められる。

抜き出し元は社内の解析リポジトリ（n-df-load-chart）で、そこで数十コミットかけて固まった
足場からプロジェクト固有の部分を除いたもの。

## 何が入っているか

| 領域 | ファイル | 役割 |
| --- | --- | --- |
| ツールチェーン | [.mise.toml](.mise.toml) | Go / Node / gh / golangci-lint / markdownlint-cli2 の固定とタスク |
| Go lint | [.golangci.yml](.golangci.yml) | golangci-lint v2。高シグナル linter を全件クリーン運用 |
| Markdown lint | [.markdownlint-cli2.yaml](.markdownlint-cli2.yaml) | エージェントが書く文書の品質ゲート |
| 改行・文字コード | [.gitattributes](.gitattributes) / [.editorconfig](.editorconfig) | LF 基本、Windows スクリプトのみ CRLF |
| CI | [.github/workflows/ci.yml](.github/workflows/ci.yml) | mise で環境再現 → build / vet / lint / test / 誤コミット検査 |
| エージェント規約 | [AGENTS.md](AGENTS.md) / [CLAUDE.md](CLAUDE.md) / [CODEX.md](CODEX.md) / [.agent/rules.md](.agent/rules.md) | 共通正典と各エージェントの入口、合意ゲート、行動規範 |
| 実装の最小化 | [AGENTS.md](AGENTS.md)「実装の最小化規範（Ponytail）」/ [docs/development/PONYTAIL.md](docs/development/PONYTAIL.md) | Ponytail（YAGNI・既存コード・標準ライブラリ優先の梯子）の常時ルールと、規約との優先順位・更新手順 |
| 手順書 | [.agent/skills/](.agent/skills/) / [.agents/skills/](.agents/skills/) | 手順本体と Codex 用入口（コミット、CLI 追加、コマンド失敗診断） |
| Claude Code 設定 | [.claude/settings.json](.claude/settings.json) / `.claude/agents/` | 共有の許可設定・Ponytail プラグインの共有有効化とサブエージェント定義の置き場 |
| 開発規約 | [docs/development/](docs/development/) | コミット規約、テスト設計指針、Git 管理外ディレクトリ命名、内部ライブラリ索引 |
| ナレッジ置き場 | [docs/knowledge/](docs/knowledge/index.md) | 索引 → 本体の階層でドメイン知識を蓄積 |
| 最小 CLI と実例テスト | [cmd/cli/scaffold-init](cmd/cli/scaffold-init/main.go) / [internal/scaffold](internal/scaffold/rename.go) | テンプレート名の一括置換ツール。規約どおりの CLI と AAA テストの実例 |
| Git 管理外 3 ディレクトリ | `data/` `output/` `temp/` | 受領物・成果物・一時ファイル。`.gitkeep` だけ追跡 |

## 使い方

### 1. 複製する

GitHub のテンプレートリポジトリとして（リポジトリ設定で Template repository を有効に
してある前提）:

```powershell
gh repo create <org>/<repo> --template standardsystem/go-agent-template --private --clone
Set-Location <repo>
```

ローカルのコピーから始める場合:

```powershell
git clone https://github.com/standardsystem/go-agent-template <repo>
Set-Location <repo>
Remove-Item -Recurse -Force .git
git init -b main
```

### 2. 名前を置き換える

```powershell
go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo> -dry-run   # 変更予定を確認
go run ./cmd/cli/scaffold-init -module github.com/<org>/<repo>            # 実行
```

モジュールパスと `go-agent-template` という名前が、go.mod・import・設定・文書で一括
置換される。置換後の scaffold-init は no-op になるので、不要なら `cmd/cli/scaffold-init`
と `internal/scaffold` を削除してよい（削除したら
[INTERNAL_LIBS.md](docs/development/INTERNAL_LIBS.md) と
[CLI_COMMANDS.md](docs/manuals/CLI_COMMANDS.md) からも節を消す）。

### 3. 環境を作って緑を確認する

```powershell
mise trust
mise run setup     # ツール導入 + go mod download
mise run check     # vet + lint + test
```

### 4. プロジェクト固有の部分を埋める

「【要記入】」を検索して埋める。主な箇所:

- [AGENTS.md](AGENTS.md): ドメインナレッジの必読資料、外部データソースの規約
- [docs/development/COMMIT_CONVENTION.md](docs/development/COMMIT_CONVENTION.md):
  チケットキーの形式、リリース運用
- [docs/development/UNTRACKED_DIR_NAMING.md](docs/development/UNTRACKED_DIR_NAMING.md):
  チケット表記、ツールが参照する固定ディレクトリ
- [docs/knowledge/index.md](docs/knowledge/index.md): 最初のナレッジ
- 本 README をプロジェクトの説明に書き換える

### 5. 最初のコミット

コミット規約は [COMMIT_CONVENTION.md](docs/development/COMMIT_CONVENTION.md)。
日本語メッセージは `-F` でファイルから渡す
（[git_commit_push](.agent/skills/git_commit_push/SKILL.md)）。

## mise タスク

| タスク | 内容 |
| --- | --- |
| `mise run setup` | ツール導入 + `go mod download` |
| `mise run build` | 全 CLI を `bin/` にビルド |
| `mise run test` | `go test -race -cover ./...` |
| `mise run test:short` | race 検出なし（cgo が使えない環境の代替） |
| `mise run vet` | `go vet ./...` |
| `mise run lint` | `lint:go`（golangci-lint）+ `lint:docs`（markdownlint-cli2） |
| `mise run fmt` | `gofmt` + `markdownlint-cli2 --fix` |
| `mise run check` | CI 相当（vet + lint + test） |
| `mise run info` | ツールのバージョン表示 |
| `mise run scaffold:init -- -module ...` | 名前の一括置換（`go run ./cmd/cli/scaffold-init` と同じ） |

## ディレクトリ構成

```text
cmd/cli/<tool>/        Go CLI (1 ツール 1 ディレクトリ。main.go)
internal/<pkg>/        内部ライブラリ (索引: docs/development/INTERNAL_LIBS.md)
docs/development/      開発規約
docs/manuals/          手順書 (CLI 一覧、PowerShell 7 導入)
docs/knowledge/        ドメインナレッジ (索引 → 本体)
.agent/                全エージェント共通の行動規範と手順書 (正本)
.agents/skills/        Codex 用の入口 (本体へリンク)
.claude/               Claude Code 用 (settings.json は共有、settings.local.json は個人)
.github/workflows/     CI
data/ output/ temp/    Git 管理外。命名規約: docs/development/UNTRACKED_DIR_NAMING.md
```

## テンプレートの育て方

- **昇格は 2 プロジェクト目から**: あるプロジェクトで入れた修正は、2 つ目のプロジェクトで
  同じ修正が要ったときにだけテンプレートへ取り込む。1 件の事情を持ち込まない
- **テンプレート自身の CI を緑に保つ**: 壊れたテンプレートは配布されるだけ害になる。
  [ci.yml](.github/workflows/ci.yml) はテンプレート自身にも回る
- **バージョンは動作確認した値に固定する**: [.mise.toml](.mise.toml) の版を上げるときは
  CI の緑を確認する
- **派生プロジェクトとの差分確認**: テンプレート側の更新を取り込むときは、足場ファイルを
  `git diff --no-index` で見比べる

```powershell
git diff --no-index --stat C:\Projects\STANDARDSYSTEM\go-agent-template\.mise.toml .\.mise.toml
```

## 関連テンプレート

- `webapp-template`: Go API + React + Docker Compose の Web アプリ向け。本テンプレートは
  CLI / 解析ツール群とエージェント運用規約に寄せている。両方に共通する足場
  （golangci-lint v2・markdownlint-cli2 の版）は揃えてある

# Codex 向けプロジェクト案内

共通ルールの正典は [AGENTS.md](AGENTS.md)。本ファイルは Codex から既存の知識・手順へ
到達するための案内であり、手順本体やドメイン知識を複製しない。

## 作業開始時の参照順

1. [AGENTS.md](AGENTS.md) と [.agent/rules.md](.agent/rules.md) を読み、今回の依頼範囲を
   確認する。必要な変更が既に依頼されている場合はその範囲で進め、未合意の外部操作を
   補わない。
2. 本ファイルのスキル対応表からタスクに合う入口を選び、参照先の手順本体を読む。
   関係のない全スキル・全ナレッジを一括で読み込む必要はない。
3. ドメインに関わる調査・判断では、結論を出す前に
   [ドメインナレッジ索引](docs/knowledge/index.md) と AGENTS.md「ドメインナレッジ」の
   資料を読む。索引でフロントマターの `type` / `status` / `stale_after` を確認し、
   該当する本体も読む。
4. 内部 API の調査は [内部ライブラリ索引](docs/development/INTERNAL_LIBS.md)、
   CLI の利用は [CLI コマンド一覧](docs/manuals/CLI_COMMANDS.md) と
   [.mise.toml](.mise.toml) から始める。テストの作成・修正時は
   [テスト設計指針](docs/development/TESTING_STANDARDS.md) を読む。
5. 報告には、実際に読んだ資料・適用した前提・検証結果・未確認点を示す。
   読めない資料や利用できないツールがあれば、その範囲を明記する。

## Claude Code と Codex の対応

| 既存資産 | Codex での扱い |
| --- | --- |
| [CLAUDE.md](CLAUDE.md) の `@AGENTS.md` | ルートの `AGENTS.md` を入口にし、Markdown リンク先は実際に開く |
| `.claude/skills/` | `.agents/skills/` の入口から手順本体を参照する |
| [.agent/skills/](.agent/skills/) | 共通手順の正本。Codex 用入口との対応は下表を参照する |
| `.claude/agents/` | エージェント定義を手順書として読む。読むだけでサブエージェントは起動しない |
| [.claude/settings.json](.claude/settings.json)・`settings.local.json` | Claude 専用の権限・環境設定。Codex の権限や実行許可として引き継がない |
| [.claude/settings.json](.claude/settings.json) の `enabledPlugins`（Ponytail） | Codex には引き継がれない。同じ常時ルールは [AGENTS.md](AGENTS.md)「実装の最小化規範」で効く。`/ponytail-review` 等が要るなら各自で `codex plugin marketplace add DietrichGebert/ponytail` → `codex plugin add ponytail@ponytail` を行う（[導入記録](docs/development/PONYTAIL.md)） |
| Claude のローカルメモリ | 共有知識は [ドメインナレッジ索引](docs/knowledge/index.md) と現行チケットから確認する |

Claude 専用の `Plan` モードや `Read`・`Grep`・`Glob`・`Bash` は、同名の Codex ツールが
存在するという意味ではない。現在のセッションにある読み取り・検索・編集・計画機能に
対応付ける。サブエージェントの起動はユーザーの依頼・承認がある場合だけ行う。

## 既存手順の読み替え

- **相対パス**: Markdown リンクは記載元ファイルを基準に解決する。手順中の `cmd/`・
  `docs/`・`scripts/` などのコマンド用パスはリポジトリルートを基準にする。リンクが
  切れている場合は `rg --files --hidden` で実在する移動先を確認する。
- **シェルと実行**: PowerShell 7 と `mise` を使う。Bash 用の環境変数指定・`/dev/null`・
  チェーンをそのまま実行せず、失敗ログも残す。新規ツールは Go で作る。WSL の存在を
  前提にしない。
- **一時ファイル**: OS の一時ディレクトリや `C:\tmp` を使う例は、
  [命名規約](docs/development/UNTRACKED_DIR_NAMING.md) に従った `temp/` に読み替える。
  受領 ZIP の展開物は `data/_extracted/`、成果物は `temp/` または `output/`。
  作業ファイルを残すか削除するかを完了時に判断する。
- **MCP**: ツール名を推測して呼び出さない。利用可能なツールの説明と引数を確認して
  対応付け、未接続なら使えない操作を報告する。Claude の設定をコピーして接続済みと
  みなさない。
- **操作の範囲**: スキルの選択自体は実行許可ではない。コミット・push・PR・投稿・
  アサイン変更・デプロイなどは、今回依頼された範囲だけ実行する。状態を変える操作は
  一つずつ結果を確認する。過去の許可リストから今回の合意を推測しない。
- **Git**: 課題対応は基準ブランチから課題専用ブランチを作成・選択してから編集する。
  無関係な変更や既存の未コミット変更を混在・破棄しない。番号未指定の作業は番号を
  捏造せず、内容が分かる専用ブランチで扱う。コミット形式は
  [コミット規約](docs/development/COMMIT_CONVENTION.md) に従い、メッセージは `temp/` の
  UTF-8 ファイルから `git commit -F` で渡す。ステージングは対象ファイルを指定し、
  別作業まで含む `git add .` を避ける。Claude 用の `Co-Authored-By` の値を Codex として
  流用しない。
- **知識の共有**: ローカルメモリや認証情報をプロジェクトへ一括コピーしない。共有すべき
  知識は無害化して [ナレッジの書き方](docs/knowledge/README.md) に従って追加する。

## スキル対応表

Codex 用入口は `.agents/skills/<名前>/SKILL.md` に置く。手順本体のアンダースコア名は
Codex 用入口ではハイフン名に対応する。CLI/IDE では `$git-commit-push` などの名前で
指定できる。

| Codex のスキル名 | 使う場面 | 手順本体 |
| --- | --- | --- |
| `git-commit-push` | 依頼されたコミット・push | [git_commit_push](.agent/skills/git_commit_push/SKILL.md) |
| `scaffold-new-tool` | Go CLI の追加 | [scaffold_new_tool](.agent/skills/scaffold_new_tool/SKILL.md) |
| `diagnose-command-failure` | 先頭バイト欠落・文字化けの診断 | [diagnose_command_failure](.agent/skills/diagnose_command_failure/SKILL.md) |

## 維持と動作確認

- 既存手順を更新するときは手順本体を更新し、Codex 入口へ本文をコピーしない。
  追加・改名・削除時には対応表と `.agents/skills/` の入口も合わせる。
- 新しい Codex スキルには `name` と `description` を持つ `SKILL.md` を置く。
  相対リンクは入口からリポジトリルートへ `../../../` で戻ることを確認する。
  Windows の通常の clone で扱えるファイルとし、シンボリックリンクを前提にしない。
- 変更後は `mise run lint:docs`、入口と正本の対応、リンク先の実在を検証する。
  Codex のスキル一覧で反映されない場合はセッションを再起動する。

配置仕様は OpenAI 公式の AGENTS.md 探索規則とスキル探索規則に基づく（2026-09-12
時点で確認）。`AGENTS.override.md` は同じ階層の `AGENTS.md` より優先されるため、
本対応では追加しない。

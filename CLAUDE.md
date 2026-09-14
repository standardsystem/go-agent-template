@AGENTS.md

# Claude Code 固有の補足

プロジェクト共通の正典は上でインポートした [AGENTS.md](AGENTS.md)。ここには
Claude Code だけに関係する対応付けを書く（共通ルールを二重に書かない）。

## 一時ファイルの置き場

リポジトリ内の一時ディレクトリは **`./temp` だけ**（正典は AGENTS.md の BASE RULES）。
セッション用スクラッチパッドではなくリポジトリ内を使うのは次のもの:

- チケット対応・調査の作業ファイル、添付ダウンロード → `./temp/<作業ID>/`
- 解析結果・他チームへの受け渡し成果物 → `./temp` または `./output`
- 受領した ZIP 等の**展開物**は一時ファイルではないので `./data/_extracted/<受領名>/`

会話中しか使わない捨てファイル（試行錯誤のワンライナー等）はスクラッチパッドでよい。
`temp/` は作業完了時に「昇格（再利用する）」か「削除（一回限り）」を判断する。
3 フォルダ共通のサブフォルダ名は `<YYYYMMDD>_<チケット>_<内容>`（正典:
[docs/development/UNTRACKED_DIR_NAMING.md](docs/development/UNTRACKED_DIR_NAMING.md)）。

## 設定ファイルの分担

- [.claude/settings.json](.claude/settings.json) はチーム共有の設定（追跡する）。
  検証系コマンド（`go build` / `go test` / `mise run` 等）の許可だけを置く
- `.claude/settings.local.json` は個人の設定（Git 管理外）。個人パス・コネクタの許可・
  autoMode の追加許可はこちらに書く
- MCP サーバの設定（`.mcp.json`）は認証情報を含みうるため Git 管理外

## サブエージェント

- チームで共有するサブエージェント定義（`<名前>.md`。フロントマターに `name` /
  `description` / `tools`）は `.claude/agents/` に置いて追跡する。個人用の定義は
  `~/.claude/agents/` に置き、リポジトリには入れない
- Claude Code はユーザーの依頼なしにサブエージェントを起動しないため、
  **該当タスクを受けたらまず起動を提案し、承認を得てから起動する**
- 起動しない場合でも、定義ファイルに書かれたチェックリストは本体の作業に適用する

## 実装着手前の合意（Plan モード）

AGENTS.md「着手前に方針を共有し、合意を得ること」に該当する作業は、Plan モードで
方針・変更範囲・影響を提示し、承認を得てから編集に入る。調査・レポート作成は
合意を待たずに完遂してよい。

## コミット

- 件名の形式は [docs/development/COMMIT_CONVENTION.md](docs/development/COMMIT_CONVENTION.md)
  が正典（`<type>: <チケット> <件名>`）
- `Co-Authored-By:` トレーラを付ける（履歴から AI の関与を追えるようにする）
- 日本語メッセージは文字化け回避のため、`-m` ではなくファイル経由（`-F`）でコミットする
  （手順は [git_commit_push](.agent/skills/git_commit_push/SKILL.md)）

## 手順書（スキル）の所在

[.agent/skills/](.agent/skills/) は全エージェント共通の手順書の正本で、**Claude Code の
スキルとしては読み込まれない**（Claude Code が探すのは `.claude/skills/`）。名前を
指定して起動することはできないため、必要な手順は該当ファイルを読んで従う参照資料として
扱う。Claude Code から `/名前` で呼びたい手順だけ `.claude/skills/<名前>/SKILL.md` に
薄い入口を置き、本体へリンクする（本文を複製しない）。

## ファイル参照の書式

回答中でファイルに言及するときは、バッククォートではなくワークスペース相対の
Markdown リンクを使う（例: `[AGENTS.md](AGENTS.md)` と書く）。

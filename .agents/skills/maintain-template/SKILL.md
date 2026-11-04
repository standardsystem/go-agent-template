---
name: maintain-template
description: テンプレートの定期メンテを頼まれたときに使う。社内プロジェクトを俯瞰して取り込む要素を選び、ツール・Actions・同梱物の古い版を更新する手順。
---

# テンプレートの定期メンテ

社内のプロジェクトを俯瞰し、(1) テンプレートに取り込む要素を選び、(2) 古くなった
コンポーネントを更新する。調査は合意なしで最後まで進め、変更は方針を示して合意を
得てから行う（[AGENTS.md](../../../AGENTS.md)「作業の進め方」）。

作業ファイルは `temp/<YYYYMMDD>_template-maintenance/` に置く。コマンドは Git Bash
（Bash ツール）の書き方で示す。社内のリポジトリは読むだけにし、fetch・checkout もしない。

## 手順 1: 社内プロジェクトを俯瞰する

GitHub の org と手元のクローンの両方を見る。手元は古いことがあるので、比較は
`origin/main` で行う。

```bash
gh repo list standardsystem --limit 200 --json name,updatedAt,primaryLanguage,description
cd /c/Projects
for d in STANDARDSYSTEM/* NHK/*; do
  [ -d "$d/.git" ] || continue
  m=""
  for f in AGENTS.md .mise.toml .golangci.yml .markdownlint-cli2.yaml .gitattributes \
           .agents/skills .claude/settings.json .github/workflows .pre-commit-config.yaml; do
    [ -e "$d/$f" ] && m="$m $f"
  done
  echo "$(git -C "$d" log -1 --format=%cs) $d :$m"
done | sort -r
```

- 足場ファイルがそろっていて最近更新されたものが比較の対象。テンプレートから作った
  リポジトリは初回コミットの件名（「go-agent-template をもとに」等）で見分けられる
- テンプレートの由来元は n-df-load-chart（README 冒頭）。CI・mise タスクの新しい工夫は
  ここに先に入ることが多い

## 手順 2: 足場ファイルを比べる

```bash
export MSYS_NO_PATHCONV=1   # Git Bash が "origin/main:<path>" をパスとして書き換えるのを防ぐ
T=/c/Projects/STANDARDSYSTEM/go-agent-template
cd <比較先のリポジトリ>
for f in .mise.toml .github/workflows/ci.yml .golangci.yml .claude/settings.json \
         .gitattributes .markdownlint-cli2.yaml .gitignore; do
  echo "#### $f"; diff -u "$T/$f" <(git show "origin/main:$f" 2>&1)
done
```

あわせて、採用の広さを grep で数える（例: `govulncheck`、`actionlint`、`run_windows`、
SHA 固定の `uses: x@<40 桁>`、settings.json の `hooks` / `deny` / `ask`）。

**取り込みの判定**: README「テンプレートの育て方」のとおり、2 つ以上のプロジェクトで
同じ修正が要ったもの。1 件だけでも、社内調査（手順 3）かセキュリティ勧告の裏付けが
あれば候補にする。プロジェクト固有の値（URL・課題キー・クラウドのプロジェクト名）は
持ち込まない。

## 手順 3: 社内調査の成果を参照する

develop-method-research（エージェンティック開発の手法調査）の
`docs/research/index.md` から各報告書の「自社での扱い」節を読む。同リポジトリは本
テンプレートから作られているので、指摘はそのまま本テンプレートに当てはまることが多い。
報告書の確認日と `stale_after` を見て、古い指摘は割り引く。

## 手順 4: 古いコンポーネントを洗い出す

| 対象 | 置き場所 | 最新の確かめ方 |
| --- | --- | --- |
| mise のツール | [.mise.toml](../../../.mise.toml) | `mise latest <tool>`（例: `mise latest npm:markdownlint-cli2`） |
| Go | `.mise.toml` と [go.mod](../../../go.mod) | `mise latest go`。上げるなら golangci-lint のリリースノートで対応版（`go1.NN support`）を確かめる |
| Node | `.mise.toml` | `https://nodejs.org/dist/index.json` の `lts` が付いた最新の偶数版（Active LTS） |
| Python | `.mise.toml` | バグ修正期間（リリースから 2 年）内の最新のマイナー版。`mise latest python@3.NN` |
| GitHub Actions | [ci.yml](../../../.github/workflows/ci.yml) | `gh api repos/<owner>/<repo>/releases/latest --jq .tag_name` |
| mise 本体（CI） | `ci.yml` の `version:` | 手元の `mise --version` で `mise run check` が通った版にする |
| yomiyasu | `.agents/skills/yomiyasu/` | [YOMIYASU.md](../../../docs/development/YOMIYASU.md) の「更新手順」 |
| Ponytail | AGENTS.md の Ponytail 節 | [PONYTAIL.md](../../../docs/development/PONYTAIL.md) の「更新手順」 |

- メジャー版・0.x のマイナー版を上げるときは、リリースノートの破壊的変更を読み、本
  テンプレートでの影響の有無を報告に書く
- Actions はタグでなくコミット SHA で固定する。SHA は次で引く。`object.type` が `tag`
  （注釈付きタグ）なら、`gh api repos/<owner>/<repo>/git/tags/<sha> --jq .object.sha` で
  コミットまでたどる

```bash
gh api repos/<owner>/<repo>/git/refs/tags/<tag> --jq '.object.type+" "+.object.sha'
```

## 手順 5: 報告して合意を得る

次の 3 つの表で報告し、Plan モード（Codex は会話）で合意を得る。

- 更新表: 対象・現在・更新後・根拠（破壊的変更の有無）
- 取り込み表: 要素・採用しているプロジェクト・変更内容
- 見送り表: 要素・理由（下の「検討候補」も見直す）

利用者の判断が要る点（規約の変更、権限の強化など）は選択肢にして聞く。

## 手順 6: 適用する

1. `main` から `chore/template-maintenance-<YYYYMMDD>` を作る（チケットがあれば規約どおり）
2. 合意した変更だけを 1 つずつ入れる。版を上げたら、文書中の版数（README の表、
   ci.yml のヘッダ、`.mise.toml` のコメント）も合わせる
3. `.agents/skills/` の追加・改名をしたら AGENTS.md の手順書の表も直す

## 手順 7: 検証する

```bash
mise install            # 新しい版を入れる
mise run info           # 版を確かめる
mise run check          # vet + lint (Go / Markdown / Actions) + vuln + test
claude plugin validate ./.agents   # yomiyasu を更新したとき
git ls-files --eol | awk '($1 == "i/crlf" || $1 == "i/mixed") && $3 != "attr/-text"'
```

markdownlint の版上げで新しいルールの違反が出たら本文を直す。設定で無効にするのは、
日本語の書き方と構造的に合わないときだけにし、理由を `.markdownlint-cli2.yaml` に書く。
コミット・push・PR は依頼を受けてから [git-commit-push](../git-commit-push/SKILL.md) で行う。

## 報告

読んだリポジトリと資料、読めなかったもの、採否とその根拠、検証結果、派生
プロジェクトへの波及（取り込みは各プロジェクトが README の `git diff --no-index` で行う）を
報告する。最後に、下の「検討候補」を更新したかを書く。

## 検討候補（毎回見直す）

前回までに見送ったもの。採用プロジェクトが増えたら取り込みを検討する。

| 候補 | 出どころ | 見送りの理由（最終確認 2026-11-04） |
| --- | --- | --- |
| mise の複数行 `run` を CI で検出 | n-df-load-chart | 1 件のみ。テンプレートは配列形式をコメントで規定済み |
| Go の対象から `temp/` を外す（`./...` を `./cmd/... ./internal/...` に） | n-df-load-chart | 1 件のみ |
| Windows で `test` から `-race` を外す（`run_windows`） | n-df-load-chart | 1 件のみ。`test:short` で代替できる |
| INTERNAL_LIBS.md と `internal/` の整合検査 | n-df-load-chart | 1 件のみ |
| 課題管理（自社 Issue Tracker・MCP）の節を AGENTS.md に | develop-method-research、keiei | 固有の値が主で、AGENTS.md が 20 KB を超えるため |
| settings.json の `ask`（commit / push）・`deny`（`.mcp.json` の読み取り） | 社内調査 DEV_RSCH-6・7 | 採用プロジェクトが無い。2026-11-04 に利用者の判断で見送り |
| Stop フックで lint を回す | 社内調査 DEV_RSCH-9 | 採用プロジェクトが無い |
| AGENTS.md の縮減（Claude Code の節を `.claude/rules/` へ） | 社内調査 DEV_RSCH-2 | 構成の変更が大きく、定期メンテの範囲外 |
| Office の一時ファイル `~$*` を `.gitignore` に | minamiurawa | 1 件のみ |

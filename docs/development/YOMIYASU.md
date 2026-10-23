# yomiyasu の導入（日本語推敲スキル）

[yomiyasu](https://github.com/nanaism/yomiyasu)（MIT）は、AI が生成した日本語の不自然な
比喩、曖昧な主述関係、不要な装飾を、意味と条件を保ったまま直す Agent Skill。本書は、
本テンプレートへの取り込み方と更新手順を記す。使い方の正典は同梱した
[SKILL.md](../../.agents/skills/yomiyasu/SKILL.md) で、ここには複製しない。

取り込み元: 上流 v2.2.0（commit `a34d94ff`、2026-10-23）の `skills/yomiyasu/` を無改変で同梱。

## 構成（正本は 1 か所、Claude Code と Codex の両方が読む）

```text
.agents/
├── .claude-plugin/marketplace.json   Claude Code 用のローカルマーケットプレイス定義
└── skills/
    └── yomiyasu/                     上流 skills/yomiyasu/ の無改変コピー (正本)
```

- **Codex**: `.agents/skills/` を自動で読む。呼び出しは `$yomiyasu:yomiyasu`
- **Claude Code**: `.agents/skills/` は読まないが、`.agents/` をマーケットプレイスの
  ルートとして登録すると、プラグイン `yomiyasu` として読む。呼び出しは
  `/yomiyasu:yomiyasu`。`marketplace.json` のエントリで `skills` を
  `./skills/yomiyasu` に絞っているので、ほかの手順書はプラグインに含まれない。
  ファイルはキャッシュに写さず、その場から読む
- シンボリックリンクと複製は使わない（[AGENTS.md](../../AGENTS.md)「手順書（スキル）」）

[.claude/settings.json](../../.claude/settings.json) の `enabledPlugins` で有効化を共有する。
マーケットプレイスの登録は各メンバーが `--scope local` で行う。CLI は登録に絶対パスを
書くので、共有設定に入れるとほかの人の環境で解決できない（Claude Code 2.1.308 で確認）。

## 導入手順（各メンバー 1 回、VS Code の統合ターミナルで）

`claude` が PATH に無い環境では、VS Code 拡張が同梱する CLI を使う（拡張の版が上がると
パスの版番号も変わる）。ドライブ文字は小文字で `cd` してから実行する（理由は
[PONYTAIL.md](PONYTAIL.md)「導入手順」）。

```powershell
$claude = "$env:USERPROFILE\.vscode\extensions\anthropic.claude-code-2.1.308-win32-x64\resources\native-binary\claude.exe"
cmd /c "cd /d c:\Projects\STANDARDSYSTEM\go-agent-template && $claude plugin marketplace add ./.agents --scope local && $claude plugin install yomiyasu@go-agent-template --scope project"
```

パスは自分のクローン先に読み替える。1 つ目のコマンドは
`.claude/settings.local.json`（Git 管理外）に登録を書き、2 つ目は共有設定の
`enabledPlugins` を使う。`git diff` で `.claude/settings.json` に差分が出ていないことを
確認する。`claude plugin details yomiyasu` に `Skills (1)  yomiyasu` と出れば登録できている。
VS Code では「Developer: Reload Window」で再起動してから、新しい会話で
`/yomiyasu:yomiyasu` を使う。

## 使い方

- 推敲を頼むときは名前を出す（例: 「`/yomiyasu:yomiyasu` でこの文書を推敲して」）。
  名前を出さない文章作成やコード・設定の更新では、スキルは働かない
- 同梱スクリプトは標準ライブラリだけで動く。上流の手順どおり `python3` で実行する
  （mise で固定した Python を使う。`uv run` は不要）
- 検査記録は git の最上位にある `.yomiyasu/` に保存される。原文と最終稿の全文が
  入るので `.gitignore` の対象にしている。個人情報を含む文書では `--no-record` を頼む
- フック（`scripts/yomiyasu_hook.py`）は登録しない。必要になったら、上流の
  `--print-config` が出す設定を確かめ、合意を取ってから個人設定に入れる

## 更新手順

1. 上流の新しいタグを `temp/<YYYYMMDD>_yomiyasu更新/` に clone し、`CHANGELOG.md` と
   `skills/yomiyasu/scripts/` の差分を読む（外部通信や外部プロセスの起動が増えていないか）
2. `.agents/skills/yomiyasu/` を削除し、新しい `skills/yomiyasu/` をそのままコピーする
3. [marketplace.json](../../.agents/.claude-plugin/marketplace.json) の `version` と、
   本書冒頭の取り込み元を上流に合わせる
4. `claude plugin validate ./.agents` と `mise run lint:docs` を通す。同梱物は
   [.markdownlint-cli2.yaml](../../.markdownlint-cli2.yaml) で lint の対象から外している
   （上流の書式を直すと、更新のたびに差分が出るため）

## 無効化

- 個人で止める: `.claude/settings.local.json` の `enabledPlugins` に
  `"yomiyasu@go-agent-template": false` を書く
- テンプレートから外す: `.agents/skills/yomiyasu/`、`.agents/.claude-plugin/`、
  `.claude/settings.json` の `enabledPlugins` の行、AGENTS.md の表の行、
  `.markdownlint-cli2.yaml` の除外を消す

## ライセンス

コードと文書は MIT（[LICENSE](../../.agents/skills/yomiyasu/LICENSE)）。同梱の Unicode
データは [UNICODE-LICENSE.txt](../../.agents/skills/yomiyasu/UNICODE-LICENSE.txt) に従う。
どちらも同梱物の中に残してある。

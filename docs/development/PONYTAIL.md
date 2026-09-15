# Ponytail の導入（実装の最小化規範）

[Ponytail](https://github.com/DietrichGebert/ponytail) は、AI エージェントに「怠け者の
シニア開発者」の判断順序（YAGNI → 既存コード → 標準ライブラリ → プラットフォーム機能 →
導入済み依存 → 1 行 → 最小実装）を常時適用させるプロンプト集。本書は、本テンプレートへ
どの範囲で取り込んだか、なぜか、どう更新するかを記す。ルール本文は
[AGENTS.md](../../AGENTS.md)「実装の最小化規範（Ponytail）」が正典で、ここには複製しない。

取り込み元: 上流 v4.10.0（commit `e3ba2aa6`、2026-09-14）の `AGENTS.md`（簡約版ルール）。

## 評価（2026-09-15 時点）

### 意義

- 本テンプレートの用途（Go の CLI・解析ツールをエージェントと共同開発）と方向が一致する。
  Go は標準ライブラリが厚く、「依存を足す前に標準で書く」「`internal/` を先に引く」は
  既存規約（[内部ライブラリ索引](INTERNAL_LIBS.md) を最初に参照）と同じ判断を求めている
- 上流の計測（Haiku 4.5、FastAPI + React の実リポジトリ、12 タスク、n=4）では、スキル
  無しに対して差分行数 -54%、トークン -22%、コスト -20%、時間 -27%、安全性（検証・
  エラー処理の維持）100%。Go の CLI で同じ値になる保証はないが、方向は再現すると見込む
- 「バグ修正は呼び出し元を全部 grep して共有関数を 1 か所で直す」「理解を省いた小さな
  差分は 2 つ目のバグ」は、エージェントが起こしやすい失敗の言語化として単独でも有用

### 本テンプレートの規約と衝突する点

AGENTS.md の「本プロジェクトの規約との優先順位」で解決している。

| 上流の記述 | 本テンプレートの規約 | 解決 |
| --- | --- | --- |
| 説明はコードの後に 3 行まで | 読んだ資料・前提・検証結果・未確認点を日本語で報告 | 報告は本テンプレートの規約が優先。上流も「明示的に求められた説明は省かない」としている |
| 不要と判断したら作らず 1 行で言う | 依頼範囲を勝手に縮小しない、変更は合意の上で | 依頼範囲を削る提案は着手前に出して合意を得る。依頼外の抽象・依存・雛形を省くのは自由 |
| テストは最小の検証 1 つ。フレームワーク・フィクスチャ無し | AAA 形式・`t.Logf`・テスト対象情報 | 「1 つは残す」を下限とし、書き方は [TESTING_STANDARDS.md](TESTING_STANDARDS.md) |
| `ponytail:` コメントで上限と引き上げ経路を残す | コメントは日本語 | 接頭辞 `ponytail:` は英語のまま（`/ponytail-debt` の grep 対象）、本文は日本語 |

### リスク

- **プラグインは第三者の Node.js フックを実行する**。v4.10.0 の `hooks/*.js` を読んだ
  範囲では、ネットワーク通信・外部プロセス起動は無く、書き込み先はホームディレクトリ
  配下（`~/.claude/.ponytail-active`、`~/.claude/.ponytail-statusline-nudged`、
  `%APPDATA%\ponytail\config.json`）だけ。ただし版が上がれば内容は変わりうる。
  単独メンテナのリポジトリである点も含め、更新は内容を確認してから取り込む
- 初回セッションで「ステータスラインを設定しないか」という提案が 1 度だけ出る
  （`~/.claude/settings.json` への `statusLine` 追加）。合意ゲートに従い、承認した場合
  だけ設定する
- 常時ルール（AGENTS.md）とプラグインの注入が両方効くと、同じ規範が 2 度読み込まれる。
  内容は同一なので害は無いが、`/ponytail off` で止まるのはプラグイン側だけで、AGENTS.md
  の規範はプロジェクト規約として残る

## 採用範囲（2 層）

| 層 | 何を | 対象 | 実体 |
| --- | --- | --- | --- |
| 1. 常時ルール | 上流 `AGENTS.md`（簡約版）を日本語化し、優先順位を付けたもの | 全エージェント（AGENTS.md を読む Claude Code / Codex / Gemini CLI / Cursor 等） | [AGENTS.md](../../AGENTS.md)「実装の最小化規範（Ponytail）」 |
| 2. プラグイン | `/ponytail-review`（差分の過剰設計レビュー）、`/ponytail-audit`（リポジトリ全体）、`/ponytail-debt`（`ponytail:` コメントの台帳化）、強度切り替え `/ponytail lite\|full\|ultra\|off` | Claude Code | [.claude/settings.json](../../.claude/settings.json) の `extraKnownMarketplaces` / `enabledPlugins` ＋ 各メンバーの CLI 導入（下記） |

層 2 の `enabledPlugins` は「このプロジェクトで使うプラグイン」の宣言であり、GitHub 配布の
プラグインはこれだけでは導入されない。Claude Code 2.1.270 では、ネットワーク上の
マーケットプレイスはユーザー設定か管理設定で宣言されたものだけを信頼し、プロジェクト
設定の `extraKnownMarketplaces` だけでは登録されない（実行ファイル内のメッセージ
「project/local scope cannot vouch for it」で確認。公式ドキュメントも「メンバーが
インストールするまで読み込まれない」としている）。未導入のまま起動すると
「Plugin "ponytail@ponytail" is enabled in project settings but isn't installed」の診断が
出る。

### 導入手順（各メンバー 1 回、VS Code の統合ターミナルで）

`node` が PATH に必要（[.mise.toml](../../.mise.toml) の Node で満たす）。`claude` が
PATH に無い環境では、VS Code 拡張が同梱する CLI を使う（拡張の版が上がるとパスの版番号も
変わる）。

**ドライブ文字の大小に注意**: VS Code はセッションの作業ディレクトリを `c:\...`（小文字）
で渡すが、PowerShell から起動した CLI は `C:\...`（大文字）で導入記録を書く。Claude Code
2.1.270 は project スコープの導入記録と作業ディレクトリを文字列で厳密比較するため、
大文字で記録すると VS Code のセッションでは「Plugin "ponytail" not cached」と誤判定されて
読み込まれない（2026-09-15 に再現確認）。導入コマンドは `cmd /c` で小文字パスに `cd`
してから実行する:

```powershell
$claude = "$env:USERPROFILE\.vscode\extensions\anthropic.claude-code-2.1.270-win32-x64\resources\native-binary\claude.exe"
& $claude plugin marketplace add DietrichGebert/ponytail
cmd /c "cd /d c:\Projects\STANDARDSYSTEM\go-agent-template && $claude plugin install ponytail@ponytail --scope project"
```

2 行目はマーケットプレイスをユーザー設定に登録する（作業ディレクトリは無関係）。3 行目の
パスは自分のクローン先に読み替え、ドライブ文字は小文字にする。`--scope project` は
`.claude/settings.json` の `enabledPlugins` を使うので、リポジトリ側の差分は出ない想定
（出たら `git diff` で内容を確認する）。既に大文字で導入済みの場合も 3 行目だけ再実行
すればよく、`claude plugin list` で記録が増えたことを確認する。

その後、コマンドパレットの「Developer: Reload Window」で再起動し、新規会話で `/` を
打って `/ponytail:ponytail-help` が一覧に出れば完了（プラグインのコマンドはプラグイン名の
名前空間付きで並ぶ）。セッション内なら `/reload-plugins` でも読み直せる。VS Code 拡張は
コマンド一覧を起動時に読み込んで保持するので、導入直後の会話には反映されない。

Codex で同じコマンドが要る場合は各自で `codex plugin marketplace add DietrichGebert/ponytail`
と `codex plugin add ponytail@ponytail` を実行する（ユーザー単位の設定で、リポジトリには
入れない）。

## 使い方

- 通常の実装依頼では何もしなくてよい。層 1 が常に効く
- 差分を「削れるところは無いか」の観点だけで見たいときは `/ponytail-review`。正しさ・
  セキュリティ・性能は対象外なので、通常のレビュー（`/code-review`）と併用する
- 意図的に手を抜いた箇所の棚卸しは `/ponytail-debt`（`grep -rnE '(#|//) ?ponytail:'` と
  同じ）

## 更新手順

1. 前回取り込み以降の上流の変更を見る:
   `https://github.com/DietrichGebert/ponytail/compare/e3ba2aa6...main`
2. 上流 `AGENTS.md` に変更があれば、AGENTS.md の該当節へ日本語で反映し、本書冒頭の版と
   commit を更新する
3. プラグイン側（`hooks/*.js`）の差分も同時に読み、ネットワーク通信や外部プロセス起動が
   増えていないか確認する。プラグインは上流 `plugin.json` の `version` が上がったときに
   配信される

## 無効化

- 層 2 だけ外す: [.claude/settings.json](../../.claude/settings.json) から
  `extraKnownMarketplaces.ponytail` と `enabledPlugins."ponytail@ponytail"` を消す。
  個人で外すだけなら `.claude/settings.local.json` に
  `"enabledPlugins": {"ponytail@ponytail": false}` を書く
- 層 1 も外す: AGENTS.md の該当節を削除し、本書と [README.md](../../README.md)・
  [docs/README.md](../README.md) の行も消す

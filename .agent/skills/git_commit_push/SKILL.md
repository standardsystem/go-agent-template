---
name: git_commit_push
description: 依頼された範囲のコミットと push を、コミット規約と合意ゲートに従って実行する。日本語メッセージはファイル経由で渡す。
---

# コミットと push

## 前提（合意ゲート）

- 実行するのは**依頼された単一の操作だけ**。「push して」と言われて未コミットなら、
  コミットせずに「未コミットです」と報告して止まる（[AGENTS.md](../../../AGENTS.md)）
- ステージングは対象ファイルを明示する。`git add .` / `git add -A` は別作業の変更を
  巻き込むので使わない
- 作業ツリーは他のセッションと共有されていることがある。コミット直前に
  `git branch --show-current` で意図したブランチにいることを確認する
- メッセージの形式は [コミット規約](../../../docs/development/COMMIT_CONVENTION.md)
  （`<type>: <チケット> <件名>`）

## 手順 1: 状態確認

```powershell
git status
git branch --show-current
git diff --stat
```

## 手順 2: 対象を明示してステージング

```powershell
git add path/to/file1 path/to/file2
git diff --cached --stat
```

## 手順 3: メッセージを temp/ の UTF-8（BOM なし）ファイルに書く

エージェントのファイル書き込みツールで `temp/commit_msg.txt` を作るのが最も確実
（UTF-8・BOM なし・LF）。PowerShell で作る場合は here-string を使い、終端の `'@` を
行頭に置く。

```powershell
@'
feat: PROJ-123 ○○を追加

変更の背景と要点を 2〜3 行で。

Co-Authored-By: <エージェント規定の名前> <エージェント規定のアドレス>
'@ | Set-Content -Path temp/commit_msg.txt -Encoding utf8NoBOM
```

`Co-Authored-By` の値は各エージェントの規定値を使う（Claude 用の値を他のエージェントが
流用しない）。

## 手順 4: ファイル経由でコミット

```powershell
git commit -F temp/commit_msg.txt
Remove-Item temp/commit_msg.txt
git log -1 --format="%h %s"
```

`git commit -m "日本語…"` は**使用禁止**。先頭バイト欠落・文字化けの実績がある
（[diagnose_command_failure](../diagnose_command_failure/SKILL.md)）。

## 手順 5: push（依頼されたときだけ）

```powershell
git push -u origin (git branch --show-current)
```

リモートが先行して拒否されたら、独断で `pull --rebase` や `--force` をせず状況を報告する。

## 報告

対象ブランチ・コミットハッシュ・push の有無を報告する。

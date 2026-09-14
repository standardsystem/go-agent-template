---
name: git-commit-push
description: 現在の変更についてコミットや push を依頼されたときに使う。日本語メッセージと対象変更の確認手順を参照する。
---

# コミットと push（Codex 用入口）

[Codex 向け案内](../../../CODEX.md) の参照順と読み替えを適用し、
[手順本体](../../../.agent/skills/git_commit_push/SKILL.md) を実際に開いてから作業する。
手順内の相対リンク・補助ファイルは手順本体の場所を基準に解決する。

push だけの依頼で未コミットなら、その状態を報告しコミットを補わない。
`Co-Authored-By` は Codex 用の値を使い、Claude 用の値を流用しない。

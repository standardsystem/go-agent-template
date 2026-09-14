# ドメインナレッジの書き方

## 置き場と階層

- [index.md](index.md) が索引。本体は `docs/knowledge/<テーマ>/<文書>.md` または直下に置く
- 索引には 1 行の要約とリンクだけを書き、本文を複製しない
- エージェントのローカルメモリで得た知見は、個人情報・認証情報・社外秘を除いてから
  ここへ昇格する（[AGENTS.md](../../AGENTS.md)「ドメイン知識の蓄積先」）

## フロントマター

各本体の先頭に最低限 `type` を付ける。索引で当たりを付けるための項目なので、
`status` と `stale_after` も付けることを推奨する。

```yaml
---
type: knowledge          # knowledge | decision | pitfall | reference
title: 文書の題名
status: active           # active | stale | superseded
stale_after: 2027-03-31  # 見直し期限 (任意)
tags: [example]
---
```

## 昇格の手順

1. 本体を書く（背景・事実・根拠・未確認点を分けて書く。推測は推測と明記する）
2. [index.md](index.md) に 1 行追記する
3. `mise run lint:docs` を通す

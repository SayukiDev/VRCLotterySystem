# VRC Lottery System API ドキュメント

VRChat向け抽選システムのHTTP API仕様書です。

## ドキュメント構成

| ドキュメント                | 対象                                                       |
| --------------------------- | ---------------------------------------------------------- |
| **index.md**（本書）        | 共通仕様・共通レスポンス・データモデル・エンドポイント早見表 |
| [user.md](./user.md)        | 一般向けAPI（フォーム取得・応募・許可リスト等。認証不要）   |
| [admin.md](./admin.md)      | 管理者向けAPI（ブラックリスト・スタッフ・抽選操作。認証必須） |

## 概要

- ベースURL: `http://<host>:8080`
- ルートプレフィックス: `/api`（管理者APIは `/api/authed`）
- 通信フォーマット: JSON（一部エンドポイントはプレーンテキスト）
- HTTPフレームワーク: [Gin](https://github.com/gin-gonic/gin)
- リクエストボディ上限: 64 KiB（超過時は `413 request entity too large`）

## 認証

- 一般向けAPI（`/api/*`）: **認証なし**
- 管理者向けAPI（`/api/authed/*`）: **トークン認証**
  - リクエストヘッダ `Authorization` に、`config.json` の `token` と **完全一致** する値を指定します。
  - 一致しない場合は `401 unauthorized` を返します。
  - 詳細は [admin.md](./admin.md) を参照してください。

```
Authorization: <config.json の token>
```

## 共通レスポンス構造

ほとんどのエンドポイントは以下の共通構造（`common.CommonResp`）でレスポンスを返します。

```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

| フィールド | 型     | JSONキー | 説明                                              |
| ---------- | ------ | -------- | ------------------------------------------------- |
| `code`     | int    | `code`   | ステータスコード（HTTPステータスと同じ値）        |
| `msg`      | string | `msg`    | 処理結果メッセージ（`success` / `bad request` 等） |
| `data`     | any    | `data`   | レスポンス本体（エンドポイントごとに異なる）       |

> 例外として、`GET /api/getAllowList` のみ共通構造ではなく **プレーンテキスト** を返します。

### 共通エラーメッセージ

| HTTPステータス | `msg`                       | 発生条件                                       |
| -------------- | --------------------------- | ---------------------------------------------- |
| 400            | `bad request`               | リクエストボディ／クエリ不正・必須項目未入力     |
| 401            | `unauthorized`              | 管理者APIで `Authorization` が不一致             |
| 404            | `not found`                 | 指定リソースが存在しない（例: フォーム `id` 不一致） |
| 413            | `request entity too large`  | リクエストボディが 64 KiB を超過                |
| 500            | `internal server error`     | サーバ内部処理エラー                           |

---

## データモデル

### SiteData（サイト情報）

サイト全体の表示情報を表します。

| フィールド   | 型     | JSONキー     | 説明                          |
| ------------ | ------ | ------------ | ----------------------------- |
| `title`      | string | `title`      | サイトのタイトル（必須）      |
| `form_title` | string | `form_title` | フォームページのタイトル（必須） |
| `terms`      | string | `terms`      | 利用規約（Markdown、必須）    |

### FormItem（フォーム項目定義）

応募フォームの各項目を表します。

| フィールド  | 型         | JSONキー   | 説明                                              |
| ----------- | ---------- | ---------- | ------------------------------------------------- |
| `is_id`     | bool       | `is_id`    | この項目が応募者の一意ID（例: VRChat ID）かどうか |
| `title`     | string     | `title`    | 項目タイトル（必須）                              |
| `desc`      | string     | `desc`     | 項目の説明文                                      |
| `required`  | bool       | `required` | 入力必須かどうか                                  |
| `options`   | []string   | `options`  | 選択肢（`type` が `options` の場合に使用）        |
| `type`      | string     | `type`     | 項目タイプ。`content` / `input` / `options` のいずれか |

#### type の種類

| type      | 説明                        |
|-----------|---------------------------|
| `content` | 表示専用のテキスト（入力欄なし。説明文の表示など） |
| `input`   | 自由入力テキスト欄                 |
| `options` | 選択肢から選ぶ項目（`options` を併用）  |

### Input（応募データ）

応募者がフォームに入力した内容を表します。フォーム項目の **インデックス番号** をキーにします。

| フィールド | 型              | JSONキー   | 説明                                                      |
| ---------- | --------------- | ---------- | --------------------------------------------------------- |
| `content`  | map[int]string  | `content`  | `input` タイプ項目の回答。キーはフォーム項目のインデックス |
| `selected` | map[int][]int   | `selected` | `options` タイプ項目の回答。キーは項目インデックス、値は選択肢インデックスの配列 |

---

## エンドポイント早見表

### 一般向け（認証なし） — 詳細は [user.md](./user.md)

| メソッド | パス               | 説明                       |
| -------- | ------------------ | -------------------------- |
| GET      | `/api/getSiteData` | サイト情報の取得           |
| GET      | `/api/getForm`     | 応募フォーム定義の取得（クエリ `id` 必須） |
| POST     | `/api/submitForm`  | 応募フォームの送信         |
| GET      | `/api/getAllowList`| 入場許可リスト（当選者＋スタッフ）の取得 |
| GET      | `/api/isActive`    | 応募の受付状態の取得       |

### 管理者向け（認証必須） — 詳細は [admin.md](./admin.md)

| メソッド | パス                          | 説明                       |
| -------- | ----------------------------- | -------------------------- |
| GET      | `/api/authed/getBlackList`    | ブラックリストの取得       |
| POST     | `/api/authed/addBlackList`    | ブラックリストへ追加       |
| POST     | `/api/authed/deleteBlackList` | ブラックリストから削除     |
| GET      | `/api/authed/getStaffList`    | スタッフリストの取得       |
| POST     | `/api/authed/addStaff`        | スタッフの追加             |
| POST     | `/api/authed/deleteStaff`     | スタッフの削除             |
| POST     | `/api/authed/setDrawing`      | 抽選の設定（締切・当選数）   |
| GET      | `/api/authed/getDrawing`      | 現在の抽選状態の取得       |
| POST     | `/api/authed/drawing`         | 抽選の実行                 |
| GET      | `/api/authed/getResults`      | 当選者リストの取得         |
| POST     | `/api/authed/removeResults`   | 当選者の削除               |

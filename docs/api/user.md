# 一般向けAPI

応募者・フロントエンドが利用する **認証不要** のAPIです。
共通仕様・データモデル（`SiteData` / `FormItem` / `Input`）は [index.md](./index.md) を参照してください。

- ルートプレフィックス: `/api`
- 認証: なし

## エンドポイント一覧

| メソッド | パス               | 説明                       |
| -------- | ------------------ | -------------------------- |
| GET      | `/api/getSiteData` | サイト情報（タイトル・フォームタイトル・利用規約）の取得 |
| GET      | `/api/getForm`     | 応募フォーム定義の取得（クエリ `id` 必須） |
| POST     | `/api/submitForm`  | 応募フォームの送信         |
| GET      | `/api/getAllowList`| 入場許可リスト（当選者＋スタッフ）の取得 |
| GET      | `/api/isActive`    | 応募の受付状態（受付中かどうか）の取得 |

---

## GET /api/getSiteData

サイト情報（サイトタイトル・フォームページタイトル・利用規約）を取得します。

### リクエスト

パラメータなし。

### レスポンス（200）

`data` に `SiteData` オブジェクトが入ります。`terms` はMarkdown文字列です。

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "title": "サンプルイベント",
    "form_title": "抽選予約フォーム",
    "terms": "本イベントの抽選参加にあたっての規約文..."
  }
}
```

### 使用例

```bash
curl http://localhost:8080/api/getSiteData
```

---

## GET /api/getForm

応募フォームの項目定義一覧を取得します。フロントエンドはこの定義をもとにフォームを描画します。

### リクエスト

クエリパラメータで `id`（現在有効な抽選のフォームを識別するランダム文字列。応募者IDではない）を指定します。
`id` は `setDrawing` の実行ごとに新しく採番され、サーバが保持している現在の値と一致する必要があります。

| パラメータ | 型     | 必須 | 制約       | 説明                                   |
| ---------- | ------ | ---- | ---------- | -------------------------------------- |
| `id`       | string | ○    | 最大10文字 | フォームを識別するランダム文字列（応募者IDではない） |

```
GET /api/getForm?id=Ab3xK9pQ2
```

### レスポンス

#### 成功（200）

`data` に `FormItem` の配列が入ります。
配列の **インデックス番号** が、応募送信時（`/api/submitForm`）のキーに対応します。

```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "is_id": false,
      "title": "注意事項",
      "desc": "応募前に必ずお読みください",
      "required": false,
      "options": null,
      "type": "content"
    },
    {
      "is_id": true,
      "title": "VRChat ID",
      "desc": "",
      "required": true,
      "options": null,
      "type": "input"
    },
    {
      "is_id": false,
      "title": "参加希望日",
      "desc": "",
      "required": true,
      "options": ["1日目", "2日目"],
      "type": "options"
    }
  ]
}
```

#### `id` 未指定／10文字超過（400）

```json
{
  "code": 400,
  "msg": "bad request",
  "data": "Key: 'GetFormReq.Id' Error:Field validation for 'Id' failed on the 'required' tag"
}
```

`data` には具体的なバリデーションエラー内容が入ります。

#### `id` がサーバ保持値と不一致（404）

形式は正しいが、現在有効な抽選のフォームIDと一致しない場合。

```json
{
  "code": 404,
  "msg": "not found",
  "data": null
}
```

### 使用例

```bash
curl "http://localhost:8080/api/getForm?id=Ab3xK9pQ2"
```

---

## POST /api/submitForm

応募フォームを送信します。

### リクエスト

- Content-Type: `application/json`

| フィールド | 型    | 必須 | 説明                       |
| ---------- | ----- | ---- | -------------------------- |
| `inputs`   | Input | ○    | 応募内容（[Input](./index.md#inputデータ) を参照） |

`content` / `selected` のキーは、`/api/getForm` で取得した項目配列のインデックスに対応します。

```json
{
  "inputs": {
    "content": {
      "1": "usr_xxxxxxxx"
    },
    "selected": {
      "2": [0]
    }
  }
}
```

上記の例は、`getForm` のインデックス `1`（VRChat ID）に `usr_xxxxxxxx` を入力し、
インデックス `2`（参加希望日）で先頭の選択肢（`1日目`）を選択した応募を表します。

### バリデーション

サーバ側で以下を検証します。

- リクエストボディが正しいJSON形式であること
- `required: true` の項目がすべて入力されていること
  - `type: input` の項目 → `content` に対応するキーが存在すること
  - `type: options` の項目 → `selected` に対応するキーが存在すること

加えて、応募の登録処理（プロバイダ層）で以下も検証されます。

- 募集期間が終了していないこと（期限切れの場合はエラー）
- 応募者IDがブラックリストに含まれていないこと（含まれる場合は登録されず、エラーにもならず黙って無視されます）
- 同一IDが二重応募でないこと

### レスポンス

#### 成功（200）

```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 必須項目未入力・JSON不正（400）

```json
{
  "code": 400,
  "msg": "bad request",
  "data": "required field(VRChat ID) not filled"
}
```

`data` には具体的なエラー内容が入ります（例: `required field(<項目タイトル>) not filled`、JSONパースエラーの内容など）。

#### サーバ内部エラー（500）

応募期限切れ・二重応募・保存失敗時など。

```json
{
  "code": 500,
  "msg": "internal server error",
  "data": "input already exists"
}
```

| `data` の値          | 意味                   |
| -------------------- | ---------------------- |
| `date expired`       | 募集期間が終了している |
| `input already exists` | 同一IDですでに応募済み |

### 使用例

```bash
curl -X POST http://localhost:8080/api/submitForm \
  -H "Content-Type: application/json" \
  -d '{
    "inputs": {
      "content": { "1": "usr_xxxxxxxx" },
      "selected": { "2": [0] }
    }
  }'
```

---

## GET /api/getAllowList

入場を許可するIDの一覧を取得します。**当選者リスト** と **スタッフリスト** を結合した結果を返します。

> このエンドポイントのみ、共通レスポンス構造（JSON）ではなく **プレーンテキスト** を返します。

### リクエスト

パラメータなし。

### レスポンス（200）

- Content-Type: `text/plain`
- 当選者IDとスタッフIDが **改行（`\n`）区切り** で並んだテキスト

```
usr_aaaaaaaa
usr_bbbbbbbb
usr_cccccccc
```

> ID の並び順は内部マップの列挙順に依存するため、保証されません。

### 使用例

```bash
curl http://localhost:8080/api/getAllowList
```

---

## GET /api/isActive

応募の **受付状態** を取得します。現在時刻が募集締切（`setDrawing` で設定した `date`）より前であれば受付中（`true`）、締切を過ぎていれば受付終了（`false`）を返します。

フロントエンドはこのエンドポイントで、応募ページを開く前に「受付中／受付終了」を判定し、`submitForm` を送信して `date expired` エラーで弾かれる前にUIを切り替えられます。

### リクエスト

パラメータなし。

### レスポンス（200）

`data` に受付状態を表す真偽値（bool）が入ります。

| `data` の値 | 意味                          |
| ----------- | ----------------------------- |
| `true`      | 受付中（現在時刻 < 締切日時） |
| `false`     | 受付終了（締切日時を過ぎた）  |

```json
{
  "code": 200,
  "msg": "success",
  "data": true
}
```

### 使用例

```bash
curl http://localhost:8080/api/isActive
```

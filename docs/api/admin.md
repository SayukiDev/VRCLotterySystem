# 管理者向けAPI

抽選システムの管理操作（ブラックリスト・スタッフ・抽選）を行うAPIです。
共通仕様・データモデルは [index.md](./index.md) を参照してください。

- ルートプレフィックス: `/api/authed`
- すべてのエンドポイントで **トークン認証が必須** です。

## 認証

リクエストヘッダ `Authorization` に、`config.json` の `token` と **完全一致** する文字列を指定します。

```
Authorization: <config.json の token>
```

一致しない場合、すべてのエンドポイントで以下を返します。

```json
{
  "code": 401,
  "msg": "unauthorized",
  "data": null
}
```

> 注意: トークンは前置詞なし（`Bearer ` などは付けない）の **生の値** をそのまま送信します。

---

## エンドポイント一覧

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

---

## ブラックリスト

ブラックリストに登録されたIDは、`/api/submitForm` での応募時に黙って無視されます（エラーにはなりません）。

### GET /api/authed/getBlackList

ブラックリストに登録されているIDの一覧を取得します。

#### リクエスト

パラメータなし。

#### レスポンス（200）

`data` にID文字列の配列が入ります。登録がない場合は空配列 `[]`。

```json
{
  "code": 200,
  "msg": "success",
  "data": ["usr_aaaaaaaa", "usr_bbbbbbbb"]
}
```

> 並び順は内部マップの列挙順に依存するため保証されません。

#### 使用例

```bash
curl http://localhost:8080/api/authed/getBlackList \
  -H "Authorization: <token>"
```

---

### POST /api/authed/addBlackList

ブラックリストにIDを追加します。

#### リクエスト

クエリパラメータで `id` を指定します。

| パラメータ | 型     | 必須 | 制約       | 説明           |
| ---------- | ------ | ---- | ---------- | -------------- |
| `id`       | string | ○    | 最大10文字 | 追加するID     |

```
POST /api/authed/addBlackList?id=usr_xxxx
```

#### レスポンス

##### 成功（200）

```json
{ "code": 200, "msg": "success", "data": null }
```

##### `id` 未指定／10文字超過（400）

```json
{
  "code": 400,
  "msg": "bad request",
  "data": "Key: 'AddBlackListReq.Id' Error:Field validation for 'Id' failed on the 'required' tag"
}
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/addBlackList?id=usr_xxxx" \
  -H "Authorization: <token>"
```

---

### POST /api/authed/deleteBlackList

ブラックリストからIDを削除します。

#### リクエスト

| パラメータ | 型     | 必須 | 制約       | 説明           |
| ---------- | ------ | ---- | ---------- | -------------- |
| `id`       | string | ○    | 最大10文字 | 削除するID     |

```
POST /api/authed/deleteBlackList?id=usr_xxxx
```

#### レスポンス（200）

```json
{ "code": 200, "msg": "success", "data": null }
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/deleteBlackList?id=usr_xxxx" \
  -H "Authorization: <token>"
```

---

## スタッフリスト

スタッフに登録されたIDは、`/api/getAllowList`（入場許可リスト）に当選者と並んで含まれます。

### GET /api/authed/getStaffList

スタッフとして登録されているIDの一覧を取得します。

#### リクエスト

パラメータなし。

#### レスポンス（200）

`data` にID文字列の配列が入ります。登録がない場合は空配列 `[]`。

```json
{
  "code": 200,
  "msg": "success",
  "data": ["usr_staff01", "usr_staff02"]
}
```

#### 使用例

```bash
curl http://localhost:8080/api/authed/getStaffList \
  -H "Authorization: <token>"
```

---

### POST /api/authed/addStaff

スタッフを追加します。

#### リクエスト

| パラメータ | 型     | 必須 | 制約       | 説明           |
| ---------- | ------ | ---- | ---------- | -------------- |
| `id`       | string | ○    | 最大10文字 | 追加するID     |

```
POST /api/authed/addStaff?id=usr_xxxx
```

#### レスポンス（200）

```json
{ "code": 200, "msg": "success", "data": null }
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/addStaff?id=usr_xxxx" \
  -H "Authorization: <token>"
```

---

### POST /api/authed/deleteStaff

スタッフを削除します。

#### リクエスト

| パラメータ | 型     | 必須 | 制約       | 説明           |
| ---------- | ------ | ---- | ---------- | -------------- |
| `id`       | string | ○    | 最大10文字 | 削除するID     |

```
POST /api/authed/deleteStaff?id=usr_xxxx
```

#### レスポンス（200）

```json
{ "code": 200, "msg": "success", "data": null }
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/deleteStaff?id=usr_xxxx" \
  -H "Authorization: <token>"
```

---

## 抽選

### POST /api/authed/setDrawing

抽選を **新規設定** します。新しいフォームID（応募受付用のランダム文字列）を採番し、締切日時と当選人数を設定します。

> 重要: この操作は **既存の応募データ・当選結果をすべてリセット** します（`Forms` と `Results` が初期化されます）。フォームID（`id`）も新しく振り直されるため、過去の応募リンクは無効になります。

#### リクエスト

クエリパラメータで指定します。

| パラメータ | 型              | 必須 | 制約     | 説明                                    |
| ---------- | --------------- | ---- | -------- | --------------------------------------- |
| `date`     | string (RFC3339) | ○    | —        | 応募締切日時。Goの `time.Time` としてパースされる |
| `max`      | int             | ○    | `>= 1`   | 当選人数の上限                          |

> `date` はGinのバインダがRFC3339形式（例: `2026-12-31T23:59:59Z`）としてパースします。

```
POST /api/authed/setDrawing?date=2026-12-31T23:59:59Z&max=50
```

#### レスポンス

##### 成功（200）

`data` に **新しく採番されたフォームID**（8文字のランダム文字列）が入ります。
このIDが、応募者向けの `GET /api/getForm?id=<id>` で使うフォームIDになります。

```json
{
  "code": 200,
  "msg": "success",
  "data": "Ab3xK9pQ"
}
```

##### パラメータ不正（400）

`date` 未指定・形式不正、`max` 未指定・`1` 未満の場合。

```json
{
  "code": 400,
  "msg": "bad request",
  "data": "..."
}
```

##### サーバ内部エラー（500）

フォームID採番（外部のEthereum RPCを用いた乱数生成）やデータ保存に失敗した場合。

```json
{ "code": 500, "msg": "internal server error", "data": null }
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/setDrawing?date=2026-12-31T23:59:59Z&max=50" \
  -H "Authorization: <token>"
```

---

### GET /api/authed/getDrawing

現在設定されている抽選の状態を取得します。フォームID・締切日時・当選人数の上限・抽選実行済みフラグを返します。

#### リクエスト

パラメータなし。

#### レスポンス（200）

`data` に抽選状態オブジェクト（`DrawingStatus`）が入ります。

| フィールド | 型               | JSONキー | 説明                                       |
| ---------- | ---------------- | -------- | ------------------------------------------ |
| `id`       | string           | `id`     | 現在のフォームID（`getForm` 用のランダム文字列） |
| `date`     | string (RFC3339) | `date`   | 応募締切日時                               |
| `max`      | int              | `max`    | 当選人数の上限                             |
| `showed`   | bool             | `showed` | 抽選を実行済みか（`drawing` 実行で `true`）    |

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": "Ab3xK9pQ",
    "date": "2026-12-31T23:59:59Z",
    "max": 50,
    "showed": false
  }
}
```

> まだ一度も `setDrawing` を実行していない場合、`id` は空文字列、`date` はゼロ値（`0001-01-01T00:00:00Z`）になります。

#### 使用例

```bash
curl http://localhost:8080/api/authed/getDrawing \
  -H "Authorization: <token>"
```

---

### POST /api/authed/drawing

抽選を **実行** します。これまでに集まった応募の中から、`setDrawing` で設定した当選人数（`max`）を上限にランダムに当選者を選び、当選者リストへ登録します。

- 応募が0件の場合は、空の当選者リストを返します（エラーにはなりません）。
- 当選者の選定にはEthereumブロックチェーン由来の乱数を使用します。
- 選定結果は当選者リストへ追加保存され、以後 `/api/getAllowList` と `/api/authed/getResults` に反映されます。

#### リクエスト

パラメータなし（ボディ・クエリ不要）。

#### レスポンス

##### 成功（200）

`data` に当選したIDの配列が入ります。

```json
{
  "code": 200,
  "msg": "success",
  "data": ["usr_winner01", "usr_winner02"]
}
```

##### サーバ内部エラー（500）

乱数生成・当選結果の保存に失敗した場合。

```json
{ "code": 500, "msg": "internal server error", "data": null }
```

#### 使用例

```bash
curl -X POST http://localhost:8080/api/authed/drawing \
  -H "Authorization: <token>"
```

---

### GET /api/authed/getResults

現在の **当選者リスト** を取得します。

#### リクエスト

パラメータなし。

#### レスポンス（200）

`data` に当選者IDの配列が入ります。未抽選・全削除済みの場合は空配列 `[]`。

```json
{
  "code": 200,
  "msg": "success",
  "data": ["usr_winner01", "usr_winner02"]
}
```

> 並び順は内部マップの列挙順に依存するため保証されません。

#### 使用例

```bash
curl http://localhost:8080/api/authed/getResults \
  -H "Authorization: <token>"
```

---

### POST /api/authed/removeResults

当選者リストから特定のIDを削除します（当選辞退・繰り上げ調整などに使用）。

#### リクエスト

| パラメータ | 型     | 必須 | 制約       | 説明               |
| ---------- | ------ | ---- | ---------- | ------------------ |
| `id`       | string | ○    | 最大10文字 | 削除する当選者ID   |

```
POST /api/authed/removeResults?id=usr_winner01
```

#### レスポンス

##### 成功（200）

```json
{ "code": 200, "msg": "success", "data": null }
```

##### `id` 未指定／10文字超過（400）

```json
{
  "code": 400,
  "msg": "bad request",
  "data": "..."
}
```

##### サーバ内部エラー（500）

保存に失敗した場合。

```json
{ "code": 500, "msg": "internal server error", "data": null }
```

#### 使用例

```bash
curl -X POST "http://localhost:8080/api/authed/removeResults?id=usr_winner01" \
  -H "Authorization: <token>"
```

# AGENTS.md — VRChat 抽選システム フロントエンド

VRChat 向け抽選イベントの **応募フォーム** SPA。Go バックエンド（`VRCLotterySystem`、`/api` を `:8080` で配信）に対する単一ページの応募フォームを提供する。

## 技術スタック

- Vue 3 + TypeScript + Vite 7
- PrimeVue 4（**Material テーマ**・ライトモード・**ピンク基調**）+ `unplugin-vue-components`（コンポーネント auto-import）
- Pinia 3（状態管理）
- Vue Router 4（**hash mode**）
- Vitest + @vue/test-utils + happy-dom（ユニットテスト）
- markdown-it + DOMPurify（規約の Discord 風 Markdown 描画）
- playwright（devDep・実機レイアウト検証用。本番では未使用）

## 開発コマンド

```bash
npm install
npm run dev         # Vite devサーバ。/api を :8080 へ proxy（バックエンドを別途起動）
npm run build       # vue-tsc 型チェック + 本番ビルド
npm run test        # Vitest（watch）
npm run test:run    # Vitest（1回実行）
npm run type-check  # vue-tsc のみ
```

> バックエンドは CORS 未設定のため、開発時は必ず Vite proxy 経由（`/api`）でアクセスする。`vite.config.ts` の `server.proxy` で `:8080` に転送している。

## ディレクトリ構成

```
src/
├── main.ts                 # Pinia / Router / PrimeVue(Material+ピンク) / ToastService 登録
├── App.vue                 # <Toast/> + <router-view/>
├── router/index.ts         # hash mode。/(Loading) と /form
├── types/api.ts            # CommonResp<T> / FormItem / SubmitInput / SubmitFormReq
├── services/
│   ├── http.ts             # fetch ラッパ。CommonResp アンラップ + ApiError
│   └── lottery.ts          # isActive / getTerms / getForm / submitForm
├── stores/lottery.ts       # 状態・isValid・buildPayload・submit の中核
├── composables/
│   ├── useMarkdown.ts      # markdown-it + DOMPurify（spoiler 後処理付き）
│   └── useScrollEnd.ts     # 「最後までスクロール」検知（短文は即達成）
├── views/
│   ├── LoadingView.vue     # 初期ロード + 受付終了/エラー Dialog
│   └── FormView.vue        # TermsDialog → FormRenderer 切替 + 送信結果
├── components/
│   ├── TermsDialog.vue     # 規約ダイアログ（スクロール末尾で確認ボタン活性化）
│   ├── FormRenderer.vue    # form を type 別カードへディスパッチ
│   ├── SubmitButton.vue
│   └── cards/{ContentCard,InputCard,OptionsCard}.vue
├── assets/main.css         # グローバル + .md-card / .field-* / .discord-md スタイル
└── test/withSetup.ts       # composable をコンポーネント setup 内で実行するテストヘルパー
```

## 画面フロー

1. **Loading**（`/`）: `isActive` と `getTerms` を並列取得。`isActive===true` なら `getForm` も取得し `/form` へ。`false` なら受付終了 Dialog、取得失敗ならエラー Dialog（再試行）。
2. **Terms**（`FormView` 内 `TermsDialog`）: 規約を Discord 風 Markdown で表示。本文を**最後までスクロール**すると「同意して続ける」が活性化（短文は即活性化）。
3. **Form**: `getForm` の `FormItem[]` を `content`/`input`/`options` の 3 種カードで 1 列フロー描画。必須項目が揃うと提出ボタン活性化 → `submitForm`。成功で完了画面、失敗で Toast。

## バックエンド API（要点）

ベース `/api`、共通レスポンス `{ code, msg, Data }`（`Data` は **大文字始まり**・JSON タグ無しのため）。

| メソッド | パス | Data |
| --- | --- | --- |
| GET | `/api/isActive` | bool（true=受付中） |
| GET | `/api/getTerms` | string（規約本文） |
| GET | `/api/getForm` | FormItem[] |
| POST | `/api/submitForm` | null（成功時） |
| GET | `/api/getAllowList` | text/plain（本フォームでは未使用） |

`FormItem = { is_id, title, desc, required, options: string[]|null, type: 'content'|'input'|'options' }`

`submitForm` ボディ: `{ inputs: { content: { "<index>": "値" }, selected: { "<index>": [選択肢index...] } } }`
キーは `getForm` 配列の**インデックス**。送信エラー文字列: `date expired` / `input already exists` / `required field(...) not filled`。

詳細は `temp/api.md` を参照。

## 実装上の重要な注意（ハマりどころ）

- **送信キーは文字列化**: Go の `map[int]` は JSON で `"1"` のようにクォート文字列キーになる。`stores/lottery.ts` の `buildPayload()` で `String(index)` に変換し、空 input は除外、`selected` はソートして送る。
- **required 検証はフロント側でも実施**: 提出ボタンの活性制御（要件）として `isValid` computed で input/options の必須を検証する。サーバ側でも `submitForm` で required 検証が行われるため二重で担保される。
- **テーマ**: PrimeVue Material をベースに `definePreset` で primary をピンク（`{pink.*}`）へ差し替え（`main.ts`）。`darkModeSelector: '.app-dark'`（未付与）で**常時ライトモード**に固定。背景の淡いピンク・カード・Markdown 色などテーマ変数で表せない部分は `assets/main.css` で調整。
- **Markdown は html:false + DOMPurify の二重防御**。spoiler（`||...||`）は markdown-it の text ルールが `|` を特殊文字扱いしないため、インラインルールではなく**レンダリング後の後処理**（`applySpoilers`）で span 化している。
- **Terms ダイアログのスクロールは 1 本に集約**: `.terms-body` を `overflow-y: auto`、PrimeVue Dialog の `.p-dialog-content` を `pt` で `overflow: hidden` にしている（両方 auto だと画面が低い時に縦スクロールバーが 2 本出る）。
- **スクロール末尾検知**: `useScrollEnd` は Dialog の遅延マウントに対応するため `@show` + `requestAnimationFrame` で再判定する。
- **options は複数選択（チェックボックス）**。`FormItem` に単一/複数のフラグが無いため全 options を複数選択として扱う。required は 1 つ以上で充足。
- **カード/ラベルのスタイルはグローバル**（`assets/main.css` の `.md-card` / `.field-title` / `.field-desc` / `.required-mark`）。scoped では共有されないため共通スタイルはここに置く。

## テスト方針

ロジック層を厚く、UI を薄く。`services/http`・`stores/lottery`・`composables/*` をユニットテストでカバー（fetch は `vi.stubGlobal` でモック）。コンポーネントの見た目確認は Playwright で実機計測する（一時スクリプトはリポジトリ外に置き、検証後に削除する）。

## バージョン留意

- Vite 7 を使うため `@vitejs/plugin-vue` は **^6**（^5 は Vite 7 非対応）。
- `@primevue/themes` は deprecated 警告が出るが PrimeVue 4 系では動作する（将来的に `@primeuix/themes` へ移行）。

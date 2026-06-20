// バックエンド（VRCLotterySystem）の API 仕様に対応する型定義。
// 共通レスポンスは { code, msg, Data }。Data は JSON タグ無しのため大文字始まり。

export interface CommonResp<T> {
  code: number
  msg: string
  Data: T
}

/** /api/getSiteData で返るサイト表示情報 */
export interface SiteData {
  /** サイトのタイトル（document.title 等） */
  title: string
  /** フォームページのタイトル（見出し） */
  form_title: string
  /** 利用規約（Markdown） */
  terms: string
}

export type FormType = 'content' | 'input' | 'options'

/** /api/getForm で返るフォーム項目定義 */
export interface FormItem {
  is_id: boolean
  title: string
  desc: string
  required: boolean
  options: string[] | null
  type: FormType
}

/**
 * /api/submitForm の応募データ。
 * キーは getForm の配列インデックス（Go の map[int] は JSON で文字列キーになる）。
 */
export interface SubmitInput {
  /** input タイプの回答。index(文字列) -> 入力値 */
  content: Record<string, string>
  /** options タイプの回答。index(文字列) -> 選択肢indexの配列 */
  selected: Record<string, number[]>
}

export interface SubmitFormReq {
  inputs: SubmitInput
}

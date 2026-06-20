import { getJson, postJson } from './http'
import type { FormItem, SiteData, SubmitFormReq } from '@/types/api'

export const lotteryApi = {
  /** 受付状態（true=受付中 / false=受付終了） */
  isActive: () => getJson<boolean>('isActive'),
  /** サイト情報（タイトル・フォームタイトル・利用規約） */
  getSiteData: () => getJson<SiteData>('getSiteData'),
  /** 応募フォームの項目定義一覧（id はフォーム識別用のランダム文字列。最大10文字） */
  getForm: (id: string) => getJson<FormItem[]>('getForm', { id }),
  /** 応募フォームの送信 */
  submitForm: (req: SubmitFormReq) => postJson<SubmitFormReq, null>('submitForm', req),
}

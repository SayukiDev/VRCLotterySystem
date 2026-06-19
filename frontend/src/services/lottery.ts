import { getJson, postJson } from './http'
import type { FormItem, SubmitFormReq } from '@/types/api'

export const lotteryApi = {
  /** 受付状態（true=受付中 / false=受付終了） */
  isActive: () => getJson<boolean>('isActive'),
  /** 利用規約のテキスト */
  getTerms: () => getJson<string>('getTerms'),
  /** 応募フォームの項目定義一覧 */
  getForm: () => getJson<FormItem[]>('getForm'),
  /** 応募フォームの送信 */
  submitForm: (req: SubmitFormReq) => postJson<SubmitFormReq, null>('submitForm', req),
}

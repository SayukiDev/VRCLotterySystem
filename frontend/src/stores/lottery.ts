import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { lotteryApi } from '@/services/lottery'
import { ApiError } from '@/services/http'
import type { FormItem, SubmitFormReq } from '@/types/api'

/** submitForm のエラーを利用者向けメッセージに変換する */
export function mapSubmitError(err: unknown): string {
  if (err instanceof ApiError) {
    const data = typeof err.data === 'string' ? err.data : ''
    if (data === 'date expired') return '受付は終了しました。'
    if (data === 'input already exists') return 'このIDはすでに応募済みです。'
    if (data.startsWith('required field')) return '必須項目が入力されていません。'
    if (data) return data
    return err.message || '送信に失敗しました。'
  }
  return '送信に失敗しました。通信環境をご確認ください。'
}

export const useLotteryStore = defineStore('lottery', () => {
  // ---- 初期ロード ----
  const loading = ref(true)
  const loadError = ref<string | null>(null)
  // フォームURLが不正（id 欠落 / getForm が 400）。専用エラー表示に使う
  const invalidForm = ref(false)
  const isActive = ref(false)
  const initialized = ref(false)

  // ---- サイト情報 ----
  const siteTitle = ref('')
  const formTitle = ref('')

  // ---- 規約 ----
  const terms = ref('')
  const termsAccepted = ref(false)

  // ---- フォーム ----
  const form = ref<FormItem[]>([])
  // input 項目: index -> 入力値
  const contentInputs = ref<Record<number, string>>({})
  // options 項目: index -> 選択された選択肢インデックスの配列
  const selectedOptions = ref<Record<number, number[]>>({})

  // ---- 提出 ----
  const submitting = ref(false)
  const submitError = ref<string | null>(null)
  const submitSuccess = ref(false)

  /**
   * 起動時: 受付状態と規約を取得し、受付中ならフォーム定義も取得する。
   * formId は URL（#/?id=...）由来のフォーム識別子。欠落/不正は invalidForm として扱う。
   */
  async function initialize(formId?: string) {
    loading.value = true
    loadError.value = null
    invalidForm.value = false

    // id が無ければ API を呼ばず即「URL不正」表示
    if (!formId) {
      invalidForm.value = true
      loading.value = false
      return
    }

    try {
      const [active, site] = await Promise.all([
        lotteryApi.isActive(),
        lotteryApi.getSiteData(),
      ])
      isActive.value = active
      siteTitle.value = site.title
      formTitle.value = site.form_title
      terms.value = site.terms
      if (site.title) document.title = site.title

      if (active) {
        // getForm の 400（id 不正）だけは専用の URL 不正表示に振り分ける
        try {
          const items = await lotteryApi.getForm(formId)
          setForm(items)
        } catch (err) {
          if (err instanceof ApiError && err.status === 400) {
            invalidForm.value = true
            return
          }
          throw err
        }
      }
      initialized.value = true
    } catch (err) {
      loadError.value =
        err instanceof ApiError
          ? err.message || '読み込みに失敗しました。'
          : '読み込みに失敗しました。通信環境をご確認ください。'
    } finally {
      loading.value = false
    }
  }

  /** フォーム定義をセットし、入力値の器を初期化する（リアクティビティ安定化） */
  function setForm(items: FormItem[]) {
    form.value = items
    const content: Record<number, string> = {}
    const selected: Record<number, number[]> = {}
    items.forEach((item, index) => {
      if (item.type === 'input') content[index] = ''
      else if (item.type === 'options') selected[index] = []
    })
    contentInputs.value = content
    selectedOptions.value = selected
  }

  function acceptTerms() {
    termsAccepted.value = true
  }

  function setContent(index: number, value: string) {
    contentInputs.value[index] = value
  }

  /** options のチェックボックスをトグルする（複数選択） */
  function toggleOption(index: number, optIndex: number) {
    const cur = selectedOptions.value[index] ?? []
    const pos = cur.indexOf(optIndex)
    if (pos === -1) selectedOptions.value[index] = [...cur, optIndex]
    else selectedOptions.value[index] = cur.filter((v) => v !== optIndex)
  }

  /** 必須項目がすべて充足しているか（提出ボタン活性条件の単一真実源） */
  const isValid = computed(() =>
    form.value.every((item, index) => {
      if (!item.required) return true
      if (item.type === 'input') {
        return (contentInputs.value[index] ?? '').trim().length > 0
      }
      if (item.type === 'options') {
        return (selectedOptions.value[index] ?? []).length > 0
      }
      return true // content は必須対象外
    }),
  )

  /** 送信ボディを組み立てる。キーは文字列化し、空 input は除外する */
  function buildPayload(): SubmitFormReq {
    const content: Record<string, string> = {}
    const selected: Record<string, number[]> = {}
    form.value.forEach((item, index) => {
      if (item.type === 'input') {
        const v = contentInputs.value[index]
        if (v != null && v.trim() !== '') content[String(index)] = v.trim()
      } else if (item.type === 'options') {
        const sel = selectedOptions.value[index]
        if (sel && sel.length > 0) {
          selected[String(index)] = [...sel].sort((a, b) => a - b)
        }
      }
    })
    return { inputs: { content, selected } }
  }

  /** フォームを送信する。成功で submitSuccess=true、失敗で submitError をセット */
  async function submit(): Promise<boolean> {
    if (submitting.value) return false
    submitting.value = true
    submitError.value = null
    try {
      await lotteryApi.submitForm(buildPayload())
      submitSuccess.value = true
      return true
    } catch (err) {
      submitError.value = mapSubmitError(err)
      return false
    } finally {
      submitting.value = false
    }
  }

  return {
    loading,
    loadError,
    invalidForm,
    isActive,
    initialized,
    siteTitle,
    formTitle,
    terms,
    termsAccepted,
    form,
    contentInputs,
    selectedOptions,
    submitting,
    submitError,
    submitSuccess,
    isValid,
    initialize,
    setForm,
    acceptTerms,
    setContent,
    toggleOption,
    buildPayload,
    submit,
  }
})

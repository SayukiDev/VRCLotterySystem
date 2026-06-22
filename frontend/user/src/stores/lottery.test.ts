import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useLotteryStore, mapSubmitError } from './lottery'
import { ApiError } from '@/services/http'
import type { FormItem } from '@/types/api'

const FORM: FormItem[] = [
  { is_id: false, title: '注意事項', desc: '読んでね', required: false, options: null, type: 'content' },
  { is_id: true, title: 'VRChat ID', desc: '', required: true, options: null, type: 'input' },
  { is_id: false, title: '参加希望日', desc: '', required: true, options: ['1日目', '2日目'], type: 'options' },
  { is_id: false, title: '備考', desc: '', required: false, options: null, type: 'input' },
]

describe('lottery store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('isValid', () => {
    it('必須未入力では false', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      expect(store.isValid).toBe(false)
    })

    it('必須 input が空白のみでは false', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.setContent(1, '   ')
      store.toggleOption(2, 0)
      expect(store.isValid).toBe(false)
    })

    it('必須 input が空 かつ options 未選択では false', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.setContent(1, 'usr_abc')
      expect(store.isValid).toBe(false)
    })

    it('必須項目がすべて充足すると true', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.setContent(1, 'usr_abc')
      store.toggleOption(2, 0)
      expect(store.isValid).toBe(true)
    })
  })

  describe('toggleOption', () => {
    it('複数選択でき、再トグルで解除される', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.toggleOption(2, 0)
      store.toggleOption(2, 1)
      expect(store.selectedOptions[2]).toEqual([0, 1])
      store.toggleOption(2, 0)
      expect(store.selectedOptions[2]).toEqual([1])
    })
  })

  describe('buildPayload', () => {
    it('文字列キー化・空input除外・selectedソートを行う', () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.setContent(1, '  usr_abc  ')
      store.setContent(3, '') // 空 input は除外される
      store.toggleOption(2, 1)
      store.toggleOption(2, 0)

      const payload = store.buildPayload()
      expect(payload).toEqual({
        inputs: {
          content: { '1': 'usr_abc' },
          selected: { '2': [0, 1] },
        },
      })
    })
  })

  describe('mapSubmitError', () => {
    it('date expired を受付終了メッセージに変換', () => {
      expect(mapSubmitError(new ApiError(500, 'date expired', 'x'))).toContain('受付は終了')
    })
    it('input already exists を応募済みメッセージに変換', () => {
      expect(mapSubmitError(new ApiError(500, 'input already exists', 'x'))).toContain('応募済み')
    })
    it('required field を未入力メッセージに変換', () => {
      expect(mapSubmitError(new ApiError(400, 'required field(VRChat ID) not filled', 'x'))).toContain(
        '必須項目',
      )
    })
    it('非 ApiError は汎用メッセージ', () => {
      expect(mapSubmitError(new Error('boom'))).toContain('通信環境')
    })
  })

  describe('initialize', () => {
    it('id が無ければ API を呼ばず invalidForm=true', async () => {
      const store = useLotteryStore()
      const fetchMock = vi.fn()
      vi.stubGlobal('fetch', fetchMock)
      await store.initialize('')
      expect(store.invalidForm).toBe(true)
      expect(store.loading).toBe(false)
      expect(fetchMock).not.toHaveBeenCalled()
      vi.unstubAllGlobals()
    })

    const SITE = { title: '宵撫旅館', form_title: '抽選予約フォーム', terms: '規約' }

    it('getForm が 400 のとき invalidForm=true', async () => {
      const store = useLotteryStore()
      vi.stubGlobal(
        'fetch',
        vi.fn().mockImplementation((url: string) => {
          const body = url.includes('isActive')
            ? { code: 200, msg: 'success', Data: true }
            : url.includes('getSiteData')
              ? { code: 200, msg: 'success', Data: SITE }
              : { code: 400, msg: 'bad request', Data: 'invalid id' }
          return Promise.resolve({ ok: true, status: 200, json: () => Promise.resolve(body) } as Response)
        }),
      )
      await store.initialize('toolongid_x')
      expect(store.invalidForm).toBe(true)
      expect(store.loadError).toBeNull()
      vi.unstubAllGlobals()
    })

    it('受付中かつ id 有効でサイト情報とフォームをセットする', async () => {
      const store = useLotteryStore()
      vi.stubGlobal(
        'fetch',
        vi.fn().mockImplementation((url: string) => {
          const body = url.includes('isActive')
            ? { code: 200, msg: 'success', Data: true }
            : url.includes('getSiteData')
              ? { code: 200, msg: 'success', Data: SITE }
              : { code: 200, msg: 'success', Data: FORM }
          return Promise.resolve({ ok: true, status: 200, json: () => Promise.resolve(body) } as Response)
        }),
      )
      await store.initialize('Ab3xK9pQ2')
      expect(store.invalidForm).toBe(false)
      expect(store.isActive).toBe(true)
      expect(store.siteTitle).toBe(SITE.title)
      expect(store.formTitle).toBe(SITE.form_title)
      expect(store.terms).toBe(SITE.terms)
      expect(store.form).toEqual(FORM)
      vi.unstubAllGlobals()
    })
  })

  describe('submit', () => {
    it('成功で submitSuccess=true', async () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      store.setContent(1, 'usr_abc')
      store.toggleOption(2, 0)
      vi.stubGlobal(
        'fetch',
        vi.fn().mockResolvedValue({
          ok: true,
          status: 200,
          json: () => Promise.resolve({ code: 200, msg: 'success', Data: null }),
        } as Response),
      )
      const ok = await store.submit()
      expect(ok).toBe(true)
      expect(store.submitSuccess).toBe(true)
      vi.unstubAllGlobals()
    })

    it('失敗で submitError がセットされる', async () => {
      const store = useLotteryStore()
      store.setForm(FORM)
      vi.stubGlobal(
        'fetch',
        vi.fn().mockResolvedValue({
          ok: true,
          status: 200,
          json: () => Promise.resolve({ code: 500, msg: 'internal', Data: 'input already exists' }),
        } as Response),
      )
      const ok = await store.submit()
      expect(ok).toBe(false)
      expect(store.submitError).toContain('応募済み')
      vi.unstubAllGlobals()
    })
  })
})

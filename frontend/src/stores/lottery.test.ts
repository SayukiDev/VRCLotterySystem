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

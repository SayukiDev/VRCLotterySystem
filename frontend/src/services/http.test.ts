import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, getJson, postJson } from './http'

function mockFetchOnce(body: unknown, ok = true, status = 200) {
  const fetchMock = vi.fn().mockResolvedValue({
    ok,
    status,
    statusText: ok ? 'OK' : 'Error',
    json: () => Promise.resolve(body),
  } as Response)
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('http service', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('CommonResp の Data をアンラップして返す', async () => {
    mockFetchOnce({ code: 200, msg: 'success', Data: true })
    await expect(getJson<boolean>('isActive')).resolves.toBe(true)
  })

  it('文字列の Data を返す', async () => {
    mockFetchOnce({ code: 200, msg: 'success', Data: '規約本文' })
    await expect(getJson<string>('getTerms')).resolves.toBe('規約本文')
  })

  it('code>=400 のとき ApiError を throw し Data をエラー情報に持つ', async () => {
    mockFetchOnce({ code: 500, msg: 'internal server error', Data: 'date expired' }, true, 200)
    await expect(getJson('getForm')).rejects.toMatchObject({
      name: 'ApiError',
      status: 500,
      data: 'date expired',
    })
  })

  it('HTTP エラー(res.ok=false)でも ApiError を throw する', async () => {
    mockFetchOnce({ code: 400, msg: 'bad request', Data: 'required field(X) not filled' }, false, 400)
    await expect(getJson('x')).rejects.toBeInstanceOf(ApiError)
  })

  it('postJson は payload を JSON で送る', async () => {
    const fetchMock = mockFetchOnce({ code: 200, msg: 'success', Data: null })
    await postJson('submitForm', { inputs: { content: { '1': 'a' }, selected: {} } })
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/submitForm',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ inputs: { content: { '1': 'a' }, selected: {} } }),
      }),
    )
  })
})

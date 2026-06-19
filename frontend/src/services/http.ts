import type { CommonResp } from '@/types/api'

// Vite の proxy 設定で /api を http://localhost:8080 へ転送する
const BASE = '/api'

/**
 * API エラーを正規化した例外。
 * `data` には CommonResp.Data（"date expired" 等のエラー文字列）が入る。
 */
export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly data: unknown,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function parseCommon<T>(res: Response): Promise<T> {
  let body: CommonResp<T>
  try {
    body = (await res.json()) as CommonResp<T>
  } catch {
    throw new ApiError(res.status, null, res.statusText || 'invalid response')
  }
  if (!res.ok || body.code >= 400) {
    throw new ApiError(body.code ?? res.status, body.Data, body.msg ?? res.statusText)
  }
  return body.Data
}

export async function getJson<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE}/${path}`)
  return parseCommon<T>(res)
}

export async function postJson<TReq, TRes>(path: string, payload: TReq): Promise<TRes> {
  const res = await fetch(`${BASE}/${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  return parseCommon<TRes>(res)
}

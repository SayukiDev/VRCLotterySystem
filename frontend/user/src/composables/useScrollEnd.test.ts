import { describe, expect, it } from 'vitest'
import { ref, type Ref } from 'vue'
import { withSetup } from '@/test/withSetup'
import { useScrollEnd } from './useScrollEnd'

function fakeEl(opts: { scrollTop: number; clientHeight: number; scrollHeight: number }) {
  return opts as unknown as HTMLElement
}

// onMounted を正しく動かすため、コンポーネントの setup 内で composable を呼ぶ
function run(el: Ref<HTMLElement | null>) {
  return withSetup(() => useScrollEnd(el))
}

describe('useScrollEnd', () => {
  it('スクロール不要な短い内容では即 reachedEnd=true', () => {
    const el = ref<HTMLElement | null>(
      fakeEl({ scrollTop: 0, clientHeight: 400, scrollHeight: 400 }),
    )
    const [{ reachedEnd, check }] = run(el)
    check()
    expect(reachedEnd.value).toBe(true)
  })

  it('未到達では false、末尾までスクロールで true', () => {
    const el = ref<HTMLElement | null>(
      fakeEl({ scrollTop: 0, clientHeight: 200, scrollHeight: 1000 }),
    )
    const [{ reachedEnd, check, onScroll }] = run(el)
    check()
    expect(reachedEnd.value).toBe(false)

    ;(el.value as unknown as { scrollTop: number }).scrollTop = 800
    onScroll()
    expect(reachedEnd.value).toBe(true)
  })

  it('el が null でも例外を投げない', () => {
    const el = ref<HTMLElement | null>(null)
    const [{ check }] = run(el)
    expect(() => check()).not.toThrow()
  })
})

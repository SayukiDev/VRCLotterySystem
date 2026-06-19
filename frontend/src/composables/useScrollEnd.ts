import { onMounted, ref, type Ref } from 'vue'

/**
 * スクロールコンテナが最後までスクロールされたかを検知する。
 * スクロール不要なほど短い内容の場合は即 reachedEnd=true にする。
 *
 * @param el       監視対象のスクロールコンテナ
 * @param threshold 末尾とみなす許容誤差(px)
 */
export function useScrollEnd(el: Ref<HTMLElement | null>, threshold = 4) {
  const reachedEnd = ref(false)

  const check = () => {
    const e = el.value
    if (!e) return
    // スクロールバーが出ないほど短ければ即達成
    if (e.scrollHeight - e.clientHeight <= threshold) {
      reachedEnd.value = true
      return
    }
    if (e.scrollTop + e.clientHeight >= e.scrollHeight - threshold) {
      reachedEnd.value = true
    }
  }

  const onScroll = () => check()

  onMounted(() => {
    // コンテンツ描画後に高さが確定するため次フレームで判定
    requestAnimationFrame(check)
  })

  return { reachedEnd, onScroll, check }
}

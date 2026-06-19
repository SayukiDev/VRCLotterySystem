import { createApp } from 'vue'

/**
 * composable をコンポーネントの setup コンテキストで実行するためのテストヘルパー。
 * onMounted などのライフサイクルフックを正しく登録できる。
 * 戻り値は [composableの戻り値, app]。app.unmount() で破棄できる。
 */
export function withSetup<T>(composable: () => T): [T, ReturnType<typeof createApp>] {
  let result!: T
  const app = createApp({
    setup() {
      result = composable()
      return () => null
    },
  })
  app.mount(document.createElement('div'))
  return [result, app]
}

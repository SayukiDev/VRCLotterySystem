<script setup lang="ts">
import { ref } from 'vue'
import { useScrollEnd } from '@/composables/useScrollEnd'
import { renderMarkdown } from '@/composables/useMarkdown'

const props = defineProps<{
  visible: boolean
  title?: string
  body: string
}>()

const emit = defineEmits<{ accept: [] }>()

const bodyEl = ref<HTMLElement | null>(null)
const { reachedEnd, onScroll, check } = useScrollEnd(bodyEl)

const rendered = ref('')

// Dialog はオープン時に内容が遅延マウントされるため、表示時に再計算する
function onShow() {
  rendered.value = renderMarkdown(props.body)
  // 描画後に高さが確定してからスクロール判定
  requestAnimationFrame(() => requestAnimationFrame(check))
}
</script>

<template>
  <Dialog
    :visible="visible"
    modal
    :closable="false"
    :draggable="false"
    :header="title || '利用規約'"
    :style="{ width: '92%', maxWidth: '800px' }"
    :pt="{
      content: { style: 'padding-bottom: 0.5rem; overflow: hidden;' },
      footer: { style: 'padding-top: 1rem;' },
    }"
    @show="onShow"
  >
    <div ref="bodyEl" class="terms-body discord-md" v-html="rendered" @scroll="onScroll" />

    <p v-if="!reachedEnd" class="scroll-hint">
      <i class="pi pi-angle-double-down" /> 最後までお読みください
    </p>

    <template #footer>
      <Button
        label="同意して続ける"
        icon="pi pi-check"
        :disabled="!reachedEnd"
        @click="emit('accept')"
      />
    </template>
  </Dialog>
</template>

<style scoped>
.terms-body {
  /* スクロールはこの要素の縦1本に集約する（Dialog content 側は overflow:hidden） */
  max-height: min(50vh, 460px);
  overflow-y: auto;
  overflow-x: hidden;
  padding: 0.25rem 0.5rem 0.5rem 0;
}

/* コードブロックなど長い要素のみ横スクロールを許可 */
.terms-body :deep(pre) {
  overflow-x: auto;
}

.scroll-hint {
  margin: 0.5rem 0 0.25rem;
  text-align: center;
  font-size: 0.85rem;
  color: var(--p-text-muted-color);
  animation: bob 1.4s ease-in-out infinite;
}

@keyframes bob {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(3px);
  }
}
</style>

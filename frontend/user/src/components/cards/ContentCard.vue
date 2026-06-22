<script setup lang="ts">
import { computed } from 'vue'
import type { FormItem } from '@/types/api'
import { renderMarkdown } from '@/composables/useMarkdown'

const props = defineProps<{ item: FormItem }>()

// content の説明文は Discord 風 Markdown で表示する
const renderedDesc = computed(() => (props.item.desc ? renderMarkdown(props.item.desc) : ''))
</script>

<template>
  <section class="md-card content-card">
    <h2 v-if="item.title" class="content-title">{{ item.title }}</h2>
    <div v-if="renderedDesc" class="discord-md content-desc" v-html="renderedDesc" />
  </section>
</template>

<style scoped>
.content-title {
  margin: 0 0 0.4rem;
  font-size: 1.15rem;
  font-weight: 600;
}

.content-desc {
  color: var(--p-text-color);
}
</style>

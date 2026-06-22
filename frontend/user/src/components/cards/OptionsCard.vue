<script setup lang="ts">
import type { FormItem } from '@/types/api'
import { useLotteryStore } from '@/stores/lottery'

const props = defineProps<{ item: FormItem; index: number }>()
const store = useLotteryStore()

function isChecked(optIndex: number): boolean {
  return (store.selectedOptions[props.index] ?? []).includes(optIndex)
}

function toggle(optIndex: number) {
  store.toggleOption(props.index, optIndex)
}
</script>

<template>
  <section class="md-card">
    <h3 class="field-title">
      {{ item.title }}
      <span v-if="item.required" class="required-mark" aria-label="必須">*</span>
    </h3>
    <p v-if="item.desc" class="field-desc">{{ item.desc }}</p>

    <ul class="option-list">
      <li v-for="(opt, optIndex) in item.options ?? []" :key="optIndex" class="option-item">
        <Checkbox
          :input-id="`opt-${index}-${optIndex}`"
          :model-value="isChecked(optIndex)"
          binary
          @update:model-value="toggle(optIndex)"
        />
        <label :for="`opt-${index}-${optIndex}`" class="option-label">{{ opt }}</label>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.option-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

/* チェックボックスは折り返し時も縮ませない */
.option-item :deep(.p-checkbox) {
  flex: 0 0 auto;
}

.option-label {
  cursor: pointer;
  user-select: none;
  /* 行高を詰めてチェックボックスと視覚的な縦中央を揃える */
  line-height: 1.25;
}
</style>

<script setup lang="ts">
import { computed } from 'vue'
import type { FormItem } from '@/types/api'
import { useLotteryStore } from '@/stores/lottery'

const props = defineProps<{ item: FormItem; index: number }>()
const store = useLotteryStore()

const value = computed({
  get: () => store.contentInputs[props.index] ?? '',
  set: (v: string) => store.setContent(props.index, v),
})
</script>

<template>
  <section class="md-card">
    <h3 class="field-title">
      {{ item.title }}
      <span v-if="item.required" class="required-mark" aria-label="必須">*</span>
    </h3>
    <p v-if="item.desc" class="field-desc">{{ item.desc }}</p>
    <InputText
      v-model="value"
      class="field-input"
      :placeholder="item.is_id ? 'xxxxxxxx' : ''"
      :aria-label="item.title"
    />
  </section>
</template>

<style scoped>
.field-input {
  width: 100%;
}
</style>

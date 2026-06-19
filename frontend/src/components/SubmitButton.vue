<script setup lang="ts">
import { useLotteryStore } from '@/stores/lottery'

const store = useLotteryStore()
const emit = defineEmits<{ submit: [] }>()
</script>

<template>
  <div class="submit-area">
    <Button
      label="応募する"
      icon="pi pi-send"
      :disabled="!store.isValid || store.submitting"
      :loading="store.submitting"
      @click="emit('submit')"
    />
    <p v-if="!store.isValid" class="submit-hint">必須項目（*）をすべて入力・選択してください</p>
  </div>
</template>
<style scoped>
.submit-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.6rem;
  padding: 0.5rem 0 2rem;
}

.submit-hint {
  margin: 0;
  font-size: 0.82rem;
  color: var(--p-text-muted-color);
}
</style>

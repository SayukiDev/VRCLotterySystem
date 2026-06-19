<script setup lang="ts">
import type { Component } from 'vue'
import { useLotteryStore } from '@/stores/lottery'
import type { FormType } from '@/types/api'
import ContentCard from './cards/ContentCard.vue'
import InputCard from './cards/InputCard.vue'
import OptionsCard from './cards/OptionsCard.vue'
import SubmitButton from './SubmitButton.vue'

const store = useLotteryStore()
const emit = defineEmits<{ submit: [] }>()

const cardMap: Record<FormType, Component> = {
  content: ContentCard,
  input: InputCard,
  options: OptionsCard,
}
</script>

<template>
  <div class="form-flow">
    <component
      :is="cardMap[item.type]"
      v-for="(item, index) in store.form"
      :key="index"
      :item="item"
      :index="index"
    />
    <SubmitButton @submit="emit('submit')" />
  </div>
</template>

<style scoped>
.form-flow {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 680px;
  margin: 0 auto;
  padding: 1.5rem 1rem;
}
</style>

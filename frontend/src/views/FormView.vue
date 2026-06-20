<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useLotteryStore } from '@/stores/lottery'
import TermsDialog from '@/components/TermsDialog.vue'
import FormRenderer from '@/components/FormRenderer.vue'

const store = useLotteryStore()
const router = useRouter()
const route = useRoute()
const toast = useToast()

// リロードや #/<id>/form 直接アクセス等で未初期化なら、id を保ったまま Loading へ戻す
onMounted(() => {
  if (!store.initialized || !store.isActive) {
    const id = typeof route.params.id === 'string' ? route.params.id : ''
    router.replace(id ? { name: 'loading', params: { id } } : { path: '/' })
  }
})

async function onSubmit() {
  const ok = await store.submit()
  if (ok) {
    toast.add({
      severity: 'success',
      summary: '応募完了',
      detail: 'ご応募ありがとうございました。',
      life: 5000,
    })
  } else {
    toast.add({
      severity: 'error',
      summary: '送信エラー',
      detail: store.submitError ?? '送信に失敗しました。',
      life: 6000,
    })
  }
}
</script>

<template>
  <main class="form-view">
    <!-- 規約ダイアログ（同意前） -->
    <TermsDialog
      :visible="!store.termsAccepted"
      title="利用規約"
      :body="store.terms"
      @accept="store.acceptTerms"
    />

    <!-- 提出完了画面 -->
    <div v-if="store.submitSuccess" class="result-state">
      <i class="pi pi-check-circle result-icon" />
      <h1 class="result-title">応募が完了しました</h1>
      <p class="result-text">ご応募ありがとうございました。<br />抽選結果をお待ちください。</p>
    </div>

    <!-- フォーム本体（同意後） -->
    <template v-else-if="store.termsAccepted">
      <header class="form-header">
        <h1 class="form-heading">{{ store.formTitle || '抽選応募フォーム' }}</h1>
      </header>
      <FormRenderer @submit="onSubmit" />
    </template>
  </main>
</template>

<style scoped>
.form-view {
  min-height: 100vh;
}

.form-header {
  max-width: 680px;
  margin: 0 auto;
  padding: 2rem 1rem 0;
}

.form-heading {
  margin: 0;
  font-size: 1.6rem;
  font-weight: 700;
}

.result-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  min-height: 100vh;
  padding: 1.5rem;
  text-align: center;
}

.result-icon {
  font-size: 3.5rem;
  color: #22c55e;
}

.result-title {
  margin: 0;
  font-size: 1.5rem;
}

.result-text {
  margin: 0;
  color: var(--p-text-muted-color);
}
</style>

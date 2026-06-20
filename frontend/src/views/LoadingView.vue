<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLotteryStore } from '@/stores/lottery'

const store = useLotteryStore()
const router = useRouter()
const route = useRoute()

async function load() {
  // フォーム識別 id は URL パス（#/<id>）から取得する
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  await store.initialize(id)
  if (!store.loadError && !store.invalidForm && store.isActive) {
    // id をパスに残したままフォームへ（リロード時に再取得できるように）
    router.replace({ name: 'form', params: { id } })
  }
}

onMounted(load)

function retry() {
  load()
}
</script>

<template>
  <main class="loading-view">
    <!-- 読み込み中 -->
    <div v-if="store.loading" class="loading-state">
      <ProgressSpinner stroke-width="4" />
      <p class="loading-text">読み込み中…</p>
    </div>

    <!-- 読み込みエラー -->
    <Dialog
      :visible="!store.loading && !!store.loadError"
      modal
      :closable="false"
      header="エラー"
      :style="{ width: '90%', maxWidth: '420px' }"
    >
      <p>{{ store.loadError }}</p>
      <template #footer>
        <Button label="再試行" icon="pi pi-refresh" @click="retry" />
      </template>
    </Dialog>

    <!-- フォームURL不正（id 欠落 / 不正） -->
    <Dialog
      :visible="!store.loading && !store.loadError && store.invalidForm"
      modal
      :closable="false"
      header="エラー"
      :style="{ width: '90%', maxWidth: '420px' }"
    >
      <p>フォームURLが正しくありません。<br />主催者から共有された正しいリンクをご確認ください。</p>
    </Dialog>

    <!-- 受付終了 -->
    <Dialog
      :visible="!store.loading && !store.loadError && !store.invalidForm && !store.isActive"
      modal
      :closable="false"
      header="受付終了"
      :style="{ width: '90%', maxWidth: '420px' }"
    >
      <div class="closed-content">
        <i class="pi pi-clock closed-icon" />
        <p>応募の受付は終了しました。<br />ご応募ありがとうございました。</p>
      </div>
    </Dialog>
  </main>
</template>

<style scoped>
.loading-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 1rem;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.loading-text {
  color: var(--p-text-muted-color);
  margin: 0;
}

.closed-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  text-align: center;
}

.closed-icon {
  font-size: 2.5rem;
  color: var(--p-primary-color);
}
</style>

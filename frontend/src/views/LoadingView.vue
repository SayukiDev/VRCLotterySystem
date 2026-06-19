<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLotteryStore } from '@/stores/lottery'

const store = useLotteryStore()
const router = useRouter()

async function load() {
  await store.initialize()
  if (!store.loadError && store.isActive) {
    router.replace('/form')
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

    <!-- 受付終了 -->
    <Dialog
      :visible="!store.loading && !store.loadError && !store.isActive"
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

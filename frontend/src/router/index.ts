import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    // id 無しアクセス（#/）。LoadingView 側で「URL不正」表示になる
    {
      path: '/',
      name: 'loading-root',
      component: () => import('@/views/LoadingView.vue'),
    },
    // フォーム識別 id をパスに持つ（#/<id>）
    {
      path: '/:id',
      name: 'loading',
      component: () => import('@/views/LoadingView.vue'),
    },
    // フォーム本体（#/<id>/form）
    {
      path: '/:id/form',
      name: 'form',
      component: () => import('@/views/FormView.vue'),
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

export default router

import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'loading',
      component: () => import('@/views/LoadingView.vue'),
    },
    {
      path: '/form',
      name: 'form',
      component: () => import('@/views/FormView.vue'),
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

export default router

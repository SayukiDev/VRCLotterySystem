import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import { definePreset } from '@primevue/themes'
import Material from '@primevue/themes/material'
import ToastService from 'primevue/toastservice'

import App from './App.vue'
import router from './router'

import 'primeicons/primeicons.css'
import './assets/main.css'

// Material をベースに primary をピンクへ差し替える
const PinkPreset = definePreset(Material, {
  semantic: {
    primary: {
      50: '{pink.50}',
      100: '{pink.100}',
      200: '{pink.200}',
      300: '{pink.300}',
      400: '{pink.400}',
      500: '{pink.500}',
      600: '{pink.600}',
      700: '{pink.700}',
      800: '{pink.800}',
      900: '{pink.900}',
      950: '{pink.950}',
    },
  },
})

createApp(App)
  .use(createPinia())
  .use(router)
  .use(PrimeVue, {
    theme: {
      preset: PinkPreset,
      options: {
        // .app-dark を付与しない限り常にライトモード（システムのダーク設定を無効化）
        darkModeSelector: '.app-dark',
      },
    },
  })
  .use(ToastService)
  .mount('#app')

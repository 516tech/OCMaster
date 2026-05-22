import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/theme-chalk/dark/css-vars.css'
import 'element-plus/theme-chalk/index.css'
import App from './App.vue'
import router from './router'
import { events } from '@/composables/useEventBus'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: undefined })

// Global error → toast
app.config.errorHandler = (err) => {
  const msg = err instanceof Error ? err.message : String(err)
  events.scanCompleted.post({ success: false, error: msg })
}

app.mount('#app')

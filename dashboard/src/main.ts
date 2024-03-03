import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { clerkPlugin } from 'vue-clerk/plugin'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Add the Clerk plugin
app.use(clerkPlugin, {
  publishableKey: "pk_test_aW50aW1hdGUta3JpbGwtNzguY2xlcmsuYWNjb3VudHMuZGV2JA",
})

app.mount('#app')

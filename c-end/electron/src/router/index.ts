import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'scan', component: () => import('@/views/ScanPage.vue') },
    { path: '/upload', name: 'upload', component: () => import('@/views/UploadPage.vue') },
    { path: '/settings', name: 'settings', component: () => import('@/views/SettingsPage.vue') },
    { path: '/about', name: 'about', component: () => import('@/views/AboutPage.vue') }
  ]
})

export default router

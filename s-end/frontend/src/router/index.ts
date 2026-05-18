import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue') },
  { path: '/register', name: 'Register', component: () => import('../views/RegisterView.vue') },
  { path: '/dashboard', name: 'Dashboard', component: () => import('../views/DashboardView.vue'), meta: { auth: true } },
  { path: '/hardware/:code', name: 'HardwareDetail', component: () => import('../views/HardwareDetailView.vue'), meta: { auth: true } },
  { path: '/suggestion/new/:code', name: 'SuggestionNew', component: () => import('../views/SuggestionEditView.vue'), meta: { auth: true } },
  { path: '/suggestion/:id', name: 'SuggestionEdit', component: () => import('../views/SuggestionEditView.vue'), meta: { auth: true } },
  { path: '/history', name: 'History', component: () => import('../views/SuggestionHistoryView.vue'), meta: { auth: true } },
  { path: '/profile', name: 'Profile', component: () => import('../views/ProfileView.vue'), meta: { auth: true } },
  { path: '/settings', name: 'Settings', component: () => import('../views/SettingsView.vue'), meta: { auth: true } },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
]

const router = createRouter({ history: createWebHashHistory(), routes })

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.token) {
    next('/login')
  } else {
    next()
  }
})

export default router

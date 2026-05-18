import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Merchant } from '../types'
import { authApi, profileApi } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const merchant = ref<Merchant | null>(null)

  async function login(phone: string, password: string) {
    const res = await authApi.login(phone, password)
    token.value = res.data.token
    localStorage.setItem('token', token.value)
    await fetchProfile()
  }

  async function register(phone: string, password: string) {
    await authApi.register(phone, password)
  }

  async function fetchProfile() {
    const res = await profileApi.get()
    merchant.value = res.data
  }

  function logout() {
    token.value = ''
    merchant.value = null
    localStorage.removeItem('token')
  }

  return { token, merchant, login, register, fetchProfile, logout }
})

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HardwareInfo } from '../types'
import { hardwareApi } from '../api/hardware'

export const useHardwareStore = defineStore('hardware', () => {
  const hardware = ref<HardwareInfo | null>(null)
  const loading = ref(false)

  async function fetchByCode(code: string) {
    loading.value = true
    const res = await hardwareApi.getByCode(code)
    hardware.value = res.data
    loading.value = false
  }

  return { hardware, loading, fetchByCode }
})

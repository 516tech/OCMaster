import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface AppConfig {
  apiUrl: string
  language: string
  autoScan: boolean
}

export const useSettingsStore = defineStore('settings', () => {
  const apiUrl = ref('https://localhost/api/v1')
  const language = ref('zh-CN')
  const autoScan = ref(true)
  const loaded = ref(false)

  async function load() {
    try {
      const config = await window.electronAPI.getConfig()
      apiUrl.value = (config.apiUrl as string) || apiUrl.value
      language.value = (config.language as string) || language.value
      autoScan.value = (config.autoScan as boolean) ?? autoScan.value
      loaded.value = true
    } catch {
      // use defaults
    }
  }

  async function save() {
    await window.electronAPI.saveConfig({
      apiUrl: apiUrl.value,
      language: language.value,
      autoScan: autoScan.value
    })
  }

  return { apiUrl, language, autoScan, loaded, load, save }
})

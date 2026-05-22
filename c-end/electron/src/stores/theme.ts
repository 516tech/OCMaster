import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ThemeMode } from '@/composables/useTheme'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>('system')

  async function load() {
    try {
      const config = await window.electronAPI.getConfig()
      mode.value = (config.theme as ThemeMode) || 'system'
    } catch {
      mode.value = 'system'
    }
  }

  async function setMode(newMode: ThemeMode) {
    mode.value = newMode
    try {
      await window.electronAPI.saveConfig({ theme: newMode })
    } catch {
      // ignore persist errors
    }
  }

  return { mode, load, setMode }
})

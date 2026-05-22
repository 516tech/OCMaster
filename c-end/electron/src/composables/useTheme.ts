// Theme composable — adapted from ImHex ThemeManager pattern
// Hot-swaps CSS variables for dark/light mode, supports system preference

import { ref, watchEffect } from 'vue'
import { useThemeStore } from '@/stores/theme'

export type ThemeMode = 'dark' | 'light' | 'system'

const SYSTEM_DARK_QUERY = window.matchMedia('(prefers-color-scheme: dark)')

export function useTheme() {
  const store = useThemeStore()
  const isDark = ref(false)

  function applyTheme(mode: ThemeMode): void {
    const effective = mode === 'system'
      ? (SYSTEM_DARK_QUERY.matches ? 'dark' : 'light')
      : mode

    isDark.value = effective === 'dark'
    document.documentElement.classList.toggle('dark', isDark.value)
    store.setMode(mode)
  }

  function toggle(): void {
    const next = isDark.value ? 'light' : 'dark'
    applyTheme(next)
  }

  // Init from store or system
  const saved = store.mode
  applyTheme(saved)

  // Follow system preference changes when in system mode
  SYSTEM_DARK_QUERY.addEventListener('change', () => {
    if (store.mode === 'system') applyTheme('system')
  })

  return { isDark, themeMode: store.mode, applyTheme, toggle }
}

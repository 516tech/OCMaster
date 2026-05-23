import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from '@/stores/theme'

vi.stubGlobal('electronAPI', {
  getConfig: vi.fn().mockResolvedValue({ theme: 'system' }),
  saveConfig: vi.fn().mockResolvedValue({ ok: true })
})
vi.stubGlobal('matchMedia', vi.fn().mockReturnValue({
  matches: false,
  addEventListener: vi.fn(),
  removeEventListener: vi.fn()
}))

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
})

describe('theme store', () => {
  it('loads saved theme preference', async () => {
    window.electronAPI.getConfig = vi.fn().mockResolvedValue({ theme: 'dark' })
    const store = useThemeStore()
    await store.load()
    expect(store.mode).toBe('dark')
  })

  it('defaults to system when no preference saved', async () => {
    window.electronAPI.getConfig = vi.fn().mockResolvedValue({})
    const store = useThemeStore()
    await store.load()
    expect(store.mode).toBe('system')
  })

  it('setMode persists to IPC', async () => {
    const store = useThemeStore()
    await store.setMode('light')
    expect(store.mode).toBe('light')
    expect(window.electronAPI.saveConfig).toHaveBeenCalledWith({ theme: 'light' })
  })
})

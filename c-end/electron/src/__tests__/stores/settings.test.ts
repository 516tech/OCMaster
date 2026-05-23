import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSettingsStore } from '@/stores/settings'

vi.stubGlobal('electronAPI', {
  getConfig: vi.fn(),
  saveConfig: vi.fn()
})

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
})

describe('settings store', () => {
  it('loads defaults when no config saved', async () => {
    window.electronAPI.getConfig = vi.fn().mockResolvedValue({})
    const store = useSettingsStore()
    await store.load()
    expect(store.apiUrl).toBe('https://localhost/api/v1')
    expect(store.language).toBe('zh-CN')
    expect(store.autoScan).toBe(true)
    expect(store.loaded).toBe(true)
  })

  it('loads saved config partially', async () => {
    window.electronAPI.getConfig = vi.fn().mockResolvedValue({ apiUrl: 'https://example.com/api/v1', language: 'en-US' })
    const store = useSettingsStore()
    await store.load()
    expect(store.apiUrl).toBe('https://example.com/api/v1')
    expect(store.language).toBe('en-US')
    expect(store.autoScan).toBe(true) // default preserved
  })

  it('save calls IPC with merged config', async () => {
    window.electronAPI.getConfig = vi.fn().mockResolvedValue({})
    window.electronAPI.saveConfig = vi.fn().mockResolvedValue({ ok: true })
    const store = useSettingsStore()
    store.apiUrl = 'https://test/api/v1'
    store.language = 'en-US'
    await store.save()
    expect(window.electronAPI.saveConfig).toHaveBeenCalledWith({
      apiUrl: 'https://test/api/v1',
      language: 'en-US',
      autoScan: true
    })
  })
})

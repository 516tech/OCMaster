import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useScanStore } from '@/stores/scan'
import { events } from '@/composables/useEventBus'

const scanResult = {
  cpu: { model: 'i7', cores: 8, threads: 16, baseFreq: '3.8 GHz', biosVersion: '' },
  motherboard: { brand: 'ASUS', model: 'Z490', chipset: '', biosVersion: '1203' },
  ram: { totalCapacity: '32 GB', stickCount: 2, channelCount: 2, sticks: [] },
  gpu: { model: 'RTX 3080', vram: '10 GB' },
  psu: { ratedWattage: '', source: 'manual' },
  cooler: { type: 'air', source: 'manual' }
}

function mockElectronAPI(result: ReturnType<typeof buildSuccess>) {
  vi.stubGlobal('electronAPI', {
    scanAll: vi.fn().mockResolvedValue(result)
  })
  // @ts-expect-error mock window
  window.electronAPI = window.electronAPI || {}
  window.electronAPI.scanAll = vi.fn().mockResolvedValue(result)
}

vi.stubGlobal('window', { electronAPI: { scanAll: vi.fn() } })

function buildSuccess() {
  return { success: true, data: JSON.stringify(scanResult) }
}

function buildError() {
  return { success: false, error: 'scanner error' }
}

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
})

describe('scan store', () => {
  it('scan success populates hardware', async () => {
    window.electronAPI.scanAll = vi.fn().mockResolvedValue(buildSuccess())
    const store = useScanStore()
    await store.scan()
    expect(store.hardware?.cpu.model).toBe('i7')
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
    expect(store.progress).toBe(100)
  })

  it('scan error sets error field', async () => {
    window.electronAPI.scanAll = vi.fn().mockResolvedValue(buildError())
    const store = useScanStore()
    await store.scan()
    expect(store.hardware).toBeNull()
    expect(store.error).toBe('scanner error')
  })

  it('scan exception is caught', async () => {
    window.electronAPI.scanAll = vi.fn().mockRejectedValue(new Error('crash'))
    const store = useScanStore()
    await store.scan()
    expect(store.error).toBe('crash')
  })

  it('cancel stops loading and clears hardware', async () => {
    window.electronAPI.scanAll = vi.fn().mockImplementation(() =>
      new Promise((resolve) => setTimeout(() => resolve(buildSuccess()), 500))
    )
    const store = useScanStore()
    const scanPromise = store.scan()
    store.cancel()
    await scanPromise
    expect(store.loading).toBe(false)
    expect(store.hardware).toBeNull()
  })

  it('reset clears all state', () => {
    const store = useScanStore()
    store.hardware = scanResult as any
    store.loading = true
    store.error = 'err'
    store.progress = 50
    store.reset()
    expect(store.hardware).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBe('')
    expect(store.progress).toBe(0)
  })

  it('posts scanStarted and scanCompleted events', async () => {
    window.electronAPI.scanAll = vi.fn().mockResolvedValue(buildSuccess())
    const store = useScanStore()
    const started = vi.fn()
    const completed = vi.fn()
    events.scanStarted.subscribe(started)
    events.scanCompleted.subscribe(completed)
    await store.scan()
    expect(started).toHaveBeenCalledOnce()
    expect(completed).toHaveBeenCalledWith({ success: true })
  })

  it('progress advances during scan', async () => {
    window.electronAPI.scanAll = vi.fn().mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve(buildSuccess()), 400))
    )
    const store = useScanStore()
    const scanPromise = store.scan()
    expect(store.loading).toBe(true)
    expect(store.progress).toBeGreaterThanOrEqual(0)
    await scanPromise
    expect(store.progress).toBe(100)
  })
})

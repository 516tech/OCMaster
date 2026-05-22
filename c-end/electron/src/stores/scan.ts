import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HardwareInfo } from '@/types/hardware'
import { events } from '@/composables/useEventBus'

export const useScanStore = defineStore('scan', () => {
  const hardware = ref<HardwareInfo | null>(null)
  const loading = ref(false)
  const error = ref('')
  const progress = ref(0)
  let abortController: AbortController | null = null

  async function scan() {
    loading.value = true
    error.value = ''
    progress.value = 0
    abortController = new AbortController()

    // Simulate progress updates (Go sidecar doesn't report real progress)
    const progressTimer = setInterval(() => {
      if (progress.value < 90) progress.value += 10
    }, 200)

    events.scanStarted.post()

    try {
      const result = await window.electronAPI.scanAll()
      clearInterval(progressTimer)

      if (abortController.signal.aborted) return

      if (result.success && result.data) {
        progress.value = 100
        hardware.value = JSON.parse(result.data)
        events.scanCompleted.post({ success: true })
      } else {
        error.value = result.error || '扫描失败'
        hardware.value = null
        events.scanCompleted.post({ success: false, error: error.value })
      }
    } catch (e) {
      clearInterval(progressTimer)
      if (abortController.signal.aborted) return
      error.value = e instanceof Error ? e.message : '扫描异常'
      hardware.value = null
      events.scanCompleted.post({ success: false, error: error.value })
    } finally {
      loading.value = false
      abortController = null
    }
  }

  function cancel() {
    abortController?.abort()
    loading.value = false
    progress.value = 0
  }

  function reset() {
    hardware.value = null
    loading.value = false
    error.value = ''
    progress.value = 0
  }

  return { hardware, loading, error, progress, scan, cancel, reset }
})

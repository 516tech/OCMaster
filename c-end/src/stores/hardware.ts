import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HardwareInfo, AppConfig } from '../types/hardware'

declare global { interface Window { ocmaster: any } }

export const useHardwareStore = defineStore('hardware', () => {
  const hardware = ref<HardwareInfo | null>(null)
  const shareCode = ref('')
  const loading = ref(false)
  const error = ref('')
  const config = ref<AppConfig>({ apiUrl: 'https://localhost/api/v1', language: 'zh-CN', autoScan: true })

  async function scan() {
    loading.value = true; error.value = ''
    try {
      hardware.value = await window.ocmaster.scanHardware()
    } catch (e: any) {
      error.value = e?.message || '扫描失败，请重试'
    } finally {
      loading.value = false
    }
  }

  async function upload() {
    if (!hardware.value) return
    loading.value = true; error.value = ''
    try {
      const result = await window.ocmaster.uploadHardware(hardware.value)
      shareCode.value = result.share_code || ''
    } catch (e: any) {
      error.value = e?.message || '上传失败，请检查网络连接'
    } finally {
      loading.value = false
    }
  }

  async function deleteData() {
    if (!shareCode.value) return
    try { await window.ocmaster.deleteHardware(shareCode.value); shareCode.value = '' }
    catch (e: any) { error.value = e?.message || '删除失败' }
  }

  async function loadConfig() {
    try { config.value = await window.ocmaster.getConfig() }
    catch { /* 使用默认配置 */ }
  }

  async function setConfig(key: string, value: unknown) {
    await window.ocmaster.setConfig(key, value)
    ;(config.value as any)[key] = value
  }

  function exportTxt() {
    if (!hardware.value) return
    const h = hardware.value
    const txt = [
      `CPU: ${h.cpu.model} | ${h.cpu.cores}C/${h.cpu.threads}T | ${h.cpu.baseFreq}`,
      `Motherboard: ${h.motherboard.brand} ${h.motherboard.model} | ${h.motherboard.chipset} | BIOS ${h.motherboard.biosVersion}`,
      `RAM: ${h.ram.totalCapacity} | ${h.ram.stickCount} sticks | ${h.ram.channelCount} channels`,
      ...h.ram.sticks.map((s, i) => `  Stick #${i + 1}: ${s.capacity} | ${s.frequency} | ${s.timings} | ${s.dieType}`),
      `GPU: ${h.gpu.model} | ${h.gpu.vram}`,
      `PSU: ${h.psu.ratedWattage} (${h.psu.source})`,
      `Cooler: ${h.cooler.type} (${h.cooler.source})`,
    ].join('\n')
    window.ocmaster.exportTxt(txt)
  }

  return { hardware, shareCode, loading, error, config, scan, upload, deleteData, loadConfig, setConfig, exportTxt }
})

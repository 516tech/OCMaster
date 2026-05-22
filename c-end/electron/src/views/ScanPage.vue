<template>
  <div class="scan-page">
    <div class="scan-controls">
      <el-button
        type="primary"
        size="large"
        :loading="store.loading"
        @click="handleScan"
      >
        {{ store.loading ? '扫描中...' : '开始扫描' }}
      </el-button>
      <el-button v-if="store.loading" type="warning" @click="handleCancel">
        取消扫描
      </el-button>
      <el-text v-if="!store.hardware && !store.loading && !store.error" type="info">
        点击按钮一键扫描您的硬件信息
      </el-text>
    </div>

    <el-progress
      v-if="store.loading"
      :percentage="store.progress"
      :indeterminate="store.progress === 0"
      style="margin: 12px 0"
    />

    <div v-if="store.hardware" class="hardware-result">
      <div class="result-header">
        <h3>硬件信息</h3>
        <el-button @click="handleExport">导出 TXT</el-button>
      </div>

      <el-collapse>
        <el-collapse-item title="CPU 处理器" name="cpu">
          <div class="info-grid">
            <span>型号: {{ store.hardware.cpu.model }}</span>
            <span>核心/线程: {{ store.hardware.cpu.cores }}C / {{ store.hardware.cpu.threads }}T</span>
            <span>基频: {{ store.hardware.cpu.baseFreq }}</span>
          </div>
        </el-collapse-item>

        <el-collapse-item title="主板 Motherboard" name="mb">
          <div class="info-grid">
            <span>品牌: {{ store.hardware.motherboard.brand }}</span>
            <span>型号: {{ store.hardware.motherboard.model }}</span>
            <span>BIOS: {{ store.hardware.motherboard.biosVersion }}</span>
          </div>
        </el-collapse-item>

        <el-collapse-item title="内存 RAM" name="ram">
          <div class="info-grid">
            <span>总容量: {{ store.hardware.ram.totalCapacity }}</span>
            <span>条数: {{ store.hardware.ram.stickCount }}</span>
            <span>通道: {{ store.hardware.ram.channelCount }}</span>
          </div>
        </el-collapse-item>

        <el-collapse-item title="显卡 GPU" name="gpu">
          <div class="info-grid">
            <span>型号: {{ store.hardware.gpu.model }}</span>
            <span>显存: {{ store.hardware.gpu.vram }}</span>
          </div>
        </el-collapse-item>

        <el-collapse-item title="电源 PSU" name="psu">
          <div class="info-grid">
            <span>电源信息需手动输入</span>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useScanStore } from '@/stores/scan'
import { useSettingsStore } from '@/stores/settings'
import { useToast } from '@/composables/useToast'
import { events } from '@/composables/useEventBus'

const store = useScanStore()
const settings = useSettingsStore()
const toast = useToast()

// Listen for scan completion from event bus
events.scanCompleted.subscribe(({ success, error }) => {
  if (success) {
    toast.success('硬件扫描完成')
  } else if (error) {
    toast.error(error)
  }
})

onMounted(async () => {
  await settings.load()
  if (settings.autoScan) handleScan()
})

async function handleScan() {
  await store.scan()
  if (store.error) {
    toast.error(store.error)
  }
}

function handleCancel() {
  store.cancel()
  toast.warning('扫描已取消')
}

async function handleExport() {
  if (!store.hardware) return
  const lines = [
    'OCMaster 硬件扫描报告',
    '========================',
    `CPU: ${store.hardware.cpu.model} | ${store.hardware.cpu.cores}C/${store.hardware.cpu.threads}T | ${store.hardware.cpu.baseFreq}`,
    `Motherboard: ${store.hardware.motherboard.brand} ${store.hardware.motherboard.model} | BIOS ${store.hardware.motherboard.biosVersion}`,
    `RAM: ${store.hardware.ram.totalCapacity} | ${store.hardware.ram.stickCount} sticks | ${store.hardware.ram.channelCount} channels`,
    `GPU: ${store.hardware.gpu.model} | ${store.hardware.gpu.vram}`
  ]
  const result = await window.electronAPI.exportTxt(lines.join('\n'))
  if (result.success) {
    toast.success('导出成功')
  } else if (!result.canceled) {
    toast.error(result.error || '导出失败')
  }
}
</script>

<style scoped>
.scan-page { max-width: 700px; margin: 0 auto; }
.scan-controls { text-align: center; margin: 24px 0 16px; display: flex; flex-direction: column; align-items: center; gap: 12px; }
.result-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.info-grid { display: flex; flex-direction: column; gap: 6px; font-size: 14px; padding: 4px 0; }
</style>

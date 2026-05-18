<script setup lang="ts">
import { useHardwareStore } from '../stores/hardware'
const store = useHardwareStore()
</script>

<template>
  <div class="hw-section">
    <h2>硬件信息 <button class="export-btn" @click="store.exportTxt">导出 TXT</button></h2>
    <table v-if="store.hardware">
      <tr><th>CPU</th><td>{{ store.hardware.cpu.model }} | {{ store.hardware.cpu.cores }}C/{{ store.hardware.cpu.threads }}T | {{ store.hardware.cpu.baseFreq }}</td></tr>
      <tr><th>主板</th><td>{{ store.hardware.motherboard.brand }} {{ store.hardware.motherboard.model }} | {{ store.hardware.motherboard.chipset }} | BIOS {{ store.hardware.motherboard.biosVersion }}</td></tr>
      <tr><th>内存</th><td>{{ store.hardware.ram.totalCapacity }} | {{ store.hardware.ram.stickCount }}条 | {{ store.hardware.ram.channelCount }}通道</td></tr>
      <tr v-for="(s, i) in store.hardware.ram.sticks" :key="i"><th>内存 #{{ i + 1 }}</th><td>{{ s.capacity }} | {{ s.frequency }} | {{ s.timings }} | {{ s.dieType }}</td></tr>
      <tr><th>显卡</th><td>{{ store.hardware.gpu.model }} | {{ store.hardware.gpu.vram }}</td></tr>
      <tr><th>电源</th><td>{{ store.hardware.psu.ratedWattage }} ({{ store.hardware.psu.source === 'smbus' ? '自动检测' : '手动输入' }})</td></tr>
      <tr><th>散热</th><td>{{ store.hardware.cooler.type }} ({{ store.hardware.cooler.source === 'auto' ? '自动检测' : '手动选择' }})</td></tr>
    </table>
  </div>
</template>

<style scoped>
.hw-section { background: #16213e; border-radius: 8px; padding: 20px; }
h2 { font-size: 18px; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center; }
.export-btn { padding: 4px 12px; font-size: 12px; background: #0f3460; color: #e0e0e0; border: none; border-radius: 4px; cursor: pointer; }
table { width: 100%; border-collapse: collapse; }
th { text-align: left; padding: 8px 12px; background: #0f3460; color: #ccc; width: 100px; }
td { padding: 8px 12px; border-bottom: 1px solid #1a1a2e; }
</style>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useHardwareStore } from './stores/hardware'
import ScanButton from './components/ScanButton.vue'
import HardwareTable from './components/HardwareTable.vue'
import UploadPanel from './components/UploadPanel.vue'
import SettingsDialog from './components/SettingsDialog.vue'

const store = useHardwareStore()
onMounted(async () => {
  await store.loadConfig()
  if (store.config.autoScan) await store.scan()
})
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>超频大师 OCMaster</h1>
      <SettingsDialog />
    </header>
    <div v-if="store.error" class="error-banner">
      <span>加载失败：{{ store.error }}</span>
      <button class="retry-btn" @click="store.scan()">重试</button>
    </div>
    <main>
      <ScanButton />
      <HardwareTable v-if="store.hardware" />
      <div v-if="store.loading" class="skeleton skeleton-table" />
      <UploadPanel v-if="store.hardware" />
    </main>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: system-ui, -apple-system, sans-serif; background: #1a1a2e; color: #e0e0e0; }
.app-container { max-width: 900px; margin: 0 auto; padding: 20px; }
.app-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.app-header h1 { font-size: 24px; color: #e0e0e0; }
main { display: flex; flex-direction: column; gap: 16px; }
.error-banner { background: #3a1515; color: #f88; padding: 10px 16px; border-radius: 6px; display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.retry-btn { padding: 4px 12px; background: #e94560; color: #fff; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; }
.skeleton { background: linear-gradient(90deg, #1a1a2e 25%, #2a2a4e 50%, #1a1a2e 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; border-radius: 8px; }
.skeleton-btn { height: 60px; width: 200px; margin: 32px auto; }
.skeleton-table { height: 300px; }
.skeleton-panel { height: 120px; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
</style>

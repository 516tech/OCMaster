<template>
  <div class="about-page">
    <h2>超频大师 OCMaster</h2>
    <p class="version">版本: {{ version }}</p>

    <div class="section">
      <h4>技术栈</h4>
      <p>· Electron + Vue 3 (跨平台 GUI)</p>
      <p>· Go 硬件扫描 (WMI / sysfs / IOKit)</p>
      <p>· Go Backend + Vue 3 商家端</p>
    </div>

    <div class="section">
      <h4>支持平台</h4>
      <p>· Windows 10/11</p>
      <p>· macOS (Apple Silicon)</p>
      <p>· Linux (amd64)</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const version = ref('---')

onMounted(async () => {
  try {
    version.value = await window.electronAPI.getVersion()
  } catch {
    version.value = '0.0.1'
  }
})
</script>

<style scoped>
.about-page { max-width: 400px; margin: 40px auto 0; text-align: center; }
.version { color: var(--el-text-color-secondary); margin: 8px 0 24px; }
.section { text-align: left; margin: 16px 0; }
.section h4 { margin-bottom: 6px; }
.section p { font-size: 14px; color: var(--el-text-color-regular); line-height: 1.8; }
</style>

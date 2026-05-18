<script setup lang="ts">
import { useHardwareStore } from '../stores/hardware'
const store = useHardwareStore()
</script>

<template>
  <div class="scan-section">
    <button class="scan-btn" :disabled="store.loading" @click="store.scan()">
      {{ store.loading ? '扫描中...' : '开始扫描' }}
    </button>
    <p v-if="!store.hardware && !store.loading" class="hint">点击按钮一键扫描您的硬件信息</p>
  </div>
</template>

<style scoped>
.scan-section { text-align: center; padding: 32px 0; }
.scan-btn {
  padding: 14px 48px; font-size: 18px; border: none; border-radius: 8px;
  background: #16213e; color: #e0e0e0; cursor: pointer; transition: background 0.2s;
}
.scan-btn:hover:not(:disabled) { background: #0f3460; }
.scan-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.hint { margin-top: 12px; color: #888; font-size: 14px; }
</style>

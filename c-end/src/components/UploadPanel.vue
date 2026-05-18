<script setup lang="ts">
import { ref } from 'vue'
import { useHardwareStore } from '../stores/hardware'
const store = useHardwareStore()
const uploading = ref(false)

async function handleUpload() { uploading.value = true; await store.upload(); uploading.value = false }
function copyCode() { (window as any).navigator.clipboard.writeText(store.shareCode) }
</script>

<template>
  <div class="upload-section">
    <template v-if="!store.shareCode">
      <p class="notice">上传前确认：您即将上传硬件信息到云端，仅用于超频服务商为您提供建议服务，数据保存 7 天后自动删除。</p>
      <button class="upload-btn" :disabled="uploading" @click="handleUpload">{{ uploading ? '上传中...' : '上传生成分享码' }}</button>
    </template>
    <template v-else>
      <div class="share-code-box">
        <p>分享码（7天有效）：</p>
        <code>{{ store.shareCode }}</code>
        <button class="copy-btn" @click="copyCode">一键复制</button>
        <button class="delete-btn" @click="store.deleteData()">删除云端数据</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.upload-section { background: #16213e; border-radius: 8px; padding: 20px; text-align: center; }
.notice { color: #888; font-size: 13px; margin-bottom: 12px; }
.upload-btn { padding: 10px 32px; font-size: 16px; background: #e94560; color: #fff; border: none; border-radius: 6px; cursor: pointer; }
.upload-btn:disabled { opacity: 0.6; }
.share-code-box { display: flex; align-items: center; gap: 12px; justify-content: center; }
.share-code-box code { font-size: 28px; letter-spacing: 4px; font-weight: bold; background: #0f3460; padding: 8px 16px; border-radius: 4px; }
.copy-btn { padding: 6px 16px; background: #0f3460; color: #e0e0e0; border: none; border-radius: 4px; cursor: pointer; }
.delete-btn { padding: 6px 16px; background: #444; color: #e0e0e0; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; }
</style>

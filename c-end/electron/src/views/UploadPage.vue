<template>
  <div class="upload-page">
    <h3>上传并获取分享码</h3>

    <el-alert
      type="info" :closable="false" show-icon
      title="隐私提示"
      description="您即将上传硬件信息到云端。数据仅用于超频服务商为您提供建议，保存 7 天后自动删除。"
    />

    <el-button
      type="primary" size="large"
      :loading="uploading"
      @click="handleUpload"
      :disabled="uploading"
    >
      {{ uploading ? '上传中...' : '上传硬件信息' }}
    </el-button>

    <el-progress
      v-if="uploading"
      :percentage="50"
      :indeterminate="true"
    />

    <div v-if="shareCode" class="code-section">
      <span>分享码（7天有效）：</span>
      <div class="code-display">{{ shareCode }}</div>
      <el-button-group>
        <el-button @click="handleCopy">一键复制</el-button>
        <el-button type="danger" @click="handleDelete">删除云端数据</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { uploadHardware, deleteHardware } from '@/api/client'
import { useSettingsStore } from '@/stores/settings'
import { useToast } from '@/composables/useToast'

const settings = useSettingsStore()
const toast = useToast()
const shareCode = ref('')
const uploading = ref(false)

onMounted(() => settings.load())

async function handleUpload() {
  uploading.value = true
  try {
    const result = await window.electronAPI.scanAll()
    if (!result.success || !result.data) {
      toast.error(result.error || '扫描硬件失败，请先到「硬件扫描」页确认扫描正常')
      return
    }
    const hw = JSON.parse(result.data)
    const code = await uploadHardware(settings.apiUrl, hw)
    if (code) {
      shareCode.value = code
      toast.success('上传成功，分享码已生成')
    } else {
      toast.error('上传失败，请检查网络连接和 API 地址')
    }
  } catch (e) {
    toast.error(`上传失败: ${e instanceof Error ? e.message : '未知错误'}`)
  } finally {
    uploading.value = false
  }
}

async function handleCopy() {
  await navigator.clipboard.writeText(shareCode.value)
  toast.success('已复制到剪贴板')
}

async function handleDelete() {
  try {
    await deleteHardware(settings.apiUrl, shareCode.value)
    shareCode.value = ''
    toast.success('已删除')
  } catch (e) {
    toast.error(`删除失败: ${e instanceof Error ? e.message : '未知错误'}`)
  }
}
</script>

<style scoped>
.upload-page { max-width: 500px; margin: 0 auto; display: flex; flex-direction: column; gap: 16px; }
.code-section { text-align: center; display: flex; flex-direction: column; gap: 12px; align-items: center; }
.code-display { font-size: 32px; font-weight: bold; font-family: Consolas, monospace; letter-spacing: 8px; background: var(--el-fill-color-light); padding: 12px 24px; border-radius: 8px; }
</style>

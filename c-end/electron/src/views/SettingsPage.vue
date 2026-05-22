<template>
  <div class="settings-page">
    <h3>设置</h3>

    <div class="form-item">
      <label>API 地址</label>
      <el-input v-model="settings.apiUrl" placeholder="https://your-server.com/api/v1" />
    </div>

    <div class="form-item">
      <label>语言</label>
      <el-select v-model="settings.language">
        <el-option label="中文" value="zh-CN" />
        <el-option label="English" value="en-US" />
      </el-select>
    </div>

    <div class="form-item">
      <label>主题</label>
      <el-radio-group v-model="themeMode" @change="onThemeChange">
        <el-radio-button value="system">跟随系统</el-radio-button>
        <el-radio-button value="light">浅色</el-radio-button>
        <el-radio-button value="dark">深色</el-radio-button>
      </el-radio-group>
    </div>

    <div class="form-item">
      <el-checkbox v-model="settings.autoScan">启动时自动扫描</el-checkbox>
    </div>

    <div class="form-actions">
      <el-button @click="$router.back()">取消</el-button>
      <el-button type="primary" @click="handleSave">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useTheme } from '@/composables/useTheme'
import { useToast } from '@/composables/useToast'

const settings = useSettingsStore()
const theme = useTheme()
const toast = useToast()
const themeMode = ref(theme.themeMode)

onMounted(() => settings.load())

function onThemeChange(val: string) {
  theme.applyTheme(val as 'dark' | 'light' | 'system')
}

async function handleSave() {
  await settings.save()
  toast.success('设置已保存')
}
</script>

<style scoped>
.settings-page { max-width: 500px; margin: 0 auto; display: flex; flex-direction: column; gap: 16px; }
.form-item { display: flex; flex-direction: column; gap: 6px; }
.form-item label { font-weight: 600; font-size: 14px; }
.form-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 8px; }
</style>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { profileApi } from '../api/auth'
const auth = useAuthStore(); const router = useRouter()
const riskTemplate = ref(auth.merchant?.risk_template || '')
async function saveTemplate() { await profileApi.update({ risk_template: riskTemplate.value }) }
</script>
<template>
  <div class="page">
    <el-button @click="router.push('/dashboard')">← 返回</el-button>
    <h2>通用模板设置</h2>
    <el-card style="margin-top:12px">
      <p style="margin-bottom:8px;color:#888">自定义风险提示模板（填写超频建议时可一键插入）</p>
      <el-input v-model="riskTemplate" type="textarea" :rows="6" placeholder="例如：超频有风险，请在专业人员指导下操作..." />
      <el-button type="primary" style="margin-top:12px" @click="saveTemplate">保存模板</el-button>
    </el-card>
  </div>
</template>
<style scoped>.page { padding:24px; max-width:700px }</style>

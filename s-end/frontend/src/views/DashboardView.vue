<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter(); const code = ref('')
function lookup() { if (code.value.length === 6) router.push(`/hardware/${code.value}`) }
</script>
<template>
  <div class="page">
    <el-card><h2>超频大师 - 商家平台</h2>
      <p>输入用户提供的 6 位分享码查看硬件信息并生成超频建议</p>
      <el-input v-model="code" placeholder="输入6位分享码" maxlength="6" style="width:200px;margin-right:8px" @keyup.enter="lookup" />
      <el-button type="primary" @click="lookup">查询</el-button>
      <div style="margin-top:16px">
        <el-button @click="router.push('/history')">查询历史</el-button>
        <el-button @click="router.push('/profile')">个人信息</el-button>
        <el-button @click="router.push('/settings')">模板设置</el-button>
      </div>
    </el-card>
  </div>
</template>
<style scoped>.page { padding:40px; background:#f5f5f5; min-height:100vh }</style>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { profileApi } from '../api/auth'
const auth = useAuthStore(); const router = useRouter()
const oldPass = ref(''), newPass = ref('')
async function changePwd() { await profileApi.changePassword(oldPass.value, newPass.value); oldPass.value = ''; newPass.value = '' }
</script>
<template>
  <div class="page">
    <el-button @click="router.push('/dashboard')">← 返回</el-button>
    <h2>个人信息</h2>
    <el-card style="margin-top:12px">
      <p>手机号：{{ auth.merchant?.phone }}</p>
      <p>状态：{{ auth.merchant?.status === 'active' ? '已激活' : '待审核' }}</p>
    </el-card>
    <el-card header="修改密码" style="margin-top:12px">
      <el-input v-model="oldPass" type="password" placeholder="旧密码" style="margin-bottom:8px" />
      <el-input v-model="newPass" type="password" placeholder="新密码" style="margin-bottom:8px" />
      <el-button type="primary" @click="changePwd">修改密码</el-button>
    </el-card>
    <el-button style="margin-top:12px" @click="auth.logout(); router.push('/login')">退出登录</el-button>
  </div>
</template>
<style scoped>.page { padding:24px; max-width:600px }</style>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore()
const router = useRouter()
const phone = ref(''), password = ref('')
async function submit() { await auth.login(phone.value, password.value); router.push('/dashboard') }
</script>
<template>
  <div class="page"><el-card style="width:360px"><h2>商家登录</h2>
    <el-form @submit.prevent="submit"><el-input v-model="phone" placeholder="手机号" /><el-input v-model="password" type="password" placeholder="密码" show-password style="margin:12px 0" /><el-button type="primary" native-type="submit" style="width:100%">登录</el-button></el-form>
    <p style="margin-top:12px"><router-link to="/register">注册新账号</router-link></p>
  </el-card></div>
</template>
<style scoped>.page { display:flex; justify-content:center; align-items:center; min-height:100vh; background:#f5f5f5 } h2 { margin-bottom:16px }</style>

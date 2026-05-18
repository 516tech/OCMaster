<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore(); const router = useRouter()
const phone = ref(''), password = ref('')
async function submit() { await auth.register(phone.value, password.value); router.push('/login') }
</script>
<template>
  <div class="page"><el-card style="width:360px"><h2>商家注册</h2>
    <el-form @submit.prevent="submit"><el-input v-model="phone" placeholder="手机号" /><el-input v-model="password" type="password" placeholder="密码" style="margin:12px 0" /><el-button type="primary" native-type="submit" style="width:100%">注册</el-button></el-form>
    <p style="margin-top:12px;color:#888;font-size:12px">注册后需管理员审核通过方可使用</p>
    <p><router-link to="/login">返回登录</router-link></p>
  </el-card></div>
</template>
<style scoped>.page { display:flex; justify-content:center; align-items:center; min-height:100vh; background:#f5f5f5 } h2 { margin-bottom:16px }</style>

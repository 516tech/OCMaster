<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSuggestionStore } from '../stores/suggestion'
const store = useSuggestionStore(); const router = useRouter()
onMounted(() => store.fetchHistory(0, 20))
</script>
<template>
  <div class="page">
    <el-button @click="router.push('/dashboard')">← 返回</el-button>
    <h2>查询历史</h2>
    <el-table :data="store.history" style="margin-top:12px">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="cpu_suggestion.frequency" label="CPU建议" />
      <el-table-column prop="ram_suggestion.frequency" label="内存建议" />
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作">
        <template #default="{ row }"><el-button size="small" @click="router.push(`/suggestion/${row.id}`)">编辑</el-button><el-button size="small" @click="store.downloadPdf(row.id)">PDF</el-button></template>
      </el-table-column>
    </el-table>
  </div>
</template>
<style scoped>.page { padding:24px }</style>

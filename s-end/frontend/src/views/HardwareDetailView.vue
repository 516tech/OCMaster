<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useHardwareStore } from '../stores/hardware'
const route = useRoute(); const router = useRouter()
const store = useHardwareStore()
onMounted(() => store.fetchByCode(route.params.code as string))
</script>
<template>
  <div class="page">
    <el-button @click="router.push('/dashboard')">← 返回</el-button>
    <h2>硬件信息</h2>
    <el-table v-if="store.hardware" :data="[store.hardware.cpu]" style="margin-top:12px">
      <el-table-column prop="model" label="CPU 型号" /><el-table-column prop="cores" label="核心" /><el-table-column prop="threads" label="线程" /><el-table-column prop="base_freq" label="基础频率" />
    </el-table>
    <el-table v-if="store.hardware" :data="[store.hardware.motherboard]" style="margin-top:12px">
      <el-table-column prop="brand" label="主板品牌" /><el-table-column prop="model" label="型号" /><el-table-column prop="chipset" label="芯片组" /><el-table-column prop="bios_version" label="BIOS" />
    </el-table>
    <el-table v-if="store.hardware" :data="store.hardware.ram.sticks" style="margin-top:12px">
      <el-table-column prop="capacity" label="容量" /><el-table-column prop="frequency" label="频率" /><el-table-column prop="timings" label="时序" /><el-table-column prop="die_type" label="颗粒" />
    </el-table>
    <el-descriptions v-if="store.hardware" :column="2" border style="margin-top:12px">
      <el-descriptions-item label="显卡">{{ store.hardware.gpu.model }} / {{ store.hardware.gpu.vram }}</el-descriptions-item>
      <el-descriptions-item label="电源">{{ store.hardware.psu.rated_wattage }} ({{ store.hardware.psu.source }})</el-descriptions-item>
      <el-descriptions-item label="散热">{{ store.hardware.cooler.type }} ({{ store.hardware.cooler.source }})</el-descriptions-item>
    </el-descriptions>
    <el-button type="primary" style="margin-top:16px" @click="router.push(`/suggestion/new/${route.params.code}`)">填写超频建议</el-button>
  </div>
</template>
<style scoped>.page { padding:24px }</style>

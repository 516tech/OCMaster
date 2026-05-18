<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSuggestionStore } from '../stores/suggestion'
const route = useRoute(); const router = useRouter()
const store = useSuggestionStore()

const cpuFreq = ref(''), cpuVoltage = ref(''), cpuNotes = ref('')
const ramFreq = ref(''), ramTimings = ref(''), ramVoltage = ref(''), ramNotes = ref('')
const testTool = ref(''), testDuration = ref(''), riskWarning = ref('')

onMounted(async () => {
  if (route.params.id) { await store.fetchById(Number(route.params.id)); const s = store.current; if (s) {
    cpuFreq.value = s.cpu_suggestion?.frequency || ''; cpuVoltage.value = s.cpu_suggestion?.voltage || ''; cpuNotes.value = s.cpu_suggestion?.notes || ''
    ramFreq.value = s.ram_suggestion?.frequency || ''; ramTimings.value = s.ram_suggestion?.timings || ''; ramVoltage.value = s.ram_suggestion?.voltage || ''; ramNotes.value = s.ram_suggestion?.notes || ''
    testTool.value = s.stability_test?.tool || ''; testDuration.value = s.stability_test?.duration || ''; riskWarning.value = s.risk_warning || ''
  }}
})

async function save() {
  await store.create({
    hardware_id: 0, cpu_suggestion: { frequency: cpuFreq.value, voltage: cpuVoltage.value, notes: cpuNotes.value },
    ram_suggestion: { frequency: ramFreq.value, timings: ramTimings.value, voltage: ramVoltage.value, notes: ramNotes.value },
    stability_test: { tool: testTool.value, duration: testDuration.value }, risk_warning: riskWarning.value,
  })
  router.push('/history')
}
</script>
<template>
  <div class="page">
    <h2>超频建议</h2>
    <el-card header="CPU 超频建议" style="margin-bottom:12px">
      <el-input v-model="cpuFreq" placeholder="建议频率 (如 5.5GHz)" style="margin-bottom:8px" />
      <el-input v-model="cpuVoltage" placeholder="建议电压 (如 1.35V)" style="margin-bottom:8px" />
      <el-input v-model="cpuNotes" type="textarea" placeholder="操作要点" />
    </el-card>
    <el-card header="内存超频建议" style="margin-bottom:12px">
      <el-input v-model="ramFreq" placeholder="建议频率 (如 6400MHz)" style="margin-bottom:8px" />
      <el-input v-model="ramTimings" placeholder="时序 (如 32-40-40-84)" style="margin-bottom:8px" />
      <el-input v-model="ramVoltage" placeholder="电压 (如 1.40V)" style="margin-bottom:8px" />
      <el-input v-model="ramNotes" type="textarea" placeholder="操作要点" />
    </el-card>
    <el-card header="稳定性测试" style="margin-bottom:12px">
      <el-input v-model="testTool" placeholder="测试工具 (如 Prime95)" style="margin-bottom:8px" />
      <el-input v-model="testDuration" placeholder="建议测试时长 (如 2小时)" />
    </el-card>
    <el-card header="风险提示">
      <el-input v-model="riskWarning" type="textarea" :rows="3" placeholder="风险提示与注意事项" />
    </el-card>
    <el-button type="primary" style="margin-top:16px" @click="save">保存建议</el-button>
  </div>
</template>
<style scoped>.page { padding:24px; max-width:800px }</style>

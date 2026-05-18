<script setup lang="ts">
import { ref } from 'vue'
import { useHardwareStore } from '../stores/hardware'
const store = useHardwareStore()
const visible = ref(false)
const apiUrl = ref(store.config.apiUrl)

function open() { apiUrl.value = store.config.apiUrl; visible.value = true }
function save() { store.setConfig('apiUrl', apiUrl.value); visible.value = false }
</script>

<template>
  <div>
    <button class="settings-btn" @click="open">设置</button>
    <div v-if="visible" class="overlay" @click.self="visible = false">
      <div class="dialog">
        <h3>设置</h3>
        <label>API 服务器地址</label>
        <input v-model="apiUrl" type="text" class="input" />
        <div class="actions">
          <button class="cancel" @click="visible = false">取消</button>
          <button class="save" @click="save">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-btn { padding: 6px 16px; background: #333; color: #ccc; border: none; border-radius: 4px; cursor: pointer; font-size: 13px; }
.overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 100; }
.dialog { background: #16213e; border-radius: 8px; padding: 24px; width: 400px; }
h3 { margin-bottom: 16px; }
label { display: block; margin-bottom: 6px; font-size: 13px; color: #888; }
.input { width: 100%; padding: 8px 12px; background: #0f3460; border: 1px solid #333; border-radius: 4px; color: #e0e0e0; font-size: 14px; }
.actions { margin-top: 16px; display: flex; gap: 8px; justify-content: flex-end; }
.cancel, .save { padding: 6px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 13px; }
.cancel { background: #444; color: #ccc; }
.save { background: #e94560; color: #fff; }
</style>

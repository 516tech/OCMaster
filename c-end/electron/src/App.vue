<template>
  <el-container class="app-container" :class="{ dark: isDark }">
    <el-header class="app-header">
      <div class="header-title">
        <h2>超频大师 OCMaster</h2>
      </div>
      <el-menu
        :default-active="activeRoute"
        mode="horizontal"
        router
        class="nav-menu"
      >
        <el-menu-item index="/">
          <el-icon><Search /></el-icon>
          <span>硬件扫描</span>
        </el-menu-item>
        <el-menu-item index="/upload">
          <el-icon><Upload /></el-icon>
          <span>上传分享</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>设置</span>
        </el-menu-item>
        <el-menu-item index="/about">
          <el-icon><InfoFilled /></el-icon>
          <span>关于</span>
        </el-menu-item>
      </el-menu>
      <el-button
        class="theme-toggle"
        :icon="isDark ? Sunny : Moon"
        circle
        @click="theme.toggle()"
      />
    </el-header>
    <el-main class="app-main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Search, Upload, Setting, InfoFilled, Sunny, Moon } from '@element-plus/icons-vue'
import { useTheme } from '@/composables/useTheme'
import { useThemeStore } from '@/stores/theme'
import { onMounted } from 'vue'

const route = useRoute()
const activeRoute = computed(() => route.path)
const theme = useTheme()
const themeStore = useThemeStore()
const { isDark } = theme

onMounted(() => themeStore.load())
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
.app-container { height: 100vh; }
.app-container.dark { background: #1a1a2e; color: #e0e0e0; }
.app-header { display: flex; align-items: center; padding: 0 16px; border-bottom: 1px solid var(--el-border-color-light); }
.header-title { margin-right: 32px; white-space: nowrap; }
.nav-menu { flex: 1; border-bottom: none !important; }
.theme-toggle { margin-left: 8px; }
.app-main { padding: 24px; overflow-y: auto; }
</style>

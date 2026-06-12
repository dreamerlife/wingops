<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Settings</p>
        <h1>系统配置</h1>
      </div>
      <RouterLink to="/dashboard" class="back-link">返回首页</RouterLink>
    </header>

    <section class="panel">
      <div class="toolbar">
        <el-input v-model="configForm.key" placeholder="配置键" />
        <el-input v-model="configForm.value" placeholder="配置值" />
        <el-button type="primary" @click="saveConfig">保存配置</el-button>
      </div>
    </section>

    <el-table :data="configs" border>
      <el-table-column prop="key" label="配置键" min-width="180" />
      <el-table-column prop="value" label="配置值" min-width="240" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

import { listSystemConfigs, saveSystemConfig, type SystemConfig } from '../api/platform'

const configs = ref<SystemConfig[]>([])
const configForm = reactive({ key: 'platform.name', value: '' })

async function loadConfigs() {
  configs.value = await listSystemConfigs()
}

async function saveConfig() {
  await saveSystemConfig(configForm)
  await loadConfigs()
  ElMessage.success('配置已保存')
}

onMounted(loadConfigs)
</script>

<style scoped>
.page-shell { min-height: 100vh; padding: 28px; background: #f6f8fb; }
.page-header, .toolbar { display: flex; align-items: center; gap: 12px; }
.page-header { justify-content: space-between; margin-bottom: 20px; }
.panel { margin-bottom: 16px; padding: 16px; border: 1px solid #d7dee8; border-radius: 8px; background: #ffffff; }
.toolbar .el-input { width: 260px; }
.eyebrow { margin: 0 0 4px; color: #2563eb; font-size: 13px; font-weight: 700; }
h1 { margin: 0; color: #111827; font-size: 24px; }
.back-link { color: #1d4ed8; font-weight: 600; text-decoration: none; }
</style>

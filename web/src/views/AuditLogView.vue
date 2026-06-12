<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Governance</p>
        <h1>审计日志</h1>
      </div>
      <RouterLink to="/dashboard" class="back-link">返回首页</RouterLink>
    </header>

    <el-table :data="logs" border>
      <el-table-column prop="actor_id" label="操作者" min-width="180" />
      <el-table-column prop="method" label="方法" width="100" />
      <el-table-column prop="path" label="路径" min-width="260" />
      <el-table-column prop="status_code" label="状态码" width="100" />
      <el-table-column prop="created_at" label="时间" min-width="200" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { listAuditLogs, type AuditLog } from '../api/platform'

const logs = ref<AuditLog[]>([])

onMounted(async () => {
  logs.value = await listAuditLogs()
})
</script>

<style scoped>
.page-shell { min-height: 100vh; padding: 28px; background: #f6f8fb; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.eyebrow { margin: 0 0 4px; color: #2563eb; font-size: 13px; font-weight: 700; }
h1 { margin: 0; color: #111827; font-size: 24px; }
.back-link { color: #1d4ed8; font-weight: 600; text-decoration: none; }
</style>

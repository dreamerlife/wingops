<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Governance</p>
        <h1>审计日志</h1>
      </div>
    </header>

    <el-table :data="logs" border>
      <el-table-column prop="actor" label="操作者" min-width="160" />
      <el-table-column prop="method" label="方法" width="100" />
      <el-table-column prop="path" label="路径" min-width="220" />
      <el-table-column prop="status" label="状态码" width="100" />
      <el-table-column prop="createdAt" label="时间" min-width="180" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const logs = [
  {
    actor: 'admin',
    method: 'PUT',
    path: '/api/v1/system/configs/platform.name',
    status: 200,
    createdAt: '实时记录'
  }
]
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px;
  background: #f6f8fb;
}

.page-header {
  margin-bottom: 20px;
}

.eyebrow {
  margin: 0 0 4px;
  color: #2563eb;
  font-size: 13px;
  font-weight: 700;
}

h1 {
  margin: 0;
  color: #111827;
  font-size: 24px;
}
</style>

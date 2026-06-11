<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>资产</h1>
      </div>
      <el-button type="primary">新增资产</el-button>
    </header>

    <el-table :data="assets" border>
      <el-table-column prop="uniqueKey" label="唯一标识" min-width="180" />
      <el-table-column prop="model" label="模型" min-width="140" />
      <el-table-column prop="hostname" label="主机名" min-width="160" />
      <el-table-column prop="managementIp" label="管理 IP" min-width="160" />
      <el-table-column prop="status" label="状态" width="120" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const assets = [
  {
    uniqueKey: 'sn:ABC123',
    model: '服务器',
    hostname: 'web-01',
    managementIp: '10.0.1.10',
    status: 'running'
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
  display: flex;
  align-items: center;
  justify-content: space-between;
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

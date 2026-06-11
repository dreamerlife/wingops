<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>资产导入</h1>
      </div>
      <el-button type="primary">上传 CSV</el-button>
    </header>

    <el-table :data="rows" border>
      <el-table-column prop="uniqueKey" label="唯一标识" min-width="180" />
      <el-table-column prop="hostname" label="主机名" min-width="160" />
      <el-table-column prop="managementIp" label="管理 IP" min-width="160" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const rows = [
  {
    uniqueKey: 'sn:ABC123',
    hostname: 'web-01',
    managementIp: '10.0.1.10'
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

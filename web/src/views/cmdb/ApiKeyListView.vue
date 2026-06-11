<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>API Key</h1>
      </div>
      <el-button type="primary">新增 Key</el-button>
    </header>

    <el-table :data="keys" border>
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="keyId" label="Key ID" min-width="180" />
      <el-table-column prop="status" label="状态" width="120" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const keys = [
  {
    name: '开发同步 Key',
    keyId: 'dev-sync-key',
    status: 'active'
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

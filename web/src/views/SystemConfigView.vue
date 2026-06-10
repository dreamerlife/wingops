<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Settings</p>
        <h1>系统配置</h1>
      </div>
      <el-button type="primary">保存配置</el-button>
    </header>

    <el-table :data="configs" border>
      <el-table-column prop="key" label="配置键" min-width="180" />
      <el-table-column prop="value" label="配置值" min-width="240" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const configs = [
  {
    key: 'platform.name',
    value: 'WingOps'
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

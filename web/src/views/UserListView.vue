<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Access</p>
        <h1>用户</h1>
      </div>
      <el-button type="primary">新增用户</el-button>
    </header>

    <el-table :data="users" border>
      <el-table-column prop="username" label="用户名" min-width="140" />
      <el-table-column prop="displayName" label="显示名称" min-width="160" />
      <el-table-column prop="role" label="角色" min-width="160" />
      <el-table-column prop="status" label="状态" width="120" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const users = [
  {
    username: 'admin',
    displayName: '管理员',
    role: '系统管理员',
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

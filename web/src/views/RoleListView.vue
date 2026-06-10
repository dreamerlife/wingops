<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Access</p>
        <h1>角色</h1>
      </div>
      <el-button type="primary">新增角色</el-button>
    </header>

    <el-table :data="roles" border>
      <el-table-column prop="name" label="角色标识" min-width="160" />
      <el-table-column prop="displayName" label="显示名称" min-width="180" />
      <el-table-column prop="permissions" label="权限数量" width="120" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const roles = [
  { name: 'system_admin', displayName: '系统管理员', permissions: 4 },
  { name: 'ops_admin', displayName: '运维管理员', permissions: 2 },
  { name: 'ops_operator', displayName: '运维操作员', permissions: 1 },
  { name: 'readonly', displayName: '只读用户', permissions: 1 }
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

<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>模型组</h1>
      </div>
      <el-button type="primary">新增模型组</el-button>
    </header>

    <el-table :data="groups" border>
      <el-table-column prop="name" label="模型组标识" min-width="160" />
      <el-table-column prop="displayName" label="显示名称" min-width="180" />
      <el-table-column prop="description" label="描述" min-width="240" />
    </el-table>
  </main>
</template>

<script setup lang="ts">
const groups = [
  {
    name: 'compute',
    displayName: '计算资源',
    description: '服务器、虚拟机和容器节点'
  },
  {
    name: 'network',
    displayName: '网络设备',
    description: '交换机、路由器和负载均衡'
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

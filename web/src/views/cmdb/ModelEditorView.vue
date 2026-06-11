<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>模型编辑</h1>
      </div>
      <el-button type="primary">保存模型</el-button>
    </header>

    <section class="editor-grid">
      <el-form label-position="top">
        <el-form-item label="模型标识">
          <el-input model-value="server" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input model-value="服务器" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input model-value="Linux/Windows 主机资产模型" />
        </el-form-item>
      </el-form>

      <el-table :data="fields" border>
        <el-table-column prop="name" label="字段名" min-width="160" />
        <el-table-column prop="displayName" label="显示名称" min-width="160" />
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="required" label="必填" width="100" />
      </el-table>
    </section>
  </main>
</template>

<script setup lang="ts">
const fields = [
  { name: 'hostname', displayName: '主机名', type: 'text', required: '是' },
  { name: 'management_ip', displayName: '管理 IP', type: 'ip', required: '是' }
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

.editor-grid {
  display: grid;
  grid-template-columns: minmax(260px, 360px) 1fr;
  gap: 20px;
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

@media (max-width: 860px) {
  .editor-grid {
    grid-template-columns: 1fr;
  }
}
</style>

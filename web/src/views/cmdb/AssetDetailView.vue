<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>资产详情</h1>
      </div>
      <el-button type="primary">保存资产</el-button>
    </header>

    <section class="detail-grid">
      <DynamicAssetForm :fields="fields" :values="values" />

      <el-table :data="changeLogs" border>
        <el-table-column prop="actor" label="操作者" min-width="150" />
        <el-table-column prop="before" label="变更前" min-width="180" />
        <el-table-column prop="after" label="变更后" min-width="180" />
        <el-table-column prop="createdAt" label="时间" min-width="160" />
      </el-table>
    </section>
  </main>
</template>

<script setup lang="ts">
import DynamicAssetForm from '../../components/DynamicAssetForm.vue'

const fields = [
  { name: 'hostname', displayName: '主机名' },
  { name: 'management_ip', displayName: '管理 IP' }
]

const values = {
  hostname: 'web-01',
  management_ip: '10.0.1.10'
}

const changeLogs = [
  {
    actor: 'admin',
    before: '{}',
    after: '{"hostname":"web-01"}',
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.detail-grid {
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
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>

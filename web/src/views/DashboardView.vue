<template>
  <main class="dashboard-page">
    <header class="page-header">
      <div>
        <p class="eyebrow">WingOps</p>
        <h1>平台首页</h1>
      </div>
    </header>

    <section class="metric-grid">
      <article v-for="metric in metrics" :key="metric.label" class="metric-card">
        <span>{{ metric.label }}</span>
        <strong>{{ metric.value }}</strong>
      </article>
    </section>

    <section class="quick-links">
      <RouterLink to="/cmdb/model-groups">CMDB 模型</RouterLink>
      <RouterLink to="/cmdb/assets">资产管理</RouterLink>
      <RouterLink to="/audit/logs">审计日志</RouterLink>
      <RouterLink to="/system/configs">系统配置</RouterLink>
    </section>
  </main>
</template>

<script setup lang="ts">
const metrics = [
  { label: '模型组', value: '0' },
  { label: '资产', value: '0' },
  { label: '告警', value: '0' },
  { label: '同步任务', value: '0' }
]
</script>

<style scoped>
.dashboard-page {
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

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(140px, 1fr));
  gap: 14px;
  margin-bottom: 22px;
}

.metric-card {
  display: grid;
  gap: 8px;
  padding: 18px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
}

.metric-card span {
  color: #64748b;
  font-size: 13px;
}

.metric-card strong {
  color: #111827;
  font-size: 28px;
}

.quick-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.quick-links a {
  padding: 10px 14px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
  color: #1d4ed8;
  font-weight: 600;
  text-decoration: none;
}

@media (max-width: 860px) {
  .metric-grid {
    grid-template-columns: repeat(2, minmax(140px, 1fr));
  }
}
</style>

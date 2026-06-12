<template>
  <main class="dashboard-page">
    <header class="hero">
      <div>
        <p class="eyebrow">M0-M2 基础平台与 CMDB</p>
        <h1>运维资产与平台治理工作台</h1>
        <p>集中维护模型、资产、API 同步、权限审计和系统配置，为后续监控与告警联动提供可信资产底座。</p>
      </div>
    </header>

    <section class="metric-grid">
      <article v-for="metric in metrics" :key="metric.label" class="metric-card">
        <span>{{ metric.label }}</span>
        <strong>{{ metric.value }}</strong>
        <small>{{ metric.hint }}</small>
      </article>
    </section>

    <section class="work-grid">
      <article class="panel">
        <div class="panel-heading">
          <h2>CMDB 核心流程</h2>
          <RouterLink to="/cmdb/assets">进入资产台账</RouterLink>
        </div>
        <div class="flow-row">
          <span>模型定义</span>
          <span>资产录入</span>
          <span>分组治理</span>
          <span>API 同步</span>
        </div>
      </article>

      <article class="panel">
        <div class="panel-heading">
          <h2>平台基础模块</h2>
          <RouterLink to="/audit/logs">查看审计</RouterLink>
        </div>
        <ul class="status-list">
          <li><span>认证与 RBAC</span><strong>已接入</strong></li>
          <li><span>审计日志</span><strong>自动记录</strong></li>
          <li><span>系统配置</span><strong>可读写</strong></li>
        </ul>
      </article>

      <article class="panel muted-panel">
        <div class="panel-heading">
          <h2>后续联动</h2>
          <span>M3-M4</span>
        </div>
        <p>监控面板、活跃告警和历史告警将在 M3/M4 后接入。当前资产详情页保留真实入口和空状态，不伪造监控数据。</p>
      </article>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { listApiKeys, listAssets, listModelGroups } from '../api/cmdb'

const modelGroupCount = ref(0)
const assetCount = ref(0)
const apiKeyCount = ref(0)

const metrics = computed(() => [
  { label: '模型组', value: String(modelGroupCount.value), hint: '计算、网络、中间件、业务应用' },
  { label: '资产', value: String(assetCount.value), hint: '支持模型、分组、状态筛选' },
  { label: 'API Key', value: String(apiKeyCount.value), hint: '外部系统签名同步凭据' },
  { label: '联动能力', value: 'M3-M4', hint: '监控与告警接入阶段' }
])

async function loadDashboard() {
  const [groups, assets, keys] = await Promise.all([listModelGroups(), listAssets(), listApiKeys()])
  modelGroupCount.value = groups.length
  assetCount.value = assets.total
  apiKeyCount.value = keys.length
}

onMounted(loadDashboard)
</script>

<style scoped>
.dashboard-page {
  min-height: 100vh;
  padding: 28px;
  background: #f3f6f8;
  color: #172033;
}

.hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  min-height: 180px;
  margin-bottom: 18px;
  padding: 28px;
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(35, 198, 184, 0.16), transparent 38%),
    linear-gradient(115deg, #ffffff 0%, #eef6ff 56%, #f8fbf2 100%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #0f7b74;
  font-size: 12px;
  font-weight: 800;
}

h1,
h2 {
  margin: 0;
  letter-spacing: 0;
}

h1 {
  font-size: 30px;
}

.hero p {
  max-width: 760px;
  margin: 10px 0 0;
  color: #667085;
}

.metric-grid,
.work-grid {
  display: grid;
  gap: 16px;
}

.metric-grid {
  grid-template-columns: repeat(4, minmax(140px, 1fr));
  margin-bottom: 16px;
}

.metric-card,
.panel {
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #fff;
}

.metric-card {
  display: grid;
  gap: 8px;
  padding: 18px;
}

.metric-card span,
.metric-card small,
.muted-panel p {
  color: #667085;
}

.metric-card span {
  font-size: 13px;
}

.metric-card strong {
  font-size: 30px;
}

.work-grid {
  grid-template-columns: 1.1fr 0.9fr 0.9fr;
}

.panel {
  padding: 18px;
}

.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.panel-heading h2 {
  font-size: 18px;
}

.panel-heading a {
  color: #1f6feb;
  font-weight: 700;
  text-decoration: none;
}

.flow-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.flow-row span {
  padding: 12px;
  border: 1px solid #d9e0ea;
  border-radius: 6px;
  background: #fafbfc;
  font-weight: 700;
  text-align: center;
}

.status-list {
  display: grid;
  gap: 10px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.status-list li {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.status-list strong {
  color: #207a4b;
}

@media (max-width: 980px) {
  .metric-grid,
  .work-grid,
  .flow-row {
    grid-template-columns: 1fr;
  }
}
</style>

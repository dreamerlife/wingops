<template>
  <main class="detail-page" v-if="asset">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB / 资产详情</p>
        <h1>{{ asset.attributes.name || asset.unique_key }}</h1>
        <p class="subtitle">{{ modelName(asset.model_id) }} · {{ asset.unique_key }}</p>
      </div>
      <RouterLink to="/cmdb/assets" class="back-link">返回资产列表</RouterLink>
    </header>

    <section class="summary-grid">
      <article class="summary-card">
        <span>生命周期</span>
        <strong>{{ statusLabel(asset.status) }}</strong>
      </article>
      <article class="summary-card">
        <span>管理 IP</span>
        <strong>{{ asset.attributes.management_ip || '-' }}</strong>
      </article>
      <article class="summary-card">
        <span>负责人</span>
        <strong>{{ asset.attributes.owner || '-' }}</strong>
      </article>
      <article class="summary-card">
        <span>资产分组</span>
        <strong>{{ groupLabels }}</strong>
      </article>
    </section>

    <section class="detail-grid">
      <section class="panel">
        <div class="panel-heading">
          <h2>资产属性</h2>
          <el-tag effect="plain">{{ modelName(asset.model_id) }}</el-tag>
        </div>
        <el-descriptions border :column="2">
          <el-descriptions-item label="唯一标识">{{ asset.unique_key }}</el-descriptions-item>
          <el-descriptions-item label="生命周期">{{ statusLabel(asset.status) }}</el-descriptions-item>
          <el-descriptions-item v-for="(value, key) in asset.attributes" :key="key" :label="fieldLabel(String(key))">
            {{ value || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </section>

      <section class="panel">
        <div class="panel-heading">
          <h2>变更历史</h2>
          <span>{{ changeLogs.length }} 条记录</span>
        </div>
        <el-timeline v-if="changeLogs.length">
          <el-timeline-item v-for="log in changeLogs" :key="log.id" :timestamp="formatTime(log.created_at)">
            <div class="change-card">
              <strong>{{ log.actor_id || '系统' }}</strong>
              <p>变更前：{{ stringify(log.before_value) }}</p>
              <p>变更后：{{ stringify(log.after_value) }}</p>
            </div>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else description="暂无变更历史" />
      </section>
    </section>

    <section class="ops-grid">
      <article class="ops-panel">
        <h2>关联监控</h2>
        <p>监控模块完成后，将按 asset_id 展示主机、网络或中间件实时指标。</p>
      </article>
      <article class="ops-panel">
        <h2>活跃告警</h2>
        <p>告警模块完成后，将展示当前未恢复告警并支持跳转处理。</p>
      </article>
      <article class="ops-panel">
        <h2>历史告警</h2>
        <p>这里预留资产维度的告警追溯入口，不伪造 M3-M4 数据。</p>
      </article>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import {
  getAsset,
  getModel,
  listAssetChangeLogs,
  listAssetGroups,
  listModelGroups,
  listModels,
  type Asset,
  type AssetChangeLog,
  type AssetGroup,
  type Model
} from '../../api/cmdb'

const route = useRoute()
const asset = ref<Asset | null>(null)
const model = ref<Model | null>(null)
const allModels = ref<Model[]>([])
const assetGroups = ref<AssetGroup[]>([])
const changeLogs = ref<AssetChangeLog[]>([])

const groupLabels = computed(() => {
  const labels = (asset.value?.group_ids ?? []).map((id) => assetGroups.value.find((group) => group.id === id)?.display_name ?? id)
  return labels.length ? labels.join(' / ') : '未分组'
})

async function loadDetail() {
  const id = String(route.params.id)
  asset.value = await getAsset(id)
  const [currentModel, groups, modelGroups, logs] = await Promise.all([
    getModel(asset.value.model_id),
    listAssetGroups(),
    listModelGroups(),
    listAssetChangeLogs(id)
  ])
  model.value = currentModel
  assetGroups.value = groups
  changeLogs.value = logs
  const modelLists = await Promise.all(modelGroups.map((group) => listModels(group.id)))
  allModels.value = modelLists.flat()
}

function fieldLabel(name: string) {
  return model.value?.fields.find((field) => field.name === name)?.display_name ?? name
}

function modelName(modelId: string) {
  return allModels.value.find((item) => item.id === modelId)?.display_name ?? model.value?.display_name ?? modelId
}

function statusLabel(value: string) {
  const labels: Record<string, string> = {
    purchased: '采购',
    racked: '上架',
    running: '运行',
    maintenance: '维修',
    retired: '报废'
  }
  return labels[value] ?? value
}

function formatTime(value: string) {
  return new Date(value).toLocaleString()
}

function stringify(value: Record<string, unknown>) {
  return Object.keys(value ?? {}).length ? JSON.stringify(value) : '空'
}

onMounted(loadDetail)
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  padding: 28px;
  background: #f4f6f8;
  color: #172033;
}

.page-header,
.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.page-header {
  margin-bottom: 20px;
}

.eyebrow {
  margin: 0 0 4px;
  color: #1f6feb;
  font-size: 12px;
  font-weight: 700;
}

h1,
h2,
p {
  margin-top: 0;
}

h1 {
  margin-bottom: 6px;
  font-size: 26px;
  letter-spacing: 0;
}

h2 {
  margin-bottom: 14px;
  font-size: 17px;
}

.subtitle,
.panel-heading span,
.ops-panel p {
  color: #667085;
}

.summary-grid,
.detail-grid,
.ops-grid {
  display: grid;
  gap: 16px;
}

.summary-grid {
  grid-template-columns: repeat(4, minmax(140px, 1fr));
  margin-bottom: 16px;
}

.detail-grid {
  grid-template-columns: minmax(0, 1.1fr) minmax(360px, 0.9fr);
}

.ops-grid {
  grid-template-columns: repeat(3, minmax(180px, 1fr));
  margin-top: 16px;
}

.summary-card,
.panel,
.ops-panel {
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #fff;
}

.summary-card {
  display: grid;
  gap: 8px;
  padding: 16px;
}

.summary-card span {
  color: #667085;
  font-size: 13px;
}

.summary-card strong {
  font-size: 20px;
}

.panel,
.ops-panel {
  padding: 16px;
}

.change-card {
  color: #344054;
}

.change-card p {
  margin: 6px 0 0;
  word-break: break-word;
}

.back-link {
  color: #1f6feb;
  font-weight: 600;
  text-decoration: none;
}

@media (max-width: 980px) {
  .summary-grid,
  .detail-grid,
  .ops-grid {
    grid-template-columns: 1fr;
  }
}
</style>

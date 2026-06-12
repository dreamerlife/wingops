<template>
  <main class="asset-page">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB / 资产台账</p>
        <h1>资产管理</h1>
        <p class="subtitle">按模型、分组、生命周期和关键字维护企业资产主数据。</p>
      </div>
      <div class="header-actions">
        <RouterLink to="/cmdb/assets/import" class="text-link">批量导入</RouterLink>
        <el-button type="primary" @click="openCreate">新增资产</el-button>
      </div>
    </header>

    <section class="workspace">
      <aside class="side-panel">
        <div class="panel-title">
          <strong>资产分组</strong>
          <el-button link type="primary" @click="groupDialogVisible = true">新增</el-button>
        </div>
        <button
          class="group-item"
          :class="{ active: !filters.group_id }"
          type="button"
          @click="selectGroup('')"
        >
          <span>全部资产</span>
          <small>{{ total }}</small>
        </button>
        <button
          v-for="group in assetGroups"
          :key="group.id"
          class="group-item"
          :class="{ active: filters.group_id === group.id }"
          type="button"
          @click="selectGroup(group.id)"
        >
          <span>{{ group.display_name }}</span>
          <small>{{ dimensionLabel(group.dimension) }}</small>
        </button>
      </aside>

      <section class="main-panel">
        <div class="filter-bar">
          <el-input v-model="filters.keyword" clearable placeholder="搜索唯一标识、IP、负责人、名称" @keyup.enter="reloadAssets" />
          <el-select v-model="filters.model_id" clearable filterable placeholder="资产模型" @change="reloadAssets">
            <el-option v-for="model in allModels" :key="model.id" :label="model.display_name" :value="model.id" />
          </el-select>
          <el-select v-model="filters.status" clearable placeholder="生命周期" @change="reloadAssets">
            <el-option v-for="status in statusOptions" :key="status.value" :label="status.label" :value="status.value" />
          </el-select>
          <el-button @click="reloadAssets">查询</el-button>
        </div>

        <el-table :data="assets" class="asset-table" border>
          <el-table-column prop="unique_key" label="唯一标识" min-width="180" />
          <el-table-column label="模型" min-width="140">
            <template #default="{ row }">{{ modelName(row.model_id) }}</template>
          </el-table-column>
          <el-table-column label="名称" min-width="160">
            <template #default="{ row }">{{ row.attributes.name || '-' }}</template>
          </el-table-column>
          <el-table-column label="管理 IP" min-width="150">
            <template #default="{ row }">{{ row.attributes.management_ip || '-' }}</template>
          </el-table-column>
          <el-table-column label="负责人" min-width="120">
            <template #default="{ row }">{{ row.attributes.owner || '-' }}</template>
          </el-table-column>
          <el-table-column label="分组" min-width="180">
            <template #default="{ row }">
              <span v-if="!row.group_ids?.length" class="muted">未分组</span>
              <el-tag v-for="groupId in row.group_ids" :key="groupId" size="small" effect="plain">
                {{ groupName(groupId) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="statusTag(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="210" fixed="right">
            <template #default="{ row }">
              <RouterLink :to="`/cmdb/assets/${row.id}`">详情</RouterLink>
              <el-button link type="primary" @click="editAsset(row)">编辑</el-button>
              <el-button link type="danger" @click="removeAsset(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-row">
          <span>共 {{ total }} 条资产</span>
          <el-pagination
            v-model:current-page="filters.page"
            v-model:page-size="filters.page_size"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="sizes, prev, pager, next"
            @change="reloadAssets"
          />
        </div>
      </section>
    </section>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑资产' : '新增资产'" size="560px">
      <el-form label-position="top">
        <el-form-item label="资产模型">
          <el-select v-model="assetForm.model_id" filterable placeholder="选择模型" @change="onModelChange">
            <el-option v-for="model in allModels" :key="model.id" :label="model.display_name" :value="model.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="唯一标识">
          <el-input v-model="assetForm.unique_key" placeholder="如 sn:ABC123 / instance:i-001" />
        </el-form-item>
        <el-form-item label="资产分组">
          <el-select v-model="assetForm.group_ids" multiple filterable placeholder="选择业务线、环境、地域分组">
            <el-option v-for="group in assetGroups" :key="group.id" :label="group.display_name" :value="group.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="生命周期状态">
          <el-select v-model="assetForm.status">
            <el-option v-for="status in statusOptions" :key="status.value" :label="status.label" :value="status.value" />
          </el-select>
        </el-form-item>

        <template v-if="selectedModel">
          <el-divider content-position="left">模型字段</el-divider>
          <el-row :gutter="12">
            <el-col v-for="field in selectedModel.fields" :key="field.name" :span="12">
              <el-form-item :label="field.display_name" :required="field.required">
                <el-select
                  v-if="field.field_type === 'enum'"
                  v-model="assetForm.attributes[field.name]"
                  :placeholder="field.display_name"
                >
                  <el-option v-for="option in field.options" :key="option" :label="option" :value="option" />
                </el-select>
                <el-input-number
                  v-else-if="field.field_type === 'number'"
                  v-model="assetForm.attributes[field.name]"
                  :controls="false"
                  class="full-control"
                />
                <el-date-picker
                  v-else-if="field.field_type === 'date'"
                  v-model="assetForm.attributes[field.name]"
                  value-format="YYYY-MM-DD"
                  class="full-control"
                />
                <el-input v-else v-model="assetForm.attributes[field.name]" :placeholder="field.display_name" />
              </el-form-item>
            </el-col>
          </el-row>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAsset">{{ editingId ? '保存资产' : '创建资产' }}</el-button>
      </template>
    </el-drawer>

    <el-dialog v-model="groupDialogVisible" title="新增资产分组" width="420px">
      <el-form label-position="top">
        <el-form-item label="分组标识">
          <el-input v-model="groupForm.name" placeholder="如 env-prod" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="groupForm.display_name" placeholder="如 生产环境" />
        </el-form-item>
        <el-form-item label="维度">
          <el-select v-model="groupForm.dimension">
            <el-option label="业务线" value="business" />
            <el-option label="环境" value="environment" />
            <el-option label="地域" value="region" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="groupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitGroup">创建分组</el-button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

import {
  createAsset,
  createAssetGroup,
  deleteAsset,
  listAssetGroups,
  listAssets,
  listModelGroups,
  listModels,
  updateAsset,
  type Asset,
  type AssetGroup,
  type AssetListQuery,
  type Model
} from '../../api/cmdb'

const statusOptions = [
  { label: '采购', value: 'purchased' },
  { label: '上架', value: 'racked' },
  { label: '运行', value: 'running' },
  { label: '维修', value: 'maintenance' },
  { label: '报废', value: 'retired' }
]

const assets = ref<Asset[]>([])
const assetGroups = ref<AssetGroup[]>([])
const allModels = ref<Model[]>([])
const total = ref(0)
const editingId = ref('')
const drawerVisible = ref(false)
const groupDialogVisible = ref(false)
const filters = reactive<AssetListQuery>({ page: 1, page_size: 20 })
const assetForm = reactive({
  model_id: '',
  unique_key: '',
  status: 'running',
  group_ids: [] as string[],
  attributes: {} as Record<string, unknown>
})
const groupForm = reactive({ name: '', display_name: '', dimension: 'environment' })

const selectedModel = computed(() => allModels.value.find((model) => model.id === assetForm.model_id))

async function loadPage() {
  const groups = await listModelGroups()
  const modelLists = await Promise.all(groups.map((group) => listModels(group.id)))
  allModels.value = modelLists.flat()
  assetGroups.value = await listAssetGroups()
  await reloadAssets()
}

async function reloadAssets() {
  const result = await listAssets(filters)
  assets.value = result.items
  total.value = result.total
}

function selectGroup(groupId: string) {
  filters.group_id = groupId || undefined
  filters.page = 1
  reloadAssets()
}

function openCreate() {
  resetForm()
  drawerVisible.value = true
}

function onModelChange() {
  assetForm.attributes = {}
}

async function saveAsset() {
  if (!assetForm.model_id || !assetForm.unique_key) {
    ElMessage.warning('请选择模型并填写唯一标识')
    return
  }
  const payload = {
    model_id: assetForm.model_id,
    unique_key: assetForm.unique_key,
    status: assetForm.status,
    group_ids: assetForm.group_ids,
    attributes: assetForm.attributes
  }
  if (editingId.value) {
    await updateAsset(editingId.value, { id: editingId.value, ...payload })
    ElMessage.success('资产已更新')
  } else {
    await createAsset(payload)
    ElMessage.success('资产已创建')
  }
  drawerVisible.value = false
  resetForm()
  await reloadAssets()
}

function editAsset(asset: Asset) {
  editingId.value = asset.id
  assetForm.model_id = asset.model_id
  assetForm.unique_key = asset.unique_key
  assetForm.status = asset.status
  assetForm.group_ids = [...(asset.group_ids ?? [])]
  assetForm.attributes = { ...asset.attributes }
  drawerVisible.value = true
}

async function removeAsset(id: string) {
  await ElMessageBox.confirm('确认删除该资产？删除后仍可通过审计日志追踪操作。', '删除资产')
  await deleteAsset(id)
  await reloadAssets()
  ElMessage.success('资产已删除')
}

async function submitGroup() {
  if (!groupForm.name || !groupForm.display_name || !groupForm.dimension) {
    ElMessage.warning('请填写分组标识、显示名称和维度')
    return
  }
  await createAssetGroup(groupForm)
  Object.assign(groupForm, { name: '', display_name: '', dimension: 'environment' })
  groupDialogVisible.value = false
  assetGroups.value = await listAssetGroups()
  ElMessage.success('资产分组已创建')
}

function resetForm() {
  editingId.value = ''
  Object.assign(assetForm, {
    model_id: '',
    unique_key: '',
    status: 'running',
    group_ids: [],
    attributes: {}
  })
}

function modelName(modelId: string) {
  return allModels.value.find((model) => model.id === modelId)?.display_name ?? modelId
}

function groupName(groupId: string) {
  return assetGroups.value.find((group) => group.id === groupId)?.display_name ?? groupId
}

function dimensionLabel(dimension: string) {
  const labels: Record<string, string> = { business: '业务线', environment: '环境', region: '地域' }
  return labels[dimension] ?? dimension
}

function statusLabel(value: string) {
  return statusOptions.find((status) => status.value === value)?.label ?? value
}

function statusTag(value: string) {
  if (value === 'running') return 'success'
  if (value === 'maintenance') return 'warning'
  if (value === 'retired') return 'info'
  return ''
}

onMounted(loadPage)
</script>

<style scoped>
.asset-page {
  min-height: 100vh;
  padding: 28px;
  background: #f4f6f8;
  color: #172033;
}

.page-header,
.header-actions,
.panel-title,
.filter-bar,
.pagination-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-header {
  justify-content: space-between;
  margin-bottom: 20px;
}

.eyebrow {
  margin: 0 0 4px;
  color: #1f6feb;
  font-size: 12px;
  font-weight: 700;
}

h1 {
  margin: 0;
  font-size: 26px;
  letter-spacing: 0;
}

.subtitle {
  margin: 6px 0 0;
  color: #667085;
}

.workspace {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 16px;
}

.side-panel,
.main-panel {
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #fff;
}

.side-panel {
  padding: 14px;
}

.panel-title {
  justify-content: space-between;
  margin-bottom: 10px;
}

.group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 42px;
  padding: 9px 10px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  color: #344054;
  cursor: pointer;
}

.group-item.active,
.group-item:hover {
  border-color: #9bb8ff;
  background: #eef4ff;
  color: #174ea6;
}

.group-item small,
.muted {
  color: #8a95a6;
}

.main-panel {
  padding: 16px;
}

.filter-bar {
  margin-bottom: 14px;
}

.filter-bar .el-input {
  max-width: 360px;
}

.filter-bar .el-select {
  width: 180px;
}

.asset-table :deep(.el-tag + .el-tag) {
  margin-left: 6px;
}

.pagination-row {
  justify-content: space-between;
  margin-top: 14px;
  color: #667085;
}

.text-link,
.asset-table a {
  color: #1f6feb;
  font-weight: 600;
  text-decoration: none;
}

.full-control {
  width: 100%;
}

@media (max-width: 960px) {
  .workspace {
    grid-template-columns: 1fr;
  }

  .page-header,
  .filter-bar {
    align-items: flex-start;
    flex-direction: column;
  }

  .filter-bar .el-input,
  .filter-bar .el-select {
    width: 100%;
    max-width: none;
  }
}
</style>

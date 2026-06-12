<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>资产导入</h1>
      </div>
      <RouterLink to="/cmdb/assets" class="back-link">返回资产列表</RouterLink>
    </header>

    <section class="panel">
      <div class="toolbar">
        <el-select v-model="modelId" placeholder="选择导入模型" filterable>
          <el-option v-for="model in models" :key="model.id" :label="model.display_name" :value="model.id" />
        </el-select>
        <el-button :disabled="!selectedModel" @click="downloadTemplate">下载导入模板</el-button>
        <input type="file" accept=".csv" @change="onFileChange" />
        <el-button type="primary" :disabled="!modelId || rows.length === 0" @click="importRows">导入预览数据</el-button>
      </div>
      <p class="hint">CSV 必须包含 unique_key 列；模板会按所选模型自动包含动态字段，自定义模型同样可下载。</p>
    </section>

    <section class="panel">
      <el-table :data="rows" border>
        <el-table-column prop="unique_key" label="唯一标识" min-width="180" />
        <el-table-column label="属性" min-width="360">
          <template #default="{ row }">{{ JSON.stringify(row.attributes) }}</template>
        </el-table-column>
      </el-table>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

import { createAsset, listModelGroups, listModels, previewCSVImport, type Model } from '../../api/cmdb'

const models = ref<Model[]>([])
const modelId = ref('')
const rows = ref<Array<{ unique_key: string; attributes: Record<string, unknown> }>>([])
const selectedModel = computed(() => models.value.find((model) => model.id === modelId.value))

async function loadModels() {
  const groups = await listModelGroups()
  const modelLists = await Promise.all(groups.map((group) => listModels(group.id)))
  models.value = modelLists.flat()
}

async function onFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  rows.value = await previewCSVImport(file)
  ElMessage.success(`预览 ${rows.value.length} 行`)
}

function downloadTemplate() {
  if (!selectedModel.value) return
  const columns = ['unique_key', 'status', ...selectedModel.value.fields.map((field) => field.name)]
  const example = columns.map((column) => {
    if (column === 'unique_key') return 'asset-001'
    if (column === 'status') return 'running'
    const field = selectedModel.value?.fields.find((item) => item.name === column)
    if (field?.field_type === 'ip') return '10.0.0.10'
    if (field?.field_type === 'number') return '1'
    if (field?.field_type === 'enum') return field.options?.[0] ?? ''
    if (field?.field_type === 'date') return '2026-01-01'
    return ''
  })
  const csv = `${columns.join(',')}\n${example.map(escapeCSVValue).join(',')}\n`
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${selectedModel.value.name}-import-template.csv`
  link.click()
  URL.revokeObjectURL(url)
}

async function importRows() {
  let success = 0
  for (const row of rows.value) {
    const { status, ...attributes } = row.attributes
    await createAsset({
      model_id: modelId.value,
      unique_key: row.unique_key,
      status: String(status || 'running'),
      attributes
    })
    success += 1
  }
  ElMessage.success(`导入完成：成功 ${success} 行`)
}

function escapeCSVValue(value: unknown) {
  const text = String(value ?? '')
  if (text.includes(',') || text.includes('"') || text.includes('\n')) {
    return `"${text.split('"').join('""')}"`
  }
  return text
}

onMounted(loadModels)
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px;
  background: #f6f8fb;
}

.page-header,
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-header {
  justify-content: space-between;
  margin-bottom: 20px;
}

.panel {
  margin-bottom: 16px;
  padding: 16px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
}

.toolbar {
  flex-wrap: wrap;
}

.toolbar .el-select {
  width: 260px;
}

.hint {
  margin: 10px 0 0;
  color: #64748b;
  font-size: 13px;
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

.back-link {
  color: #1d4ed8;
  font-weight: 600;
  text-decoration: none;
}
</style>

<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>模型组与模型</h1>
      </div>
      <RouterLink to="/dashboard" class="back-link">返回首页</RouterLink>
    </header>

    <section class="toolbar">
      <el-input v-model="groupForm.name" placeholder="模型组标识" />
      <el-input v-model="groupForm.display_name" placeholder="显示名称" />
      <el-input v-model="groupForm.description" placeholder="描述" />
      <el-button type="primary" @click="submitGroup">新增模型组</el-button>
    </section>

    <el-row :gutter="16">
      <el-col :span="8">
        <el-table :data="groups" border highlight-current-row @current-change="selectGroup">
          <el-table-column prop="name" label="模型组标识" min-width="130" />
          <el-table-column prop="display_name" label="显示名称" min-width="130" />
        </el-table>
      </el-col>
      <el-col :span="16">
        <section v-if="selectedGroup" class="panel">
          <div class="panel-title">
            <strong>{{ selectedGroup.display_name }}</strong>
            <span>{{ selectedGroup.description }}</span>
          </div>
          <div class="toolbar compact">
            <el-input v-model="modelForm.name" placeholder="模型标识" />
            <el-input v-model="modelForm.display_name" placeholder="显示名称" />
            <el-button type="primary" @click="submitModel">新增模型</el-button>
          </div>
          <el-table :data="models" border>
            <el-table-column prop="name" label="模型标识" min-width="140" />
            <el-table-column prop="display_name" label="显示名称" min-width="160" />
            <el-table-column label="字段数" width="90">
              <template #default="{ row }">{{ row.fields.length }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160">
              <template #default="{ row }">
                <RouterLink :to="`/cmdb/models/${row.id}`">编辑字段</RouterLink>
              </template>
            </el-table-column>
          </el-table>
        </section>
      </el-col>
    </el-row>
  </main>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

import {
  createModel,
  createModelGroup,
  listModelGroups,
  listModels,
  type Model,
  type ModelGroup
} from '../../api/cmdb'

const groups = ref<ModelGroup[]>([])
const models = ref<Model[]>([])
const selectedGroup = ref<ModelGroup | null>(null)
const groupForm = reactive({ name: '', display_name: '', description: '' })
const modelForm = reactive({ name: '', display_name: '', description: '', fields: [], relations: [] })

async function loadGroups() {
  groups.value = await listModelGroups()
  if (!selectedGroup.value && groups.value.length > 0) {
    await selectGroup(groups.value[0])
  }
}

async function selectGroup(group: ModelGroup | null) {
  selectedGroup.value = group
  models.value = group ? await listModels(group.id) : []
}

async function submitGroup() {
  if (!groupForm.name || !groupForm.display_name) {
    ElMessage.warning('请填写模型组标识和显示名称')
    return
  }
  await createModelGroup(groupForm)
  Object.assign(groupForm, { name: '', display_name: '', description: '' })
  await loadGroups()
  ElMessage.success('模型组已创建')
}

async function submitModel() {
  if (!selectedGroup.value || !modelForm.name || !modelForm.display_name) {
    ElMessage.warning('请选择模型组并填写模型信息')
    return
  }
  await createModel(selectedGroup.value.id, modelForm)
  Object.assign(modelForm, { name: '', display_name: '', description: '', fields: [], relations: [] })
  await selectGroup(selectedGroup.value)
  ElMessage.success('模型已创建')
}

onMounted(loadGroups)
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px;
  background: #f6f8fb;
}

.page-header,
.panel-title,
.toolbar {
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
  color: #2563eb;
  font-size: 13px;
  font-weight: 700;
}

h1 {
  margin: 0;
  color: #111827;
  font-size: 24px;
}

.toolbar {
  margin-bottom: 16px;
}

.toolbar.compact {
  margin: 14px 0;
}

.panel {
  padding: 16px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
}

.panel-title {
  justify-content: space-between;
  color: #475569;
}

.back-link {
  color: #1d4ed8;
  font-weight: 600;
  text-decoration: none;
}
</style>

<template>
  <main class="page-shell" v-if="model">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>模型编辑</h1>
      </div>
      <div class="actions">
        <RouterLink to="/cmdb/model-groups" class="back-link">返回模型组</RouterLink>
        <el-button type="danger" plain @click="removeModel">删除模型</el-button>
        <el-button type="primary" @click="saveModel">保存模型</el-button>
      </div>
    </header>

    <section class="editor-grid">
      <el-form label-position="top" class="panel">
        <el-form-item label="模型标识">
          <el-input v-model="model.name" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="model.display_name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="model.description" type="textarea" />
        </el-form-item>
      </el-form>

      <section class="panel">
        <div class="field-form">
          <el-input v-model="fieldForm.name" placeholder="字段名" />
          <el-input v-model="fieldForm.display_name" placeholder="显示名称" />
          <el-select v-model="fieldForm.field_type" placeholder="类型">
            <el-option label="文本" value="text" />
            <el-option label="数字" value="number" />
            <el-option label="枚举" value="enum" />
            <el-option label="日期" value="date" />
            <el-option label="IP 地址" value="ip" />
            <el-option label="关联引用" value="relation" />
          </el-select>
          <el-checkbox v-model="fieldForm.required">必填</el-checkbox>
          <el-input v-model="optionsText" placeholder="枚举选项，逗号分隔" />
          <el-button type="primary" @click="addField">添加字段</el-button>
        </div>

        <el-table :data="model.fields" border>
          <el-table-column prop="name" label="字段名" min-width="140" />
          <el-table-column prop="display_name" label="显示名称" min-width="140" />
          <el-table-column prop="field_type" label="类型" width="110" />
          <el-table-column label="必填" width="90">
            <template #default="{ row }">{{ row.required ? '是' : '否' }}</template>
          </el-table-column>
          <el-table-column label="枚举选项" min-width="180">
            <template #default="{ row }">{{ row.options?.join(', ') }}</template>
          </el-table-column>
          <el-table-column label="操作" width="90">
            <template #default="{ $index }">
              <el-button link type="danger" @click="model?.fields.splice($index, 1)">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </section>

      <section class="panel relation-panel">
        <h2>模型关系</h2>
        <div class="field-form">
          <el-select v-model="relationForm.target_model_id" filterable placeholder="目标模型">
            <el-option v-for="item in relatedModelOptions" :key="item.id" :label="item.display_name" :value="item.id" />
          </el-select>
          <el-select v-model="relationForm.relation_type" placeholder="关系类型">
            <el-option label="运行于 / 运行" value="runs" />
            <el-option label="依赖" value="depends_on" />
            <el-option label="包含" value="contains" />
            <el-option label="入口" value="exposes_by" />
          </el-select>
          <el-input v-model="relationForm.display_name" placeholder="显示名称" />
          <el-button type="primary" @click="addRelation">添加关系</el-button>
        </div>
        <el-table :data="model.relations" border>
          <el-table-column label="源模型" min-width="140">
            <template #default>{{ model?.display_name }}</template>
          </el-table-column>
          <el-table-column label="关系" min-width="120">
            <template #default="{ row }">{{ row.display_name || row.relation_type }}</template>
          </el-table-column>
          <el-table-column label="目标模型" min-width="160">
            <template #default="{ row }">{{ modelName(row.target_model_id) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="90">
            <template #default="{ $index }">
              <el-button link type="danger" @click="model?.relations.splice($index, 1)">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </section>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import {
  deleteModel,
  getModel,
  listModelGroups,
  listModels,
  updateModel,
  type FieldDefinition,
  type Model,
  type ModelRelation
} from '../../api/cmdb'

const route = useRoute()
const router = useRouter()
const model = ref<Model | null>(null)
const allModels = ref<Model[]>([])
const optionsText = ref('')
const fieldForm = reactive<FieldDefinition>({
  name: '',
  display_name: '',
  field_type: 'text',
  required: false,
  unique_value: false,
  options: [],
  sort_order: 0
})
const relationForm = reactive<ModelRelation>({
  target_model_id: '',
  relation_type: 'depends_on',
  display_name: ''
})

const relatedModelOptions = computed(() => allModels.value.filter((item) => item.id !== model.value?.id))

async function loadModel() {
  const [current, groups] = await Promise.all([getModel(String(route.params.id)), listModelGroups()])
  model.value = { ...current, relations: current.relations ?? [] }
  const modelLists = await Promise.all(groups.map((group) => listModels(group.id)))
  allModels.value = modelLists.flat()
}

function addField() {
  if (!model.value || !fieldForm.name || !fieldForm.display_name) {
    ElMessage.warning('请填写字段名和显示名称')
    return
  }
  model.value.fields.push({
    ...fieldForm,
    options: optionsText.value ? optionsText.value.split(',').map((item) => item.trim()).filter(Boolean) : [],
    sort_order: model.value.fields.length + 1
  })
  Object.assign(fieldForm, {
    name: '',
    display_name: '',
    field_type: 'text',
    required: false,
    unique_value: false,
    options: [],
    sort_order: 0
  })
  optionsText.value = ''
}

function addRelation() {
  if (!model.value || !relationForm.target_model_id || !relationForm.relation_type) {
    ElMessage.warning('请选择目标模型和关系类型')
    return
  }
  model.value.relations.push({
    target_model_id: relationForm.target_model_id,
    relation_type: relationForm.relation_type,
    display_name: relationForm.display_name || relationForm.relation_type
  })
  Object.assign(relationForm, { target_model_id: '', relation_type: 'depends_on', display_name: '' })
}

async function saveModel() {
  if (!model.value) return
  model.value = await updateModel(model.value.id, model.value)
  ElMessage.success('模型已保存')
}

async function removeModel() {
  if (!model.value) return
  await ElMessageBox.confirm('确认删除该模型？有资产引用时后端会拒绝删除。', '删除模型')
  await deleteModel(model.value.id)
  ElMessage.success('模型已删除')
  router.push('/cmdb/model-groups')
}

function modelName(modelId: string) {
  return allModels.value.find((item) => item.id === modelId)?.display_name ?? modelId
}

onMounted(loadModel)
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px;
  background: #f6f8fb;
}

.page-header,
.actions,
.field-form {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-header {
  justify-content: space-between;
  margin-bottom: 20px;
}

.actions {
  flex-wrap: wrap;
}

.editor-grid {
  display: grid;
  grid-template-columns: minmax(260px, 360px) 1fr;
  gap: 20px;
}

.relation-panel {
  grid-column: 1 / -1;
}

.panel {
  padding: 16px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
}

.field-form {
  flex-wrap: wrap;
  margin-bottom: 14px;
}

.field-form .el-input,
.field-form .el-select {
  width: 180px;
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

@media (max-width: 960px) {
  .editor-grid {
    grid-template-columns: 1fr;
  }
}
</style>

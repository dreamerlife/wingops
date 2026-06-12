<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">CMDB</p>
        <h1>API Key</h1>
      </div>
      <RouterLink to="/dashboard" class="back-link">返回首页</RouterLink>
    </header>

    <section class="panel">
      <div class="toolbar">
        <el-input v-model="name" placeholder="Key 名称" />
        <el-button type="primary" @click="addKey">新增 Key</el-button>
      </div>
      <p v-if="createdSecret" class="secret">新 Key Secret：{{ createdSecret }}</p>
    </section>

    <section class="panel">
      <el-table :data="keys" border>
        <el-table-column prop="name" label="名称" min-width="180" />
        <el-table-column prop="key_id" label="Key ID" min-width="220" />
        <el-table-column prop="status" label="状态" width="120" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="danger" :disabled="row.status === 'revoked'" @click="revoke(row.id)">吊销</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'

import { createApiKey, listApiKeys, revokeApiKey, type ApiKey } from '../../api/cmdb'

const keys = ref<ApiKey[]>([])
const name = ref('')
const createdSecret = ref('')

async function loadKeys() {
  keys.value = await listApiKeys()
}

async function addKey() {
  const key = await createApiKey({ name: name.value || '同步 Key' })
  createdSecret.value = key.secret ?? ''
  name.value = ''
  await loadKeys()
  ElMessage.success('API Key 已创建')
}

async function revoke(id: string) {
  await revokeApiKey(id)
  await loadKeys()
  ElMessage.success('API Key 已吊销')
}

onMounted(loadKeys)
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

.toolbar .el-input {
  width: 260px;
}

.secret {
  margin: 12px 0 0;
  color: #0f766e;
  font-weight: 700;
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

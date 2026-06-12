<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Access / 用户治理</p>
        <h1>用户管理</h1>
        <p>维护平台登录账号、状态和角色授权。</p>
      </div>
      <el-button type="primary" @click="openCreate">新增用户</el-button>
    </header>

    <section class="panel">
      <el-table :data="users" border>
        <el-table-column prop="username" label="用户名" min-width="140" />
        <el-table-column prop="display_name" label="显示名称" min-width="160" />
        <el-table-column label="角色" min-width="260">
          <template #default="{ row }">
            <el-tag v-for="role in row.roles" :key="role.name" size="small" effect="plain">
              {{ role.display_name || role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status === 'active' ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="removeUser(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑用户' : '新增用户'" size="480px">
      <el-form label-position="top">
        <el-form-item label="用户名">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item :label="editingId ? '密码（留空不修改）' : '密码'">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status">
            <el-option label="启用" value="active" />
            <el-option label="停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色授权">
          <el-select v-model="form.role_names" multiple filterable>
            <el-option v-for="role in roles" :key="role.name" :label="role.display_name || role.name" :value="role.name" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">保存</el-button>
      </template>
    </el-drawer>
  </main>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

import { createUser, deleteUser, listRoles, listUsers, updateUser, type Role, type User } from '../api/platform'

const users = ref<User[]>([])
const roles = ref<Role[]>([])
const drawerVisible = ref(false)
const editingId = ref('')
const form = reactive({
  username: '',
  password: '',
  display_name: '',
  status: 'active',
  role_names: [] as string[]
})

async function loadPage() {
  const [userRows, roleRows] = await Promise.all([listUsers(), listRoles()])
  users.value = userRows
  roles.value = roleRows
}

function openCreate() {
  editingId.value = ''
  Object.assign(form, { username: '', password: '', display_name: '', status: 'active', role_names: [] })
  drawerVisible.value = true
}

function openEdit(user: User) {
  editingId.value = user.id
  Object.assign(form, {
    username: user.username,
    password: '',
    display_name: user.display_name,
    status: user.status || 'active',
    role_names: user.roles?.map((role) => role.name) ?? []
  })
  drawerVisible.value = true
}

async function saveUser() {
  if (!form.username || !form.display_name || (!editingId.value && !form.password)) {
    ElMessage.warning('请填写用户名、显示名称和密码')
    return
  }
  if (editingId.value) {
    await updateUser(editingId.value, { ...form, password: form.password || undefined })
    ElMessage.success('用户已更新')
  } else {
    await createUser(form)
    ElMessage.success('用户已创建')
  }
  drawerVisible.value = false
  await loadPage()
}

async function removeUser(id: string) {
  await ElMessageBox.confirm('确认删除该用户？相关审计记录会保留。', '删除用户')
  await deleteUser(id)
  await loadPage()
  ElMessage.success('用户已删除')
}

onMounted(loadPage)
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 28px;
  background: #f3f6f8;
  color: #172033;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.page-header p {
  margin: 6px 0 0;
  color: #667085;
}

.eyebrow {
  margin: 0;
  color: #0f7b74;
  font-size: 12px;
  font-weight: 800;
}

h1 {
  margin: 4px 0 0;
  font-size: 26px;
  letter-spacing: 0;
}

.panel {
  padding: 16px;
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #fff;
}

:deep(.el-tag + .el-tag) {
  margin-left: 6px;
}
</style>

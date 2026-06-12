<template>
  <main class="page-shell">
    <header class="page-header">
      <div>
        <p class="eyebrow">Access / RBAC</p>
        <h1>角色权限</h1>
        <p>维护角色和权限授权，权限码由系统模块预置。</p>
      </div>
      <el-button type="primary" @click="openCreate">新增角色</el-button>
    </header>

    <section class="panel">
      <el-table :data="roles" border>
        <el-table-column prop="name" label="角色标识" min-width="160" />
        <el-table-column prop="display_name" label="显示名称" min-width="160" />
        <el-table-column label="权限" min-width="420">
          <template #default="{ row }">
            <el-tag v-for="permission in row.permissions" :key="permission.code" size="small" effect="plain">
              {{ permission.description || permission.code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="removeRole(row.name)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>

    <el-drawer v-model="drawerVisible" :title="editingName ? '编辑角色' : '新增角色'" size="560px">
      <el-form label-position="top">
        <el-form-item label="角色标识">
          <el-input v-model="form.name" :disabled="Boolean(editingName)" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="权限授权">
          <el-checkbox-group v-model="form.permission_codes" class="permission-grid">
            <el-checkbox v-for="permission in permissions" :key="permission.code" :label="permission.code">
              <strong>{{ permission.description }}</strong>
              <span>{{ permission.code }}</span>
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRole">保存</el-button>
      </template>
    </el-drawer>
  </main>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

import {
  createRole,
  deleteRole,
  listPermissions,
  listRoles,
  updateRole,
  type Permission,
  type Role
} from '../api/platform'

const roles = ref<Role[]>([])
const permissions = ref<Permission[]>([])
const drawerVisible = ref(false)
const editingName = ref('')
const form = reactive({
  name: '',
  display_name: '',
  permission_codes: [] as string[]
})

async function loadPage() {
  const [roleRows, permissionRows] = await Promise.all([listRoles(), listPermissions()])
  roles.value = roleRows
  permissions.value = permissionRows
}

function openCreate() {
  editingName.value = ''
  Object.assign(form, { name: '', display_name: '', permission_codes: [] })
  drawerVisible.value = true
}

function openEdit(role: Role) {
  editingName.value = role.name
  Object.assign(form, {
    name: role.name,
    display_name: role.display_name,
    permission_codes: role.permissions.map((permission) => permission.code)
  })
  drawerVisible.value = true
}

async function saveRole() {
  if (!form.name || !form.display_name) {
    ElMessage.warning('请填写角色标识和显示名称')
    return
  }
  if (editingName.value) {
    await updateRole(editingName.value, form)
    ElMessage.success('角色已更新')
  } else {
    await createRole(form)
    ElMessage.success('角色已创建')
  }
  drawerVisible.value = false
  await loadPage()
}

async function removeRole(name: string) {
  await ElMessageBox.confirm('确认删除该角色？已有用户的该角色授权会被移除。', '删除角色')
  await deleteRole(name)
  await loadPage()
  ElMessage.success('角色已删除')
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

.permission-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.permission-grid :deep(.el-checkbox) {
  align-items: flex-start;
  height: auto;
  min-height: 54px;
  padding: 10px;
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #fbfcfe;
  white-space: normal;
}

.permission-grid strong,
.permission-grid span {
  display: block;
  line-height: 1.35;
}

.permission-grid span {
  color: #667085;
  font-size: 12px;
}

:deep(.el-tag + .el-tag) {
  margin-left: 6px;
  margin-top: 4px;
}

@media (max-width: 760px) {
  .permission-grid {
    grid-template-columns: 1fr;
  }
}
</style>

import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/users',
      component: () => import('../views/UserListView.vue')
    },
    {
      path: '/roles',
      component: () => import('../views/RoleListView.vue')
    },
    {
      path: '/audit/logs',
      component: () => import('../views/AuditLogView.vue')
    },
    {
      path: '/system/configs',
      component: () => import('../views/SystemConfigView.vue')
    },
    {
      path: '/cmdb/model-groups',
      component: () => import('../views/cmdb/ModelGroupListView.vue')
    },
    {
      path: '/cmdb/models/:id',
      component: () => import('../views/cmdb/ModelEditorView.vue')
    },
    {
      path: '/cmdb/assets',
      component: () => import('../views/cmdb/AssetListView.vue')
    },
    {
      path: '/cmdb/assets/import',
      component: () => import('../views/cmdb/AssetImportView.vue')
    },
    {
      path: '/cmdb/assets/:id',
      component: () => import('../views/cmdb/AssetDetailView.vue')
    },
    {
      path: '/cmdb/api-keys',
      component: () => import('../views/cmdb/ApiKeyListView.vue')
    }
  ]
})

export default router

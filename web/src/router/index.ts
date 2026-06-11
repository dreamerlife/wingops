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
    }
  ]
})

export default router

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
    }
  ]
})

export default router

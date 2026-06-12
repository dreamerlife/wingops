<template>
  <div class="app-shell" :class="{ collapsed, 'menu-open': mobileOpen }">
    <button class="mobile-toggle" type="button" @click="mobileOpen = true">菜单</button>
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">WO</div>
        <div class="brand-copy">
          <strong>WingOps</strong>
          <span>统一运维管理平台</span>
        </div>
      </div>

      <nav class="nav-groups">
        <section v-for="group in navGroups" :key="group.title" class="nav-group">
          <h2>{{ group.title }}</h2>
          <RouterLink
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            class="nav-item"
            :class="{ active: isActive(item.to) }"
            :title="item.label"
            @click="mobileOpen = false"
          >
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-label">{{ item.label }}</span>
          </RouterLink>
        </section>
      </nav>

      <div class="sidebar-actions">
        <button type="button" class="collapse-button" @click="collapsed = !collapsed">
          {{ collapsed ? '展开' : '折叠' }}
        </button>
        <button type="button" class="logout-button" @click="logout">退出登录</button>
      </div>
    </aside>
    <button class="scrim" type="button" aria-label="关闭菜单" @click="mobileOpen = false" />
    <main class="shell-content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const collapsed = ref(false)
const mobileOpen = ref(false)

const navGroups = [
  {
    title: '总览',
    items: [{ label: '平台首页', to: '/dashboard', icon: 'OV' }]
  },
  {
    title: 'CMDB',
    items: [
      { label: '模型管理', to: '/cmdb/model-groups', icon: 'MD' },
      { label: '资产管理', to: '/cmdb/assets', icon: 'AS' },
      { label: '资产导入', to: '/cmdb/assets/import', icon: 'IM' },
      { label: 'API Key', to: '/cmdb/api-keys', icon: 'AK' }
    ]
  },
  {
    title: '平台治理',
    items: [
      { label: '用户管理', to: '/users', icon: 'US' },
      { label: '角色权限', to: '/roles', icon: 'RB' },
      { label: '审计日志', to: '/audit/logs', icon: 'AU' }
    ]
  },
  {
    title: '系统设置',
    items: [{ label: '系统配置', to: '/system/configs', icon: 'CF' }]
  }
]

function isActive(path: string) {
  if (path === '/dashboard') return route.path === path
  return route.path === path || route.path.startsWith(`${path}/`)
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-shell {
  --sidebar-width: 264px;
  min-height: 100vh;
  background: #f3f6f8;
}

.sidebar {
  position: fixed;
  inset: 0 auto 0 0;
  z-index: 20;
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  padding: 18px 14px;
  border-right: 1px solid #d9e0ea;
  background: #ffffff;
  color: #172033;
  box-shadow: 8px 0 28px rgba(55, 84, 115, 0.06);
  transition: width 160ms ease, transform 180ms ease;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 8px 18px;
  border-bottom: 1px solid #edf1f6;
}

.brand-mark,
.nav-icon {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border: 1px solid #cfe4ff;
  background: #eef7ff;
  color: #1f6feb;
  font-weight: 800;
}

.brand-mark {
  width: 40px;
  height: 40px;
  border-radius: 8px;
}

.brand-copy {
  display: grid;
  gap: 2px;
  min-width: 0;
}

.brand-copy strong {
  color: #172033;
  font-size: 18px;
}

.brand-copy span {
  color: #667085;
  font-size: 12px;
}

.nav-groups {
  display: grid;
  gap: 16px;
  margin-top: 18px;
  overflow-y: auto;
}

.nav-group h2 {
  margin: 0 8px 7px;
  color: #7b8794;
  font-size: 12px;
  font-weight: 800;
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 40px;
  padding: 8px 10px;
  border-radius: 8px;
  color: #344054;
  text-decoration: none;
  transition: background 140ms ease, color 140ms ease, transform 140ms ease;
}

.nav-item:hover {
  background: #f1f7ff;
  color: #174ea6;
  transform: translateX(2px);
}

.nav-item.active {
  background: #e9f7ff;
  color: #0d3b66;
  box-shadow: inset 3px 0 0 #23c6b8;
}

.nav-icon {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  font-size: 11px;
  letter-spacing: 0;
}

.nav-item.active .nav-icon {
  border-color: #1f6feb;
  background: #dceeff;
  color: #174ea6;
}

.nav-label {
  white-space: nowrap;
}

.sidebar-actions {
  display: grid;
  gap: 8px;
  margin-top: auto;
  padding-top: 14px;
}

.collapse-button,
.logout-button,
.mobile-toggle {
  min-height: 38px;
  border: 1px solid #d9e0ea;
  border-radius: 8px;
  background: #f7fafc;
  color: #344054;
  cursor: pointer;
}

.logout-button {
  background: #ffffff;
}

.shell-content {
  min-height: 100vh;
  margin-left: var(--sidebar-width);
  transition: margin-left 160ms ease;
}

.mobile-toggle,
.scrim {
  display: none;
}

.app-shell.collapsed {
  --sidebar-width: 76px;
}

.app-shell.collapsed .brand-copy,
.app-shell.collapsed .nav-group h2,
.app-shell.collapsed .nav-label,
.app-shell.collapsed .logout-button {
  display: none;
}

.app-shell.collapsed .brand {
  justify-content: center;
  padding-inline: 0;
}

.app-shell.collapsed .nav-item {
  justify-content: center;
}

@media (max-width: 860px) {
  .app-shell {
    --sidebar-width: 264px;
  }

  .mobile-toggle {
    position: fixed;
    top: 14px;
    left: 14px;
    z-index: 15;
    display: block;
    padding: 0 14px;
    background: #ffffff;
    color: #174ea6;
    box-shadow: 0 8px 22px rgba(55, 84, 115, 0.14);
  }

  .sidebar {
    transform: translateX(-100%);
  }

  .app-shell.menu-open .sidebar {
    transform: translateX(0);
  }

  .app-shell.menu-open .scrim {
    position: fixed;
    inset: 0;
    z-index: 10;
    display: block;
    border: 0;
    background: rgba(23, 32, 51, 0.28);
  }

  .shell-content {
    margin-left: 0;
    padding-top: 52px;
  }

  .app-shell.collapsed .brand-copy,
  .app-shell.collapsed .nav-group h2,
  .app-shell.collapsed .nav-label,
  .app-shell.collapsed .logout-button {
    display: initial;
  }

  .app-shell.collapsed .nav-item {
    justify-content: flex-start;
  }
}
</style>

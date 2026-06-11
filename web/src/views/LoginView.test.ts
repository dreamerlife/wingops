import { flushPromises, mount } from '@vue/test-utils'
import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import LoginView from './LoginView.vue'

vi.mock('../api/auth', () => ({
  login: vi.fn().mockResolvedValue({
    access_token: 'test-token',
    token_type: 'Bearer'
  })
}))

function createTestRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', redirect: '/login' },
      { path: '/login', component: LoginView },
      { path: '/dashboard', component: { template: '<main>平台首页</main>' } }
    ]
  })
}

describe('LoginView', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders login button', () => {
    const router = createTestRouter()
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createPinia(), router, ElementPlus]
      }
    })

    expect(wrapper.text()).toContain('登录')
  })

  it('navigates to dashboard after login succeeds', async () => {
    const router = createTestRouter()
    router.push('/login')
    await router.isReady()

    const wrapper = mount(LoginView, {
      global: {
        plugins: [createPinia(), router, ElementPlus]
      }
    })

    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })
})

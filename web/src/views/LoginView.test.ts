import { mount } from '@vue/test-utils'
import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import { describe, expect, it } from 'vitest'

import LoginView from './LoginView.vue'

describe('LoginView', () => {
  it('renders login button', () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createPinia(), ElementPlus]
      }
    })

    expect(wrapper.text()).toContain('登录')
  })
})

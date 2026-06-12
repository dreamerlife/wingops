import { defineStore } from 'pinia'
import { ref } from 'vue'

import { login as loginApi } from '../api/auth'

const tokenStorageKey = 'wingops.access_token'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(localStorage.getItem(tokenStorageKey) ?? '')

  async function login(username: string, password: string) {
    const result = await loginApi({ username, password })
    accessToken.value = result.access_token
    localStorage.setItem(tokenStorageKey, result.access_token)
  }

  function logout() {
    accessToken.value = ''
    localStorage.removeItem(tokenStorageKey)
  }

  return {
    accessToken,
    login,
    logout
  }
})

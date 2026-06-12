<template>
  <main class="login-page">
    <section class="ops-visual">
      <p class="eyebrow">WingOps Playground</p>
      <h1>统一运维管理平台</h1>
      <p>资产、权限、审计和后续监控告警链路统一接入，形成可信运维数据底座。</p>
      <div class="mascot-board" aria-hidden="true">
        <div class="board-card server-card">
          <span class="dot green"></span>
          <strong>资产在线</strong>
          <small>CMDB</small>
        </div>
        <div class="board-card pulse-card">
          <span class="spark"></span>
          <strong>同步队列</strong>
          <small>API</small>
        </div>
        <div class="board-card shield-card">
          <span class="dot blue"></span>
          <strong>权限守护</strong>
          <small>RBAC</small>
        </div>
      </div>
      <div class="signal-grid">
        <span>CMDB</span>
        <span>RBAC</span>
        <span>Audit</span>
        <span>Sync API</span>
      </div>
    </section>

    <section class="login-panel" aria-label="登录">
      <div class="brand">
        <p class="eyebrow">Secure Access</p>
        <h2>登录控制台</h2>
      </div>

      <el-form class="login-form" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="用户名">
          <el-input v-model="username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
        <el-button class="login-button" type="primary" native-type="submit" :loading="submitting">登录</el-button>
      </el-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const username = ref('')
const password = ref('')
const submitting = ref(false)
const errorMessage = ref('')

async function handleLogin() {
  errorMessage.value = ''
  submitting.value = true
  try {
    await authStore.login(username.value, password.value)
    await router.push('/dashboard')
  } catch {
    errorMessage.value = '用户名或密码错误'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: relative;
  display: grid;
  grid-template-columns: minmax(360px, 1fr) 420px;
  align-items: center;
  min-height: 100vh;
  padding: 56px;
  overflow: hidden;
  background:
    linear-gradient(90deg, rgba(31, 111, 235, 0.07) 1px, transparent 1px),
    linear-gradient(180deg, rgba(31, 111, 235, 0.06) 1px, transparent 1px),
    radial-gradient(circle at 16% 20%, rgba(255, 217, 102, 0.38), transparent 24%),
    radial-gradient(circle at 72% 18%, rgba(72, 214, 191, 0.28), transparent 28%),
    linear-gradient(135deg, #f9fcff 0%, #eff7ff 48%, #f7fbef 100%);
  background-size: 32px 32px, 32px 32px, auto, auto, auto;
  color: #172033;
}

.login-page::after {
  position: absolute;
  right: -120px;
  bottom: -130px;
  width: 420px;
  height: 420px;
  border-radius: 50%;
  pointer-events: none;
  content: "";
  background: repeating-linear-gradient(45deg, rgba(31, 111, 235, 0.1) 0 12px, transparent 12px 24px);
  opacity: 0.72;
}

.ops-visual {
  position: relative;
  z-index: 1;
  max-width: 680px;
}

.eyebrow {
  margin: 0 0 10px;
  color: #0f7b74;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

h1,
h2 {
  margin: 0;
  letter-spacing: 0;
}

h1 {
  max-width: 560px;
  font-size: 46px;
  line-height: 1.12;
}

.ops-visual p {
  max-width: 620px;
  margin: 18px 0 0;
  color: #5d6b7a;
  font-size: 16px;
  line-height: 1.7;
}

.mascot-board {
  position: relative;
  width: min(100%, 560px);
  height: 190px;
  margin-top: 28px;
}

.board-card {
  position: absolute;
  display: grid;
  gap: 5px;
  min-width: 150px;
  padding: 16px;
  border: 1px solid #d9e4ef;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 16px 38px rgba(55, 84, 115, 0.12);
}

.server-card {
  left: 0;
  top: 14px;
}

.pulse-card {
  left: 190px;
  top: 58px;
}

.shield-card {
  left: 380px;
  top: 10px;
}

.board-card strong {
  color: #172033;
  font-size: 17px;
}

.board-card small {
  color: #667085;
  font-weight: 800;
}

.dot,
.spark {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.dot.green {
  background: #21b66f;
  box-shadow: 0 0 0 6px rgba(33, 182, 111, 0.12);
}

.dot.blue {
  background: #1f6feb;
  box-shadow: 0 0 0 6px rgba(31, 111, 235, 0.12);
}

.spark {
  background: #ffb020;
  box-shadow: 0 0 0 6px rgba(255, 176, 32, 0.16);
}

.signal-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(90px, 1fr));
  gap: 10px;
  max-width: 560px;
  margin-top: 34px;
}

.signal-grid span {
  padding: 14px 12px;
  border: 1px solid #d9e4ef;
  border-radius: 8px;
  background: #ffffff;
  color: #24506d;
  font-weight: 800;
  text-align: center;
  box-shadow: 0 8px 20px rgba(55, 84, 115, 0.08);
}

.login-panel {
  position: relative;
  z-index: 1;
  width: min(100%, 420px);
  padding: 30px;
  border: 1px solid #d8e4ef;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 24px 70px rgba(55, 84, 115, 0.18);
  backdrop-filter: blur(10px);
}

.brand {
  margin-bottom: 24px;
}

h2 {
  font-size: 24px;
}

.login-form {
  display: grid;
  gap: 4px;
}

.login-form :deep(.el-form-item__label) {
  color: #344054;
}

.login-button {
  width: 100%;
  margin-top: 8px;
}

.error-message {
  margin: 0;
  color: #b42318;
  font-size: 13px;
}

@media (max-width: 860px) {
  .login-page {
    grid-template-columns: 1fr;
    gap: 28px;
    padding: 28px;
  }

  h1 {
    font-size: 34px;
  }

  .signal-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .mascot-board {
    height: auto;
  }

  .board-card {
    position: static;
    margin-bottom: 10px;
  }
}
</style>

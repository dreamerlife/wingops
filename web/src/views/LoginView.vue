<template>
  <main class="login-page">
    <section class="login-panel" aria-label="登录">
      <div class="brand">
        <p class="eyebrow">WingOps</p>
        <h1>统一运维管理平台</h1>
      </div>

      <el-form class="login-form" label-position="top">
        <el-form-item label="用户名">
          <el-input v-model="username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <el-button class="login-button" type="primary">登录</el-button>
      </el-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const username = ref('')
const password = ref('')
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 32px;
  background:
    linear-gradient(90deg, rgba(31, 41, 55, 0.06) 1px, transparent 1px),
    linear-gradient(180deg, rgba(31, 41, 55, 0.06) 1px, transparent 1px),
    #f6f8fb;
  background-size: 28px 28px;
}

.login-panel {
  width: min(100%, 360px);
  padding: 28px;
  border: 1px solid #d7dee8;
  border-radius: 8px;
  background: #ffffff;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
}

.brand {
  margin-bottom: 24px;
}

.eyebrow {
  margin: 0 0 6px;
  color: #2563eb;
  font-size: 13px;
  font-weight: 700;
}

h1 {
  margin: 0;
  color: #111827;
  font-size: 22px;
  font-weight: 700;
}

.login-form {
  display: grid;
  gap: 4px;
}

.login-button {
  width: 100%;
  margin-top: 8px;
}
</style>

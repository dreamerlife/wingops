<template>
  <RouterView v-if="$route.path === '/login'" />
  <AppShell v-else />
</template>

<script setup lang="ts">
import AppShell from './components/AppShell.vue'
</script>

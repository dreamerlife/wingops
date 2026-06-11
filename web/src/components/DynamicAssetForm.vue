<template>
  <el-form label-position="top">
    <el-form-item v-for="field in fields" :key="field.name" :label="field.displayName">
      <el-input :model-value="values[field.name]" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
interface Field {
  name: string
  displayName: string
}

defineProps<{
  fields: Field[]
  values: Record<string, string>
}>()
</script>

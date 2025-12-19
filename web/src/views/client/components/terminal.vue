<script setup lang="ts">
import Terminal from "@/components/terminal/index.vue"
const props = defineProps<{
  id: string
}>()
</script>

<template>
    <Terminal :id="id"/>
</template>
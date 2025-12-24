<script setup lang="ts">
import Terminal from "@/components/terminal/index.vue"
const props = defineProps<{
  id: string
}>()

const terminalRef = ref(null)
function closeTerminal() {
  terminalRef.value?.close();
}

onUnmounted(() => {
  closeTerminal();
})
</script>

<template>
    <Terminal :id="id" ref="terminalRef" />
</template>
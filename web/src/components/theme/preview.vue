<template>
  <div class="btn group btn-sm btn-ghost px-1" :data-theme="theme">
    <div class="bg-base-100 grid grid-cols-2 gap-0.5 rounded-md p-1">
      <div class="bg-base-content size-2 rounded-full"></div>
      <div class="bg-primary size-2 rounded-full"></div>
      <div class="bg-secondary size-2 rounded-full"></div>
      <div class="bg-accent size-2 rounded-full"></div>
    </div>
  </div>
</template>

<script setup lang="ts" name="ThemePreview">

defineProps<{
  theme: string
}>()

</script>

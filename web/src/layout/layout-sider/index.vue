<template>
  <div class="menu p-4 w-54 h-full bg-base-100 flex flex-col border-r border-base-200">
    <!-- Logo -->
    <Logo class="mb-8 px-4 lg:hidden" />

    <!-- 菜单项 -->
    <SiderMenu class="flex-1" />

    <!-- 底部固定区域 -->
    <WorkSpace />
  </div>
</template>

<script setup lang="ts">
import SiderMenu from './menu.vue';
import WorkSpace from './work-space.vue';
</script>

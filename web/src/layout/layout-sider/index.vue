<template>
  <div class="menu p-4 w-56 h-full bg-base-100 flex flex-col border-r border-t border-base-300">
    <!-- Logo -->
    <Logo class="mb-8 px-4 lg:hidden" />

    <!-- 菜单项 -->
    <SiderMenu />

    <!-- 底部固定区域 -->
    <WorkSpace />
  </div>
</template>

<script setup lang="ts">
import SiderMenu from './menu.vue';
import WorkSpace from './work-space.vue';
</script>

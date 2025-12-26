<template>
  <div class="navbar bg-base-100 border-b border-base-200 h-16">
    <div class="flex-1 flex items-center">
      <!-- Drawer toggle button for mobile -->
      <label v-if="!route.meta.hideSidebar" for="main-drawer" class="btn btn-ghost lg:hidden">
        <Icon icon="mdi:menu" class="text-xl" />
      </label>
      <Logo class="hidden lg:block ml-2" />
    </div>
    <div class="flex-none flex items-center gap-2 pr-4">
      <Theme />
      <Locales />
      <div class="divider divider-horizontal mx-0"></div>
      <UserInfo />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();
</script>

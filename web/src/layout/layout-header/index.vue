<template>
  <div class="navbar bg-base-100 shadow-sm">
    <div class="flex-1 flex items-center">
      <!-- Drawer toggle button for mobile -->
      <label v-if="!route.meta.hideSidebar" for="main-drawer" class="btn btn-ghost lg:hidden">
        <Icon icon="mdi:menu" />
      </label>
      <Logo class="hidden lg:block" />
    </div>
    <div class="flex-none flex items-center gap-1">
      <Theme />
      <Locales />
      <UserInfo />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();
</script>

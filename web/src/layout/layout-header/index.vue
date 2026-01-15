<template>
  <div class="navbar bg-base-100 border-b border-base-200 h-16">
    <div class="flex-1 flex items-center">
      <!-- Drawer toggle button for mobile -->
      <label v-if="!route.meta.hideSidebar" for="main-drawer" class="btn btn-ghost lg:hidden">
        <Icon icon="mdi:menu" class="text-xl" />
      </label>
      <Logo class="hidden lg:block ml-2" />
    </div>
    <div class="flex-none flex items-center gap-2">
      <Theme />
      <Locales />
      <div class="divider divider-horizontal mx-0"></div>
      <button
        v-if="userStore.isLogined"
        class="btn btn-ghost btn-sm gap-2 normal-case"
        @click="router.push({name: 'Dashboard'})"
      >
        <Icon icon="mdi:view-dashboard" class="text-lg" />
        <span>工作台</span>
      </button>
      <UserInfo />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from "@/stores/auth"

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
</script>

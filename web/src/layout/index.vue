<template>
  <div class="h-screen flex flex-col">
    <Header class="sticky top-0 z-10" />
    <div class="flex flex-1 overflow-hidden">
      <!-- Drawer for mobile and fixed sidebar for desktop -->
      <div class="drawer lg:drawer-open flex-1">
        <input id="main-drawer" type="checkbox" class="drawer-toggle" />
        <div class="drawer-content flex flex-col flex-1 overflow-y-auto bg-base-200/50">
          <!-- Main content area -->
          <main class="p-4 sm:p-6 flex-1 w-full max-w-[1920px] mx-auto">
            <RouterView />
          </main>
        </div>
        <!-- Sidebar -->
        <div v-if="!route.meta.hideSidebar" class="drawer-side h-full">
          <label for="main-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
          <Sidebar />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import Header from './layout-header/index.vue';
import Sidebar from './layout-sider/index.vue'; // 引入侧边栏组件

const route = useRoute();
</script>

<template>
  <div class="mt-auto pt-4 border-t border-base-300">
    <div class="dropdown dropdown-top w-full">
      <div tabindex="0" role="button" class="btn btn-ghost w-full justify-between hover:bg-base-200">
        <div class="flex items-center gap-3">
          <div class="avatar placeholder">
            <Icon icon="material-symbols:group-outline" />
          </div>
          <div>
            <span class="text-gray-500">{{ currentWorkSpaceName }}</span>
          </div>
        </div>
      </div>
      <ul tabindex="0"
        class="dropdown-content z-1 menu p-2 shadow-lg bg-base-100 rounded-box w-50 mb-2 border border-base-200">
        <li class="menu-title">切换工作空间</li>
        <li v-for="space in workSpaceList">
          <a>{{ space.name }}</a>
        </li>
        <li class="border-t mt-2 pt-2"><a class="text-primary">创建空间</a></li>
      </ul>
    </div>

    <div class="dropdown dropdown-top w-full">
      <div tabindex="0" role="button" class="btn btn-ghost w-full justify-between hover:bg-base-300">
        <div class="flex items-center gap-3">
          <div class="avatar placeholder">
            <Icon icon="majesticons:applications" />
          </div>
          <div>
            <span class="text-gray-500">{{ currentWorkAppName }}</span>
          </div>
        </div>
      </div>
      <ul tabindex="0"
        class="dropdown-content z-1 menu p-2 shadow-lg bg-base-100 rounded-box w-64 mb-2 border border-base-200">
        <li class="menu-title">切换应用</li>
        <li v-for="app in workAppList">
          <a>{{ app.name }}</a>
        </li>
        <li class="border-t mt-2 pt-2"><a class="text-primary">创建应用 +</a></li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/auth";
import { storeToRefs } from 'pinia';

const userStore = useUserStore();
const { currentWorkSpace, currentWorkApp, workSpaceList, workAppList } = storeToRefs(userStore);
const currentWorkSpaceName = computed(() => workSpaceList.value.find(space => space.id === currentWorkSpace.value)?.name)
const currentWorkAppName = computed(() => workAppList.value?.find(app => app.id === currentWorkApp.value)?.name)
</script>

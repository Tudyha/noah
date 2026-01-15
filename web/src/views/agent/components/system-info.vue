<script setup lang="ts">
import { useRequest } from 'vue-hooks-plus';
import { getAgentDetail } from '@/api/agent';
import { clientStatusMap, clientOsTypeIconMap } from '@/map';
const router = useRouter();

const props = defineProps<{
  id: string
}>()

const { data: client } = useRequest(() => getAgentDetail(props.id));
</script>

<template>
  <div class="bg-base-100 border-b border-base-200 px-4 py-2 flex items-center justify-between">
    <div class="flex items-center gap-4">
      <button class="btn btn-ghost btn-sm btn-square" @click="router.back()" title="返回列表">
        <Icon icon="mdi:arrow-left" class="w-5 h-5" />
      </button>

      <div class="flex items-center gap-3">
        <div class="flex flex-col">
          <div class="flex items-center gap-2">
            <h1 class="text-sm font-bold truncate max-w-[120px] sm:max-w-xs">
              {{ client?.hostname || '加载中...' }}
            </h1>
            <div v-if="client" class="badge badge-xs font-normal"
              :class="client.status === 1 ? 'badge-success text-white' : 'badge-error text-white'">
              {{ clientStatusMap[client.status] }}
            </div>
          </div>
          <div class="flex items-center gap-2 text-[10px] text-base-content/50">
            <span class="flex items-center gap-0.5">
              <Icon icon="mdi:ip-outline" class="w-2.5 h-2.5" />
              {{ client?.remote_ip }}
            </span>
            <span class="hidden sm:inline opacity-50">•</span>
            <span class="hidden sm:flex items-center gap-0.5">
              <Icon :icon="client ? clientOsTypeIconMap[client.os_type] : ''" class="w-2.5 h-2.5" />
              {{ client?.os_name }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 详细信息气泡 -->
    <div class="dropdown dropdown-end">
      <label tabindex="0" class="btn btn-ghost btn-xs gap-1 opacity-70 hover:opacity-100">
        <Icon icon="mdi:information-outline" class="w-4 h-4" />
        <span class="hidden sm:inline">详情</span>
      </label>
      <div tabindex="0"
        class="dropdown-content z-[30] card card-compact w-64 p-2 shadow-xl bg-base-100 border border-base-200 mt-2">
        <div class="card-body">
          <h3 class="font-bold text-xs mb-2 border-b border-base-200 pb-1">主机详情</h3>
          <div class="space-y-2">
            <div class="flex justify-between items-center text-[11px]">
              <span class="opacity-50">用户名</span>
              <span class="font-medium">{{ client?.username }}</span>
            </div>
            <div class="flex justify-between items-center text-[11px]">
              <span class="opacity-50">内核版本</span>
              <span class="font-medium">{{ client?.kernel_version }}</span>
            </div>
            <div class="flex justify-between items-center text-[11px]">
              <span class="opacity-50">主机程序版本</span>
              <span class="font-medium">{{ client?.version_name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { clientStatusMap, clientOsTypeIconMap, clientOsTypeColorMap } from '@/map'
import type { ClientResponse } from '@/types'
import { formatBytes, formatDateTime, formatUptime } from '@/utils'
import { deleteClient as deleteClientApi} from '@/api/client'

const props = defineProps<{
  item: ClientResponse;
}>()

const emit = defineEmits(['refresh']);

const rowInfo = computed(() => [
  {
    Icon: "mdi:account",
    label: "用户",
    value: `${props.item.username} (UID: ${props.item.uid})`
  },
  {
    Icon: "mdi:ip-outline",
    label: "IP地址",
    value: `${props.item.remote_ip}:${props.item.port}`
  },
  {
    Icon: "mdi:clock-outline",
    label: "运行时间",
    value: formatUptime(props.item.uptime)
  },
  {
    Icon: "material-symbols:online-prediction",
    label: "上次在线",
    value: formatDateTime(props.item.last_online_time)
  },
]);

const statInfo = computed(() => [
  {
    Icon: "mdi:cpu-64-bit",
    label: "CPU",
    value: `${props.item.cpu_num} 核`,
    color: "text-primary"
  },
  {
    Icon: "mdi:memory",
    label: "内存",
    value: formatBytes(props.item.mem_total),
    color: "text-secondary"
  },
  {
    Icon: "mdi:harddisk",
    label: "磁盘",
    value: formatBytes(props.item.disk_total),
    color: "text-accent"
  },
]);

const deleteClient = () => {
  deleteClientApi(props.item.id)
  emit('refresh')
}
</script>

<template>
  <div class="card relative overflow-hidden bg-gradient-to-br from-base-100 via-base-100 to-base-200/50 shadow-sm hover:shadow-xl transition-all duration-300 border border-base-200 group hover:-translate-y-1">
    <!-- 装饰性背景光斑 -->
    <div class="absolute -top-10 -right-10 w-32 h-32 bg-primary/10 rounded-full blur-3xl group-hover:bg-primary/20 transition-colors duration-500 pointer-events-none"></div>
    <div class="absolute -bottom-10 -left-10 w-32 h-32 bg-secondary/10 rounded-full blur-3xl group-hover:bg-secondary/20 transition-colors duration-500 pointer-events-none"></div>
    
    <!-- 装饰性顶部光条 -->
    <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-primary/80 via-secondary/80 to-accent/80 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>

    <div class="card-body p-3 sm:p-4 relative z-10">
      <!-- 头部信息 -->
      <div class="flex justify-between items-start mb-2">
        <div class="flex items-center gap-2">
          <!-- OS 图标容器 -->
          <div class="w-10 h-10 rounded-lg bg-primary/5 group-hover:bg-primary/10 flex items-center justify-center text-primary transition-all duration-300 group-hover:scale-110">
            <Icon :icon="clientOsTypeIconMap[item.os_type]" class="w-6 h-6" />
          </div>
          
          <div class="flex flex-col">
            <h2 class="font-bold text-base leading-tight truncate max-w-[140px]" :title="item.hostname">
              {{ item.hostname }}
            </h2>
            <div class="flex items-center gap-1 mt-0.5">
              <span class="text-[10px] text-base-content/60 font-medium px-1 py-0.5 rounded bg-base-200/50">
                {{ item.platform_family }}
              </span>
              <span class="text-[10px] text-base-content/40">
                {{ item.platform_version }}
              </span>
            </div>
          </div>
        </div>

        <!-- 状态标签 -->
        <div class="badge border-0 gap-1 py-2 px-2 shadow-sm min-h-0 h-auto" 
             :class="item.status === 1 ? 'badge-success/10 text-success' : 'badge-error/10 text-error'">
          <span class="relative flex h-1.5 w-1.5">
            <span v-if="item.status === 1" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
            <span class="relative inline-flex rounded-full h-1.5 w-1.5" :class="item.status === 1 ? 'bg-success' : 'bg-error'"></span>
          </span>
          <span class="font-semibold text-[10px]">{{ clientStatusMap[item.status] }}</span>
        </div>
      </div>

      <!-- 核心指标 Grid -->
      <div class="grid grid-cols-3 gap-2 my-1">
        <div v-for="(stat, index) in statInfo" :key="index" 
             class="flex flex-col items-center justify-center p-2 rounded-lg bg-base-200/30 hover:bg-base-200/60 transition-colors duration-200">
          <Icon :icon="stat.Icon" :class="['w-4 h-4 mb-0.5', stat.color]" />
          <span class="text-xs font-bold text-base-content">{{ stat.value }}</span>
          <span class="text-[10px] text-base-content/50 uppercase tracking-wide scale-90 origin-center">{{ stat.label }}</span>
        </div>
      </div>

      <!-- 详细信息列表 -->
      <div class="mt-2 space-y-1">
        <div v-for="(info, index) in rowInfo" :key="index" class="flex items-center justify-between text-xs group/row hover:bg-base-200/20 p-1 rounded -mx-1 transition-colors">
          <div class="flex items-center gap-1.5 text-base-content/60">
            <Icon :icon="info.Icon" class="w-3.5 h-3.5 opacity-70" />
            <span class="text-[11px]">{{ info.label }}</span>
          </div>
          <span class="font-medium text-base-content/90 truncate max-w-[150px] text-right select-all" :title="info.value">
            {{ info.value }}
          </span>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="p-3 pt-0 mt-auto flex gap-2">
      <router-link
        :to="{ name: 'ClientConsole', params: { id: item.id } }"
        class="btn btn-primary btn-sm flex-1 gap-2 font-normal shadow-primary/20 hover:shadow-primary/40 hover:-translate-y-0.5 transition-all duration-200"
      >
        <Icon icon="mdi:console" class="w-4 h-4" />
        控制台
      </router-link>

      <div class="dropdown dropdown-end dropdown-top">
        <div tabindex="0" role="button" class="btn btn-square btn-sm btn-ghost hover:bg-base-200 transition-colors">
          <Icon icon="mdi:dots-vertical" class="w-5 h-5 text-base-content/60" />
        </div>
        <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow-lg bg-base-100 rounded-xl w-32 border border-base-200">
          <li>
            <button class="text-xs font-medium hover:text-warning hover:bg-warning/10 active:bg-warning/20">
              <Icon icon="mdi:arrow-up-bold-hexagon-outline" class="w-4 h-4" />
              升级
            </button>
          </li>
          <li>
            <button class="text-xs font-medium text-error hover:bg-error/10 active:bg-error/20" @click="deleteClient">
              <Icon icon="mdi:trash-can-outline" class="w-4 h-4" />
              解绑
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

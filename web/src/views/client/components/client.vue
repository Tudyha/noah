<script setup lang="ts">
import { clientStatusMap, clientOsTypeIconMap, clientOsTypeColorMap } from '@/map'
import type { ClientResponse } from '@/types'
import { formatBytes, formatDateTime, formatUptime } from '@/utils'
import { deleteClient as deleteClientApi} from '@/api/client'

const props = defineProps<{
  item: ClientResponse;
}>()

const emit = defineEmits(['refresh']);

const statusColor = computed(() => {
  return props.item.status === 1 ? 'status-success' : 'status-error'
})

const badgeColor = computed(() => {
  return props.item.status === 1 ? 'badge-success' : 'badge-error'
})

const rowInfo = computed(() => [
  {
    Icon: "mdi:account",
    label: "用户",
    value: `${props.item.username} (UID: ${props.item.uid})`
  },
  {
    Icon: "mdi:ip-outline",
    label: "ip地址",
    value: `${props.item.remote_ip}:${props.item.port}`
  },
  {
    Icon: "mdi:clock",
    label: "系统运行时间",
    value: formatUptime(props.item.uptime)
  },
  {
    Icon: "material-symbols:online-prediction",
    label: "上次在线时间",
    value: formatDateTime(props.item.last_online_time)
  },
]);

const statInfo = computed(() => [
  {
    Icon: "mdi:cpu-64-bit",
    label: "CPU",
    value: `${props.item.cpu_num} 核`
  },
  {
    Icon: "mdi:memory",
    label: "内存总量",
    value: formatBytes(props.item.mem_total)
  },
  {
    Icon: "mdi:disk",
    label: "磁盘总量",
    value: formatBytes(props.item.disk_total)
  },

]);

const deleteClient = () => {
  deleteClientApi(props.item.id)
  emit('refresh')
}

</script>

<template>
  <div class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow border border-base-300">
    <div class="card-body p-1">
      <!-- header -->
      <div class="flex justify-between items-start">
        <div class="flex items-center gap-3">
          <div :class="['avatar', 'placeholder', clientOsTypeColorMap[item.os_type], 'text-white']">
            <Icon :icon="clientOsTypeIconMap[item.os_type]" class="w-10 h-10" />
          </div>
          <div class="max-w-42">
            <h2 class="text-base font-bold truncate">
              {{ item.hostname }}
            </h2>
            <span class="font-semibold text-base-content/70 text-xs">
              {{ item.platform_family }}-{{ item.platform_version }}
            </span>
          </div>
        </div>

        <!-- 在线状态 -->
        <div class="flex items-center justify-between gap-1 badge px-1" :class="badgeColor">
          <div class="inline-grid *:[grid-area:1/1]">
            <div :class="['status animate-ping', statusColor]"></div>
            <div :class="['status', statusColor]"></div>
          </div>
          {{ clientStatusMap[item.status] }}
        </div>
      </div>
      <!-- Body -->
      <div class="p-1">
        <div class="stats w-full shadow-inner bg-base-200 border border-base-300">
          <div class="stat px-4 py-2" v-for="(item, index) in statInfo" :key="index">
            <div class="stat-title text-xs">{{ item.label }}</div>
            <div class="stat-value text-lg">{{ item.value }}</div>
          </div>
        </div>

        <div class="space-y-2 text-sm mt-4">
          <template v-for="(item, index) in rowInfo" :key="index">
            <div class="flex justify-between border-b border-base-300">
              <div class="flex items-center gap-2">
                <Icon :icon="item.Icon" class="text-blue-500 text-xs" />
                <span class="font-mono text-xs">{{ item.label }}</span>
              </div>
              <div class="text-xs">
                {{ item.value }}
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <div class="card-actions flex justify-between gap-1 p-1">
      <router-link :to="{name: 'ClientConsole', params: {id: item.id}}" class="btn btn-primary btn-sm flex-1">
        <Icon icon="mdi:console" />
        进入控制台
      </router-link>
      <!-- 更多操作下拉菜单 -->
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-secondary btn-sm px-2 tooltip tooltip-left" data-tip="更多操作">
          <Icon icon="mdi:dots-vertical" />
        </div>
        <ul tabindex="0" class="dropdown-content z-1 menu p-2 shadow-lg bg-base-300 rounded-box w-30 text-sm">
          <li>
            <button class="text-warning">
              <Icon icon="mdi:arrow-up-bold-hexagon-outline" />
              升级
            </button>
          </li>
          <li>
            <button class="text-error" @click="deleteClient">
              <Icon icon="mdi:trash-can" />
              解绑
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
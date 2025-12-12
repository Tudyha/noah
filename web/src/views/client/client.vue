<script setup lang="ts">
import { clientStatusMap, clientOsTypeIconMap } from '@/map'
import type { ClientResponse } from '@/types'

const props = defineProps<{
  item: ClientResponse;
}>()

const statusColor = computed(() => {
  return props.item.status === 1 ? 'status-success' : 'badge-error'
})

const rowInfo = computed(() => [
  {
    Icon: "mdi:account",
    label: "用户",
    value: `${props.item.username} (UID: ${props.item.uid})`
  },
  {
    Icon: "mdi:cpu-64-bit",
    label: "CPU 核心数",
    value: `${props.item.cpu_num} Cores`
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

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatUptime = (totalSeconds: number) => {
  const days = Math.floor(totalSeconds / (3600 * 24));
  const hours = Math.floor((totalSeconds % (3600 * 24)) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = Math.floor(totalSeconds % 60);

  let parts = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (minutes > 0) parts.push(`${minutes}m`);
  if (seconds > 0 || parts.length === 0) parts.push(`${seconds}s`);

  return parts.join(' ');
};

const formatDateTime = (date: Date) => {
  try {
    if (date.getFullYear() <= 1) return 'N/A';
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    });
  } catch (e) {
    return 'Invalid Date';
  }
};

</script>

<template>
  <div class="card border border-base-400">
    <!-- header -->
    <div class="card-title flex justify-between items-start p-1 border-b border-base-300">
      <div class="flex flex-col">
        <h2 class="text-xl font-extrabold">
          {{ item.hostname }}
        </h2>
        <div class="text-xs">
          <span class="font-semibold text-base-content/70">
            {{ item.platform_family }}-{{ item.platform_version }}-{{ item.kernel_arch }}-{{ item.kernel_version }}
          </span>
        </div>
      </div>
      <div class="flex flex-col items-center w-12">
        <div class="flex items-center justify-between gap-1">
          <div class="inline-grid *:[grid-area:1/1]">
            <div :class="['status animate-ping', statusColor]"></div>
            <div :class="['status', statusColor]"></div>
          </div>
          <div class="text-sm">{{ clientStatusMap[item.status] }}</div>
        </div>
        <div>
          <Icon :icon="clientOsTypeIconMap[item.os_type]" />
        </div>
      </div>
    </div>

    <!-- Body -->
    <div class="p-1">
      <template v-for="(item, index) in rowInfo" :key="index">
        <div class="border-b border-base-300 flex items-center justify-between py-1 px-2">
          <div class="flex items-center">
            <Icon :icon="item.Icon" class="w-6 h-6 text-blue-500" />
            {{ item.label }}
          </div>
          <div class="">
            {{ item.value }}
          </div>
        </div>
      </template>
    </div>

    <!-- Footer -->
    <div class="flex flex-wrap justify-end p-1">
      <button class="btn btn-sm btn-primary btn-outline">
        <Icon icon="mdi:console" />
        控制台
      </button>
    </div>

  </div>
</template>
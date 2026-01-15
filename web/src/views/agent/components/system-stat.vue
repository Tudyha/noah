<template>
  <div class="flex flex-col h-full space-y-4">
    <div class="flex items-center gap-2 shrink-0">
      <div class="join">
        <input type="datetime-local" class="input input-bordered input-sm join-item w-44" v-model="filter.start"/>
        <div class="join-item btn btn-sm btn-disabled bg-base-200 border-base-300 px-2">-</div>
        <input type="datetime-local" class="input input-bordered input-sm join-item w-44" v-model="filter.end"/>
      </div>
      <button class="btn btn-primary btn-sm gap-2" @click="refresh" :disabled="loading">
        <Icon v-if="loading" icon="mdi:loading" class="animate-spin" />
        <Icon v-else icon="mdi:magnify" />
        查询
      </button>
      <button class="btn btn-ghost btn-sm btn-square" @click="refresh" title="刷新">
        <Icon icon="mdi:refresh" :class="{ 'animate-spin': loading }" />
      </button>
    </div>

    <div class="flex-1 grid grid-cols-1 lg:grid-cols-2 gap-6 min-h-0 overflow-y-auto pr-2">
      <div class="card bg-base-100 border border-base-200 shadow-sm">
        <div class="card-body p-4">
          <h3 class="card-title text-sm opacity-70">CPU 使用率</h3>
          <div class="h-[250px] sm:h-[300px]">
            <Line v-if="!loading" :data="cpuData" :options="chartOptions" />
            <div v-else class="h-full flex items-center justify-center">
              <span class="loading loading-dots loading-md text-primary"></span>
            </div>
          </div>
        </div>
      </div>

      <div class="card bg-base-100 border border-base-200 shadow-sm">
        <div class="card-body p-4">
          <h3 class="card-title text-sm opacity-70">内存使用详情</h3>
          <div class="h-[250px] sm:h-[300px]">
            <Line v-if="!loading" :data="memoryData" :options="chartOptions" />
            <div v-else class="h-full flex items-center justify-center">
              <span class="loading loading-dots loading-md text-primary"></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import Line from '@/components/chart/line.vue'
import { getAgentSystemInfo } from '@/api/agent'
import type { AgentSystemInfoResponse } from '@/types'
import { formatBytesToGB } from '@/utils'

const props = defineProps<{
  id: string
}>()

// 获取今日日期范围
const today = new Date();
const startOfToday = new Date(today.setHours(0,0,0,0)).toISOString().slice(0, 16);
const endOfToday = new Date(today.setHours(23,59,59,999)).toISOString().slice(0, 16);

const filter = ref({
  start: startOfToday,
  end: endOfToday
})

const data = ref<AgentSystemInfoResponse[]>([])
const loading = ref(false)

const labels = computed(() => {
  return data.value.map(item => {
    const date = new Date(item.created_at);
    return date.getHours().toString().padStart(2, '0') + ':' +
           date.getMinutes().toString().padStart(2, '0');
  })
})

const cpuData = computed(() => {
  return {
    labels: labels.value,
    datasets: [
      {
        label: 'CPU使用率(%)',
        data: data.value.map(item => item.cpu_percent),
        fill: true,
        backgroundColor: 'rgba(75, 192, 192, 0.1)',
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.4,
        pointRadius: 2,
      }
    ]
  }
})

const memoryData = computed(() => {
  return {
    labels: labels.value,
    datasets: [
      {
        label: '使用率(%)',
        data: data.value.map(item => item.mem_used_percent),
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.1)',
        fill: true,
        tension: 0.4,
        pointRadius: 2,
        yAxisID: 'y',
      },
      {
        label: '已用(GB)',
        data: data.value.map(item => formatBytesToGB(item.mem_used)),
        borderColor: 'rgb(255, 99, 132)',
        tension: 0.4,
        pointRadius: 2,
        yAxisID: 'y1',
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        boxWidth: 12,
        font: { size: 11 }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      max: 100,
      grid: { color: 'rgba(0,0,0,0.05)' }
    },
    y1: {
      position: 'right' as const,
      beginAtZero: true,
      grid: { display: false }
    },
    x: {
      grid: { display: false },
      ticks: {
        maxRotation: 0,
        font: { size: 10 }
      }
    }
  }
}

const refresh = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getAgentSystemInfo(Number(props.id))
    data.value = res
    console.log(data.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>

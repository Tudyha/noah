<template>
  <div class="space-y-6 bg-base-200/50 min-h-screen">
    <!-- 欢迎栏 -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-base-content">仪表盘</h1>
        <p class="text-base-content/60 text-sm mt-1">欢迎回来，这是您的系统实时运行状态。</p>
      </div>
      <div class="flex items-center gap-3">
        <div
          class="hidden sm:flex items-center gap-2 px-3 py-1.5 bg-success/10 text-success rounded-full border border-success/20 text-xs font-medium">
          <div class="w-2 h-2 rounded-full bg-success animate-pulse"></div>
          系统运行中
        </div>
        <div class="text-xs text-base-content/40 font-mono">
          最后更新: {{ lastUpdateTime }}
        </div>
      </div>
    </div>

    <!-- 核心统计指标 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in stats" :key="stat.title"
        class="group relative overflow-hidden card bg-base-100 border border-base-200 shadow-sm hover:shadow-md transition-all duration-300">
        <div class="card-body p-5">
          <div class="flex justify-between items-start">
            <div class="space-y-1">
              <p class="text-sm font-medium text-base-content/60">{{ stat.title }}</p>
              <h3 class="text-2xl font-bold tracking-tight">{{ stat.value }}</h3>
            </div>
            <div :class="`p-2.5 rounded-xl bg-opacity-10 ${stat.color}`">
              <Icon :icon="stat.icon" class="w-6 h-6" :class="stat.textColor" />
            </div>
          </div>
        </div>
        <!-- 背景装饰 -->
        <div class="absolute -right-2 -bottom-2 opacity-[0.03] group-hover:scale-110 transition-transform duration-500">
          <Icon :icon="stat.icon" class="w-24 h-24" />
        </div>
      </div>
    </div>

    <!-- 图表与系统状态 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 流量趋势图 -->
      <div class="lg:col-span-2 card bg-base-100 border border-base-200 shadow-sm overflow-hidden">
        <div class="card-body p-0">
          <div class="flex items-center justify-between p-6 pb-0">
            <div>
              <h2 class="card-title text-lg font-bold flex items-center gap-2">
                <Icon icon="mdi:chart-line" class="w-5 h-5 text-primary" />
                流量趋势
              </h2>
              <p class="text-xs text-base-content/50 mt-1">实时监测入站与出站流量波动</p>
            </div>
            <div class="tabs tabs-boxed bg-base-200/50 p-1">
              <button v-for="range in ['7d', '30d']" :key="range" @click="timeRange = range"
                class="tab tab-xs rounded-md px-4 transition-all"
                :class="{ 'tab-active bg-base-100 shadow-sm': timeRange === range }">
                {{ range === '7d' ? '最近7天' : '最近30天' }}
              </button>
            </div>
          </div>
          <div class="p-6 pt-4">
            <div class="h-[320px] w-full">
              <LineChart :data="chartData" :options="chartOptions" />
            </div>
          </div>
        </div>
      </div>

      <!-- 系统状态与资源 -->
      <div class="card bg-base-100 border border-base-200 shadow-sm">
        <div class="card-body p-6">
          <h2 class="card-title text-lg font-bold mb-6 flex items-center gap-2">
            <Icon icon="mdi:server" class="w-5 h-5 text-primary" />
            系统资源
          </h2>

          <div class="space-y-7">
            <div v-for="resource in systemResources" :key="resource.name" class="space-y-2.5">
              <div class="flex justify-between items-end">
                <div class="flex items-center gap-2">
                  <div :class="`w-1.5 h-1.5 rounded-full ${resource.dotClass}`"></div>
                  <span class="text-sm font-semibold">{{ resource.name }}</span>
                </div>
                <div class="flex items-baseline gap-1">
                  <span class="text-lg font-bold tracking-tight" :class="resource.textColor">{{ resource.value }}</span>
                  <span class="text-[10px] text-base-content/40 font-medium">%</span>
                </div>
              </div>
              <div class="relative h-2 w-full bg-base-200 rounded-full overflow-hidden">
                <div class="absolute top-0 left-0 h-full transition-all duration-1000 ease-out rounded-full"
                  :class="resource.progressClass" :style="{ width: `${resource.value}%` }"></div>
              </div>
              <div class="flex justify-between text-[10px] text-base-content/40 font-medium px-0.5">
                <span>{{ resource.desc }}</span>
                <span>{{ resource.usage }}</span>
              </div>
            </div>
          </div>

          <div class="divider my-8 opacity-50"></div>

          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-base-200/30 rounded-xl border border-base-200/50">
              <div class="flex items-center gap-3">
                <div class="p-2 bg-info/10 text-info rounded-lg">
                  <Icon icon="mdi:clock-outline" class="w-4 h-4" />
                </div>
                <span class="text-sm font-medium text-base-content/70">运行时间</span>
              </div>
              <span class="text-sm font-mono font-bold">{{ dashboard.sys_info?.uptime || '' }}</span>
            </div>
            <div class="flex items-center justify-between p-3 bg-base-200/30 rounded-xl border border-base-200/50">
              <div class="flex items-center gap-3">
                <div class="p-2 bg-warning/10 text-warning rounded-lg">
                  <Icon icon="mdi:package-variant-closed" class="w-4 h-4" />
                </div>
                <span class="text-sm font-medium text-base-content/70">版本号</span>
              </div>
              <span class="badge badge-ghost font-mono text-[10px] font-bold tracking-wider">V1.2.4-STABLE</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LineChart from '@/components/chart/line.vue'
import { getDashboard } from '@/api/dashboard'
import dayjs from 'dayjs'
import type { DashboardResponse } from '@/types'

const timeRange = ref('7d')
const lastUpdateTime = ref(dayjs().format('HH:mm:ss'))
const dashboard = ref<DashboardResponse>({})

// 统计指标数据
const stats = computed(() => [
  {
    title: '在线主机数',
    value: dashboard.value.agent_stats?.online || 0,
    icon: 'material-symbols:computer-outline',
    color: 'bg-primary',
    textColor: 'text-primary',
  },
  {
    title: '离线主机数',
    value: dashboard.value.agent_stats?.offline || 0,
    icon: 'material-symbols:computer-outline',
    color: 'bg-secondary',
    textColor: 'text-secondary',
  },
])

// 系统资源数据
const systemResources = computed(() => [
  {
    name: 'CPU 使用率',
    value: dashboard.value.cpu_usage || 24,
    desc: '8 核处理架构',
    usage: '1.2GHz / 3.4GHz',
    progressClass: (dashboard.value.cpu_usage || 24) > 80 ? 'bg-error' : 'bg-primary',
    textColor: (dashboard.value.cpu_usage || 24) > 80 ? 'text-error' : 'text-primary',
    dotClass: (dashboard.value.cpu_usage || 24) > 80 ? 'bg-error' : 'bg-primary'
  },
  {
    name: '内存占用',
    value: dashboard.value.ram_usage || 45,
    desc: '系统总计 16GB',
    usage: '7.2GB / 16GB',
    progressClass: (dashboard.value.ram_usage || 45) > 85 ? 'bg-error' : 'bg-secondary',
    textColor: (dashboard.value.ram_usage || 45) > 85 ? 'text-error' : 'text-secondary',
    dotClass: (dashboard.value.ram_usage || 45) > 85 ? 'bg-error' : 'bg-secondary'
  },
  {
    name: '磁盘空间',
    value: dashboard.value.disk_usage || 12,
    desc: 'NVMe SSD 512GB',
    usage: '62GB / 512GB',
    progressClass: (dashboard.value.disk_usage || 12) > 90 ? 'bg-error' : 'bg-accent',
    textColor: (dashboard.value.disk_usage || 12) > 90 ? 'text-error' : 'text-accent',
    dotClass: (dashboard.value.disk_usage || 12) > 90 ? 'bg-error' : 'bg-accent'
  }
])

const fetchData = async () => {
  try {
    const res = await getDashboard()
    dashboard.value = { ...res }
    lastUpdateTime.value = dayjs().format('HH:mm:ss')
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  }
}

onMounted(() => {
  fetchData()
})

// 定时刷新 (30s)
useIntervalFn(fetchData, 30000)

const chartData = computed(() => ({
  labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '23:59'],
  datasets: [
    {
      label: '入站流量 (MB)',
      backgroundColor: 'rgba(99, 102, 241, 0.1)',
      borderColor: '#6366f1',
      borderWidth: 3,
      pointBackgroundColor: '#6366f1',
      pointBorderColor: '#fff',
      pointHoverRadius: 6,
      data: [450, 380, 520, 780, 610, 890, 720],
      fill: true,
      tension: 0.4
    },
    {
      label: '出站流量 (MB)',
      backgroundColor: 'rgba(236, 72, 153, 0.1)',
      borderColor: '#ec4899',
      borderWidth: 3,
      pointBackgroundColor: '#ec4899',
      pointBorderColor: '#fff',
      pointHoverRadius: 6,
      data: [210, 180, 290, 420, 350, 510, 430],
      fill: true,
      tension: 0.4
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index' as const,
  },
  plugins: {
    legend: {
      position: 'top' as const,
      align: 'end' as const,
      labels: {
        usePointStyle: true,
        padding: 20,
        boxWidth: 8,
        font: {
          size: 12,
          weight: '600'
        }
      }
    },
    tooltip: {
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      titleColor: '#1f2937',
      bodyColor: '#4b5563',
      borderColor: '#e5e7eb',
      borderWidth: 1,
      padding: 12,
      boxPadding: 6,
      usePointStyle: true,
      callbacks: {
        label: function (context: any) {
          return ` ${context.dataset.label}: ${context.parsed.y} MB`
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      border: { display: false },
      grid: {
        color: 'rgba(0, 0, 0, 0.04)',
      },
      ticks: {
        font: { size: 11 },
        padding: 10
      }
    },
    x: {
      border: { display: false },
      grid: { display: false },
      ticks: {
        font: { size: 11 },
        padding: 10
      }
    }
  }
}
</script>

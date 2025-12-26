<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold text-base-content">仪表盘</h2>
      <div class="text-sm text-base-content/60">{{ new Date().toLocaleDateString() }}</div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="stats shadow-sm border border-base-200 bg-base-100">
        <div class="stat">
          <div class="stat-figure text-primary bg-primary/10 p-2 rounded-full">
            <Icon icon="mdi:server" size="24" />
          </div>
          <div class="stat-title text-base-content/70">在线主机</div>
          <div class="stat-value text-primary text-3xl">25.6K</div>
          <div class="stat-desc text-success font-medium">↗︎ 21% 比上月</div>
        </div>
      </div>
      <div class="stats shadow-sm border border-base-200 bg-base-100">
        <div class="stat">
          <div class="stat-figure text-secondary bg-secondary/10 p-2 rounded-full">
            <Icon icon="mdi:bag-personal" size="24" />
          </div>
          <div class="stat-title text-base-content/70">活跃用户</div>
          <div class="stat-value text-secondary text-3xl">2,600</div>
          <div class="stat-desc text-success font-medium">↗︎ 12% 比上周</div>
        </div>
      </div>
      <div class="stats shadow-sm border border-base-200 bg-base-100">
        <div class="stat">
          <div class="stat-figure text-accent bg-accent/10 p-2 rounded-full">
            <Icon icon="mdi:chart-line" size="24" />
          </div>
          <div class="stat-title text-base-content/70">系统负载</div>
          <div class="stat-value text-accent text-3xl">86%</div>
          <div class="stat-desc text-warning font-medium">需关注</div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="card bg-base-100 shadow-sm border border-base-200">
      <div class="card-body p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="card-title text-lg">流量趋势</h2>
          <select class="select select-bordered select-xs">
            <option>最近7天</option>
            <option>最近30天</option>
          </select>
        </div>
        <div class="h-80 w-full">
          <LineChart :data="chartData" :options="chartOptions" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LineChart from '@/components/chart/line.vue'

const chartData = {
  labels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
  datasets: [
    {
      label: '入站流量 (MB)',
      backgroundColor: 'rgba(99, 102, 241, 0.2)',
      borderColor: 'rgb(99, 102, 241)',
      borderWidth: 2,
      pointBackgroundColor: 'rgb(99, 102, 241)',
      data: [40, 39, 10, 40, 39, 80, 40],
      fill: true,
      tension: 0.4
    },
    {
      label: '出站流量 (MB)',
      backgroundColor: 'rgba(236, 72, 153, 0.2)',
      borderColor: 'rgb(236, 72, 153)',
      borderWidth: 2,
      pointBackgroundColor: 'rgb(236, 72, 153)',
      data: [20, 25, 30, 28, 35, 40, 38],
      fill: true,
      tension: 0.4
    }
  ]
}

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
      align: 'end' as const,
      labels: {
        usePointStyle: true,
        boxWidth: 8
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: {
        display: true,
        color: 'rgba(0, 0, 0, 0.05)'
      }
    },
    x: {
      grid: {
        display: false
      }
    }
  }
}
</script>

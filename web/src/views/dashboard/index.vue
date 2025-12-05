<template>
  <div>
    <h2 class="text-3xl font-bold mb-6">仪表盘</h2>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-figure text-primary">
            <Icon icon="mdi:server" size="32" />
          </div>
          <div class="stat-title">在线主机</div>
          <div class="stat-value text-primary">25.6K</div>
          <div class="stat-desc">比上月增长 21%</div>
        </div>
      </div>
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-figure text-secondary">
            <Icon icon="mdi:bag-personal" size="32" />
          </div>
          <div class="stat-title">活跃用户</div>
          <div class="stat-value text-secondary">2,600</div>
          <div class="stat-desc">↗︎ 400 (22%)</div>
        </div>
      </div>
      <div class="stats shadow bg-base-100">
        <div class="stat">
          <div class="stat-figure text-accent">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
              class="inline-block w-8 h-8 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z">
              </path>
            </svg>
          </div>
          <div class="stat-title">系统负载</div>
          <div class="stat-value text-accent">86%</div>
          <div class="stat-desc">需关注</div>
        </div>
      </div>
    </div>

    <!-- 模拟图表区域 -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">流量趋势</h2>
        <div
          class="w-full h-64 bg-base-200 rounded-box flex items-center justify-center border-2 border-dashed border-base-300 text-base-content/50">
          [ 此处集成 ECharts 图表 ]
        </div>
      </div>
    </div>
  </div>
</template>

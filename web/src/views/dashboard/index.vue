<template>
  <el-container>
    <el-header>
      <el-row :gutter="12">
        <el-col :span="6">
          <el-card shadow="always" class="card">
            <div slot="header">
              <span><b>服务器信息</b></span>
            </div>
            <div>
              主机名称: {{ dashboardData.hostname }}<br><br>
              启动时间：{{ parseTime(dashboardData.bootTime) }}<br><br>
              OS: {{ dashboardData.platform }} {{ dashboardData.platformFamily }}-{{ dashboardData.platformVersion }} {{ dashboardData.kernelArch }}-{{ dashboardData.kernelVersion }}<br><br>
<!--              CPU:-->
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="always" class="card">
            <div slot="header">
              <span><b>服务器状态</b></span>
            </div>
            <el-row>
              <el-col :span="12">
                <div ref="memoryChart" class="chart" />
              </el-col>
              <el-col :span="12">
                <div ref="diskChart" class="chart" />
              </el-col>
            </el-row>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="always" class="card">
            <div slot="header">
              <span><b>客户端统计</b></span>
            </div>
            <div ref="clientChart" class="chart" />
          </el-card>
        </el-col>
      </el-row>
    </el-header>
  </el-container>
</template>

<script>
import * as echarts from 'echarts'
import { getDashboardData } from '@/api/dashboard'
import { parseTime } from '@/utils'

export default {
  name: 'Dashboard',
  data() {
    return {
      parseTime,
      cpuChart: null,
      memoryChart: null,
      diskChart: null,
      clientChart: null,
      memoryData: [],
      diskData: [],
      clientData: [],
      dashboardData: {}
    }
  },
  mounted() {
    this.initCharts()
  },
  created() {
    this.getDashboardData()
  },
  methods: {
    getDashboardData() {
      getDashboardData().then(response => {
        if (response.code === 0) {
          this.dashboardData = response.data
          this.memoryData = [
            { value: response.data.memoryUsed, name: '已用(GB)' },
            { value: response.data.memoryRemain, name: '剩余(GB)' }
          ]
          this.diskData = [
            { value: response.data.diskUsed, name: '已用(GB)' },
            { value: response.data.diskFree, name: '剩余(GB)' }
          ]
          this.clientData = [
            { value: response.data.clientOnlineCount, name: '在线' },
            { value: response.data.clientOfflineCount, name: '离线' }
          ]
        }
      })
    },
    initCharts() {
      setTimeout(() => {
        this.memoryChart = echarts.init(this.$refs.memoryChart)
        this.diskChart = echarts.init(this.$refs.diskChart)
        this.clientChart = echarts.init(this.$refs.clientChart)

        this.updateChart(this.clientChart, this.clientData, '')
        this.updateChart(this.memoryChart, this.memoryData, '内存')
        this.updateChart(this.diskChart, this.diskData, '磁盘')
      }, 500)
    },
    updateChart(c, data, title) {
      // 配置项
      const option = {
        title: {
          text: title
        },
        tooltip: {
          trigger: 'item'
        },
        label: {
          show: true,
          formatter: '{b}:{c}'
        },
        // grid: {
        //   top: 1200,
        //   containLabel: true
        // },
        series: [
          {
            type: 'pie',
            radius: ['40%', '80%'],
            startAngle: 180,
            endAngle: 360,
            data: data
          }
        ]
      }
      c.setOption(option)
    }
  }
}
</script>

<style lang="scss" scoped>
.card {
  height: 200px;
}

.chart {
  height: 180px;
  width: 100%;
}
</style>

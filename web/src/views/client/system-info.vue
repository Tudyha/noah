<template>
  <div>
    <el-row>
      <el-col :span="24">
        <el-time-picker
          is-range
          v-model="timeRange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          @change="onTimeRangeChange"
        ></el-time-picker>
      </el-col>
    </el-row>
  <div class="dashboard-container">
    <div class="chart-container">
      <h4>CPU</h4>
      <div ref="cpuChart" class="chart"></div>
    </div>
    <div class="chart-container">
      <h4>内存</h4>
      <div ref="memoryChart" class="chart"></div>
    </div>
  </div>
  <div class="dashboard-container">
    <div class="chart-container">
      <h4>磁盘</h4>
      <div ref="diskChart" class="chart"></div>
    </div>
    <div class="chart-container">
      <h4>带宽</h4>
      <div ref="bandwidthChart" class="chart"></div>
    </div>
  </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { systemInfo } from '@/api/client'
import { parseTime } from '@/utils'

export default {
  name: 'ClientSystemInfo',
  data() {
    return {
      id: null,
      timeRange: [],
      cpuChart: null,
      memoryChart: null,
      diskChart: null,
      bandwidthChart: null,
      intervalId: null,
      cpuData: [],
      memoryData: [],
      diskData: [],
      bandwidthData: [],
      cpuCharSeries: [{ name: 'CPU(%)', key: 'cpuUsage'}],
      memoryCharSeries: [
        { name: '已用内存(GB)', key: 'memoryUsed'},
        { name: '空闲内存(GB)', key: 'memoryFree'},
        { name: '可用内存(GB)', key: 'memoryAvailable'},
        { name: '内存占用百分比(%)', key: 'memoryPercent'}
      ],
      diskCharSeries: [
        { name: '已用空间(GB)', key: 'diskUsed'},
        { name: '剩余空间(GB)', key: 'diskFree'}
      ],
      bandwidthCharSeries: [
        { name: '入口带宽(MB/s)', key: 'bandwidthIn'},
        { name: '出口带宽(MB/s)', key: 'bandwidthOut'}
      ]
    };
  },
  created() {
    //query上获取id
    this.id = this.$route.query.id;
    //时间组件默认5分钟
    this.timeRange = [new Date(new Date().getTime() - 5 * 60 * 1000), new Date()];
  },
  mounted() {
    this.initCharts();
    this.fetchData();
  },
  methods: {
    initCharts() {
      this.cpuChart = echarts.init(this.$refs.cpuChart);
      this.memoryChart = echarts.init(this.$refs.memoryChart);
      this.diskChart = echarts.init(this.$refs.diskChart);
      this.bandwidthChart = echarts.init(this.$refs.bandwidthChart);

      this.updateCharts();
    },
    fetchData() {
      this.clearData()
      systemInfo({ id: this.id, start: parseTime(this.timeRange[0]), end: parseTime(this.timeRange[1]) })
        .then(response => {
          const data = response.data;
          data.forEach(item => {
            const t = new Date(item.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
            this.cpuData.push({ time: t, cpuUsage: item.cpuUsage });
            this.memoryData.push({
              time: t, memoryUsed: item.memoryUsed,
              memoryFree: item.memoryFree,
              memoryAvailable: item.memoryAvailable,
              memoryPercent: item.memoryPercent });
            this.diskData.push({
              time: t,
              diskUsed: item.diskUsed,
              diskFree: item.diskFree
            });
            this.bandwidthData.push({
              time: t,
              bandwidthIn: item.bandwidthIn,
              bandwidthOut: item.bandwidthOut
            });
          })
          this.updateCharts();
        })
        .catch(error => {
          console.error('Error fetching data:', error);
        });
    },
    clearData() {
      this.cpuData = [];
      this.memoryData = [];
      this.diskData = [];
      this.bandwidthData = [];
    },
    updateCharts() {
      this.updateChart(this.cpuChart, this.cpuCharSeries, this.cpuData);
      this.updateChart(this.memoryChart, this.memoryCharSeries, this.memoryData);
      this.updateChart(this.diskChart, this.diskCharSeries, this.diskData);
      this.updateChart(this.bandwidthChart, this.bandwidthCharSeries, this.bandwidthData);
    },
    updateChart(chart, series, data) {
      const so = series.map(s => {
        return {
          name: s.name,
          data: data.map((item) => item[s.key]),
          type: 'line',
          smooth: true
        }
      })
      chart.setOption({
        legend: {
          data: series.map(s => s.name)
        },
        tooltip: {
          trigger: 'axis',
          formatter: function (params) {
            let result = '';
            params.forEach(function (param) {
              result += `${param.seriesName}: ${param.data}<br>`;
            });
            return result;
          }
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: data.map(item => item.time)
        },
        yAxis: {
          type: 'value'
        },
        series: so
      });
    },
    onTimeRangeChange(range) {
      if (!range || !range.length || range[0] === '' || range[1] === '') {
        return;
      }
      this.fetchData()
    },
  },
  beforeDestroy() {
    if (this.cpuChart) {
      this.cpuChart.dispose();
    }
    if (this.memoryChart) {
      this.memoryChart.dispose();
    }
    if (this.diskChart) {
      this.diskChart.dispose();
    }
    if (this.bandwidthChart) {
      this.bandwidthChart.dispose();
    }
  }
}
</script>

<style>
.dashboard-container {
  display: flex;
  justify-content: space-around;
  margin-bottom: 10px; /* 添加一些间隔 */
}

.chart-container {
  flex: 1; /* 每个图表容器都有相同的宽度 */
  padding: 10px; /* 内部间隔 */
  box-sizing: border-box; /* 包括padding在内的宽度计算 */
}

.chart {
  height: 400px; /* 图表高度 */
}

h3 {
  margin: 0 0 0 0; /* 标题的外边距 */
}
</style>

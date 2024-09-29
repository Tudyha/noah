<template>
  <el-container class="container">
    <el-header class="header" height="80px">
      <el-card>
        <el-row>
          <el-col :span="8">{{ hostname }}</el-col>
          <el-col :span="8">CPU: {{ cpuCores }}核<br>内存: {{ memoryTotal }}</el-col>
          <el-col :span="8">内网IP: {{ localIp }}<br>公网IP: {{ remoteIp }}</el-col>
        </el-row>
      </el-card>
    </el-header>
    <el-container style="padding: 20px">
      <el-main class="content">
        <el-tabs v-model="currentComponent" tab-position="left" @tab-click="handleTagClick">
          <el-tab-pane
            v-for="(item, index) in menuItems"
            :key="index"
            :label="item.label"
            :name="item.index"
            :icon="item.icon"
          />
          <keep-alive>
            <component :is="currentComponent" :id="id" />
          </keep-alive>
        </el-tabs>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import Load from './components/system-info.vue'
import Status from './components/system-status.vue'
import { getClient } from '@/api/client'
import Terminal from '@/components/Terminal'
import File from '@/components/FileTree'
import Channel from './components/chennel.vue'
import App from './components/app.vue'

export default {
  name: 'Console',
  components: {
    Load,
    Status,
    Terminal,
    File,
    Channel,
    App
  },
  data() {
    return {
      id: null,
      hostname: '',
      cpuCores: '',
      memoryTotal: '',
      localIp: '',
      remoteIp: '',
      currentComponent: 'load', // 默认显示资源负载
      menuItems: [
        { index: 'load', label: '资源负载', icon: 'el-icon-cpu' },
        { index: 'status', label: '系统状态', icon: 'el-icon-monitor' },
        { index: 'terminal', label: '在线终端', icon: 'el-icon-s-promotion' },
        { index: 'file', label: '文件管理', icon: 'el-icon-folder-opened' },
        { index: 'channel', label: '隧道代理', icon: 'el-icon-connection' },
        { index: 'app', label: '应用管理', icon: 'el-icon-s-grid' }
      ]
    }
  },
  created() {
    this.id = +this.$route.query.id
    this.fetchSystemInfo()
  },
  methods: {
    async fetchSystemInfo() {
      try {
        const response = await getClient(this.id)
        const { hostname, cpuCores, memoryTotal, localIp, remoteIp } = response.data
        this.hostname = hostname
        this.cpuCores = cpuCores
        this.memoryTotal = memoryTotal
        this.localIp = localIp
        this.remoteIp = remoteIp
      } catch (error) {
        console.error('Error fetching system info:', error)
      }
    },
    handleTagClick(tab) {
      this.changeContent(tab.name)
    },
    changeContent(selected) {
      this.currentComponent = selected
    }
  }
}
</script>

<style scoped>
.container {
  height: 94vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content {
  padding: 0;
  overflow-y: auto;
  height: calc(100vh - 150px);
}
</style>

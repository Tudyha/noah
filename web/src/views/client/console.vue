<template>
  <el-container class="container">
    <el-header class="header" height="80px">
      <el-card>
        <el-row>
              <el-col :span="8">{{ hostname }}</el-col>
              <el-col :span="8">CPU: {{ cpuCores }}核<br />内存: {{ memoryTotal }}G</el-col>
              <el-col :span="8">内网IP: {{ localIp }}<br />公网IP: {{ remoteIp }}</el-col>
        </el-row>
      </el-card>
    </el-header>
    <el-container style="padding: 20px">
      <el-aside width="200px" class="sidebar">
<!--        <el-card>-->
          <el-card class="menu-card">
            <el-button v-for="(item, index) in menuItems"
                       :key="index"
                       :class="{ active: item.index === currentComponent }"
                       :icon="item.icon"
                       @click="changeContent(item.index)">
              {{ item.label }}
            </el-button>
          </el-card>
<!--        </el-card>-->
      </el-aside>
      <el-main class="content">
        <el-card class="main-card">
          <keep-alive>
            <component :is="currentComponent" :id="id"></component>
          </keep-alive>
        </el-card>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import Load from './components/system-info.vue'
import Status from './components/system-status.vue'
import { getClient } from '@/api/client'
import Shell from '@/components/Shell'
import File from '@/components/FileTree'
import Channel from './components/chennel.vue'
import App from './components/app.vue'

export default {
  name: 'Console',
  components: {
    Load,
    Status,
    Shell,
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
        { index: 'shell', label: '在线终端', icon: 'el-icon-s-promotion' },
        { index: 'file', label: '文件管理', icon: 'el-icon-folder-opened' },
        { index: 'channel', label: '隧道管理', icon: 'el-icon-connection' },
        { index: 'app', label: '应用管理', icon: 'el-icon-s-grid' }
      ],
    };
  },
  created() {
    this.id = +this.$route.query.id;
    this.fetchSystemInfo();
  },
  methods: {
    async fetchSystemInfo() {
      try {
        const response = await getClient(this.id);
        const { hostname, cpuCores, memoryTotal, localIp, remoteIp } = response.data;
        this.hostname = hostname;
        this.cpuCores = cpuCores;
        this.memoryTotal = memoryTotal;
        this.localIp = localIp;
        this.remoteIp = remoteIp;
      } catch (error) {
        console.error('Error fetching system info:', error);
      }
    },
    changeContent(selected) {
      this.currentComponent = selected;
    },
  },
};
</script>

<style scoped>
.container {
  height: 94vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar {
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: stretch;
}

.menu-card {
  padding: 10px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.menu-card .el-button {
  width: 100%;
  margin-bottom: 10px;
  background-color: aliceblue;
  color: rgba(3, 1, 1, 0.77);
  border: none;
  font-size: 16px;
  font-weight: 500;
  text-align: center;
  margin-left: 0;
}

.menu-card .el-button.active {
  background-color: #78c6ee;
  color: #fff;
}

.content {
  padding: 0;
  overflow-y: auto;
  height: calc(100vh - 150px);
}
</style>

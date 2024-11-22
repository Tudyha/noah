<template>
  <div class="app-container">
    <div class="menu-with-refresh">
      <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect">
        <el-menu-item index="1">进程清单</el-menu-item>
        <el-menu-item index="2">网络端口</el-menu-item>
<!--        <el-menu-item index="3">系统账号</el-menu-item>-->
<!--        <el-menu-item index="4">登录历史</el-menu-item>-->
      </el-menu>
      <el-button plain type="text" icon="el-icon-refresh" @click="refreshData">刷新</el-button>
    </div>
    <keep-alive>
      <component :is="currentComponent" :id="id" ref="currentComponent" />
    </keep-alive>
  </div>
</template>

<script>
import Process from './process.vue'; // 引入进程组件
import Network from './network.vue';

export default {
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      activeIndex: '1',
      currentComponent: 'Process', // 初始组件
    };
  },
  methods: {
    handleSelect(key, keyPath) {
      console.log(key, keyPath);
      this.currentComponent = this.getComponentByKey(key);
    },
    getComponentByKey(key) {
      switch (key) {
        case '1':
          return 'Process';
        case '2':
          return 'Network';
        // 添加更多 case 以支持其他组件
        default:
          return 'Process';
      }
    },
    refreshData() {
      this.$refs.currentComponent.refresh();
    }
  },
  components: {
    Process, // 注册进程组件
    Network
  }
}
</script>

<style scoped>
.app-container {
  padding: 2px;
}

.menu-with-refresh {
  display: flex; /* 使用 Flexbox */
  justify-content: space-between; /* 水平方向上两端对齐 */
  align-items: center; /* 垂直居中对齐 */
}

.el-menu-demo {
  border-bottom: none;
  flex-grow: 1; /* 让菜单占据尽可能多的空间 */
}
</style>

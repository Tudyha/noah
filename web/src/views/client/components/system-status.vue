<template>
  <div class="app-container">
    <el-menu :default-active="activeIndex" class="el-menu-demo" mode="horizontal" @select="handleSelect">
      <el-menu-item index="1">进程清单</el-menu-item>
      <!-- 可以添加更多菜单项 -->
    </el-menu>
    <keep-alive>
      <component :is="currentComponent" :id="id"></component>
    </keep-alive>
  </div>
</template>

<script>
import Process from './process.vue'; // 引入进程组件

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
        // 添加更多 case 以支持其他组件
        default:
          return 'Process';
      }
    }
  },
  components: {
    Process // 注册进程组件
  }
}
</script>

<style scoped>
.app-container {
  padding: 2px;
}

.el-menu-demo {
  border-bottom: none;
}
</style>

<template>
  <div>
    <el-tabs v-model="selectedTabIndex" @tab-click="handleTabClick" @tab-remove="removeTab">
      <el-tab-pane
        v-for="(tab) in tabs"
        :key="tab.key"
        :label="tab.label"
        :name="tab.name"
        closable
      >
        <terminal :id="id" :ref="tab.terminalId" :terminal-id="tab.terminalId" :shell-type="shellType" />
      </el-tab-pane>
      <el-tab-pane key="add" name="add">
        <span slot="label" style="padding: 8PX;font-size:20PX;font-weight:bold;">
          +
        </span>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import Terminal from './terminal.vue'

export default {
  components: { Terminal },
  props: {
    id: {
      type: Number,
      default: null
    },
    shellType: {
      type: Number,
      default: 2
    }
  },
  data() {
    return {
      tabs: [], // 保存所有终端标签页的信息
      selectedTabIndex: '', // 保存当前选中的终端标签页索引
      tabIndex: 0
    }
  },
  created() {
    this.addTab()
  },
  methods: {
    addTab() {
      const idx = this.tabIndex++
      this.tabs.push({ key: idx, label: 'Terminal' + idx, name: idx.toString(), terminalId: 'terminal' + idx })
      this.selectedTabIndex = idx.toString()
    },
    handleTabClick(tab) {
      if (tab.name === 'add') {
        this.addTab()
        return false
      } else {
        this.selectedTabIndex = tab.name
      }
    },
    removeTab(idx) {
      for (let i = 0; i < this.tabs.length; i++) {
        if (this.tabs[i].name === idx) {
          const terminals = this.$refs[this.tabs[i].terminalId]
          for (const v in terminals) {
            terminals[v].close()
          }

          this.tabs.splice(i, 1)
        }
      }

      if (this.tabs.length > 0) {
        this.selectedTabIndex = this.tabs[this.tabs.length - 1].name
      }
    }
  }
}
</script>

<style scoped>
.el-tabs__nav-wrap::after {
  display: none;
}

/deep/ .el-tabs__header {
  margin: 0 0 3px;
}
</style>

<template>
  <el-table :data="processList" height="70vh">
    <el-table-column prop="pid" label="进程id" width="70px"></el-table-column>
    <el-table-column prop="name" label="进程名" width="100px"></el-table-column>
    <el-table-column prop="username" label="用户" width="100px"></el-table-column>
    <el-table-column prop="createTime" label="启动时间" :formatter="(row, column, cellValue, index) => parseTime(cellValue)" width="160px"></el-table-column>
    <el-table-column prop="cpu" label="cpu" width="100px"></el-table-column>
    <el-table-column prop="memory" label="内存" width="100px"></el-table-column>
    <el-table-column prop="command" label="命令参数" show-overflow-tooltip="true"></el-table-column>
    <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="50">
      <template slot-scope="{row}">
        <el-dropdown trigger="click" placement="bottom-start" size="small">
          <el-button type="text" size="medium">
            <i class="el-icon-more" />
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item icon="el-icon-close" @click.native="renameItem(row)">结束进程</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </template>
    </el-table-column>
  </el-table>
</template>
<script>
import { fetchProcessList } from '@/api/client'
import { parseTime } from '@/utils'

export default {
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      processList: [],
    };
  },
  methods: {
    parseTime,
    fetchProcessList() {
      fetchProcessList(this.id)
        .then(response => {
          this.processList = response.data;
        })
        .catch(error => {
          console.error('Error fetching process list:', error);
        });
    },
  },
  components: {
  },
  created() {
    this.fetchProcessList();
  },
}
</script>
<style scoped>
</style>

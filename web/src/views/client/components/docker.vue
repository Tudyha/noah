<template>
      <el-table :data="list" height="70vh">
        <el-table-column prop="Id" label="容器id" :show-overflow-tooltip="true" width="130px"></el-table-column>
        <el-table-column prop="Names" label="容器名称" :formatter="(row, column, cellValue, index) => cellValue.map(name => name.replace('/', ''))"></el-table-column>
        <el-table-column prop="ImageID" label="镜像id" :show-overflow-tooltip="true" width="130px"></el-table-column>
        <el-table-column prop="Image" label="镜像名称" :show-overflow-tooltip="true"></el-table-column>
        <el-table-column prop="Created" label="创建时间" :formatter="(row, column, cellValue, index) => formatTime(cellValue, '{y}-{m}-{d} {h}:{i}:{s}')" width="160px"></el-table-column>
        <el-table-column prop="State" label="状态" width="100px"></el-table-column>
        <el-table-column prop="Status" label="Last started" width="100px"></el-table-column>
        <el-table-column prop="Ports" label="端口"  width="180px">
          <template slot-scope="scope">
            <div style="white-space: pre-line;">{{ formatPorts(scope.row.Ports) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="Command" label="命令" :show-overflow-tooltip="true"></el-table-column>
        <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="50">
          <template slot-scope="{row}">
            <el-dropdown trigger="click" placement="bottom-start" size="small">
              <el-button type="text" size="medium">
                <i class="el-icon-more" />
              </el-button>
<!--              <el-dropdown-menu slot="dropdown">-->
<!--                <el-dropdown-item icon="el-icon-close" @click.native="killProcess(row.pid)">结束进程</el-dropdown-item>-->
<!--              </el-dropdown-menu>-->
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
</template>
<script>
import { fetchDockerContainerList } from '@/api/client'
import { formatTime } from '@/utils'

export default {
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      list: [],
    };
  },
  methods: {
    formatTime,
    fetchList() {
      fetchDockerContainerList(this.id)
        .then(response => {
          this.list = response.data;
        })
        .catch(error => {
          console.error('Error fetching process list:', error);
        });
    },
    refresh() {
      this.fetchList();
    },
    formatPorts(ports) {
      //0.0.0.0:1337->1337/tcp
      if (!ports || ports.length === 0) {
        return '';
      }
      return ports.map(port => `${port.IP || ""}:${port.PublicPort || ""}->${port.PrivatePort}/${port.Type}`).join('\n');
    },
  },
  components: {
  },
  created() {
    this.fetchList();
  },
}
</script>
<style scoped>
</style>

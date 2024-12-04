<template>
  <div class="app-container">
    <div class="flex_row date">
      <el-input v-model="listQuery.hostname" style="width: 320px" placeholder="主机名称" class="input-with-select round" />
      <el-select v-model="listQuery.status" clearable placeholder="在线状态">
        <el-option v-for="item in map.statusOptions" :key="item.value" :label="item.label" :value="+item.value" />
      </el-select>
      <el-button class="margin-left-10" type="primary" @click="fetchData">搜索</el-button>
      <el-button class="right-aligned" type="success" @click="handleBind">绑定主机</el-button>
    </div>

    <el-table v-loading="listLoading" class="table-container" :data="list" element-loading-text="Loading" border fit
      highlight-current-row>
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column label="主机">
        <template slot-scope="scope">
          <a @click="handleConsole(scope.row)" class="hostname-link">
            {{ scope.row.hostname }}
          </a>
        </template>
      </el-table-column>

      <el-table-column label="用户名称" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.username }}</span>
        </template>
      </el-table-column>
      <el-table-column label="主机类型" align="center">
        <template slot-scope="scope">
          {{ scope.row.osName }}
        </template>
      </el-table-column>
      <el-table-column label="IP" align="center">
        <template slot-scope="scope">
          内: {{ scope.row.localIp }}<br />外: {{ scope.row.remoteIp }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="在线状态" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status | statusTagTypeFilter">{{ scope.row.status | statusFilter }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="上次在线时间">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ parseTime(scope.row.lastOnlineTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="100">
        <template slot-scope="{row}">
          <el-dropdown trigger="click" placement="bottom-start" size="small">
            <el-button type="text" size="medium">
              <i class="el-icon-more" />
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <!-- <el-dropdown-item @click.native="handleUpdateClient(row)">
                <i class="el-icon-s-tools" /> 更新客户端
              </el-dropdown-item> -->
              <el-dropdown-item @click.native="handleDelete(row)">
                <i class="el-icon-delete" /> 删除
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size"
      @pagination="fetchData" />

    <el-dialog title="绑定主机" :visible.sync="bindDialogShow" width="50%">
      <el-tabs type="border-card">
        <el-tab-pane label="Linux">
          <el-input v-model="command" type="textarea" :rows="5" readonly />
        </el-tab-pane>
        <!-- <el-tab-pane label="Windows">Windows</el-tab-pane> -->
      </el-tabs>
      <div slot="footer" class="dialog-footer">
        <el-button @click="bindDialogShow = false">关闭</el-button>
      </div>
    </el-dialog>

    <!-- <update-client :v-show="updateDialogShow" :client-id="selectedRow.id" :visible="updateDialogShow" :os-type="selectedRow.osType" @hide="updateDialogShow = false" /> -->
  </div>
</template>

<script>
import { fetchList, deleteClient, getClientInstallScript } from '@/api/client'
import Pagination from '@/components/Pagination'
import formMixin from '@/mixins/form-father'
import * as map from '@/map/client'
import { parseTime } from '@/utils'
// import UpdateClient from './components/update-client.vue'

export default {
  name: 'Client',

  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: '离线',
        2: '在线'
      }
      return statusMap[status]
    },
    statusTagTypeFilter(status) {
      const statusMap = {
        1: 'info',
        2: 'success'
      }
      return statusMap[status]
    }
  },
  mixins: [formMixin],
  data() {
    return {
      map,
      list: null,
      listLoading: true,
      selectedRow: {},
      listQuery: {
        page: 1,
        size: 10,
        hostname: '',
        status: 2
      },
      total: 0,
      dialogTitle: '',
      updateDialogShow: false,
      bindDialogShow: false,
      command: ''
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    parseTime,
    fetchData() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        this.listLoading = false
      })
    },
    handleUpdateClient(row) {
      this.selectedRow = row
      this.updateDialogShow = true
    },
    handleConsole(row) {
      this.$router.push({ path: '/client/console', query: { id: row.id } })
    },
    handleDelete(row) {
      this.$confirm('是否确认删除：' + row.hostname, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteClient(row.id).then((res) => {
          if (res.code === 0) {
            this.$message.success('删除成功')
            this.fetchData()
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      })
    },
    async handleBind() {
      const res = await getClientInstallScript()
      const tempToken = res.data
      let ip, port

      ip = window.location.hostname
      port = window.location.port

      this.command = `curl -kfsSL 'http://${ip}:${port}/api/file/download/install-cli?token=${tempToken}' | bash -s -- ${ip} ${port} ${tempToken}`
      this.bindDialogShow = true
    }
  }
}
</script>

<style scoped>
.hostname-link {
  color: #0e96dc;
}

.right-aligned {
  float: right;
  /* 使用 float 属性 */
  /* 或者使用 Flexbox */
  /* display: flex;
     justify-content: flex-end; */
}
</style>

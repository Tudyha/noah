<template>
  <div class="app-container">
    <collapse-filter>
      <template v-slot:collapse-title-left>
        <div class="flex_row date">
          <el-input
            v-model="listQuery.hostname"
            style="width: 320px"
            placeholder="主机名称"
            class="input-with-select round"
          />
          <el-select v-model="listQuery.status" clearable placeholder="在线状态">
            <el-option
              v-for="item in map.statusOptions"
              :key="item.value"
              :label="item.label"
              :value="+item.value"
            />
          </el-select>
          <el-button
            class="margin-left-10"
            type="primary"
            @click="fetchData"
          >搜索</el-button>
        </div>
      </template>
    </collapse-filter>
    <el-table v-loading="listLoading" class="table-container" :data="list" element-loading-text="Loading" border fit highlight-current-row>
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.row.id }}
        </template>
      </el-table-column>
      <el-table-column label="主机名称">
        <template slot-scope="scope">
          {{ scope.row.hostname }}
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
      <el-table-column label="ip" align="center">
        <template slot-scope="scope">
          {{ scope.row.ipAddress }}
        </template>
      </el-table-column>
      <!-- <el-table-column label="端口号" align="center">
        <template slot-scope="scope">
          {{ scope.row.port }}
        </template>
      </el-table-column> -->
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
              <el-dropdown-item @click.native="handleFile(row)">
                <i class="el-icon-document" /> 文件管理
              </el-dropdown-item>
              <el-dropdown-item @click.native="handleCmd(row)">
                <i class="el-icon-s-promotion" /> 执行命令
              </el-dropdown-item>
              <el-dropdown-item @click.native="handlePtyShell(row)">
                <i class="el-icon-setting" /> Shell
              </el-dropdown-item>
              <el-dropdown-item @click.native="handleUpdateClient(row)">
                <i class="el-icon-s-tools" /> 更新客户端
              </el-dropdown-item>
              <el-dropdown-item @click.native="handleDelete(row)">
                <i class="el-icon-delete" /> 删除
              </el-dropdown-item>
<!--              <el-dropdown-item @click.native="handleSshShell(row)">-->
<!--                <i class="el-icon-edit" /> SSH Shell-->
<!--              </el-dropdown-item>-->
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size" @pagination="fetchData" />
    <common-dialog
      v-if="shellDialogShow"
      :title="dialogTitle"
      :visible.sync="shellDialogShow"
      @closed="closeShellDialog()">
      <shell :id="selectedRow.id" ref="shell" :shell-type="shellType" />
    </common-dialog>

    <common-dialog
      v-if="fileDialogShow"
      :title="dialogTitle"
      :visible.sync="fileDialogShow"
      @closed="fileDialogShow = false">
      <file-manager :id="selectedRow.id" ref="file" />
    </common-dialog>

    <cmd :v-show="cmdDialogShow" :client-id="selectedRow.id" :visible="cmdDialogShow" :title="dialogTitle" @hide="cmdDialogShow = false" />
    <update-client :v-show="updateDialogShow" :client-id="selectedRow.id" :visible="updateDialogShow" :os-type="selectedRow.osType" @hide="updateDialogShow = false" />
  </div>
</template>

<script>
import { fetchList, deleteClient } from '@/api/client'
import Shell from '@/components/Shell'
import Pagination from '@/components/Pagination'
import CollapseFilter from '@/components/CollapseFilter/index.vue'
import Cmd from './components/cmd.vue'
import formMixin from '@/mixins/form-father'
import * as map from '@/map/client'
import { parseTime } from '@/utils'
import FileManager from '@/components/FileManager/index.vue'
import UpdateClient from './components/update-client.vue'
import CommonDialog from '@/components/CommonDialog/index.vue'

export default {
  name: 'Client',

  components: { CollapseFilter, Shell, Pagination, Cmd, FileManager, UpdateClient, CommonDialog },
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
      shellDialogShow: false,
      shellType: 2,
      listQuery: {
        page: 1,
        size: 10,
        hostname: '',
        status: 2
      },
      total: 0,
      cmdDialogShow: false,
      fileDialogShow: false,
      dialogTitle: '',
      updateDialogShow: false
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
    handleFile(row) {
      this.selectedRow = row
      this.setDialogTitle()
      this.fileDialogShow = true
    },
    handlePtyShell(row) {
      this.selectedRow = row
      this.setDialogTitle()
      this.shellDialogShow = true
      this.shellType = 2
    },
    handleSshShell(row) {
      this.selectedRow = row
      this.setDialogTitle()
      this.shellDialogShow = true
      this.shellType = 1
    },
    handleCmd(row) {
      this.selectedRow = row
      this.setDialogTitle()
      this.cmdDialogShow = true
    },
    setDialogTitle() {
      const d = this.selectedRow
      this.dialogTitle = `IP: ${d.ipAddress} Hostname: ${d.hostname} Username: ${d.username}`
    },
    closeShellDialog() {
      this.$refs.shell.close()
    },
    handleUpdateClient(row) {
      this.selectedRow = row
      this.updateDialogShow = true
    },
    handleDelete(row) {
      this.$confirm('是否确认删除：' + row.hostname, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
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

      });
    }
  }
}
</script>

<style scoped>
</style>

<template>
  <div>
    <el-container>
      <el-header>
        <el-button type="primary" @click="dialogVisible = true">新增</el-button>
      </el-header>
      <el-main>
        <el-table :data="list">
          <el-table-column label="模式">
            <template slot-scope="scope">
              {{ m.tunnelType[scope.row.tunnelType] }}
            </template>
          </el-table-column>
          <el-table-column label="端口">
            <template slot-scope="scope">
              {{ scope.row.serverPort }}
            </template>
          </el-table-column>
          <el-table-column label="目标IP">
            <template slot-scope="scope">
              {{ scope.row.clientIp }}
            </template>
          </el-table-column>
          <el-table-column label="目标端口">
            <template slot-scope="scope">
              {{ scope.row.clientPort }}
            </template>
          </el-table-column>
          <el-table-column label="服务端状态">
            <template slot-scope="scope">
              <span :title="scope.row.failReason">
                {{ m.status[scope.row.status] }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="50">
            <template slot-scope="{row}">
              <el-dropdown trigger="click" placement="bottom-start" size="small">
                <el-button type="text" size="medium">
                  <i class="el-icon-more" />
                </el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item icon="el-icon-close" @click.native="deleteTunnel(row)">删除</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
      </el-main>

    </el-container>
    <el-dialog title="新增隧道" :visible.sync="dialogVisible" width="30%">
      <el-form :model="form">
        <el-form-item label="模式">
          <el-select v-model="form.tunnelType" placeholder="请选择">
            <el-option
              v-for="item in m.tunnelTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="+item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="端口">
          <el-input v-model.number="form.serverPort" autocomplete="off" />
        </el-form-item>
        <el-form-item v-if="form.tunnelType === 1" label="目标IP">
          <el-input v-model="form.clientIp" autocomplete="off" />
        </el-form-item>
        <el-form-item v-if="form.tunnelType === 1" label="目标端口">
          <el-input v-model.number="form.clientPort" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="newTunnel">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { newTunnel, fetchList, deleteTunnel } from '@/api/tunnel'
import * as m from '@/map/tunnel'

export default {
  components: {
  },
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      m,
      list: [],
      dialogVisible: false,
      form: {
        tunnelType: null,
        serverPort: null,
        clientIp: null,
        clientPort: null
      }
    }
  },
  created() {
    this.fetchList()
  },
  methods: {
    fetchList() {
      fetchList(this.id)
        .then(response => {
          this.list = response.data
        })
        .catch(error => {
          console.error('Error fetching process list:', error)
        })
    },
    refresh() {
      this.fetchList()
    },
    newTunnel() {
      this.form.id = this.id
      newTunnel(this.form).then(res => {
        if (res.code === 0) {
          this.$message({
            message: '新增隧道成功',
            type: 'success'
          })
          this.dialogVisible = false
          this.fetchList()
        } else {
          this.$message({
            message: res.msg,
            type: 'error'
          })
        }
      })
    },
    deleteTunnel(row) {
      this.$confirm('是否确认删除', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteTunnel({ id: this.id, tunnelId: row.id }).then((res) => {
          if (res.code === 0) {
            this.$message.success('删除成功')
            this.fetchList()
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      })
    }
  }
}
</script>

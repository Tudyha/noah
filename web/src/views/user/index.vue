<template>
  <div class="app-container">
    <el-table v-loading="listLoading" class="table-container" :data="list" element-loading-text="Loading" fit
      highlight-current-row>
      <el-table-column label="用户名称">
        <template slot-scope="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="头像" align="center">
        <template slot-scope="scope">
          <img :src="scope.row.avatar" alt="头像" style="width: 40px; height: 40px; border-radius: 50%;">
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="最近登录时间">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ parseTime(scope.row.loginTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="100">
        <template slot-scope="{row}">
          <el-dropdown trigger="click" placement="bottom-start" size="small">
            <el-button type="text" size="medium">
              <i class="el-icon-more" />
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item @click.native="handleUpdatePassword(row)">
                <i class="el-icon-s-tools" /> 修改密码
              </el-dropdown-item>
              <!-- <el-dropdown-item @click.native="handleDelete(row)">
                <i class="el-icon-delete" /> 删除
              </el-dropdown-item> -->
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total > 0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.size"
      @pagination="fetchData" />
  </div>
</template>

<script>
import { fetchList, updatePassword } from '@/api/user'
import Pagination from '@/components/Pagination'
import formMixin from '@/mixins/form-father'
import * as map from '@/map/client'
import { parseTime } from '@/utils'

export default {
  name: 'User',

  components: { Pagination },
  mixins: [formMixin],
  data() {
    return {
      map,
      list: null,
      listLoading: true,
      listQuery: {
        page: 1,
        size: 10
      },
      total: 0
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
    handleUpdatePassword(row) {
      this.$prompt('请输入密码', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        updatePassword({ id: row.userId, password: value }).then((res) => {
          if (res.code === 0) {
            this.$message.success('修改成功')
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '取消输入'
        });
      });
    }
  }
}
</script>

<style scoped>
.hostname-link {
  color: #0e96dc;
}
</style>

<template>
  <div>
    <el-header>
      <el-row type="flex" justify="space-between" align="middle">
        <el-col :span="24">
          <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item v-for="(item, index) in breadcrumbItems" :key="index">
              <a href="#" @click.prevent="pathChange(item.path)">{{ item.filename }}</a>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </el-col>
        <el-col :span="6">
          <el-row type="flex" justify="start" align="middle">
            <el-col :span="10">
              <el-button icon="el-icon-folder-add" size="mini" @click="newDir()">新建目录</el-button>
            </el-col>
            <el-col :span="10">
              <el-upload ref="upload"
                         class="upload"
                         action
                         :data="selectItem"
                         :http-request="handleUpload"
                         :headers="{'Authorization': 'Bearer ' + getToken()}"
                         :show-file-list="false">
                <el-button icon="el-icon-upload" size="mini">上传文件</el-button>
              </el-upload>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
    </el-header>
    <el-table :data="files" height="70vh">
      <el-table-column prop="filename" label="名称">
        <template #default="{ row }">
          <a @click="fileClick(row)">
            <i :class="getIconClass(row.type)" style="margin-right: 5px;"></i> {{ row.filename }}
          </a>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小" width="100px"></el-table-column>
      <el-table-column prop="user" label="用户" width="100px"></el-table-column>
      <el-table-column prop="mod" label="权限" width="120px"></el-table-column>
      <el-table-column prop="modTime" label="修改时间" :formatter="(row, column, cellValue, index) => parseTime(cellValue)" width="160px"></el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width" width="50">
        <template slot-scope="{row}">
          <el-dropdown trigger="click" placement="bottom-start" size="small">
            <el-button type="text" size="medium">
              <i class="el-icon-more" />
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item icon="el-icon-edit" @click.native="renameItem(row)">重命名</el-dropdown-item>
              <el-dropdown-item icon="el-icon-delete" @click.native="deleteItem(row)">删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <file-preview v-show="filePreviewDialogVisible" :id="id" :file-path="selectItem.path" :visible="filePreviewDialogVisible" @hide="filePreviewDialogVisible = false" />
  </div>
</template>

<script>
import {deleteFile, fetchList, newDir, renameFile, uploadFile} from '@/api/file'
import { parseTime } from '@/utils'
import {getToken} from "@/utils/auth"
import FilePreview from "@/components/FileManager/FilePreview.vue";

export default {
  components: {FilePreview},
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      path: '/',
      files: [],
      breadcrumbItems: [],
      selectItem: {},
      filePreviewDialogVisible: false,
    };
  },
  created() {
    this.fetchFiles();
    this.freshBreadcrumb()
  },
  methods: {
    parseTime,
    getToken,
    fileClick(item) {
      if (item.type === 1) {
        this.filePreviewDialogVisible = true;
        this.selectItem = item;
      }
      if (item.type === 2) {
        this.pathChange(item.path)
      }
    },
    pathChange(path) {
      this.path = path;
      this.fetchFiles();
    },
    fetchFiles() {
      fetchList({ id: this.id, path: this.path })
        .then(response => {
          if (Array.isArray(response.data)) {
            this.files = response.data.map(file => ({
              ...file,
              icon: this.getIconClass(file.type)
            }));
          } else {
            this.files = [];
          }
        })
        .catch(error => {
          console.error('Error fetching files:', error);
          // 可以在这里显示错误提示信息
        });
    },
    getIconClass(type) {
      const iconMap = {
        2: 'el-icon-folder',
        1: 'el-icon-document',
        3: 'el-icon-tickets'
      };
      return iconMap[type] || iconMap.default;
    },
    renameItem(item) {
      this.$prompt('', '重命名', {
        inputValue: item.filename,
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        renameFile({id: this.id, path: item.path, filename: value}).then((res) => {
          if (res.code === 0) {
            this.$message.success('重命名成功')
            this.fetchFiles();
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      });
    },
    deleteItem(item) {
      this.$confirm('是否确认删除：' + item.filename, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        deleteFile({id: this.id, path: item.path, type: item.type}).then((res) => {
          if (res.code === 0) {
            if (res.data === '') {
              this.$message.success('删除成功')
            } else {
              this.$message.info(res.data)
            }
            this.fetchFiles();
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      });
    },
    freshBreadcrumb() {
      this.breadcrumbItems = []
      this.breadcrumbItems.push({filename: '/', path: '/'})
      const paths = this.path.split('/').filter((path) => path !== '')
      if (paths.length > 0) {
        for (let i = 0; i < paths.length; i++) {
          this.breadcrumbItems.push({filename: paths[i], path: `${this.breadcrumbItems[i].path}${paths[i]}/`})
        }
      }
    },
    handleUpload(data) {
      const formData = new FormData()
      formData.append('file', data.file)
      formData.append('path', this.path)
      formData.append('id', this.id)
      uploadFile(formData).then((res) => {
        if (res.code === 0) {
          this.$message.success('上传成功')
          this.fetchFiles()
        } else {
          this.$message.error(res.msg)
        }
      })
    },
    newDir() {
      this.$prompt('', '新建目录', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        newDir({id: this.id, path: this.path + '/' + value}).then((res) => {
          if (res.code === 0) {
            this.$message.success('新建目录成功')
            this.fetchFiles()
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      });
    },
  },
  watch: {
    path(newPath, oldPath) {
      this.freshBreadcrumb()
    }
  }
};
</script>

<style scoped>
</style>

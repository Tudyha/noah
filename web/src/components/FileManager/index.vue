<template>
  <div class="file-manager">
    <div class="path-bar">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item v-for="(item, index) in breadcrumbItems" :key="index">
          <a href="#" @click.prevent="handlerPathClick(item.path)">{{ item.name }}</a>
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="file-list" style="height: 500px;overflow: auto">
      <div v-for="(item, index) in items" :key="index" class="file-item" @click="fileClick(item)">
        <img v-if="item.type === 2" src="./文件夹.png" alt="" />
        <img v-else src="./文件.png" alt="" />
        <div class="file-name" @mouseover="showTooltip = true" @mouseout="showTooltip = false">
          {{ item.name }}
          <el-tooltip v-if="showTooltip" effect="dark" :content="item.name" placement="top-start">
            <div></div>
          </el-tooltip>
        </div>
        <div class="actions">
          <el-tooltip effect="dark" content="重命名" placement="top">
            <el-button class="action-button" icon="el-icon-edit" type="primary" v-if="item.type === 1 || item.type === 2" size="mini" @click.stop="renameItem(item)"></el-button>
          </el-tooltip>
          <el-tooltip effect="dark" content="上传" placement="top">
            <el-button class="action-button" icon="el-icon-upload" type="primary" v-if="item.type === 2" @click.stop="uploadFile(item)"></el-button>
          </el-tooltip>
          <el-tooltip effect="dark" content="删除" placement="top">
            <el-button class="action-button" icon="el-icon-delete" type="danger" v-if="item.type === 1 || item.type === 2" size="mini" @click.stop="deleteItem(item)"></el-button>
          </el-tooltip>
        </div>
      </div>
    </div>

      <file-preview v-show="filePreviewDialogVisible" :id="id" :file-path="selectedItem.path" :visible="filePreviewDialogVisible" @hide="filePreviewDialogVisible = false" />
  </div>
</template>

<script>
import { fetchList, renameFile, deleteFile } from '@/api/file'
import FilePreview from './FilePreview.vue'

export default {
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  components: {
    FilePreview
  },
  data() {
    return {
      currentPath: '/',
      items: [],
      folders: [],
      breadcrumbItems: [],
      showTooltip: false,
      filePreviewDialogVisible: false,
      selectedItem: {}
    };
  },
  methods: {
    fileClick(item) {
      if (item.type === 1) {
        this.selectedItem = item
        this.filePreviewDialogVisible = true
      } else if (item.type === 2) {
        this.selectItem(item)
      }
    },
    selectItem(item) {
      this.currentPath = item.path;
      this.freshBreadcrumb()
      if (item.type === 2) {
        this.loadItems();
      }
    },
    handlerPathClick(path) {
      this.currentPath = path;
      this.freshBreadcrumb()
      this.loadItems();
    },
    loadItems() {
      fetchList({id: this.id, path: this.currentPath}).then((res) => {
        this.items = res.data;
      })
    },
    freshBreadcrumb() {
      this.breadcrumbItems = []
      this.breadcrumbItems.push({name: '/', path: '/'})
      const paths = this.currentPath.split('/').filter((path) => path !== '')
      if (paths.length > 0) {
        for (let i = 0; i < paths.length; i++) {
          this.breadcrumbItems.push({name: paths[i], path: `${this.breadcrumbItems[i].path}${paths[i]}/`})
        }
      }
    },
    renameItem(item) {
      this.$prompt('', '重命名', {
        inputValue: item.name,
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        console.log(item);
        renameFile({id: this.id, path: item.path, name: value}).then((res) => {
          if (res.code === 0) {
            if (res.data === '') {
              this.$message.success('重命名成功')
            } else {
              this.$message.info(res.data)
            }
            this.loadItems();
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      });
    },
    deleteItem(item) {
      this.$confirm('是否确认删除：' + item.name, '提示', {
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
            this.loadItems();
          } else {
            this.$message.error(res.msg)
          }
        })
      }).catch(() => {

      });
    },
    uploadFile(folder) {
      // Upload file to folder logic
      console.log(`Uploading file to ${folder.name}`);
    }
  },
  mounted() {
    this.freshBreadcrumb()
    this.loadItems()
  }
};
</script>

<style scoped>
.file-manager {
  display: flex;
  flex-direction: column;
}

.file-list {
  display: flex;
  flex-wrap: wrap;
  padding: 10px;
}

.file-item {
  border: 3px solid #ccc;
  padding: 5px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  width: 120px; /* 固定宽度 */
  height: 140px; /* 固定高度 */
}

.file-item img {
  max-width: 80px;
  max-height: 80px;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 1px; /* 添加间距 */
}

.action-button {
  width: 24px; /* 固定宽度 */
  height: 24px; /* 固定高度 */
  padding: 0; /* 移除默认内边距 */
  line-height: 24px; /* 确保图标居中 */
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100px; /* 可以根据需要调整 */
}
</style>

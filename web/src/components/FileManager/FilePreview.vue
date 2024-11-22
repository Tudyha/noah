<template>
  <div>
    <el-dialog :visible.sync="dialogVisible" append-to-body @close="handleClose">

      <el-input v-model="content" type="textarea" rows="20" :disabled="textAreaDisabled">{{ content }}</el-input>

      <span slot="footer" class="dialog-footer">
        <el-button v-if="editButtonShow" type="primary" @click="handleUpdate()">编辑</el-button>
        <el-button v-if="saveButtonShow" type="primary" @click="handleSave()">保存</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { fetchFileContent, saveFile } from '@/api/file'

export default {
  name: 'FilePreview',
  props: {
    filePath: {
      type: String,
      required: true,
      default: ''
    },
    id: {
      type: Number,
      required: true
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      dialogVisible: false,
      content: '',
      textAreaDisabled: true,
      confirmText: '编辑',
      saveButtonShow: false,
      editButtonShow: true
    };
  },
  watch: {
    visible(newVal) {
      this.dialogVisible = newVal;
    },
    filePath(newVal) {
      this.getFileContent()
    }
  },
  mounted() {
    this.getFileContent()
  },
  computed: {
  },
  methods: {
    getFileContent() {
      if (this.filePath !== '') {
        fetchFileContent({ id: this.id, path: this.filePath })
          .then((res) => {
            this.content = res.data;
          })
          .catch(error => {
            console.error('Failed to fetch file content:', error);
            this.content = 'Failed to load file content.';
          });
      }
    },
    handleClose() {
      this.dialogVisible = false;
      this.$emit('hide');
    },
    handleUpdate() {
      this.textAreaDisabled = false;
      this.saveButtonShow = true;
      this.editButtonShow = false;
    },
    handleSave() {
      this.textAreaDisabled = true;
      this.saveButtonShow = false;
      this.editButtonShow = true;
      //请求接口保存数据
      saveFile({ id: this.id, path: this.filePath, content: this.content }).then((res) => {
        if (res.code === 0) {
          this.$message({
            message: "保存成功",
            type: 'success'
          });
        } else {
          this.$message({
            message: "保存失败：" + res.msg,
            type: 'error'
          });
        }
      })
    }
  }
}
</script>

<style scoped>
/* 样式代码 */
</style>

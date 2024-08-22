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
import { fetchFileContent } from '@/api/file'

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
    fileType() {
      const extension = this.filePath.split('.').pop();
      return this.getFileTypeByExtension(extension);
    },
    formattedJson() {
      try {
        const parsed = JSON.parse(this.content);
        return JSON.stringify(parsed, null, 2);
      } catch (error) {
        console.error('Failed to parse JSON:', error);
        return this.content + '\n\n[Warning: Failed to parse JSON. Please check the file content.]';
      }
    }
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
    getFileTypeByExtension(extension) {
      const fileTypes = {
        json: ['json'],
        yaml: ['yaml', 'yml'],
        txt: ['txt'],
        // 其他文件类型
        image: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg'],
        pdf: ['pdf'],
        // ... 更多文件类型
      };

      for (const [type, extensions] of Object.entries(fileTypes)) {
        if (extensions.includes(extension.toLowerCase())) {
          return type;
        }
      }

      return 'unknown';
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
    }
  }
}
</script>

<style scoped>
/* 样式代码 */
</style>

<template>
  <el-dialog
    :title="title"
    :visible.sync="isVisible"
    width="60%"
    @close="onClose"
    center
    custom-class="common-dialog"
  >
    <div class="dialog-content">
      <slot></slot>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'CommonDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      isVisible: this.visible
    };
  },
  watch: {
    visible(newVal) {
      this.isVisible = newVal;
    }
  },
  methods: {
    onClose() {
      this.$emit('update:visible', false);
      this.$emit('closed');
    }
  }
}
</script>

<style scoped>
/deep/.common-dialog .el-dialog__body {
  padding: 0;
}

/deep/.common-dialog .el-dialog__header {
  padding: 5px;
}

/deep/.common-dialog .el-dialog__headerbtn {
  top: 10px;
}

.dialog-content {
  border: none; /* 去掉内容区域的边框 */
  padding: 5px; /* 可以根据需要调整内边距 */
}
</style>

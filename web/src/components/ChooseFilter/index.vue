<template>
  <div class="choose">
    <div class="txt"><span>{{ title }}</span></div>
    <el-select v-model="newValue" :filterable="filterable" clearable placeholder="请选择" :disabled="isDisabled" style="min-width:120px" @clear="setValueNull" @focus="visibleChange">
      <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
    </el-select>
  </div>
</template>

<script>
/**
 * setNull 置空选择输入框
 * @setValueNull 置空回调
 * @visibleChange 选中回调
 */
export default {
  model: {
    prop: 'value',
    event: 'parent-event'
  },
  props: {
    title: { type: String, default: '' },
    options: { type: Array, default: null },
    setvalue: { type: String, default: null },
    setNull: { type: Boolean, default: false },
    filterable: { type: Boolean, default: false },
    isDisabled: { type: Boolean, default: false },
    value: {
      type: [Number, String, Array],
      default: ''
    }
  },
  data() {
    return {
      newValue: null
    }
  },
  watch: {
    newValue(newVlue, oldValue) {
      this.$emit('parent-event', this.newValue)
    },
    setNull() {
      this.newValue = null
    },
    setvalue: {
      handler(nv) {
        if (nv !== null) {
          this.newValue = nv
        }
      },
      immediate: true
    },
    // 监听父组件传回来的v-model的值
    value: {
      handler(nv) {
        if (nv !== null) {
          this.newValue = nv
        }
      },
      immediate: true
    }
  },
  methods: {
    setValueNull() {
      this.newValue = null
      this.$emit('setValueNull')
    },
    visibleChange(value) {
      this.$emit('visibleChange')
    }
  }
}
</script>
<style lang="scss" scoped>
.choose {
  display: flex;
  flex-direction: row;
  margin: 0 10px;
  .txt {
    margin-right: 10px;
    font-size: 12px;
    word-break: keep-all;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}
</style>

<template>
  <!-- 折叠搜索表单容器 -->
  <div class="collapse-filter">
    <el-collapse v-model="activeName" accordion v-on="$listeners">
      <el-collapse-item name="filter">
        <template v-slot:title>
          <!-- 折叠标题左边部分 -->
          <div style="flex: 1; font-size: 0" @click.stop="() => {}" @keyup.stop="() => {}">
            <slot name="collapse-title-left" />
          </div>
          <!-- 折叠标题右边部分 -->
          <slot name="collapse-title-right" />
          <span style="padding: 0 5px;">{{ activeName ? '收起常用筛选' : '展开常用筛选' }}</span>
        </template>
        <!-- 折叠内容 -->
        <slot />
      </el-collapse-item>
    </el-collapse>
  </div>
</template>
<script>
export default {
  props: {
    default: { type: String, default: '' }
  },
  data() {
    return {
      activeName: this.default || ''
    }
  }
}
</script>
<style lang="scss">
.collapse-filter {
  overflow-y: scroll;
  span{
    white-space: nowrap;
  }
  &.no-margin {
    .el-collapse-item__header .el-form-item {
      margin-bottom: 0;
    }
  }
  .el-collapse-item:last-child {
    margin-bottom: 0;
  }
  .el-collapse-item > div:first-child {
    margin-bottom: 10px;
  }
  .el-collapse-item__wrap {
    border-bottom: 0;
  }
  .el-collapse-item__header {
    border-bottom: 0
  }
  .el-collapse {
    border: 0
  }
  .el-collapse-item__content {
    padding-bottom: 0px;
    margin-top: 0px;
  }
}
</style>

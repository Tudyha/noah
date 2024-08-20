<template>
  <!-- 表格容器 -->
  <div class="table-box-container">
    <!-- 容器顶部部分，一般是搜索表单 -->
    <div class="table-box__top">
      <slot name="table-box-top" />
    </div>
    <!-- 表格部分 -->
    <div class="table-box__table">
      <slot />
    </div>
    <!-- 容器底部部分 -->
    <div class="table-box__bottom">
      <slot name="table-box-bottom" />
    </div>
  </div>
</template>
<style lang="scss" scoped>
.table-box-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .table-box__table {
    flex: 1;
    overflow: hidden;
  }
}
</style>

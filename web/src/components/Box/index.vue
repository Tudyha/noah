<template>
  <!-- 统一容器 -->
  <div class="app-container" :class="classObj">
    <div class="mod box-shadow bg-white" style="padding: 0">
      <el-card shadow="never" style="border: 0" v-bind="$attrs">
        <template v-if="title" v-slot:header>
          <div style="line-height: 24px">
            <!-- 容器顶部标题 -->
            <span>{{ title }}</span>
            <!-- 容器顶部右边部分 -->
            <div style="float: right; font-size: 0px;">
              <slot name="header-btns" />
            </div>
          </div>
        </template>
        <slot />
      </el-card>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    title: { type: String, default: '' }, // 标题,
    customClass: { type: String, default: '' }, // 自定义class
    isFlexLayout: { type: Boolean, default: true } // flex 布局
  },
  computed: {
    classObj() {
      return {
        is_flex_layout: this.isFlexLayout,
        [this.customClass]: true
      }
    }
  }
}
</script>
<style lang="scss">
.app-container.is_flex_layout {
  height: 100%;
  .mod {
    height: 100%;
  }
  .el-card {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  .el-card__header {
    padding: 12px 20px
  }
  .el-card__body {
    padding: 10px 20px;
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}
</style>

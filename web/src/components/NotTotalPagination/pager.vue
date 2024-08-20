<template>
  <ul class="el-pager">
    <li
      v-for="pager in pagers"
      :key="pager"
      :class="{ active: currentPage === pager, disabled }"
      class="number"
    >{{ pager }}</li>
  </ul>
</template>

<script type="text/babel">
export default {
  name: 'NotTotalPager',

  props: {
    currentPage: Number,

    pageCount: Number,

    pagerCount: Number,

    disabled: Boolean
  },

  computed: {
    pagers() {
      const currentPage = Number(this.currentPage)
      return [currentPage]
    }
  }
}
</script>

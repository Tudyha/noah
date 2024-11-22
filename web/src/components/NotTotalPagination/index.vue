<template>
  <div :class="['el-pagination', 'pagination-container', { 'is-background': background, 'el-pagination--small': small }]">
    <span class="el-pagination__sizes">
      <el-select
        v-model="internalPageSize"
        :popper-class="popperClass || ''"
        size="mini"
        @change="handleSizeChange"
      >
        <el-option
          v-for="item in pageSizes"
          :key="item"
          :value="item"
          :label="item + '条/页'"
        />
      </el-select>
    </span>
    <button
      type="button"
      class="btn-prev"
      :disabled="internalCurrentPage <= 1"
      @click="prev"
    >
      <i class="el-icon el-icon-arrow-left" />
    </button>
    <pager
      :current-page="internalCurrentPage"
      :page-count="internalPageCount"
      :pager-count="pagerCount"
    />
    <button
      type="button"
      class="btn-next"
      :disabled="total < internalPageSize"
      @click="next"
    >
      <i class="el-icon el-icon-arrow-right" />
    </button>
    <span class="el-pagination__jump">
      前往
      <el-input
        class="el-pagination__editor is-in-pagination"
        :min="1"
        :max="total < internalPageSize ? internalCurrentPage : Number.MAX_SAFE_INTEGER"
        :value="userInput !== null ? userInput : internalCurrentPage"
        type="number"
        @keyup.native="handleKeyup"
        @input="handleInput"
        @change="handleChange"
      />
      页
    </span>
  </div>
</template>

<script>
import Pager from './pager.vue'

const valueEquals = (a, b) => {
  // see: https://stackoverflow.com/questions/3115982/how-to-check-if-two-arrays-are-equal-with-javascript
  if (a === b) return true
  if (!(a instanceof Array)) return false
  if (!(b instanceof Array)) return false
  if (a.length !== b.length) return false
  for (let i = 0; i !== a.length; ++i) {
    if (a[i] !== b[i]) return false
  }
  return true
}

export default {
  name: 'NotTotalPagination',

  components: {
    Pager
  },

  props: {
    pageSize: {
      type: Number,
      default: 100
    },

    small: Boolean,

    total: Number,

    pageCount: Number,

    pagerCount: {
      type: Number,
      validator(value) {
        return (value | 0) === value && value > 4 && value < 22 && (value % 2) === 1
      },
      default: 7
    },

    page: {
      type: Number,
      default: 1
    },

    layout: {
      default: 'prev, pager, next, jumper, ->, total'
    },

    pageSizes: {
      type: Array,
      default() {
        return [100, 200, 500, 1000, 2000]
      }
    },

    popperClass: String,

    prevText: String,

    nextText: String,

    background: { type: Boolean, default: true },

    disabled: Boolean,

    hideOnSinglePage: Boolean
  },

  data() {
    return {
      internalCurrentPage: 1,
      internalPageSize: 0,
      lastEmittedPage: -1,
      userChangePageSize: false,
      userInput: null
    }
  },

  computed: {
    components() {
      return this.layout.split(',').map((item) => item.trim())
    },
    internalPageCount() {
      if (typeof this.total === 'number') {
        return Math.max(1, Math.ceil(this.total / this.internalPageSize))
      } else if (typeof this.pageCount === 'number') {
        return Math.max(1, this.pageCount)
      }
      return null
    }
  },

  watch: {
    page: {
      immediate: true,
      handler(val) {
        this.internalCurrentPage = this.getValidCurrentPage(val)
      }
    },

    pageSize: {
      immediate: true,
      handler(val) {
        this.internalPageSize = isNaN(val) ? 100 : val
      }
    },

    internalCurrentPage: {
      immediate: true,
      handler(newVal) {
        this.$emit('update:page', newVal)
        this.lastEmittedPage = -1
      }
    },

    internalPageCount(newVal) {
      /* istanbul ignore if */
      this.userInput = null
      const oldPage = this.internalCurrentPage
      if (newVal > 0 && oldPage === 0) {
        this.internalCurrentPage = 1
      } else if (oldPage > newVal) {
        this.internalCurrentPage = newVal === 0 ? 1 : newVal
        this.userChangePageSize && this.emitChange()
      }
      this.userChangePageSize = false
    },

    pageSizes: {
      immediate: true,
      handler(newVal, oldVal) {
        if (valueEquals(newVal, oldVal)) return
        if (Array.isArray(newVal)) {
          this.internalPageSize = newVal.indexOf(this.pageSize) > -1
            ? this.pageSize
            : this.pageSizes[0]
        }
      }
    }
  },

  methods: {
    handleCurrentChange(val) {
      this.internalCurrentPage = this.getValidCurrentPage(val)
      this.userChangePageSize = true
      this.emitChange()
    },

    prev() {
      const newVal = this.internalCurrentPage - 1
      this.internalCurrentPage = this.getValidCurrentPage(newVal)
      this.$emit('prev-click', this.internalCurrentPage)
      this.emitChange()
    },

    next() {
      const newVal = this.internalCurrentPage + 1
      this.internalCurrentPage = this.getValidCurrentPage(newVal)
      this.$emit('next-click', this.internalCurrentPage)
      this.emitChange()
    },

    getValidCurrentPage(value) {
      value = parseInt(value, 10)

      const havePageCount = typeof this.internalPageCount === 'number'

      let resetValue
      if (!havePageCount) {
        if (isNaN(value) || value < 1) resetValue = 1
      } else {
        if (value < 1) {
          resetValue = 1
        }
      }

      if (resetValue === undefined && isNaN(value)) {
        resetValue = 1
      } else if (resetValue === 0) {
        resetValue = 1
      }

      return resetValue === undefined ? value : resetValue
    },

    emitChange() {
      this.$nextTick(() => {
        if (this.internalCurrentPage !== this.lastEmittedPage || this.userChangePageSize) {
          this.$emit('pagination', { page: this.internalCurrentPage, limit: this.internalPageSize })
          this.lastEmittedPage = this.internalCurrentPage
          this.userChangePageSize = false
        }
      })
    },
    handleKeyup({ keyCode, target }) {
      // Chrome, Safari, Firefox triggers change event on Enter
      // Hack for IE: https://github.com/ElemeFE/element/issues/11710
      // Drop this method when we no longer supports IE
      if (keyCode === 13) {
        this.handleChange(target.value)
      }
    },
    handleInput(value) {
      this.userInput = value
    },
    handleSizeChange() {
      this.internalCurrentPage = 1
      this.emitChange()
    },
    handleChange(value) {
      this.internalCurrentPage = this.getValidCurrentPage(value)
      this.emitChange()
      this.userInput = null
    }
  }
}

</script>

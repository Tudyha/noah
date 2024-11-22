import CommonTable from '@/components/CommonTable'
import moduleMixin from './module'

export default {
  components: { CommonTable },
  mixins: [moduleMixin],
  data() {
    return {
      selected: [],
      searchLabel: '搜索',
      resetLabel: '重置',
      commonFilter: {},
      rowKey: '',
      dateStart: '00:00:00',
      dateEnd: '23:59:59',
      rowHeight: 50
    }
  },
  props: {
    parentFilter: { type: Object, default: null }
  },
  computed: {
    tableScrollTop() {
      const table = (this.cached && this.cached.table) || {}
      return table.scrollTop || 0
    },
    tableScrollLeft() {
      const table = (this.cached && this.cached.table) || {}
      return table.scrollLeft || 0
    },
    fields() {
      return [
        ...(this.filterFields || []),
        { type: 'btns', key: 'btns', btns: [
          { key: 'search', type: 'primary', label: this.searchLabel, click: this.handleSearch },
          { key: 'reaset', label: this.resetLabel, click: this.handleReset }
        ] }
      ]
    },
    commonFields() {
      return [
        ...(this.commonFilterFields || []),
        { type: 'btns', key: 'btns', btns: [
          { key: 'search', type: 'primary', label: '搜索', click: this.handleSearch }
        ] }
      ]
    },
    filter() {
      return this.customMapState('filter') || { page: 1, size: 100 }
    },
    list() {
      return this.customMapState('list')
    },
    total() {
      return this.customMapState('total') || 0
    },
    loading() {
      return this.customMapState('loading') || false
    },
    cached() {
      return this.customMapState('cached') || {}
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    getSelectedIds() {
      return (this.selected || []).map(item => item[this.rowKey])
    },
    init() {
      this.handleSearch()
    },
    _beforeBeforeDesorty() {
      this.resetListState()
    },
    _setTableScroll(data) {
      this._commit('setCache', { table: { ...(this.cached && this.cached.table) || {}, ...data }})
    },
    fetchList(data) {
      if (this.moduleName) {
        return this._dispatch('fetchList', data)
      }
      throw new Error('请重写fetchList函数或指定moduleName')
    },
    resetListState() {
      if (this.moduleName) {
        return this._commit('resetListState')
      }
      throw new Error('请重写resetListState函数或指定moduleName')
    },
    _formateFilter(filter) {
      return filter
    },
    _handleSearch() {
      const filter = this.$refs.filter ? this.$refs.filter.getValue() : {}
      if (this.searchFilterValidate && !this.searchFilterValidate(filter)) {
        // 如果有需要校验filter的hook，而且没有校验通过，则不往下继续执行
        return
      }
      const allFilter = {
        ...filter,
        ...this.commonFilter,
        ...(this.moreFilter || {}),
        page: 1,
        size: this.filter.size
      }
      if (this._getPage) {
        allFilter._getPage = this._getPage
      }

      return this.fetchList(this._formateFilter(allFilter))
    },
    handleSearch() {
      return this._handleSearch()
    },
    _handleReset() {
      if (this.$refs.filter) {
        this.$refs.filter.resetFields()
      }
      if (this.$refs.commonFilter) {
        this.$refs.commonFilter.resetFields()
      }
      // reset 额外操作
      this.handleExtraReset && this.handleExtraReset()
      this.handleSearch()
    },
    handleReset() {
      this._handleReset()
    },
    handlePagination({ page, limit: size }) {
      const filter = {
        ...this.filter,
        ...(this.moreFilter || {}),
        page: this.filter.size !== size ? 1 : page, size
      }
      if (this._getPage) {
        filter._getPage = this._getPage
      }
      this.fetchList(filter)
    },
    handleSelectionChange(selected) {
      this.selected = selected
    },
    async handleConfirmBeforeClose({ action, instance, api, data, cb }) {
      if (instance.confirmButtonLoading) {
        return
      }
      if (action === 'confirm') {
        instance.confirmButtonLoading = true
        const res = await api()
        instance.confirmButtonLoading = false
        if (res && res.code === 0) {
          // 只有confirm action 会res参数
          cb(res)
        }
      } else {
        cb()
      }
    },
    _handleBatchDelete(data) {
      if (this.moduleName) {
        return this._dispatch('batchDelete', data)
      }
      throw new Error('请重写_handleBatchDelete函数或指定moduleName')
    },
    _handleDelete(data) {
      if (this.moduleName) {
        return this._dispatch('delete', data)
      }
      throw new Error('请重写_handleDelete函数或指定moduleName')
    },
    isBatchDelete(item) {
      return item instanceof MouseEvent || item == null
    },
    _validateSelected() {
      return true
    },
    handleBatchDelete(item) {
      const { moduleChineseName = '记录', primaryKey = 'id' } = this
      const num = this.selected.length
      const isBatchDelete = this.isBatchDelete(item)
      if (isBatchDelete && num === 0) {
        return this.$message(`未选择${moduleChineseName}`)
      }
      if (!this._validateSelected()) {
        return
      }
      this.$confirm(`您确认要${isBatchDelete ? '批量' : ''}删除所选中的${moduleChineseName}吗？`, '提示', {
        type: 'warning',
        beforeClose: async(action, instance, done) => {
          const data = { ids: isBatchDelete ? this.selected.map(s => s[primaryKey]) : [item[primaryKey]] }
          this.handleConfirmBeforeClose({
            action,
            instance,
            api: () => this._handleBatchDelete(data),
            cb: (res) => {
              done()
              if (res && res.code === 0) {
                this.$message({ type: 'success', message: '删除成功' })
              }
            }
          })
        }
      }).catch(() => {})
    },
    handleDelete(item) {
      const { moduleChineseName = '记录', primaryKey = 'id' } = this
      this.$confirm(`您确认要删除所选中的${moduleChineseName}吗？`, '提示', {
        type: 'warning',
        beforeClose: async(action, instance, done) => {
          this.handleConfirmBeforeClose({
            action,
            instance,
            api: () => this._handleDelete(item[primaryKey]),
            cb: (res) => {
              done()
              if (res && res.code === 0) {
                this.$message({ type: 'success', message: '删除成功' })
              }
            }
          })
        }
      }).catch(() => {})
    },
    $cf(title, msg = '提示', config = {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }) {
      return new Promise((resolve, reject) => {
        this.$confirm(title, msg, config).then(() => {
          resolve(true)
        }).catch(() => {})
      })
    },
    handleTableScroll(e) {
      if (e && e.target) {
        const { scrollTop, scrollLeft } = e.target
        this._setTableScroll({ scrollTop, scrollLeft })
      }
    }
  }
}

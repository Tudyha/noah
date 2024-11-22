import { mapState } from 'vuex'
import modules, { moduleFactory } from '@/store/base_modules'

export default {
  props: {
    parentContainerName: { type: String, default: '' },
    unstableNeedParentContainerName: { type: Boolean, default: false } // todo: 兼容详情返回的是非相同管理模块的页面
  },
  data() {
    return {
      pageActivated: true
    }
  },
  computed: {
    _moduleName() {
      return `${this.moduleName}${this.extraModuleName ? '-' : ''}${this.extraModuleName}`
    },
    parentModuleName() {
      return this.unstableNeedParentContainerName
        ? `${this.moduleName}${this.parentContainerName ? '-' : ''}${this.parentContainerName || ''}`
        : this.moduleName
    },
    extraModuleName() {
      return this.parentContainerName
    }
  },
  created() {
    this.registerModule(this._moduleName)
  },
  beforeDestroy() {
    this._beforeBeforeDesorty && this._beforeBeforeDesorty()
    this.unregisterModule(this._moduleName)
  },
  activated() {
    this.pageActivated = true
  },
  deactivated() {
    this.pageActivated = false
  },
  methods: {
    registerModule(name) {
      if (!this.$store.state[name]) {
        const moduleApi = this.moduleApi || (modules[this.moduleName] && modules[this.moduleName]._api)
        moduleApi && this.$store.registerModule(name, moduleFactory(moduleApi))
      } else {
        this._commit('setState', { _count: this.$store.state[name]._count + 1 })
      }
    },
    unregisterModule(name) {
      if (this.$store.state[name]) {
        if (this.$store.state[name]._count <= 1) {
          this.$store.unregisterModule(name)
        } else {
          this._commit('setState', { _count: this.$store.state[name]._count - 1 })
        }
      }
    },
    customMapState(field) {
      if (!this.$store.state[this._moduleName]) {
        return null
      }
      const map = mapState(this._moduleName, {
        [field]: state => state[field]
      })
      return map[field].call(this)
    },
    async _dispatch(actionName, data, ...args) {
      const isUpdateData = ['create', 'update', 'batchUpdate'].includes(actionName)
      const dontReload = data && !!data.dontReload
      if (isUpdateData && this.__isDetailPage__) {
        data.dontReload = true
      }
      const res = await this.$store.dispatch(`${this._moduleName}/${actionName}`, data, ...args)
      if (res && res.code === 0) {
        isUpdateData && data.dontReload && !dontReload && this.$store.dispatch(`${this.parentModuleName}/fetchList`)
      }
      return res
    },
    _commit(actionName, data, ...args) {
      return this.$store.commit(`${this._moduleName}/${actionName}`, data, ...args)
    }
  }
}

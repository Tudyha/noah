import MutationTypes from './mutation_types'

const mutations = {
  [MutationTypes.resetListState](state) {
    state.list = null
    state.total = 0
    state.loading = false
    state.filter = { page: 1, size: 100 }
    state.cached = {}
  },
  [MutationTypes.updateFilter](state, payload) {
    state.filter = { ...state.filter, ...payload }
  },
  [MutationTypes.setFilter](state, filter) {
    state.filter = filter
  },
  [MutationTypes.setList](state, { data, _getPage }) {
    const page = _getPage ? _getPage(data) : data
    const isArray = page instanceof Array
    state.list = Object.freeze((isArray ? page : ((page && page.list) || [])).map((item, i) => ({
      ...item,
      _number: (state.filter.page - 1) * state.filter.size + i + 1
    })))
    state.total = isArray ? page.length : ((page && page.total) || 0)
  },
  [MutationTypes.updateListItem](state, { id, data, idKey = 'id' }) {
    state.list = state.list.map((item) => {
      if (item[idKey] === id) {
        return { ...item, ...data }
      }
      return item
    })
  },
  [MutationTypes.updateLoading](state, loading) {
    state.loading = loading
  },
  [MutationTypes.updateDetail](state, detail) {
    state.detail = detail
  },
  [MutationTypes.updateDetailLoading](state, detailLoading) {
    state.detailLoading = detailLoading
  },
  [MutationTypes.setState](state, data) {
    if (data) {
      Object.keys(data).forEach((key) => {
        state[key] = data[key]
      })
    }
  },
  [MutationTypes.setCache](state, data) {
    state.cached = { ...state.cached, ...data }
  },
  [MutationTypes.resetCache](state, data) {
    const { type, keys } = data || {}
    if (type === 'list') {
      state.list = null
      state.total = 0
      state.loading = false
      state.filter = { page: 1, size: 100 }
    } else if (type === 'detail') {
      state.detail = null
      state.detailLoading = false
    }
    (keys || []).forEach((key) => {
      state.cached[key] = undefined
    })
  }
}

export default mutations

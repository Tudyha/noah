const state = {
  list: null,
  total: 0,
  loading: false,
  filter: {
    page: 1,
    size: 100
  },
  detail: null,
  detailLoading: false,
  hasUpdate: false, // 用于是否使用缓存中的数据
  cached: {},
  _count: 1, // module被引用次数
  _detailCount: 0 // 详情数据被被引用次数
}

export default state

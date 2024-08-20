import ActionTypes from './action_types'
import MutationTypes from './mutation_types'

/**
 * 生成通用action
 * @param {Object} param
 * @param {Function} param.fetchList
 * @param {Function} param.create
 * @param {Function} param.update
 * @param {Function} param.batchUpdate
 * @param {Function} param.deleteItem
 * @param {Function} param.batchDelete
 * @param {Function} param.fetchDetail
 * @returns {Object} 根据参数返回对应的action
 */
function baseActionsGenerator({
  fetchList, create, update, batchUpdate, deleteItem, batchDelete, fetchDetail
}) {
  const base = {}
  fetchList && (
    base[ActionTypes.fetchList] = async({ commit, state }, filter) => {
      const { _getPage } = filter || {}
      if (filter) {
        delete filter._getPage
      }
      const { page } = filter || state.filter
      const params = { ...(filter || state.filter), page }
      commit(MutationTypes.updateLoading, true)
      const res = await fetchList(params).catch(err => ({ err }))
      commit(MutationTypes.updateLoading, false)
      if (res && res.code === 0) {
        filter && commit(MutationTypes.setFilter, filter)
        commit(MutationTypes.setList, { data: res.data, _getPage })
      }
      return res
    }
  )
  create && (
    base[ActionTypes.create] = async({ dispatch, commit }, { dontReload, ...data }) => {
      const res = await create(data).catch(err => ({ err }))
      if (res && res.code === 0 && !dontReload) {
        commit(MutationTypes.setState, { hasUpdate: true })
        dispatch(ActionTypes.fetchList)
      }
      return res
    }
  )
  update && (
    base[ActionTypes.update] = async({ dispatch, commit }, { dontReload, ...data }) => {
      const res = await update(data).catch(err => ({ err }))
      if (res && res.code === 0) {
        commit(MutationTypes.setState, { hasUpdate: true })
        !dontReload && dispatch(ActionTypes.fetchList)
      }
      return res
    }
  )
  batchUpdate && (
    base[ActionTypes.batchUpdate] = async({ dispatch, commit }, { dontReload, ...data }) => {
      const res = await batchUpdate(data).catch(err => ({ err }))
      commit(MutationTypes.setState, { hasUpdate: true })
      if (res && res.code === 0 && !dontReload) {
        await dispatch(ActionTypes.fetchList)
      }
      return res
    }
  )
  deleteItem && (
    base[ActionTypes.delete] = async({ dispatch, commit }, id) => {
      const res = await deleteItem(id).catch(err => ({ err }))
      if (res && res.code === 0) {
        commit(MutationTypes.updateFilter, { page: 1 })
        dispatch(ActionTypes.fetchList)
      }
      return res
    }
  )
  batchDelete && (
    base[ActionTypes.batchDelete] = async({ dispatch, commit }, data) => {
      const res = await batchDelete(data).catch(err => ({ err }))
      if (res && res.code === 0) {
        commit(MutationTypes.updateFilter, { page: 1 })
        dispatch(ActionTypes.fetchList)
      }
      return res
    }
  )
  fetchDetail && (
    base[ActionTypes.fetchDetail] = async({ commit }, id) => {
      commit(MutationTypes.updateDetailLoading, true)
      const res = await fetchDetail(id).catch(err => ({ err }))
      commit(MutationTypes.updateDetailLoading, false)
      if (res && res.code === 0) {
        commit(MutationTypes.updateDetail, res.data)
      }
      return res
    }
  )
  return base
}

export default baseActionsGenerator

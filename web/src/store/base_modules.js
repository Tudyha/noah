import baseState from './base_state'
import baseMutations from './base_mutations'
import baseActionsGenerator from './base_actions'

const autoGenModuleNames = [
  { name: 'permission-common-role', file: 'permission-role' },
  { name: 'permission-system-role', file: 'permission-role' }, 'sys-user', 'login-log',
  'client',
  { name: 'records', file: 'transcation-records' }
]

/**
 * Vuex Module生成工厂方法
 *
 * @param {Object} api
 * @returns
 */
export function moduleFactory(api) {
  const state = { ...baseState }
  const mutations = { ...baseMutations }
  const actions = { ...(baseActionsGenerator(api)) }
  return {
    _api: api,
    namespaced: true,
    state,
    mutations,
    actions
  }
}

/**
 * 转成驼峰命名
 *
 * @param {string} s
 */
function toCamelCase(s, parents) {
  return [...parents, ...s.split('/')].reduce((res, item) => {
    return [...res, ...item.split('-')]
  }, []).map((item, i) => (
    i === 0 ? item : `${item[0].toUpperCase()}${item.slice(1).toLowerCase()}`
  )).join('')
}

/**
 * 获取module名字和api文件路径
 *
 * @param {String|Object} module
 * @param {string[]} parents 父目录文件名数组
 * @param {Object} pathMap api文件路径映射
 * @returns
 */
function formateMoudle(module, parents, pathMap) {
  const isString = typeof module === 'string'
  const name = toCamelCase(isString ? module : (module.name || module.file), parents)
  const paths = [...parents, isString ? module : module.file]
  // 如果是子文件，则代表是目录，使用子文件index(可能不存在)
  if (module.children) {
    paths.push('index')
  }
  const path = pathMap[paths.join('/')]
  return { name, path }
}

// 遍历需要生成vue module的文件名数组
function reduceModules(modulesFiles, namePathMap, names, parents = []) {
  return names.reduce((modules, name) => {
    let childrenModules = []
    // 是否有子文件
    if (name.children && Array.isArray(name.children)) {
      childrenModules = reduceModules(modulesFiles, namePathMap, name.children, [...parents, name.file])
    }
    // 获取module名字和文件路径
    const module = formateMoudle(name, parents, namePathMap)
    if (module.path) {
      const api = modulesFiles(module.path)
      return { ...modules, ...childrenModules, [module.name]: moduleFactory(api) }
    }
    return { ...modules, ...childrenModules }
  }, {})
}

/**
 * 根据api文件生成Vuex module
 *
 * @param {Array} names 文件数组
 * @returns
 */
function genModules(names) {
  // 获取api目录下获取js后缀文件
  const modulesFiles = require.context('../api', true, /\.js$/)
  // api文件路径映射map
  const namePathMap = modulesFiles.keys().reduce((namePathMap, modulePath) => {
    const moduleName = modulePath.replace(/^\.\/(.*)\.\w+$/, '$1')
    namePathMap[moduleName] = modulePath
    return namePathMap
  }, {})
  return reduceModules(modulesFiles, namePathMap, names)
}

const modules = genModules(autoGenModuleNames)

export default modules

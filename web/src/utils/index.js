/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * Parse the time to string
 * @param {(Object|string|number)} time
 * @param {string} cFormat
 * @returns {string}
 */
export function parseTime(time, cFormat) {
  if (arguments.length === 0) {
    return null
  }
  const format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if ((typeof time === 'string') && (/^[0-9]+$/.test(time))) {
      time = parseInt(time)
    }
    if ((typeof time === 'number') && (time.toString().length === 10)) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }
  const time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
    let value = formatObj[key]
    // Note: getDay() returns 0 on Sunday
    if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value] }
    if (result.length > 0 && value < 10) {
      value = '0' + value
    }
    return value || 0
  })
  return time_str
}
export function formatterDate(cellValue) {
  // 秒转时分秒格式
  let _time = Number(cellValue)
  if (!_time) return '00:00:00'
  let secondTime = 0// 秒
  let minuteTime = 0// 分
  let hourTime = 0// 小时

  // 获取小时
  hourTime = parseInt(_time / 3600)
  if (hourTime < 10) {
    hourTime = '0' + hourTime
  }
  // console.log('小时：', hourTime, _time)
  _time = _time - (+hourTime) * 60 * 60
  // 获取分钟
  minuteTime = parseInt(_time / 60)
  if (minuteTime < 10) {
    minuteTime = '0' + minuteTime
  }
  // console.log('分钟：', minuteTime, _time)
  // 获取秒
  _time = _time - (+minuteTime) * 60
  secondTime = _time
  if (secondTime < 10) {
    secondTime = '0' + secondTime
  }
  // console.log('秒：', secondTime, _time)
  return hourTime + ':' + minuteTime + ':' + secondTime
}

/**
 * @param {number} time
 * @param {string} option
 * @returns {string}
 */
export function formatTime(time, option) {
  if (('' + time).length === 10) {
    time = parseInt(time) * 1000
  } else {
    time = +time
  }
  const d = new Date(time)
  const now = Date.now()

  const diff = (now - d) / 1000

  if (diff < 30) {
    return '刚刚'
  } else if (diff < 3600) {
    // less 1 hour
    return Math.ceil(diff / 60) + '分钟前'
  } else if (diff < 3600 * 24) {
    return Math.ceil(diff / 3600) + '小时前'
  } else if (diff < 3600 * 24 * 2) {
    return '1天前'
  }
  if (option) {
    return parseTime(time, option)
  } else {
    return (
      d.getMonth() +
      1 +
      '月' +
      d.getDate() +
      '日' +
      d.getHours() +
      '时' +
      d.getMinutes() +
      '分'
    )
  }
}

/**
 * 根据指定日期和指定天数，获取开始和结束日期
 *
 * @export
 * @param {Date} [start=new Date()]
 * @param {Object} [options={ days: 7 }]
 * @returns
 */
export function getStartAndEndDate(start, options = { days: 7 }) {
  const { days = 7, format } = options
  const endDate = start || new Date()
  endDate.setHours(23)
  endDate.setMinutes(59)
  endDate.setSeconds(59)
  const startDate = new Date(endDate.getTime() - days * 24 * 60 * 60 * 1000 + 1000)
  return [parseTime(startDate, format), parseTime(endDate, format)]
}

/**
 * @param {string} url
 * @returns {Object}
 */
export function getQueryObject(url) {
  url = url == null ? window.location.href : url
  const search = url.substring(url.lastIndexOf('?') + 1)
  const obj = {}
  const reg = /([^?&=]+)=([^?&=]*)/g
  search.replace(reg, (rs, $1, $2) => {
    const name = decodeURIComponent($1)
    let val = decodeURIComponent($2)
    val = String(val)
    obj[name] = val
    return rs
  })
  return obj
}

/**
 * @param {string} input value
 * @returns {number} output value
 */
export function byteLength(str) {
  // returns the byte length of an utf8 string
  let s = str.length
  for (var i = str.length - 1; i >= 0; i--) {
    const code = str.charCodeAt(i)
    if (code > 0x7f && code <= 0x7ff) s++
    else if (code > 0x7ff && code <= 0xffff) s += 2
    if (code >= 0xDC00 && code <= 0xDFFF) i--
  }
  return s
}

/**
 * @param {Array} actual
 * @returns {Array}
 */
export function cleanArray(actual) {
  const newArray = []
  for (let i = 0; i < actual.length; i++) {
    if (actual[i]) {
      newArray.push(actual[i])
    }
  }
  return newArray
}

/**
 * @param {Object} json
 * @returns {Array}
 */
export function param(json) {
  if (!json) return ''
  return cleanArray(
    Object.keys(json).map(key => {
      if (json[key] === undefined) return ''
      return encodeURIComponent(key) + '=' + encodeURIComponent(json[key])
    })
  ).join('&')
}

/**
 * @param {string} url
 * @returns {Object}
 */
export function param2Obj(url) {
  const search = url.split('?')[1]
  if (!search) {
    return {}
  }
  return JSON.parse(
    '{"' +
    decodeURIComponent(search)
      .replace(/"/g, '\\"')
      .replace(/&/g, '","')
      .replace(/=/g, '":"')
      .replace(/\+/g, ' ') +
    '"}'
  )
}

/**
 * @param {string} val
 * @returns {string}
 */
export function html2Text(val) {
  const div = document.createElement('div')
  div.innerHTML = val
  return div.textContent || div.innerText
}

/**
 * Merges two objects, giving the last one precedence
 * @param {Object} target
 * @param {(Object|Array)} source
 * @returns {Object}
 */
export function objectMerge(target, source) {
  if (typeof target !== 'object') {
    target = {}
  }
  if (Array.isArray(source)) {
    return source.slice()
  }
  Object.keys(source).forEach(property => {
    const sourceProperty = source[property]
    if (typeof sourceProperty === 'object') {
      target[property] = objectMerge(target[property], sourceProperty)
    } else {
      target[property] = sourceProperty
    }
  })
  return target
}

/**
 * @param {HTMLElement} element
 * @param {string} className
 */
export function toggleClass(element, className) {
  if (!element || !className) {
    return
  }
  let classString = element.className
  const nameIndex = classString.indexOf(className)
  if (nameIndex === -1) {
    classString += '' + className
  } else {
    classString =
      classString.substr(0, nameIndex) +
      classString.substr(nameIndex + className.length)
  }
  element.className = classString
}

/**
 * @param {string} type
 * @returns {Date}
 */
export function getTime(type) {
  if (type === 'start') {
    return new Date().getTime() - 3600 * 1000 * 24 * 90
  } else {
    return new Date(new Date().toDateString())
  }
}

/**
 * @param {Function} func
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*}
 */
export function debounce(func, wait, immediate) {
  let timeout, args, context, timestamp, result

  const later = function() {
    // 据上一次触发时间间隔
    const last = +new Date() - timestamp

    // 上次被包装函数被调用时间间隔 last 小于设定时间间隔 wait
    if (last < wait && last > 0) {
      timeout = setTimeout(later, wait - last)
    } else {
      timeout = null
      // 如果设定为immediate===true，因为开始边界已经调用过了此处无需调用
      if (!immediate) {
        result = func.apply(context, args)
        if (!timeout) context = args = null
      }
    }
  }

  return function() {
    args = [].slice.call(arguments)
    context = this
    timestamp = +new Date()
    const callNow = immediate && !timeout
    // 如果延时不存在，重新设定延时
    if (!timeout) timeout = setTimeout(later, wait)
    if (callNow) {
      result = func.apply(context, args)
      context = args = null
    }

    return result
  }
}

/**
 * 节流函数
 *
 * @export
 * @param {Function} fn
 * @param {Number} delay 延迟多少毫秒执行
 * @param {Number} mustRunDelay 延迟多少毫秒必须执行
 * @returns {Function}
 */
export function throttle(fn, delay, mustRunDelay) {
  let timer = null
  let t_start
  return function(...args) {
    const context = this
    const t_curr = +new Date()
    clearTimeout(timer)
    if (!t_start) {
      t_start = t_curr
    }
    if (t_curr - t_start >= mustRunDelay) {
      fn.apply(context, args)
      t_start = t_curr
    } else {
      timer = setTimeout(function() {
        fn.apply(context, args)
      }, delay)
    }
  }
}

/**
 * This is just a simple version of deep copy
 * Has a lot of edge cases bug
 * If you want to use a perfect deep copy, use lodash's _.cloneDeep
 * @param {Object} source
 * @returns {Object}
 */
export function deepClone(source) {
  if (!source && typeof source !== 'object') {
    throw new Error('error arguments', 'deepClone')
  }
  const targetObj = source.constructor === Array ? [] : {}
  Object.keys(source).forEach(keys => {
    if (source[keys] && typeof source[keys] === 'object') {
      targetObj[keys] = deepClone(source[keys])
    } else {
      targetObj[keys] = source[keys]
    }
  })
  return targetObj
}

/**
 * @param {Array} arr
 * @returns {Array}
 */
export function uniqueArr(arr) {
  return Array.from(new Set(arr))
}

/**
 * @returns {string}
 */
export function createUniqueString() {
  const timestamp = +new Date() + ''
  const randomNum = parseInt((1 + Math.random()) * 65536) + ''
  return (+(randomNum + timestamp)).toString(32)
}

/**
 * Check if an element has a class
 * @param {HTMLElement} elm
 * @param {string} cls
 * @returns {boolean}
 */
export function hasClass(ele, cls) {
  return !!ele.className.match(new RegExp('(\\s|^)' + cls + '(\\s|$)'))
}

/**
 * Add class to element
 * @param {HTMLElement} elm
 * @param {string} cls
 */
export function addClass(ele, cls) {
  if (!hasClass(ele, cls)) ele.className += ' ' + cls
}

/**
 * Remove class from element
 * @param {HTMLElement} elm
 * @param {string} cls
 */
export function removeClass(ele, cls) {
  if (hasClass(ele, cls)) {
    const reg = new RegExp('(\\s|^)' + cls + '(\\s|$)')
    ele.className = ele.className.replace(reg, ' ')
  }
}

/*
 * @author: houlong
 * @usage: 计算主容器高度
 * @param: 无
 * @return: {int}height
 */
export function getContainerHeight() {
  // 计算当前main容器的高度
  const height = document.body.clientHeight
  return height - 84
}

/*
 * @author: houlong
 * @usage: 格式化对象数组
 * @param: {Array} origin, 原对象的键值
 * @param: {Array} target, 目标键值， 要与origin的位置对应
 * @return: {Array} res, 格式化后的对象数组
 */
export function formatItems(items, origin, target) {
  return items.map(item => {
    const res = { ...item }
    origin.forEach((o, index) => {
      res[target[index]] = item[o]
    })
    return res
  })
}

/*
 * @author: houlong
 * @usage: 时间戳格式化日期格式（YYYY-MM-dd HH:mm:ss）
 * @param:
 * @return:
 */
export function formatTimestamp(timestamp) {
  if (!timestamp) {
    return ''
  }
  const unixtimestamp = new Date(timestamp)
  const year = 1900 + unixtimestamp.getYear()
  const month = '0' + (unixtimestamp.getMonth() + 1)
  const date = '0' + unixtimestamp.getDate()
  const hour = '0' + unixtimestamp.getHours()
  const minute = '0' + unixtimestamp.getMinutes()
  const second = '0' + unixtimestamp.getSeconds()
  return year + '-' + month.substring(month.length - 2, month.length) + '-' + date.substring(date.length - 2, date.length) +
            ' ' + hour.substring(hour.length - 2, hour.length) + ':' +
            minute.substring(minute.length - 2, minute.length) + ':' +
            second.substring(second.length - 2, second.length)
}

function getFieldValue(field, data) {
  const keys = field.split('.')
  return keys.reduce((value, key) => {
    if (!value) {
      return value
    }
    return value[key]
  }, data)
}

export function formatterTypeField(field, data, value) {
  switch (field.propType) {
    case 'boolean':
      return value ? (field.yes || '是') : (field.no || '否')
    case 'float':
      if (!value) {
        return `${value != null && value !== '' ? `${value}` : ''}`
      }
      return `${(field.valueTime ? (value * field.valueTime) : value).toFixed(field.places == null ? 2 : field.places)}`
    case 'price':
      return value != null ? value.toFixed(field.places == null ? 2 : field.places) : value
    case 'date':
      return value && parseTime(value, field.format)
    case 'date-array':
      return (value && value.length > 0) ? value.reduce((res, item, i) => {
        if (i === 0) {
          return item
        }
        if (i < 3) {
          return `${res}-${item < 10 ? '0' : ''}${item}`
        }
        if (i === 3) {
          return `${res} ${item < 10 ? '0' : ''}${item}`
        }
        return `${res}:${item < 10 ? '0' : ''}${item}`
      }, '') : value
    case 'daterange': {
      const { startKey, endKey } = field
      if (!startKey || !endKey) {
        return ''
      }
      const startValue = data[startKey]
      const endValue = data[endKey]
      if (!startValue || !endValue) {
        return ''
      }
      return `${parseTime(startValue, field.format)} 至 ${parseTime(endValue, field.format)}`
    }
    case 'minute': {
      if (value == null) {
        return value
      }
      const hours = Math.floor(value / 60)
      const minute = value % 60
      return `${hours < 10 ? '0' : ''}${hours}:${minute < 10 ? '0' : ''}${minute}:00`
    }
    case 'second': {
      if (value == null) {
        return value
      }
      const hours = Math.floor(value / 60 / 60)
      const rest = value - hours * 60 * 60
      const minute = Math.floor(rest / 60)
      const second = value % 60
      return `${hours < 10 ? '0' : ''}${hours}:${minute < 10 ? '0' : ''}${minute}:${second < 10 ? '0' : ''}${second}`
    }
    case 'region': {
      const { provinceKey = 'provinceName', cityKey = 'cityName', areaKey = 'areaName' } = field
      return `${data[provinceKey] || ''}${data[cityKey] || ''}${data[areaKey] || ''}`
    }
    default:
      return value
  }
}

/**
 * 格式化对应字段
 *
 * @export
 * @param {Object} field
 * @param {Object} data
 * @param {String} [propKey]
 * @returns
 */
export function formatterField(field, data, propKey) {
  if (!data) {
    return ''
  }
  const prop = propKey || (field.prop || field.key)
  // 多个属性决定对应的值
  const props = prop.split('|')
  let value = null
  if (typeof field.getValue === 'function') {
    value = field.getValue(field, data)
  } else {
    for (let i = 0; i < props.length; i++) {
      value = getFieldValue(props[i], data)
      if (value != null) {
        break
      }
    }
  }
  if (value == null && field.defaultValue != null) {
    return field.defaultValue
  }
  if (field.map) {
    return field.map[value] || value
  }
  if (field.propType) {
    const v = formatterTypeField(field, data, value)
    if (v == null || v === '') {
      return v
    }
    return `${field.prefix || ''}${v}${field.unit || ''}`
  }
  if (field.formatter) {
    return field.formatter(data, field, value)
  }
  if (value != null && value !== '') {
    return `${field.prefix || ''}${value}${field.unit || ''}`
  }
  return value
}

// 8-12位数字字母
export function randomPassword() {
  var s = ''
  var randomchar = function() {
    var n = Math.floor(Math.random() * 62)
    if (n < 10) return n // 1-10
    if (n < 36) return String.fromCharCode(n + 55) // A-Z
    return String.fromCharCode(n + 61) // a-z
  }
  const L = Math.round(Math.random() * 4 + 8)
  while (s.length < L) s += randomchar()
  return s
}


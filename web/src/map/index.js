/**
 * map对象生成options
 *
 * @export
 * @param {Object} map
 * @param sortFn
 * @returns
 */
export function mapToOptions(map, sortFn) {
  const options = Object.keys(map).map(key => ({
    value: key, label: map[key]
  }))
  if (!sortFn) {
    return options
  }
  return options.sort(sortFn)
}


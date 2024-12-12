import { mapToOptions } from '@/map/index'

export const status = { 1: '离线', 2: '在线' }
export const statusOptions = mapToOptions(status)

export const osType = { 'windows': 'Windows', 'linux': 'Linux', 'darwin': 'Mac'}
export const osTypeOptions = mapToOptions(osType)

export const goarch = { 'amd64': 'amd64', 'arm64': 'arm64'}
export const goarchOptions = mapToOptions(goarch)

export const compress = { 0: '无', 1: 'gzip', 2: 'flate'}
export const compressOptions = mapToOptions(compress)
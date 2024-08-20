import { mapToOptions } from '@/map/index'

export const status = { 0: '离线', 1: '在线' }
export const statusOptions = mapToOptions(status)

export const osType = { 1: 'Windows', 2: 'Linux', 3: 'Mac' }
export const osTypeOptions = mapToOptions(osType)

import { mapToOptions } from '@/map/index'

export const tunnelType = { 1: 'TCP', 2: 'HTTP' }
export const tunnelTypeOptions = mapToOptions(tunnelType)

export const status = { 0: '等待运行', 1: '运行中', 2: '运行失败' }
export const statusOptions = mapToOptions(status)

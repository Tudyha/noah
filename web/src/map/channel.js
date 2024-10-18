import { mapToOptions } from '@/map/index'

export const channelType = { 2: 'TCP', 4: 'HTTP' }
export const channelTypeOptions = mapToOptions(channelType)

export const status = { 0: '等待运行', 1: '运行中', 2: '运行失败' }
export const statusOptions = mapToOptions(status)

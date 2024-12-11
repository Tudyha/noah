import { mapToOptions } from '@/map/index'

export const tunnelType = { 1: 'tcp', 2: 'shadowsocks' }
export const tunnelTypeOptions = mapToOptions(tunnelType)

export const status = { 0: '等待运行', 1: '运行中', 2: '运行失败' }
export const statusOptions = mapToOptions(status)

export const cipher = { 'chacha20-ietf-poly1305': 'chacha20-ietf-poly1305', 'aes-128-gcm': 'aes-128-gcm', 'aes-256-gcm': 'aes-256-gcm' }
export const cipherOptions = mapToOptions(cipher)

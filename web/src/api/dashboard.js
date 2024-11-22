import request from '@/utils/request'

export function getDashboardData() {
  return request({
    url: '/dashboard',
    method: 'get'
  })
}

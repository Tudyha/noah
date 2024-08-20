import request from '@/utils/request'

export function generate(data) {
  return request({
    url: '/client/generate',
    method: 'post',
    data
  })
}

export function cmd(data) {
  return request({
    url: '/client/cmd',
    method: 'post',
    data
  })
}

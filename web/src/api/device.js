import request from '@/utils/request'

export function fetchList(params) {
  return request({
    url: '/device',
    method: 'get',
    params
  })
}

export function cmd(data) {
  return request({
    url: '/client/cmd',
    method: 'post',
    data
  })
}

export function deleteDevice(id) {
  return request({
    url: `/device/${id}`,
    method: 'delete'
  })
}

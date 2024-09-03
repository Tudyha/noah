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

export function update(data) {
  return request({
    url: `/client/${data.id}/update`,
    method: 'post',
    data
  })
}

export function fetchList(params) {
  return request({
    url: '/client',
    method: 'get',
    params
  })
}

export function deleteClient(id) {
  return request({
    url: `/client/${id}`,
    method: 'delete'
  })
}

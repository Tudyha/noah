import request from '@/utils/request'

export function newTunnel(data) {
    return request({
      url: `/client/${data.id}/tunnel`,
      method: 'post',
      data
    })
}

export function fetchList(clientId) {
  return request({
    url: `/client/${clientId}/tunnel`,
    method: 'get',
  })
}

export function deleteTunnel(data) {
  return request({
    url: `/client/${data.id}/tunnel/${data.tunnelId}`,
    method: 'delete',
  })
}

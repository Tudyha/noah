import request from '@/utils/request'

export function newChannel(data) {
    return request({
      url: `/client/${data.id}/channel`,
      method: 'post',
      data
    })
}

export function fetchList(clientId) {
  return request({
    url: `/client/${clientId}/channel`,
    method: 'get',
  })
}

export function deleteChannel(data) {
  return request({
    url: `/client/${data.id}/channel/${data.channelId}`,
    method: 'delete',
  })
}

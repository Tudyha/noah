import request from '@/utils/request'

export function newChannel(data) {
    return request({
      url: `/client/${data.id}/channel`,
      method: 'post',
      data
    })
  }
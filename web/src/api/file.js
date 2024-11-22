import request from '@/utils/request'

export function fetchList(params) {
  return request({
    url: `/client/${params.id}/file`,
    method: 'get',
    params: {path: params.path}
  })
}

export function fetchFileContent(params) {
  return request({
    url: `/client/${params.id}/file/content`,
    method: 'get',
    params: {path: params.path}
  })
}

export function renameFile(data) {
  return request({
    url: `/client/${data.id}/file/rename`,
    method: 'post',
    data
  })
}

export function deleteFile(data) {
  return request({
    url: `/client/${data.id}/file`,
    method: 'delete',
    data
  })
}

export function saveFile(data) {
  return request({
    url: `/client/${data.id}/file/content`,
    method: 'put',
    data
  })
}

export function uploadFile(data) {
  return request({
    url: `/client/${data.get('id')}/file`,
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function newDir(data) {
  return request({
    url: `/client/${data.id}/file/dir`,
    method: 'post',
    data
  })
}

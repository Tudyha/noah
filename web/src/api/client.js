import request from '@/utils/request'
import requestFile from '@/utils/request-file'

export function generate(params) {
  return requestFile({
    url: '/client/build',
    method: 'get',
    params: params,
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
    url: '/client/page',
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

export function systemInfo(params) {
  return request({
    url: `/client/${params.id}/systemInfo`,
    method: 'get',
    params: {
      start: params.start,
      end: params.end
    }
  })
}

export function fetchProcessList(id) {
  return request({
    url: `/client/${id}/process`,
    method: 'get'
  })
}

export function killProcess(params) {
  return request({
    url: `/client/${params.id}/process/${params.pid}`,
    method: 'delete'
  })
}

export function fetchNetworkList(id) {
  return request({
    url: `/client/${id}/network`,
    method: 'get'
  })
}

export function fetchDockerContainerList(id) {
  return request({
    url: `/client/${id}/docker/container`,
    method: 'get'
  })
}

export function getClient(id) {
  return request({
    url: `/client/${id}`,
    method: 'get'
  })
}

export function getClientInstallScript() {
  return request({
    url: `/client/install_script`,
    method: 'get'
  })
}

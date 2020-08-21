import request from '@/utils/request'

export function getPermission(params) {
  return request({
    url: '/permission',
    method: 'get',
    params: params
  })
}

export function editPermission(data) {
  return request({
    url: `/permission`,
    method: 'put',
    data
  })
}


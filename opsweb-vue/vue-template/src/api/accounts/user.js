import request from '@/utils/request'

export function getUserList(params) {
  return request({
    url: '/user',
    method: 'get',
    params: params
  })
}

export function addUser(data) {
  return request({
    url: '/user',
    method: 'post',
    data
  })
}

export function editUser(data) {
  return request({
    url: `/user`,
    method: 'put',
    data
  })
}

export function editUserStatus(data) {
  return request({
    url: `/user/status`,
    method: 'put',
    data
  })
}

export function editUserPasswd(data) {
  return request({
    url: `/user/pwd`,
    method: 'put',
    data
  })
}

export function editUserAdmin(data) {
  return request({
    url: `/user/admin`,
    method: 'put',
    data
  })
}

export function editUserGroups(data) {
  return request({
    url: `/user/groups`,
    method: 'put',
    data
  })
}

export function deleteUser(data) {
  return request({
    url: `/user`,
    method: 'delete',
    data
  })
}

export function getUserInfo() {
  return request({
    url: `/user/info`,
    method: 'get'
  })
}


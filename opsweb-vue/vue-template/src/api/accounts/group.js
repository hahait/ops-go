import request from '@/utils/request'

export function getGroupList(data) {
  return request({
    url: '/group',
    method: 'get',
    data
  })
}

export function addGroup(data) {
  return request({
    url: '/group',
    method: 'post',
    data
  })
}

export function editGroup(data) {
  return request({
    url: `/group`,
    method: 'put',
    data
  })
}

export function deleteGroup(data) {
  return request({
    url: `/group`,
    method: 'delete',
    data
  })
}

export function editGroupUsers(data) {
  return request({
    url: `/group/users`,
    method: 'put',
    data
  })
}

export function editGroupPerms(data) {
  return request({
    url: `/group/perms`,
    method: 'put',
    data
  })
}

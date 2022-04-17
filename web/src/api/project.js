import request from '@/utils/request'

const projectApi = {
  ListProjects: '/server/project/list'
}

export function listProjects (parameter) {
  return request({
    url: projectApi.ListProjects,
    method: 'get',
    data: parameter
  })
}

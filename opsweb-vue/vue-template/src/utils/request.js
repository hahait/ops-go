import axios from 'axios'
import { Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'
import router from '@/router'

// 创建axios实例
const service = axios.create({
  baseURL: process.env.BASE_API, // api的base_url
  timeout: 5000 // 请求超时时间
})

// request拦截器
service.interceptors.request.use(
  config => {
    if (store.getters.token) {
      config.headers['Authorization'] = 'Bearer ' + getToken() // 让每个请求携带自定义token 请根据实际情况自行修改
    }
    return config
  },
  error => {
    // Do something with request error
    console.log('file: src/utils/request.js request 过滤器报错: ', error) // for debug
    Promise.reject(error)
  }
)

// respone拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    console.log('file: src/utils/request.js 后端返回的 response.data: ', response.data)
    if (res.code === 1) {
      Message.error(res.msg)
      return Promise.reject(res.errmsg)
    }
    return res
  },
  error => {
    // console.log('file: src/utils/request.js 从后端收到的响应 error : ', error)
    // console.log('file: src/utils/request.js 从后端收到的响应 error.response : ', error.response)
    const erd = error.response.data
    if (!error.response) {
      Message.error('系统错误')
      return Promise.reject(error.response)
    }
    if (error.response.status === 401) {
      Message({
        message: '认证已过期或未认证,请输入用户名和密码进行认证...',
        type: 'error',
        duration: 2000,
        onClose: function() {
          router.push({ path: '/login' })
        }
      })
    } else if (error.response.status === 403) {
      Message({
        message: '权限拒绝',
        type: 'error'
        // duration: 1500,
        // onClose: function() {
        //   router.push({ path: '/dashboard' })
        // }
      })
    } else if (error.response.status === 400) {
      Message({
        type: 'error',
        message: '认证失效，请重新登陆'
      })
    } else if (error.response.status === 500) {
      Message({
        type: 'error',
        message: '服务器内部错误: ' + erd.msg
      })
      return Promise.reject(erd.errmsg)
    } else {
      Message({
        type: 'error',
        message: error.response.data.detail
      })
    }
    return Promise.reject(error)
  }
)

export default service

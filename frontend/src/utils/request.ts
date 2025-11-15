/**
 * 统一的 HTTP 请求工具
 * 基于 axios 封装，类似 Web 端的实现
 */

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { GetAPIURL, GetToken } from '../services/auth'

interface ApiResponse<T = any> {
  code: number
  data: T
  message?: string
  msg?: string
}

// 创建 axios 实例
let axiosInstance: AxiosInstance | null = null

/**
 * 获取或创建 axios 实例
 */
async function getAxiosInstance(): Promise<AxiosInstance> {
  if (axiosInstance) {
    return axiosInstance
  }

  // 获取 API 基础 URL
  const baseURL = await GetAPIURL()

  // 创建实例
  axiosInstance = axios.create({
    baseURL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
      'accept': 'application/json, text/plain, */*',
      'clienttype': '1',
      'language': 'zh',
    },
  })

  // 请求拦截器 - 添加 token
  axiosInstance.interceptors.request.use(
    async (config) => {
      try {
        const token = await GetToken()
        if (token && config.headers) {
          config.headers['accesstoken'] = token
        }
      } catch (error) {
        console.warn('获取 token 失败:', error)
      }
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )

  // 响应拦截器 - 统一处理响应
  axiosInstance.interceptors.response.use(
    (response) => {
      const data: ApiResponse = response.data
      
      // 检查业务状态码
      if (data.code === 0 || data.code === 200) {
        return data.data
      } else {
        const errorMsg = data.message || data.msg || '请求失败'
        console.error('API 错误:', errorMsg)
        return Promise.reject(new Error(errorMsg))
      }
    },
    (error) => {
      console.error('请求失败:', error.message)
      return Promise.reject(error)
    }
  )

  return axiosInstance
}

/**
 * 统一的请求对象，类似 Web 端的 api 对象
 */
export const request = {
  /**
   * 获取原始 axios 实例（如需要自定义配置）
   */
  getAxios: getAxiosInstance,

  /**
   * GET 请求
   */
  get: async <T = any>(
    url: string,
    params?: any,
    config?: AxiosRequestConfig
  ): Promise<T> => {
    const instance = await getAxiosInstance()
    return instance.get(url, { params, ...config })
  },

  /**
   * POST 请求
   */
  post: async <T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> => {
    const instance = await getAxiosInstance()
    return instance.post(url, data, config)
  },

  /**
   * PUT 请求
   */
  put: async <T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> => {
    const instance = await getAxiosInstance()
    return instance.put(url, data, config)
  },

  /**
   * DELETE 请求
   */
  delete: async <T = any>(
    url: string,
    params?: any,
    config?: AxiosRequestConfig
  ): Promise<T> => {
    const instance = await getAxiosInstance()
    return instance.delete(url, { params, ...config })
  },
}

// 默认导出
export default request

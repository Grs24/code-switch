/**
 * 自动配置服务
 * 负责在用户登录后自动获取 API Key 并执行一键配置
 */

import { request } from '../utils/request'
import { ConfigureAll, GetDefaultBaseURL } from './autoConfig'

interface AuthUser {
  email?: string
  userId?: string
  [key: string]: any
}

export interface ApiKeyData {
  id: number
  token: string
  status: number
}

/**
 * 获取用户的 API Key
 * @param userInfo 用户信息
 * @returns API Key 数据，如果失败返回 null
 */
export async function fetchUserApiKey(userInfo: AuthUser | null): Promise<ApiKeyData | null> {
  if (!userInfo) return null
  
  try {
    const result = await request.get<{ list: ApiKeyData[] }>('/aicoding0011/xapi/user/token')
    return result?.list?.[0] || null
  } catch (error) {
    console.error('[AutoConfig] 获取 API Key 失败:', error)
    return null
  }
}

/**
 * 后台静默执行自动配置
 * @param userInfo 用户信息
 * @returns 配置成功则返回 API Key 数据，否则返回 null
 */
export async function runAutoConfigInBackground(userInfo: AuthUser | null): Promise<ApiKeyData | null> {
  if (!userInfo) return null
  
  try {
    const apiKeyData = await fetchUserApiKey(userInfo)
    if (!apiKeyData?.token) return null
    
    const baseURL = await GetDefaultBaseURL().catch(() => '')
    const configResult = await ConfigureAll(apiKeyData.token, baseURL)
    
    if (configResult.success) {
      console.log('[AutoConfig] 配置成功')
      return apiKeyData
    } else {
      console.warn('[AutoConfig] 配置失败:', configResult.errors)
      return null
    }
  } catch (error) {
    console.error('[AutoConfig] 配置异常:', error)
    return null
  }
}

/**
 * 检查用户是否已经配置过
 * @param userInfo 用户信息
 * @returns 是否已配置
 */
export function hasAutoConfigured(userInfo: AuthUser | null): boolean {
  if (!userInfo) return false
  
  const userId = userInfo.userId || userInfo.email
  const configKey = `auto_config_done_${userId}`
  return localStorage.getItem(configKey) === 'true'
}

/**
 * 标记用户已完成自动配置
 * @param userInfo 用户信息
 */
export function markAutoConfigDone(userInfo: AuthUser | null): void {
  if (!userInfo) return
  
  const userId = userInfo.userId || userInfo.email
  const configKey = `auto_config_done_${userId}`
  localStorage.setItem(configKey, 'true')
}

/**
 * 清除用户的自动配置标记
 * @param userInfo 用户信息
 */
export function clearAutoConfigMark(userInfo: AuthUser | null): void {
  if (!userInfo) return
  
  const userId = userInfo.userId || userInfo.email
  const configKey = `auto_config_done_${userId}`
  localStorage.removeItem(configKey)
}

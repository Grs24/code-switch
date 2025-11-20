import { Call } from '@wailsio/runtime'

export interface UserInfo {
  token: string
  userId?: string
  username?: string
  email?: string
  avatar?: string
}

const SERVICE_NAME = 'codeswitch/services.AuthService'

/**
 * 获取登录页 URL
 */
export async function GetLoginURL(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetLoginURL`)
}

/**
 * 设置用户信息
 */
export async function SetUserInfo(userInfoJSON: string): Promise<void> {
  return Call.ByName(`${SERVICE_NAME}.SetUserInfo`, userInfoJSON)
}

/**
 * 获取用户信息
 */
export async function GetUserInfo(): Promise<UserInfo> {
  return Call.ByName(`${SERVICE_NAME}.GetUserInfo`)
}

/**
 * 检查是否已登录
 */
export async function IsLogin(): Promise<boolean> {
  return Call.ByName(`${SERVICE_NAME}.IsLogin`)
}

/**
 * 退出登录
 */
export async function Logout(): Promise<void> {
  return Call.ByName(`${SERVICE_NAME}.Logout`)
}

/**
 * 获取用户 token
 */
export async function GetToken(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetToken`)
}

/**
 * 获取 API 基础 URL
 */
export async function GetAPIURL(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetAPIURL`)
}

/**
 * 获取 AICoding Base URL
 */
export async function GetAICodeBaseURL(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetAICodeBaseURL`)
}

/**
 * 刷新 0011 供应商的 API Key
 * 用于用户在界面上更新 API Key 后同步到供应商卡片
 */
export async function RefreshProviderApiKey(): Promise<void> {
  return Call.ByName(`${SERVICE_NAME}.RefreshProviderApiKey`)
}

/**
 * 获取订阅页面 URL
 */
export async function GetSubscriptionURL(): Promise<string> {
  return Call.ByName(`${SERVICE_NAME}.GetSubscriptionURL`)
}

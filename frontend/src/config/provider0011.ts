/**
 * 0011 供应商配置
 * 集中管理 0011 供应商的所有配置信息
 */

import type { AutomationCard } from '../data/cards'
import { GetAICodeBaseURL } from '../services/auth'

// 默认 API Base URL（作为 fallback）
export const API_BASE_URL = 'https://aicoding.2233.ai'

// 缓存从后端获取的 Base URL
let cachedBaseURL: string | null = null

/**
 * 获取 API Base URL（优先从后端获取，失败则使用默认值）
 */
export async function getAPIBaseURL(): Promise<string> {
  if (cachedBaseURL) {
    return cachedBaseURL
  }
  
  try {
    cachedBaseURL = await GetAICodeBaseURL()
    return cachedBaseURL
  } catch (error) {
    console.warn('[Provider0011] 获取 Base URL 失败，使用默认值:', error)
    return API_BASE_URL
  }
}

/**
 * 0011 供应商的基础配置
 */
export const PROVIDER_0011_CONFIG = {
  // Claude Code 供应商 ID
  CLAUDE_ID: 100,
  
  // Codex 供应商 ID
  CODEX_ID: 200,
  
  // 供应商名称
  NAME: '0011',
  
  // API Base URL
  API_URL: API_BASE_URL,
  
  // 官网地址
  OFFICIAL_SITE: 'https://0011.ai',
  
  // 图标
  ICON: 'aicoding',
  
  // 主题色
  TINT: 'rgba(10, 132, 255, 0.14)',
  ACCENT: '#0aff5cff',
} as const

/**
 * 创建 0011 供应商卡片配置
 * @param type 供应商类型 ('claude' 或 'codex')
 * @param apiKey API Key（可选）
 * @param enabled 是否启用（可选）
 * @param apiUrl API URL（可选，如果不提供则使用默认值）
 * @returns 供应商卡片配置
 */
export function create0011Provider(
  type: 'claude' | 'codex',
  apiKey: string = '',
  enabled: boolean = false,
  apiUrl?: string
): AutomationCard {
  const id = type === 'claude' ? PROVIDER_0011_CONFIG.CLAUDE_ID : PROVIDER_0011_CONFIG.CODEX_ID
  
  return {
    id,
    name: PROVIDER_0011_CONFIG.NAME,
    apiUrl: apiUrl || PROVIDER_0011_CONFIG.API_URL,
    apiKey,
    officialSite: PROVIDER_0011_CONFIG.OFFICIAL_SITE,
    icon: PROVIDER_0011_CONFIG.ICON,
    tint: PROVIDER_0011_CONFIG.TINT,
    accent: PROVIDER_0011_CONFIG.ACCENT,
    enabled,
  }
}

/**
 * 创建 0011 供应商卡片配置（异步版本，从后端获取 Base URL）
 * @param type 供应商类型 ('claude' 或 'codex')
 * @param apiKey API Key（可选）
 * @param enabled 是否启用（可选）
 * @returns 供应商卡片配置
 */
export async function create0011ProviderAsync(
  type: 'claude' | 'codex',
  apiKey: string = '',
  enabled: boolean = false
): Promise<AutomationCard> {
  const apiUrl = await getAPIBaseURL()
  return create0011Provider(type, apiKey, enabled, apiUrl)
}

/**
 * 检查是否为 0011 供应商
 * @param id 供应商 ID
 * @returns 是否为 0011 供应商
 */
export function is0011Provider(id: number): boolean {
  return id === PROVIDER_0011_CONFIG.CLAUDE_ID || id === PROVIDER_0011_CONFIG.CODEX_ID
}

/**
 * 获取 0011 供应商的 ID（根据类型）
 * @param type 供应商类型
 * @returns 供应商 ID
 */
export function get0011ProviderId(type: 'claude' | 'codex'): number {
  return type === 'claude' ? PROVIDER_0011_CONFIG.CLAUDE_ID : PROVIDER_0011_CONFIG.CODEX_ID
}

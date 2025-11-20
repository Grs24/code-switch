import { request } from '../utils/request'

// Coding Total Tokens 响应（与 Web 端保持一致）
export interface CodingTotalTokensResponse {
  today_used: number
  total_available: number
  total_daily_limit: number
  /** 总额度金额（美元） */
  total_amount: string
  /** 每日额度金额（美元） */
  total_daily_limit_amount: string
  /** 当日已使用金额（美元） */
  today_used_amount: string
  uid: number
}

// 余额信息（用于前端显示）
export interface QuotaInfo {
  totalAmount: number // 总金额
  dailyQuota: number // 每日额度
  usedToday: number // 今日已用
  remaining: number // 剩余
  isLowQuota: boolean // 是否低额度
  isCriticalQuota: boolean // 是否极低额度
}

/**
 * 获取 Coding 总 Tokens 信息
 */
export async function fetchCodingTotalTokens(): Promise<CodingTotalTokensResponse | null> {
  try {
    const data = await request.get<CodingTotalTokensResponse>('/aicoding0011/xapi/orders/total_tokens')
    return data
  } catch (error) {
    console.error('[Service0011] 获取总 Tokens 失败:', error)
    return null
  }
}

/**
 * 获取用户余额信息
 */
export async function fetchQuotaInfo(): Promise<QuotaInfo | null> {
  try {
    const data = await fetchCodingTotalTokens()
    if (!data) {
      return null
    }

    // 解析金额字符串为数字
    const totalAmount = parseFloat(data.total_amount) || 0
    const dailyQuota = parseFloat(data.total_daily_limit_amount) || 0
    const usedToday = parseFloat(data.today_used_amount) || 0
    const remaining = totalAmount - usedToday

    // 判断是否低额度（基于剩余金额百分比）
    const usagePercentage = totalAmount > 0 ? (usedToday / totalAmount) * 100 : 0
    const remainingPercentage = 100 - usagePercentage
    const isLowQuota = remainingPercentage < 10
    const isCriticalQuota = remainingPercentage < 5

    return {
      totalAmount,
      dailyQuota,
      usedToday,
      remaining,
      isLowQuota,
      isCriticalQuota,
    }
  } catch (error) {
    console.error('[Service0011] 获取余额信息失败:', error)
    return null
  }
}

/**
 * 刷新余额信息
 */
export async function refreshQuota(): Promise<QuotaInfo | null> {
  return fetchQuotaInfo()
}



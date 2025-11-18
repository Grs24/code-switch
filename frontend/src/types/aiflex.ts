/**
 * AiFlex 套餐相关类型定义
 */

/** 使用状态枚举 */
export enum AiflexUseStatus {
  /** 未使用 */
  NOT_USED = 0,
  /** 使用中 */
  IN_USE = 1,
  /** 已过期 */
  EXPIRED = 2,
}

/** 获取使用信息-响应参数 */
export interface GetAiflexUseInfoResponse {
  /** 倒计时，字符串或数字 */
  count_down: string | number
  /** 过期时间戳，毫秒 */
  expiration_time: number
  /** 是否支持o1模型 */
  o1: boolean
  /** 随心用套餐 */
  share_type: string | number
  /** 开始时间戳，毫秒 */
  start_time: number
  /** 使用状态 */
  use_status: AiflexUseStatus
}

/** 获取使用信息-请求参数 */
export interface GetAiflexUseInfoParams {
  /** 平台类型 */
  platform: string
}

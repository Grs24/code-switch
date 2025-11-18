/**
 * AiFlex 套餐相关服务
 */

import { request } from '../utils/request'
import type { GetAiflexUseInfoResponse, GetAiflexUseInfoParams } from '../types/aiflex'

/**
 * 获取 AiFlex 使用信息
 */
export const fetchAiflexUseInfo = async (
  params: GetAiflexUseInfoParams
): Promise<GetAiflexUseInfoResponse> => {
  return request.get<GetAiflexUseInfoResponse>('/aiflex/shareUseInfo', params)
}

/**
 * 退出登录 Composable
 * 封装退出登录的完整逻辑，可在任何组件中复用
 */

import { ref } from 'vue'
import { Logout, GetLoginURL } from '../services/auth'
import { clearAutoConfigMark } from '../services/autoConfigService'
import { Events } from '@wailsio/runtime'

export function useLogout() {
  const logoutLoading = ref(false)

  const logout = async (userInfo: any = null) => {
    logoutLoading.value = true
    try {
      // 通知 Web 端清除登录状态
      const loginURL = await GetLoginURL()
      const url = new URL(loginURL)
      const logoutURL = `${url.protocol}//${url.host}/logout`
      
      const logoutFrame = document.createElement('iframe')
      logoutFrame.style.display = 'none'
      
      // 添加错误处理，避免 iframe 加载失败时的错误提示
      logoutFrame.onerror = () => {
        console.warn('[Logout] Web 端登出请求失败，但不影响桌面端退出')
      }
      
      logoutFrame.src = logoutURL
      document.body.appendChild(logoutFrame)
      
      // 等待 Web 端清除完成（给足够时间让事件处理和 localStorage 清除）
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      try {
        document.body.removeChild(logoutFrame)
      } catch (e) {
        // 忽略移除 iframe 时的错误
      }
      
      // 清除桌面端状态
      await Logout()
      clearAutoConfigMark(userInfo)
      
      // 触发退出事件，通知其他组件
      Events.Emit('auth:logout')
      
      console.log('[useLogout] 退出登录成功')
      return true
    } catch (error) {
      console.error('[useLogout] 退出登录失败:', error)
      return false
    } finally {
      logoutLoading.value = false
    }
  }

  return {
    logout,
    logoutLoading
  }
}

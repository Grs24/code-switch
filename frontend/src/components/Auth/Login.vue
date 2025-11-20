<template>
  <div class="login-container">
    <div class="login-header">
      <span></span>
      <button @click="closeWindow" class="close-btn">×</button>
    </div>
    
    <div class="webview-container">
      <div v-if="loading" class="loading-overlay">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
      
      <iframe
        v-if="loginURL"
        :key="iframeKey"
        ref="loginIframe"
        :src="loginURL"
        class="login-iframe"
        :class="{ 'iframe-hidden': loading }"
        @load="onIframeLoad"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { GetLoginURL, SetUserInfo } from '../../services/auth'

const { t } = useI18n()
const router = useRouter()
const loginIframe = ref<HTMLIFrameElement | null>(null)
const loginURL = ref('')
const loading = ref(true)
const iframeKey = ref(Date.now())

let checkInterval: number | null = null

onMounted(async () => {
  
  // 生成新的 key，确保 iframe 完全重新创建
  iframeKey.value = Date.now()
  
  try {
    const baseURL = await GetLoginURL()
    // 添加时间戳参数，强制 iframe 重新加载，避免使用缓存
    loginURL.value = `${baseURL}?t=${iframeKey.value}`
    
    // 开始轮询检查 localStorage
    startCheckingAuth()
  } catch (error) {
    console.error('[Login] 获取登录 URL 失败:', error)
    loading.value = false
  }
})

onUnmounted(() => {
  stopCheckingAuth()
})

onUnmounted(() => {
  stopCheckingAuth()
})

const onIframeLoad = () => {
  loading.value = false
}

const startCheckingAuth = () => {
  // 每 500ms 检查一次 iframe 的 localStorage
  checkInterval = window.setInterval(() => {
    checkAuthStatus()
  }, 500)
}

const stopCheckingAuth = () => {
  if (checkInterval) {
    clearInterval(checkInterval)
    checkInterval = null
  }
}

const checkAuthStatus = () => {
  try {
    if (!loginIframe.value?.contentWindow) return

    // 尝试从 iframe 的 localStorage 读取用户信息
    const iframeWindow = loginIframe.value.contentWindow
    
    // 注入脚本到 iframe 检查 localStorage
    iframeWindow.postMessage({ type: 'CHECK_AUTH' }, '*')
  } catch (error) {
    // 跨域限制，无法直接访问
    console.debug('[Login] 无法直接访问 iframe localStorage')
  }
}

// 监听来自 iframe 的消息
const handleMessage = async (event: MessageEvent) => {
  // 验证消息来源
  if (!loginURL.value || !event.origin.includes(new URL(loginURL.value).origin)) {
    return
  }

  if (event.data.type === 'AUTH_SUCCESS' && event.data.userInfo) {
    stopCheckingAuth()
    loading.value = true

    try {
      // 将用户信息发送到 Go 后端
      await SetUserInfo(JSON.stringify(event.data.userInfo))
      // 登录成功后，静默执行自动配置
      const { runAutoConfigInBackground } = await import('../../services/autoConfigService')
      runAutoConfigInBackground(event.data.userInfo).catch((error) => {
        // console.warn('[Login] 后台自动配置失败:', error)
      })

      // 触发登录成功事件，通知主界面更新
      const { Events } = await import('@wailsio/runtime')
      Events.Emit('auth:login-success', event.data.userInfo)

      // 跳转到主页
      setTimeout(() => {
        router.push('/')
      }, 500)
    } catch (error) {
      loading.value = false
    }
  }
}

const closeWindow = () => {
  router.push('/')
}

// 添加消息监听
onMounted(() => {
  window.addEventListener('message', handleMessage)
})

onUnmounted(() => {
  window.removeEventListener('message', handleMessage)
})
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #1a1a1a;
}

.login-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #2a2a2a;
  border-bottom: 1px solid #3a3a3a;
}

.login-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: #fff;
}

.close-btn {
  background: none;
  border: none;
  color: #999;
  font-size: 32px;
  line-height: 1;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #3a3a3a;
  color: #fff;
}

.webview-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.login-iframe {
  width: 100%;
  height: 100%;
  border: none;
  background: #fff;
  opacity: 1;
  transition: opacity 0.3s ease;
}

.login-iframe.iframe-hidden {
  opacity: 0;
  pointer-events: none;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 100;
  pointer-events: none;
}

.loading-spinner {
  width: 56px;
  height: 56px;
  border: 5px solid #2a2a2a;
  border-top-color: #0a84ff;
  border-right-color: #0a84ff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-overlay p {
  margin-top: 24px;
  color: #ccc;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 0.5px;
}
</style>

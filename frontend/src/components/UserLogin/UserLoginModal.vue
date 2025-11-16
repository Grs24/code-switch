<template>
  <BaseModal :open="open" title="用户信息" @close="handleClose">
    <div class="user-info-modal">
      <!-- 用户基本信息 -->
      <div class="user-info-section">
        <label class="user-info-label">邮箱</label>
        <div class="user-info-value">{{ userInfo?.email || '-' }}</div>
      </div>

      <!-- API Key 信息 -->
      <div class="user-info-section">
        <div class="api-key-header">
          <label class="user-info-label">API Key</label>
          <BaseButton 
            v-if="!apiKeyLoading"
            variant="outline" 
            @click="updateApiKey"
            style="padding: 4px 12px; font-size: 12px;"
          >
            {{ apiKeyData ? '更新密钥' : '创建密钥' }}
          </BaseButton>
        </div>
        
        <div v-if="apiKeyLoading" class="api-key-loading">
          加载中...
        </div>
        
        <div v-else-if="!apiKeyData" class="api-key-empty">
          还没有 API 密钥，点击上方按钮创建
        </div>
        
        <div v-else class="api-key-display">
          <code class="api-key-code">{{ displayApiKey }}</code>
          <button class="icon-btn" @click="toggleApiKeyVisibility">
            <svg v-if="showApiKey" viewBox="0 0 24 24" width="16" height="16">
              <path d="M13.73 4.5A6.5 6.5 0 0 0 5.5 12c0 .88.18 1.71.5 2.47M9.88 9.88a3 3 0 1 0 4.24 4.24M9.88 9.88L5.5 5.5m4.38 4.38l4.24 4.24m0 0L18.5 18.5M3 3l18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" width="16" height="16">
              <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
          </button>
          <button class="icon-btn" @click="copyApiKey">
            <svg viewBox="0 0 24 24" width="16" height="16">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2" stroke="currentColor" stroke-width="2" fill="none"/>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" stroke="currentColor" stroke-width="2" fill="none"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- 自动配置区域 -->
      <AutoConfigSection
        v-if="apiKeyData"
        ref="autoConfigRef"
        :api-key-data="apiKeyData"
        :user-info="userInfo"
      />

      <!-- 退出登录 -->
      <div class="user-info-actions">
        <BaseButton variant="danger" @click="handleLogout" :disabled="logoutLoading">
          <span v-if="logoutLoading" class="loading-spinner"></span>
          {{ logoutLoading ? '退出中...' : '退出登录' }}
        </BaseButton>
      </div>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import BaseModal from '../common/BaseModal.vue'
import BaseButton from '../common/BaseButton.vue'
import AutoConfigSection from './AutoConfigSection.vue'
import { fetchUserApiKey } from '../../services/autoConfigService'
import { request } from '../../utils/request'

interface Props {
  open: boolean
  userInfo: any
  logoutLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  logoutLoading: false
})
const emit = defineEmits(['close', 'logout'])

// 获取 AutoConfigSection 的引用
const autoConfigRef = ref<InstanceType<typeof AutoConfigSection> | null>(null)

// 状态
const showApiKey = ref(false)
const apiKeyData = ref<{ id: number; token: string; status: number } | null>(null)
const apiKeyLoading = ref(false)

// 计算属性
const apiKey = computed(() => apiKeyData.value?.token || '')

const displayApiKey = computed(() => {
  if (!apiKey.value) return '-'
  if (showApiKey.value) return apiKey.value
  // 脱敏显示
  if (apiKey.value.length < 13) return apiKey.value
  return `${apiKey.value.slice(0, 8)}****${apiKey.value.slice(-4)}`
})

// 监听弹窗打开
watch(() => props.open, async (isOpen) => {
  if (isOpen && props.userInfo) {
    await loadApiKey()
  }
})

// 加载 API Key（复用 autoConfigService 中的函数）
const loadApiKey = async () => {
  if (!props.userInfo) return
  
  apiKeyLoading.value = true
  try {
    const result = await fetchUserApiKey(props.userInfo)
    if (result) {
      apiKeyData.value = result
    }
  } catch (error) {
    console.error('获取 API Key 失败:', error)
  } finally {
    apiKeyLoading.value = false
  }
}

// 更新 API Key
const updateApiKey = async () => {
  if (!props.userInfo) return
  
  apiKeyLoading.value = true
  try {
    if (apiKeyData.value) {
      const result = await request.put('/aicoding0011/xapi/user/token', {
        id: apiKeyData.value.id,
        token: apiKeyData.value.token,
        status: apiKeyData.value.status,
      })
      apiKeyData.value = result
    } else {
      const result = await request.post('/aicoding0011/xapi/user/token', { status: 1 })
      apiKeyData.value = result
    }
  } catch (error) {
    console.error('更新 API Key 失败:', error)
  } finally {
    apiKeyLoading.value = false
  }
}

// 切换 API Key 可见性
const toggleApiKeyVisibility = () => {
  showApiKey.value = !showApiKey.value
}

// 复制 API Key
const copyApiKey = async () => {
  if (!apiKey.value) return
  try {
    await navigator.clipboard.writeText(apiKey.value)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 退出登录
const handleLogout = () => {
  // 触发父组件的退出登录处理
  emit('logout')
}

// 关闭弹窗
const handleClose = () => {
  showApiKey.value = false
  emit('close')
}

// 暴露给父组件的方法：触发自动配置
const triggerAutoConfig = async () => {
  await loadApiKey()
  
  if (apiKeyData.value && autoConfigRef.value) {
    await autoConfigRef.value.triggerAutoConfig()
  }
}

// 暴露方法给父组件
defineExpose({
  triggerAutoConfig
})
</script>

<style scoped>
.user-info-modal {
  padding: 8px 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-info-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-info-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--mac-text-secondary);
  margin-bottom: 4px;
}

.user-info-value {
  font-size: 0.9375rem;
  color: var(--mac-text-primary);
  padding: 10px 12px;
  background: var(--mac-bg-secondary);
  border-radius: 6px;
  border: 1px solid var(--mac-border);
}

.api-key-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.api-key-loading,
.api-key-empty {
  padding: 12px;
  text-align: center;
  color: var(--mac-text-secondary);
  background: var(--mac-bg-secondary);
  border-radius: 6px;
  font-size: 0.875rem;
}

.api-key-display {
  display: flex;
  gap: 8px;
  align-items: center;
}

.api-key-code {
  flex: 1;
  padding: 10px 12px;
  background: var(--mac-bg-secondary);
  border: 1px solid var(--mac-border);
  border-radius: 6px;
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-size: 0.8125rem;
  color: var(--mac-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-btn {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--mac-bg-secondary);
  border: 1px solid var(--mac-border);
  border-radius: 6px;
  color: var(--mac-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.icon-btn:hover {
  background: var(--mac-surface);
  border-color: rgba(99, 102, 241, 0.3);
  color: rgba(99, 102, 241, 0.8);
}

.icon-btn:active {
  transform: scale(0.95);
}

.user-info-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  border-top: 1px solid var(--mac-border);
}

.loading-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  margin-right: 8px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>

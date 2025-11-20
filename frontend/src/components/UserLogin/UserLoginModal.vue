<template>
  <BaseModal :open="open" title="用户信息" @close="handleClose">
    <div class="user-info-modal">
      <div class="user-info-row">
        <div class="user-info-section user-info-email">
          <label class="user-info-label">邮箱</label>
          <div class="user-info-value">{{ userInfo?.email || '-' }}</div>
        </div>
        
        <div class="user-info-section user-info-subscription">
          <label class="user-info-label">套餐</label>
          <SubscriptionStatus 
            ref="subscriptionRef"
            :user-info="userInfo"
            compact
          />
        </div>
      </div>

      <div class="user-info-section">
        <div class="api-key-header">
          <label class="user-info-label">API Key</label>
          <BaseButton 
            v-if="!apiKeyLoading"
            variant="primary" 
            @click="updateApiKey"
            class="create-api-key-btn"
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

      <AutoConfigSection
        v-if="apiKeyData"
        ref="autoConfigRef"
        :api-key-data="apiKeyData"
        :user-info="userInfo"
      />

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
import SubscriptionStatus from '../Subscription/SubscriptionStatus.vue'
import { fetchUserApiKey } from '../../services/autoConfigService'
import { request } from '../../utils/request'
import { useLogout } from '../../composables/useLogout'

interface Props {
  open: boolean
  userInfo: any
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'logout-success'])

const { logout, logoutLoading } = useLogout()

const autoConfigRef = ref<InstanceType<typeof AutoConfigSection> | null>(null)
const subscriptionRef = ref<InstanceType<typeof SubscriptionStatus> | null>(null)

const showApiKey = ref(false)
const apiKeyData = ref<{ id: number; token: string; status: number } | null>(null)
const apiKeyLoading = ref(false)

const apiKey = computed(() => apiKeyData.value?.token || '')

const displayApiKey = computed(() => {
  if (!apiKey.value) return '-'
  if (showApiKey.value) return apiKey.value
  if (apiKey.value.length < 13) return apiKey.value
  return `${apiKey.value.slice(0, 8)}****${apiKey.value.slice(-4)}`
})

watch(() => props.open, async (isOpen) => {
  if (isOpen && props.userInfo) {
    await loadApiKey()
    subscriptionRef.value?.refresh()
  }
})

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

    // 更新成功后，同步刷新 0011 供应商卡片的 API Key
    const { RefreshProviderApiKey } = await import('../../services/auth')
    await RefreshProviderApiKey()
    
    // 触发事件，通知主界面重新加载供应商列表
    const { Events } = await import('@wailsio/runtime')
    Events.Emit('provider:refresh-needed', { timestamp: Date.now() })
  } catch (error) {
    console.error('更新 API Key 失败:', error)
  } finally {
    apiKeyLoading.value = false
  }
}

const toggleApiKeyVisibility = () => {
  showApiKey.value = !showApiKey.value
}

const copyApiKey = async () => {
  if (!apiKey.value) return
  try {
    await navigator.clipboard.writeText(apiKey.value)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

const handleLogout = async () => {
  const success = await logout(props.userInfo)
  if (success) {
    emit('logout-success')
    handleClose()
  }
}

const handleClose = () => {
  showApiKey.value = false
  emit('close')
}

const triggerAutoConfig = async () => {
  await loadApiKey()
  
  if (apiKeyData.value && autoConfigRef.value) {
    await autoConfigRef.value.triggerAutoConfig()
  }
}

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

.user-info-row {
  display: flex;
  gap: 16px;
}

.user-info-email {
  flex: 0 0 40%;
  min-width: 0;
}

.user-info-subscription {
  flex: 1;
  min-width: 0;
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

.create-api-key-btn{
  font-size: 12px;
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

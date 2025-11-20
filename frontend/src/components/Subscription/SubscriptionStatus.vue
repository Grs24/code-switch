<template>
  <div v-if="loading" class="subscription-loading" :class="{ compact }">
    <div class="loading-spinner"></div>
    <span>{{ t('components.subscription.loading') }}</span>
  </div>
  
  <div v-else-if="error" class="subscription-error" :class="{ compact }">
    <svg viewBox="0 0 24 24" width="20" height="20">
      <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
      <path d="M12 8v4m0 4h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    </svg>
    <span>{{ error }}</span>
  </div>

  <div v-else-if="useInfo" class="subscription-status" :class="{ compact }">
    <div v-if="statusTip && !compact" class="status-tip" :class="statusTipClass">
      {{ statusTip }}
    </div>

    <div v-if="showExpireInfo || isExpired" class="expire-info" :class="{ compact, expired: isExpired }">
      <svg v-if="!compact" viewBox="0 0 24 24" width="16" height="16">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
        <path d="M12 6v6l4 2" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <span>{{ isExpired ? expiredInfoText : expireInfoText }}</span>
    </div>

    <BaseButton 
      :variant="buttonVariant" 
      @click="handleButtonClick"
      class="subscription-button"
      :class="{ compact }"
    >
      {{ buttonText }}
    </BaseButton>
  </div>

  <div v-else-if="!userInfo" class="subscription-empty" :class="{ compact }">
    <span>{{ t('components.subscription.pleaseLogin') }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Browser } from '@wailsio/runtime'
import BaseButton from '../common/BaseButton.vue'
import { fetchAiflexUseInfo } from '../../services/aiflex'
import { AiflexUseStatus, type GetAiflexUseInfoResponse } from '../../types/aiflex'

const { t } = useI18n()

interface Props {
  userInfo?: any
  platform?: string
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  platform: 'aicoding_0011',
  compact: false
})

const emit = defineEmits(['refresh'])

const loading = ref(false)
const error = ref<string | null>(null)
const useInfo = ref<GetAiflexUseInfoResponse | null>(null)

const remainDays = computed(() => {
  if (!useInfo.value?.expiration_time) return 0
  // expiration_time 是毫秒时间戳
  const now = Date.now()
  const expireDate = new Date(useInfo.value.expiration_time)
  const nowDate = new Date(now)
  
  const expireStartOfDay = new Date(expireDate.getFullYear(), expireDate.getMonth(), expireDate.getDate())
  const nowStartOfDay = new Date(nowDate.getFullYear(), nowDate.getMonth(), nowDate.getDate())
  
  return Math.ceil((expireStartOfDay.getTime() - nowStartOfDay.getTime()) / (24 * 3600 * 1000))
})

const expireDate = computed(() => {
  if (!useInfo.value?.expiration_time) return ''
  // expiration_time 是毫秒时间戳
  const date = new Date(useInfo.value.expiration_time)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}.${month}.${day}`
})

const isPlatformNewUser = computed(() => {
  return useInfo.value?.use_status === AiflexUseStatus.NOT_USED
})

const isExpired = computed(() => {
  return useInfo.value?.use_status === AiflexUseStatus.EXPIRED
})

const isInUse = computed(() => {
  return useInfo.value?.use_status === AiflexUseStatus.IN_USE
})

const buttonText = computed(() => {
  if (isPlatformNewUser.value) {
    return t('components.subscription.subscribe')
  }
  return t('components.subscription.renew')
})

const buttonVariant = computed(() => {
  if (isExpired.value) {
    return 'danger'
  }
  return 'primary'
})

const statusTip = computed(() => {
  if (isExpired.value) {
    return t('components.subscription.expiredTip')
  }
  return null
})

const statusTipClass = computed(() => {
  if (isExpired.value) {
    return 'tip-danger'
  }
  return ''
})

const showExpireInfo = computed(() => {
  return isInUse.value && remainDays.value >= 0
})

const expireInfoText = computed(() => {
  if (!isInUse.value) return ''
  return t('components.subscription.expireInfo', {
    days: remainDays.value,
    date: expireDate.value
  })
})

const expiredInfoText = computed(() => {
  if (!isExpired.value) return ''
  return t('components.subscription.expiredOn', {
    date: expireDate.value
  })
})

const loadSubscriptionInfo = async () => {
  if (!props.userInfo) {
    useInfo.value = null
    return
  }

  loading.value = true
  error.value = null

  try {
    const data = await fetchAiflexUseInfo({ platform: props.platform })
    useInfo.value = data
  } catch (err: any) {
    console.error('[SubscriptionStatus] 获取套餐信息失败:', err)
    error.value = t('components.subscription.loadError')
    useInfo.value = null
  } finally {
    loading.value = false
  }
}

const handleButtonClick = async () => {
  try {
    const { GetSubscriptionURL } = await import('../../services/auth')
    const url = await GetSubscriptionURL()
    Browser.OpenURL(url)
  } catch (error) {
    console.error('[SubscriptionStatus] 获取订阅链接失败:', error)
    Browser.OpenURL('https://0011.ai/subscribe')
  }
}

const refresh = () => {
  loadSubscriptionInfo()
}

// 监听用户信息变化
watch(() => props.userInfo, () => {
  loadSubscriptionInfo()
}, { immediate: false })

onMounted(() => {
  loadSubscriptionInfo()
})

defineExpose({
  refresh
})
</script>

<style scoped>
.subscription-loading,
.subscription-error,
.subscription-empty {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  border-radius: 8px;
  background: var(--mac-bg-secondary);
  font-size: 0.875rem;
  color: var(--mac-text-secondary);
}

.subscription-loading.compact,
.subscription-error.compact,
.subscription-empty.compact {
  padding: 8px 10px;
  font-size: 0.8125rem;
}

.subscription-error {
  color: rgb(239, 68, 68);
  background: rgba(239, 68, 68, 0.1);
}

.subscription-error svg {
  flex-shrink: 0;
  color: rgb(239, 68, 68);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(99, 102, 241, 0.3);
  border-top-color: rgba(99, 102, 241, 0.8);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.subscription-status {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  min-height: 46px;
}

.subscription-status.compact {
  gap: 6px;
}

.status-tip {
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  text-align: center;
}

.tip-danger {
  background: rgba(239, 68, 68, 0.1);
  color: rgb(239, 68, 68);
}

.expire-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: var(--mac-bg-secondary);
  border-radius: 6px;
  font-size: 14px;
  color: var(--mac-text-secondary);
  min-height: 46px;
}

.expire-info.compact {
  padding: 6px 8px;
  gap: 6px;
}

.expire-info svg {
  flex-shrink: 0;
  color: rgba(99, 102, 241, 0.8);
}

.expire-info.expired {
  background: rgba(239, 68, 68, 0.1);
  color: rgb(239, 68, 68);
}

.expire-info.expired svg {
  color: rgb(239, 68, 68);
}

.subscription-button {
  width: 100%;
}

.subscription-button.compact {
  padding: 4px 14px;
  font-size: 12px;
  min-height: auto;
  width: 60px;
  min-width: 48px;
}
</style>

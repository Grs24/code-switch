<template>
  <button 
    v-if="userInfo && subscriptionStatus"
    class="ghost-icon subscription-btn"
    :class="{ 'subscription-expired': subscriptionStatus.use_status === 2 }"
    :data-tooltip="tooltipText"
    @click="handleClick"
  >
    <span class="subscription-text">{{ buttonText }}</span>
  </button>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Browser } from '@wailsio/runtime'
import { fetchAiflexUseInfo } from '../../services/aiflex'
import { GetSubscriptionURL } from '../../services/auth'
import { AiflexUseStatus, type GetAiflexUseInfoResponse } from '../../types/aiflex'

const { t } = useI18n()

interface Props {
  userInfo?: any
}

const props = defineProps<Props>()

const subscriptionStatus = ref<GetAiflexUseInfoResponse | null>(null)
const loading = ref(false)

// 按钮文本
const buttonText = computed(() => {
  if (!subscriptionStatus.value) return ''
  const { use_status } = subscriptionStatus.value
  if (use_status === AiflexUseStatus.NOT_USED) {
    return t('components.subscription.subscribe')
  }
  return t('components.subscription.renew')
})

// Tooltip 文本
const tooltipText = computed(() => {
  if (!subscriptionStatus.value) return ''
  const { use_status, expiration_time } = subscriptionStatus.value
  
  if (use_status === AiflexUseStatus.NOT_USED) {
    return t('components.subscription.buyTip')
  }
  
  if (use_status === AiflexUseStatus.EXPIRED) {
    const date = new Date(expiration_time)
    const dateStr = `${date.getFullYear()}.${String(date.getMonth() + 1).padStart(2, '0')}.${String(date.getDate()).padStart(2, '0')}`
    return t('components.subscription.expiredOn', { date: dateStr })
  }
  
  if (use_status === AiflexUseStatus.IN_USE) {
    const now = Date.now()
    const expireDate = new Date(expiration_time)
    const nowDate = new Date(now)
    
    const expireStartOfDay = new Date(expireDate.getFullYear(), expireDate.getMonth(), expireDate.getDate())
    const nowStartOfDay = new Date(nowDate.getFullYear(), nowDate.getMonth(), nowDate.getDate())
    const remainDays = Math.ceil((expireStartOfDay.getTime() - nowStartOfDay.getTime()) / (24 * 3600 * 1000))
    
    const dateStr = `${expireDate.getFullYear()}.${String(expireDate.getMonth() + 1).padStart(2, '0')}.${String(expireDate.getDate()).padStart(2, '0')}`
    return t('components.subscription.expireInfo', { days: remainDays, date: dateStr })
  }
  
  return ''
})

// 加载套餐状态
const loadStatus = async () => {
  if (!props.userInfo) {
    subscriptionStatus.value = null
    return
  }
  
  loading.value = true
  try {
    const data = await fetchAiflexUseInfo({ platform: 'aicoding_0011' })
    subscriptionStatus.value = data
  } catch (error) {
    console.error('[SubscriptionButton] 获取套餐状态失败:', error)
    subscriptionStatus.value = null
  } finally {
    loading.value = false
  }
}

// 处理点击
const handleClick = async () => {
  try {
    const url = await GetSubscriptionURL()
    Browser.OpenURL(url)
  } catch (error) {
    console.error('[SubscriptionButton] 获取订阅链接失败:', error)
    Browser.OpenURL('https://0011.ai/subscribe')
  }
}

// 监听用户信息变化
watch(() => props.userInfo, () => {
  loadStatus()
}, { immediate: true })

// 暴露刷新方法
defineExpose({
  refresh: loadStatus
})
</script>

<style scoped>
</style>

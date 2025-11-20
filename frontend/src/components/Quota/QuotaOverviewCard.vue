<template>
  <section v-if="userInfo && quota" class="quota-overview-section">
    <div class="quota-overview-card" :class="statusClass">
      <div class="quota-header">
        <div class="quota-title-row">
          <div class="icon-wrapper">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path
                d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14z"
                fill="currentColor"
              />
              <path
                d="M7 10h2v7H7zm4-3h2v10h-2zm4 6h2v4h-2z"
                fill="currentColor"
              />
            </svg>
          </div>
          <h3 class="quota-title">
            {{ t('components.quota.overview.title') }}
          </h3>
          <span v-if="showWarning" class="warning-badge-inline">
            {{ warningText }}
          </span>
        </div>
        <div class="header-actions">
          <SubscriptionButton :user-info="userInfo" />
          <button class="refresh-btn" @click="handleRefresh" :disabled="loading" :class="{ spinning: loading }">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path
                d="M17.65 6.35A7.958 7.958 0 0012 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0112 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"
                fill="currentColor"
              />
            </svg>
          </button>
        </div>
      </div>

      <div class="quota-stats">
        <div class="stat-item stat-primary">
          <div class="stat-content">
            <span class="stat-label">{{ t('components.quota.overview.totalAmount') }}</span>
            <span class="stat-value">${{ formatAmount(totalAmount) }}</span>
          </div>
          <div class="stat-decoration"></div>
        </div>
        
        <div class="stat-item">
          <span class="stat-label">{{ t('components.quota.overview.dailyQuota') }}</span>
          <span class="stat-value">${{ formatAmount(dailyQuota) }}</span>
        </div>
        
        <div class="stat-item stat-highlight">
          <span class="stat-label">{{ t('components.quota.overview.usedToday') }}</span>
          <span class="stat-value stat-used">${{ formatAmount(usedToday) }}</span>
        </div>
        
        <div class="stat-item stat-highlight">
          <span class="stat-label">{{ t('components.quota.overview.remaining') }}</span>
          <span class="stat-value stat-remaining">${{ formatAmount(remaining) }}</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Events } from '@wailsio/runtime'
import { fetchQuotaInfo, refreshQuota, type QuotaInfo } from '../../services/service0011'
import SubscriptionButton from '../Subscription/SubscriptionButton.vue'

const { t } = useI18n()

interface Props {
  userInfo?: any
}

const props = defineProps<Props>()

const loading = ref(false)
const quota = ref<QuotaInfo | null>(null)

// 计算金额相关数据
const totalAmount = computed(() => {
  if (!quota.value) return 0
  return quota.value.totalAmount
})

const dailyQuota = computed(() => {
  if (!quota.value) return 0
  return quota.value.dailyQuota
})

const usedToday = computed(() => {
  if (!quota.value) return 0
  return quota.value.usedToday
})

const remaining = computed(() => {
  if (!quota.value) return 0
  return quota.value.remaining
})

const statusClass = computed(() => {
  if (!quota.value) return ''
  if (quota.value.isCriticalQuota) return 'status-critical'
  if (quota.value.isLowQuota) return 'status-warning'
  return 'status-normal'
})

const showWarning = computed(() => {
  if (!quota.value) return false
  return quota.value.isLowQuota || quota.value.isCriticalQuota
})

const warningText = computed(() => {
  if (!quota.value) return ''
  if (quota.value.isCriticalQuota) return t('components.quota.overview.criticalLow')
  if (quota.value.isLowQuota) return t('components.quota.overview.low')
  return ''
})

const formatAmount = (amount: number): string => {
  return amount.toFixed(2)
}

const loadQuota = async () => {
  if (!props.userInfo) {
    quota.value = null
    return
  }

  try {
    const data = await fetchQuotaInfo()
    quota.value = data
  } catch (error) {
    console.error('[QuotaOverviewCard] 获取额度信息失败:', error)
    quota.value = null
  }
}

const handleRefresh = async () => {
  loading.value = true
  try {
    const data = await refreshQuota()
    quota.value = data
  } catch (error) {
    console.error('[QuotaOverviewCard] 刷新额度失败:', error)
  } finally {
    loading.value = false
  }
}

// 监听用户信息变化
watch(() => props.userInfo, () => {
  loadQuota()
}, { immediate: true })

// 监听登录事件
let removeLoginListener: (() => void) | undefined
let removeLogoutListener: (() => void) | undefined

onMounted(() => {
  removeLoginListener = Events.On('auth:login-success', () => {
    loadQuota()
  })

  removeLogoutListener = Events.On('auth:logout', () => {
    quota.value = null
  })
})

onUnmounted(() => {
  removeLoginListener?.()
  removeLogoutListener?.()
})
</script>

<style scoped>
.quota-overview-section {
  margin: 24px 0;
}

.quota-overview-card {
  background: var(--mac-bg-secondary);
  border: 1px solid var(--mac-border);
  border-radius: 16px;
  padding: 24px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.quota-overview-card:hover {
  border-color: var(--mac-accent);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 0 0 1px var(--mac-accent);
  transform: translateY(-2px);
}

.quota-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.quota-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon-wrapper {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--mac-accent), var(--success-color, #34c759));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.2);
}

.icon-wrapper svg {
  width: 20px;
  height: 20px;
  color: white;
}

.quota-title {
  margin: 0;
  font-size: 17px;
  font-weight: 600;
  color: var(--mac-text-primary);
  letter-spacing: -0.3px;
}

.warning-badge-inline {
  padding: 4px 10px;
  font-size: 11px;
  font-weight: 600;
  border-radius: 8px;
  background: var(--warning-bg);
  color: var(--warning-text);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.refresh-btn {
  padding: 8px;
  background: var(--mac-bg-secondary);
  border: 1px solid var(--mac-border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.refresh-btn:hover:not(:disabled) {
  background: var(--mac-bg-tertiary);
  border-color: var(--mac-accent);
  transform: scale(1.05);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn.spinning svg {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.refresh-btn svg {
  width: 18px;
  height: 18px;
  color: var(--mac-text-secondary);
  transition: color 0.2s ease;
}

.refresh-btn:hover:not(:disabled) svg {
  color: var(--mac-accent);
}

.quota-stats {
  display: grid;
  grid-template-columns: 1.8fr 1fr 1fr 1fr;
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px;
  border-radius: 12px;
  background: var(--mac-bg-secondary);
  transition: all 0.2s ease;
}

.stat-item:hover {
  background: var(--mac-bg-tertiary);
  transform: translateY(-2px);
}

.stat-primary {
  position: relative;
  background: linear-gradient(135deg, rgba(0, 122, 255, 0.08), rgba(52, 199, 89, 0.08));
  border: 1px solid rgba(0, 122, 255, 0.2);
}

.stat-primary .stat-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-decoration {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 60px;
  height: 60px;
  background: radial-gradient(circle at center, rgba(0, 122, 255, 0.15), transparent);
  border-radius: 50%;
  pointer-events: none;
}

.stat-highlight {
  transition: all 0.2s ease;
}

.stat-label {
  font-size: 11px;
  color: var(--mac-text-tertiary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--mac-text-primary);
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.stat-primary .stat-value {
  font-size: 32px;
  background: linear-gradient(135deg, var(--mac-accent), var(--success-color, #34c759));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-used {
  color: var(--warning-color, #ff9500);
}

.stat-remaining {
  color: var(--success-color, #34c759);
}

/* 状态颜色 */
.status-normal {
  --warning-color: #ff9500;
  --success-color: #34c759;
}

.status-warning {
  border-color: rgba(255, 149, 0, 0.4);
  --warning-bg: rgba(255, 149, 0, 0.15);
  --warning-text: #ff9500;
  --warning-color: #ff9500;
}

.status-warning:hover {
  box-shadow: 0 8px 24px rgba(255, 149, 0, 0.2), 0 0 0 1px rgba(255, 149, 0, 0.4);
}

.status-critical {
  border-color: rgba(255, 59, 48, 0.4);
  --warning-bg: rgba(255, 59, 48, 0.15);
  --warning-text: #ff3b30;
  --warning-color: #ff3b30;
}

.status-critical:hover {
  box-shadow: 0 8px 24px rgba(255, 59, 48, 0.2), 0 0 0 1px rgba(255, 59, 48, 0.4);
}

/* 暗色模式 */
:global(.dark) .quota-overview-card {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.12);
}

:global(.dark) .quota-overview-card:hover {
  border-color: var(--mac-accent);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 0 1px var(--mac-accent);
}

:global(.dark) .icon-wrapper {
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.3);
}

:global(.dark) .stat-item {
  background: rgba(255, 255, 255, 0.04);
}

:global(.dark) .stat-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

:global(.dark) .stat-primary {
  background: linear-gradient(135deg, rgba(0, 122, 255, 0.12), rgba(52, 199, 89, 0.12));
  border-color: rgba(0, 122, 255, 0.3);
}

:global(.dark) .refresh-btn {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.12);
}

:global(.dark) .refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
}
</style>

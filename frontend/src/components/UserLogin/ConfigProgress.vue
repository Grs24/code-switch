<template>
  <div class="config-progress">
    <div class="progress-title">配置进度</div>
    <div class="progress-steps">
      <div 
        v-for="(step, index) in steps" 
        :key="step.id"
        class="progress-step"
        :class="{ 
          'step-success': step.status === 'success',
          'step-error': step.status === 'error',
          'step-running': step.status === 'running',
          'step-pending': step.status === 'pending'
        }"
      >
        <div class="step-icon">
          <svg v-if="step.status === 'success'" viewBox="0 0 24 24" width="20" height="20">
            <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else-if="step.status === 'error'" viewBox="0 0 24 24" width="20" height="20">
            <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else-if="step.status === 'running'" class="spinner" viewBox="0 0 24 24" width="20" height="20">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" fill="none" opacity="0.25"/>
            <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round"/>
          </svg>
          <div v-else class="step-number">{{ index + 1 }}</div>
        </div>
        <div class="step-content">
          <div class="step-title">{{ step.title }}</div>
          <div v-if="step.description" class="step-description">{{ step.description }}</div>
          <div v-if="step.error" class="step-error-msg">{{ step.error }}</div>
          <div v-if="step.details && step.details.length > 0" class="step-details">
            <div v-for="(detail, idx) in step.details" :key="idx" class="detail-item">
              {{ detail }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Step {
  id: string
  title: string
  description: string
  status: 'pending' | 'running' | 'success' | 'error'
  error?: string
  details?: string[]
}

interface Props {
  steps: Step[]
}

defineProps<Props>()
</script>

<style scoped>
.config-progress {
  margin: 16px 0;
  padding: 16px;
  background: var(--mac-bg-secondary);
  border-radius: 8px;
}

.progress-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--mac-text-primary);
  margin-bottom: 12px;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.progress-step {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.step-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--mac-surface);
  color: var(--mac-text-secondary);
}

.step-success .step-icon {
  background: rgba(34, 197, 94, 0.1);
  color: rgb(34, 197, 94);
}

.step-error .step-icon {
  background: rgba(239, 68, 68, 0.1);
  color: rgb(239, 68, 68);
}

.step-running .step-icon {
  background: rgba(99, 102, 241, 0.1);
  color: rgba(99, 102, 241, 0.8);
}

.step-number {
  font-size: 0.875rem;
  font-weight: 600;
}

.spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.step-content {
  flex: 1;
  padding-top: 4px;
}

.step-title {
  font-weight: 500;
  font-size: 0.875rem;
  color: var(--mac-text-primary);
  margin-bottom: 2px;
}

.step-description {
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
}

.step-error-msg {
  font-size: 0.8125rem;
  color: rgb(239, 68, 68);
  margin-top: 4px;
}

.step-details {
  margin-top: 6px;
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
}

.detail-item {
  padding: 4px 0;
}
</style>

<template>
  <div class="user-info-section auto-config-section">
    <label class="user-info-label">一键配置 AICoding</label>
    <div class="auto-config-description">
      让 Claude Code 和 Codex 在你的电脑上直接可用
    </div>

    <!-- API Key 配置 -->
    <div class="api-key-config">
      <label class="config-label">API Key</label>
      <div class="url-input-group">
        <input 
          v-model="customApiKey"
          :disabled="!isEditingApiKey"
          :type="showApiKeyValue ? 'text' : 'password'"
          class="url-input"
          :class="{ 'editing': isEditingApiKey }"
          placeholder="sk-..."
        />
        <button 
          v-if="customApiKey"
          class="toggle-visibility-btn"
          @click="showApiKeyValue = !showApiKeyValue"
          :title="showApiKeyValue ? '隐藏' : '显示'"
        >
          <svg v-if="showApiKeyValue" viewBox="0 0 24 24" width="16" height="16">
            <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <line x1="1" y1="1" x2="23" y2="23" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" width="16" height="16">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
        </button>
        <button 
          class="edit-btn"
          @click="toggleEditApiKey"
          :title="isEditingApiKey ? '保存' : '编辑'"
        >
          <svg v-if="!isEditingApiKey" viewBox="0 0 24 24" width="16" height="16">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" width="16" height="16">
            <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
      <div class="url-hint">
        从账户中获取的 API Key，用于配置工具
      </div>
    </div>

    <!-- Base URL 配置 -->
    <div class="base-url-config">
      <label class="config-label">API Base URL</label>
      <div class="url-input-group">
        <input 
          v-model="customBaseUrl"
          :disabled="!isEditingBaseUrl"
          type="text"
          class="url-input"
          :class="{ 'editing': isEditingBaseUrl }"
          :placeholder="DEFAULT_BASE_URL"
        />
        <button 
          v-if="isEditingBaseUrl && customBaseUrl !== DEFAULT_BASE_URL"
          class="reset-btn"
          @click="resetBaseUrl"
          title="重置为默认值"
        >
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M21 3v5h-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M3 21v-5h5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <button 
          class="edit-btn"
          @click="toggleEditBaseUrl"
          :title="isEditingBaseUrl ? '保存' : '编辑'"
        >
          <svg v-if="!isEditingBaseUrl" viewBox="0 0 24 24" width="16" height="16">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" width="16" height="16">
            <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </div>
    
    <!-- 配置按钮 -->
    <div class="config-buttons">
      <BaseButton 
        variant="primary" 
        @click="handleAutoConfig"
        :disabled="autoConfigLoading"
        class="config-btn config-btn-full"
      >
        <span v-if="!autoConfigLoading">一键配置</span>
        <span v-else>配置中，请稍候...</span>
      </BaseButton>
    </div>

    <!-- 配置进度展示 -->
    <ConfigProgress 
      v-if="autoConfigLoading && configSteps.length > 0"
      :steps="configSteps"
    />
    
    <!-- 配置结果显示 -->
    <ConfigResult
      v-if="autoConfigResult && !autoConfigLoading"
      :result="autoConfigResult"
      :api-key="customApiKey"
      :custom-base-url="customBaseUrl"
      :default-base-url="DEFAULT_BASE_URL"
      :verification-info="verificationInfo"
      :show-detailed-info="showDetailedInfo"
      @show-details="showDetailedInfo = !showDetailedInfo"
      @copy-command="copyCommand"
      @retry="handleAutoConfig"
      @open-tutorial="openTutorial"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import BaseButton from '../common/BaseButton.vue'
import ConfigProgress from './ConfigProgress.vue'
import ConfigResult from './ConfigResult.vue'
import {
  ConfigureAll,
  CheckDependencies,
  CheckConfigurationStatus,
  GetDefaultBaseURL,
  GetToolVersions,
  type AutoConfigResult,
  type DependencyCheckResult
} from '../../services/autoConfig'

interface Props {
  apiKeyData: { id: number; token: string; status: number }
  userInfo: any
}

const props = defineProps<Props>()

// 默认 Base URL
const DEFAULT_BASE_URL = 'https://aicoding.2233.ai'

// 状态
const customApiKey = ref('')
const isEditingApiKey = ref(false)
const showApiKeyValue = ref(false)
const customBaseUrl = ref('')
const isEditingBaseUrl = ref(false)
const autoConfigLoading = ref(false)
const autoConfigResult = ref<AutoConfigResult | null>(null)
const dependencyCheckResult = ref<DependencyCheckResult | null>(null)
const showDetailedInfo = ref(false)
const configSteps = ref<Array<{
  id: string
  title: string
  description: string
  status: 'pending' | 'running' | 'success' | 'error'
  error?: string
  details?: string[]
}>>([])
const verificationInfo = ref({
  claudeVersion: '',
  codexVersion: ''
})

// 计算属性
const apiKey = computed(() => customApiKey.value || props.apiKeyData?.token || '')

const defaultDomain = computed(() => {
  try {
    const url = new URL(DEFAULT_BASE_URL)
    return url.hostname
  } catch {
    return 'aicoding.2233.ai'
  }
})

// 初始化
onMounted(async () => {
  customApiKey.value = props.apiKeyData?.token || ''
  await loadDefaultBaseURL()
  await checkExistingConfiguration()
})

// 监听 apiKeyData 变化
watch(() => props.apiKeyData, (newData) => {
  if (newData) {
    customApiKey.value = newData.token || ''
  }
}, { deep: true })

// 加载默认 Base URL
const loadDefaultBaseURL = async () => {
  try {
    const defaultURL = await GetDefaultBaseURL()
    if (defaultURL && !customBaseUrl.value) {
      customBaseUrl.value = defaultURL
    }
  } catch {
    if (!customBaseUrl.value) {
      customBaseUrl.value = DEFAULT_BASE_URL
    }
  }
}

// 检查现有配置
const checkExistingConfiguration = async () => {
  try {
    const status = await CheckConfigurationStatus()
    if (status.success) {
      autoConfigResult.value = status
      await detectToolVersions()
    }
  } catch (error) {
    console.error('检查配置状态失败:', error)
  }
}

// 检测工具版本
const detectToolVersions = async () => {
  verificationInfo.value = {
    claudeVersion: '',
    codexVersion: ''
  }
  
  try {
    const versions = await GetToolVersions()
    verificationInfo.value.claudeVersion = versions.claude || ''
    verificationInfo.value.codexVersion = versions.codex || ''
  } catch {
    // 忽略版本检测失败
  }
}

// 初始化配置步骤
const initConfigSteps = () => {
  configSteps.value = [
    {
      id: 'check-nodejs',
      title: '检查 Node.js',
      description: '验证 Node.js 是否已安装',
      status: 'pending'
    },
    {
      id: 'configure-claude',
      title: '配置 Claude Code',
      description: '创建配置文件和环境变量',
      status: 'pending'
    },
    {
      id: 'configure-codex',
      title: '配置 Codex',
      description: '创建配置文件',
      status: 'pending'
    },
    {
      id: 'verify',
      title: '验证配置',
      description: '确认配置是否正确',
      status: 'pending'
    }
  ]
}

// 更新步骤状态
const updateStepStatus = (
  stepId: string,
  status: 'pending' | 'running' | 'success' | 'error',
  error?: string,
  details?: string[]
) => {
  const step = configSteps.value.find(s => s.id === stepId)
  if (step) {
    step.status = status
    if (error) step.error = error
    if (details) step.details = details
  }
}

// 一键配置
const handleAutoConfig = async () => {
  if (!props.userInfo || !apiKey.value) return

  autoConfigLoading.value = true
  autoConfigResult.value = null
  initConfigSteps()

  try {
    // 步骤 1: 检查 Node.js
    updateStepStatus('check-nodejs', 'running')
    const depCheck = await CheckDependencies()
    
    if (!depCheck.passed) {
      updateStepStatus('check-nodejs', 'error', depCheck.errors.join(', '))
      autoConfigResult.value = {
        success: false,
        claudeStatus: '未配置',
        codexStatus: '未配置',
        errors: depCheck.errors,
        warnings: [],
        os: '',
        configPaths: {},
        verificationCommands: {}
      }
      return
    }
    
    updateStepStatus('check-nodejs', 'success', undefined, [
      `Node.js ${depCheck.nodejs.version} 已安装`
    ])

    // 步骤 2-3: 配置
    updateStepStatus('configure-claude', 'running')
    updateStepStatus('configure-codex', 'running')
    await new Promise(resolve => setTimeout(resolve, 300))
    
    const result = await ConfigureAll(apiKey.value, customBaseUrl.value)
    
    if (result.claudeStatus.includes('成功')) {
      updateStepStatus('configure-claude', 'success', undefined, result.configPaths.claude || [])
    } else {
      updateStepStatus('configure-claude', 'error', result.claudeStatus)
    }
    
    if (result.codexStatus.includes('成功')) {
      updateStepStatus('configure-codex', 'success', undefined, result.configPaths.codex || [])
    } else {
      updateStepStatus('configure-codex', 'error', result.codexStatus)
    }

    // 步骤 4: 验证
    updateStepStatus('verify', 'running')
    await new Promise(resolve => setTimeout(resolve, 200))
    
    if (result.success) {
      updateStepStatus('verify', 'success', undefined, ['所有配置已验证通过'])
      detectToolVersions()
    } else {
      updateStepStatus('verify', 'error', '部分配置失败')
    }
    
    autoConfigResult.value = result
  } catch (error) {
    console.error('自动配置出错:', error)
    configSteps.value.forEach(step => {
      if (step.status === 'running' || step.status === 'pending') {
        step.status = 'error'
        step.error = String(error)
      }
    })
    autoConfigResult.value = {
      success: false,
      claudeStatus: '配置失败',
      codexStatus: '配置失败',
      errors: [String(error)],
      warnings: [],
      os: '',
      configPaths: {},
      verificationCommands: {}
    }
  } finally {
    autoConfigLoading.value = false
  }
}

// 切换 API Key 编辑状态
const toggleEditApiKey = () => {
  if (isEditingApiKey.value) {
    const key = customApiKey.value.trim()
    if (!key) {
      customApiKey.value = props.apiKeyData?.token || ''
    }
    showApiKeyValue.value = false
  }
  isEditingApiKey.value = !isEditingApiKey.value
}

// 切换 Base URL 编辑状态
const toggleEditBaseUrl = async () => {
  if (isEditingBaseUrl.value) {
    const url = customBaseUrl.value.trim()
    if (!url) {
      await loadDefaultBaseURL()
    } else if (!url.startsWith('http://') && !url.startsWith('https://')) {
      customBaseUrl.value = `https://${url}`
    }
  }
  isEditingBaseUrl.value = !isEditingBaseUrl.value
}

// 重置 Base URL
const resetBaseUrl = async () => {
  try {
    const defaultURL = await GetDefaultBaseURL()
    customBaseUrl.value = defaultURL || DEFAULT_BASE_URL
  } catch {
    customBaseUrl.value = DEFAULT_BASE_URL
  }
}

// 复制命令
const copyCommand = async (command: string) => {
  try {
    await navigator.clipboard.writeText(command)
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 打开教程
const openTutorial = () => {
  const baseUrl = customBaseUrl.value || DEFAULT_BASE_URL
  const tutorialUrl = `${baseUrl}/claude-code/tutorial`
  window.open(tutorialUrl, '_blank')
}

// 暴露给父组件的方法：触发自动配置
const triggerAutoConfig = async () => {
  await handleAutoConfig()
}

// 暴露方法给父组件
defineExpose({
  triggerAutoConfig
})
</script>

<style scoped>
.auto-config-section {
  padding-top: 16px;
  border-top: 1px solid var(--mac-border);
}

.auto-config-description {
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
  margin-top: -4px;
  margin-bottom: 12px;
}

.api-key-config,
.base-url-config {
  margin-bottom: 16px;
}

.config-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--mac-text-secondary);
  margin-bottom: 8px;
}

.url-input-group {
  display: flex;
  gap: 8px;
  align-items: center;
}

.url-input {
  flex: 1;
  padding: 10px 12px;
  background: var(--mac-bg-secondary);
  border: 1px solid var(--mac-border);
  border-radius: 6px;
  font-size: 0.875rem;
  color: var(--mac-text-primary);
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  transition: all 0.2s ease;
}

.url-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.url-input.editing {
  border-color: rgba(99, 102, 241, 0.5);
  background: var(--mac-surface);
}

.url-input:focus {
  outline: none;
  border-color: rgba(99, 102, 241, 0.5);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.edit-btn,
.reset-btn,
.toggle-visibility-btn {
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

.edit-btn:hover,
.reset-btn:hover,
.toggle-visibility-btn:hover {
  background: var(--mac-surface);
  border-color: rgba(99, 102, 241, 0.3);
  color: rgba(99, 102, 241, 0.8);
}

.edit-btn:active,
.reset-btn:active,
.toggle-visibility-btn:active {
  transform: scale(0.95);
}

.reset-btn:hover {
  border-color: rgba(245, 158, 11, 0.3);
  color: rgba(245, 158, 11, 0.8);
}

.url-hint {
  margin-top: 6px;
  font-size: 0.75rem;
  color: var(--mac-text-secondary);
  line-height: 1.4;
}

.config-buttons {
  display: flex;
  gap: 10px;
}

.config-btn {
  height: 56px !important;
  font-size: 18px !important;
  font-weight: 500 !important;
  border-radius: 40px !important;
  transition: all 0.2s ease !important;
}

.config-btn-full {
  width: 100%;
}
</style>

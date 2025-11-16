<template>
  <div class="auto-config-result" :class="{ 'success': result.success, 'error': !result.success }">
    <!-- 成功状态 -->
    <div v-if="result.success" class="result-success">
      <div class="success-icon">
        <svg viewBox="0 0 24 24" width="48" height="48" style="color: #10b981;">
          <circle cx="12" cy="12" r="11" fill="rgba(16, 185, 129, 0.1)" stroke="currentColor" stroke-width="2"/>
          <path d="M8 12l3 3 5-5" stroke="currentColor" stroke-width="2.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
      <div class="success-title">配置成功！🎉</div>
      <div class="success-message">
        Claude Code 和 Codex 已经可以使用了
      </div>
      
      <!-- 快速测试按钮 -->
      <div class="quick-actions">
        <BaseButton 
          variant="outline" 
          @click="$emit('show-details')"
          class="action-btn"
        >
          {{ showDetailedInfo ? '隐藏详情' : '📋 查看详情' }}
        </BaseButton>
      </div>

      <!-- 详细信息（可折叠） -->
      <div v-if="showDetailedInfo" class="detailed-info">
        <div class="info-section">
          <div class="info-label">📁 配置文件位置：</div>
          <div v-for="(paths, tool) in result.configPaths" :key="tool" class="info-content">
            <div v-for="(path, idx) in paths" :key="idx" class="path-item">{{ path }}</div>
          </div>
        </div>
        
        <div class="info-section">
          <div class="info-label">🔍 验证信息：</div>
          <div class="verification-output">
            <VerificationBlock
              title="Claude Code 版本："
              command="claude --version"
              :result="verificationInfo.claudeVersion"
              @copy="$emit('copy-command', $event)"
            />

            <VerificationBlock
              title="环境变量配置："
              command="echo $ANTHROPIC_BASE_URL"
              :result="customBaseUrl || defaultBaseUrl"
              @copy="$emit('copy-command', $event)"
            />

            <VerificationBlock
              command="echo $ANTHROPIC_AUTH_TOKEN"
              :result="maskApiKey(apiKey)"
              @copy="$emit('copy-command', $event)"
            />

            <VerificationBlock
              title="Codex 版本："
              command="codex -V"
              :result="verificationInfo.codexVersion"
              @copy="$emit('copy-command', $event)"
            />

            <VerificationBlock
              title="Codex 配置："
              command="cat ~/.codex/config.toml | grep base_url"
              :result="`base_url = &quot;${customBaseUrl || defaultBaseUrl}&quot;`"
              @copy="$emit('copy-command', $event)"
            />

            <div class="verification-tip">
              💡 提示：环境变量已写入 ~/.zshrc 或 ~/.bashrc<br/>
              请重新打开终端，或运行 <code class="inline-code">source ~/.zshrc</code> 使配置生效
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 失败状态 -->
    <div v-else class="result-failure">
      <div class="failure-icon">
        <svg viewBox="0 0 24 24" width="48" height="48" style="color: #ef4444;">
          <circle cx="12" cy="12" r="11" fill="rgba(239, 68, 68, 0.1)" stroke="currentColor" stroke-width="2"/>
          <path d="M12 8v4m0 4h.01" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
        </svg>
      </div>
      <div class="failure-title">配置遇到问题</div>
      
      <!-- 检查是否是 Node.js 问题 -->
      <div v-if="result.errors.some(e => e.includes('Node.js'))" class="failure-reason">
        <div class="reason-title">❌ 缺少 Node.js</div>
        <div class="reason-text">
          Claude Code 和 Codex 需要 Node.js 才能运行。
        </div>
        <div class="solution-steps">
          <div class="step-title">📝 安装步骤：</div>
          <div class="step-item">
            <strong>1. 下载 Node.js</strong><br/>
            访问 <a href="https://nodejs.org" target="_blank" class="link">nodejs.org</a>，下载 LTS 版本
          </div>
          <div class="step-item">
            <strong>2. 安装 Node.js</strong><br/>
            双击安装包，按照提示完成安装
          </div>
          <div class="step-item">
            <strong>3. 验证安装</strong><br/>
            打开终端，运行命令：<code class="inline-code">node --version</code>
          </div>
          <div class="step-item">
            <strong>4. 重新配置</strong><br/>
            安装完成后，点击下方"重试配置"按钮
          </div>
        </div>
      </div>

      <!-- 其他错误 -->
      <div v-else class="failure-reason">
        <div class="reason-title">错误详情：</div>
        <div v-for="(error, index) in result.errors" :key="index" class="error-text">
          • {{ error }}
        </div>
        <div class="solution-steps">
          <div class="step-title">💡 解决建议：</div>
          <div class="step-item">
            <a href="#" @click.prevent="$emit('open-tutorial')" class="link">📖 查看详细配置教程</a>
            （包含完整的手动配置步骤）
          </div>
        </div>
      </div>

      <!-- 重试按钮 -->
      <div class="quick-actions">
        <BaseButton 
          variant="primary" 
          @click="$emit('retry')"
          class="action-btn"
        >
          🔄 重试配置
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '../common/BaseButton.vue'
import VerificationBlock from './VerificationBlock.vue'

interface Props {
  result: any
  apiKey: string
  customBaseUrl: string
  defaultBaseUrl: string
  verificationInfo: {
    claudeVersion: string
    codexVersion: string
  }
  showDetailedInfo?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showDetailedInfo: false
})
defineEmits(['show-details', 'copy-command', 'retry', 'open-tutorial'])

const maskApiKey = (key: string | undefined) => {
  if (!key || key.length < 13) return key || 'sk-****'
  return `${key.slice(0, 8)}****${key.slice(-4)}`
}
</script>

<style scoped>
.auto-config-result {
  margin-top: 16px;
}

.result-success,
.result-failure {
  text-align: center;
  padding: 24px;
  border-radius: 12px;
}

.result-success {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(16, 185, 129, 0.1) 100%);
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.result-failure {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.success-icon,
.failure-icon {
  margin-bottom: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.success-title,
.failure-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--mac-text-primary);
  margin-bottom: 8px;
}

.success-message {
  font-size: 0.875rem;
  color: var(--mac-text-secondary);
  margin-bottom: 16px;
}

.quick-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 16px;
}

.action-btn {
  min-width: 120px;
}

.detailed-info {
  margin-top: 20px;
  text-align: left;
  padding: 16px;
  background: var(--mac-surface);
  border-radius: 8px;
}

.info-section {
  margin-bottom: 20px;
}

.info-section:last-child {
  margin-bottom: 0;
}

.info-label {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--mac-text-primary);
  margin-bottom: 8px;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.path-item {
  padding: 8px 12px;
  background: var(--mac-bg-secondary);
  border-radius: 6px;
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
}

.verification-output {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.verification-tip {
  margin-top: 16px;
  padding: 12px;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 6px;
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
  line-height: 1.6;
}

.inline-code {
  padding: 2px 6px;
  background: var(--mac-bg-secondary);
  border-radius: 3px;
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-size: 0.75rem;
}

.failure-reason {
  text-align: left;
  margin-top: 16px;
}

.reason-title {
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--mac-text-primary);
  margin-bottom: 8px;
}

.reason-text {
  font-size: 0.875rem;
  color: var(--mac-text-secondary);
  margin-bottom: 12px;
}

.solution-steps {
  margin-top: 16px;
}

.step-title {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--mac-text-primary);
  margin-bottom: 8px;
}

.step-item {
  padding: 12px;
  background: var(--mac-surface);
  border-radius: 6px;
  margin-bottom: 8px;
  font-size: 0.8125rem;
  color: var(--mac-text-secondary);
  line-height: 1.6;
}

.error-text {
  padding: 8px 12px;
  background: var(--mac-surface);
  border-radius: 6px;
  margin-bottom: 6px;
  font-size: 0.8125rem;
  color: rgb(239, 68, 68);
}

.link {
  color: rgba(99, 102, 241, 0.8);
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}
</style>

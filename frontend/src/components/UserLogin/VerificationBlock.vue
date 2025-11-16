<template>
  <div class="output-block">
    <div v-if="title" class="output-title">{{ title }}</div>
    <div class="terminal-output-wrapper">
      <div class="terminal-output">
        <div class="terminal-line">
          <span class="prompt">$ </span>
          <span class="command">{{ command }}</span>
        </div>
        <div class="terminal-result" :class="{ 'loading': !result }">
          <span v-if="result">{{ result }}</span>
          <span v-else class="loading-text">检测中...</span>
        </div>
      </div>
      <button class="copy-btn" @click="$emit('copy', command)" title="复制命令">
        <svg viewBox="0 0 24 24" width="14" height="14">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
  command: string
  result?: string
}

defineProps<Props>()
defineEmits(['copy'])
</script>

<style scoped>
.output-block {
  margin-bottom: 12px;
}

.output-title {
  font-weight: 500;
  font-size: 0.8125rem;
  color: var(--mac-text-primary);
  margin-bottom: 6px;
}

.terminal-output-wrapper {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.terminal-output {
  flex: 1;
  padding: 12px;
  background: var(--mac-bg-secondary);
  border-radius: 6px;
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-size: 0.8125rem;
}

.terminal-line {
  color: var(--mac-text-secondary);
  margin-bottom: 6px;
}

.prompt {
  color: rgba(99, 102, 241, 0.8);
  font-weight: 600;
}

.command {
  color: var(--mac-text-primary);
}

.terminal-result {
  color: rgb(34, 197, 94);
  word-break: break-all;
}

.terminal-result.loading {
  color: var(--mac-text-secondary);
  font-style: italic;
}

.loading-text {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.copy-btn {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
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

.copy-btn:hover {
  background: var(--mac-surface);
  border-color: rgba(99, 102, 241, 0.3);
  color: rgba(99, 102, 241, 0.8);
}

.copy-btn:active {
  transform: scale(0.95);
}
</style>

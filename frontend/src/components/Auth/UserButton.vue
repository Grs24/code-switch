<template>
  <div>
  <button 
    :class="['ghost-icon', { 'user-avatar-btn': userInfo }]"
    :data-tooltip="userInfo ? t('components.main.controls.profile') : t('components.main.controls.login')"
    @click="handleClick"
  >
    <div v-if="userInfo" class="user-avatar" :style="{ background: userAvatarGradient }">
      <span class="user-initial">{{ userInitial }}</span>
    </div>
    <svg v-else viewBox="0 0 24 24" aria-hidden="true">
      <circle cx="12" cy="8" r="4" stroke="currentColor" stroke-width="1.5" fill="none" />
      <path
        d="M4 20c0-4 3.5-6 8-6s8 2 8 6"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
      />
    </svg>
  </button>

  <!-- 用户信息弹窗 -->
  <UserLoginModal
    :open="showUserModal"
    :user-info="userInfo"
    :logout-loading="logoutLoading"
    @close="closeUserModal"
    @logout="handleLogoutFromModal"
  />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import UserLoginModal from '../UserLogin/UserLoginModal.vue'

const { t } = useI18n()
const router = useRouter()

interface Props {
  userInfo?: any
}

const props = defineProps<Props>()
const emit = defineEmits(['logout'])

const showUserModal = ref(false)
const logoutLoading = ref(false)

// 用户头像渐变色
const gradientColors = [
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
  'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
  'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
  'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)',
  'linear-gradient(135deg, #ff6e7f 0%, #bfe9ff 100%)',
]

// 根据邮箱生成一致的渐变色
const userAvatarGradient = computed(() => {
  const email = props.userInfo?.email || ''
  if (!email) return gradientColors[0]
  
  let hash = 0
  for (let i = 0; i < email.length; i++) {
    hash = email.charCodeAt(i) + ((hash << 5) - hash)
  }
  const index = Math.abs(hash) % gradientColors.length
  return gradientColors[index]
})

// 用户首字母
const userInitial = computed(() => {
  if (!props.userInfo) return 'U'
  const username = props.userInfo.username || props.userInfo.email || ''
  return username.charAt(0).toUpperCase() || 'U'
})

const handleClick = () => {
  if (props.userInfo) {
    showUserModal.value = true
  } else {
    router.push('/login')
  }
}

const closeUserModal = () => {
  showUserModal.value = false
}

const handleLogoutFromModal = async () => {
  logoutLoading.value = true
  try {
    emit('logout')
  } finally {
    logoutLoading.value = false
    showUserModal.value = false
  }
}
</script>

<style scoped>
.user-avatar {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 12px;
}

.user-avatar-btn {
  padding: 0;
}
</style>

<template>
  <div v-if="isLogin" class="user-profile">
    <div class="user-avatar">
      <img v-if="userInfo?.avatar" :src="userInfo.avatar" alt="avatar" />
      <div v-else class="avatar-placeholder">
        {{ userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
      </div>
    </div>
    <div class="user-info">
      <div class="username">{{ userInfo?.username || userInfo?.email }}</div>
      <button @click="handleLogout" class="logout-btn">
        {{ t('components.main.user.logout') }}
      </button>
    </div>
  </div>
  <button v-else @click="goToLogin" class="login-btn">
    {{ t('app.login') }}
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { IsLogin, GetUserInfo, Logout, type UserInfo } from '../../services/auth'

const { t } = useI18n()
const router = useRouter()
const isLogin = ref(false)
const userInfo = ref<UserInfo | null>(null)

onMounted(async () => {
  await checkLoginStatus()
})

const checkLoginStatus = async () => {
  try {
    isLogin.value = await IsLogin()
    if (isLogin.value) {
      userInfo.value = await GetUserInfo()
    }
  } catch (error) {
    console.error('[UserProfile] 检查登录状态失败:', error)
  }
}

const goToLogin = () => {
  router.push('/login')
}

const handleLogout = async () => {
  try {
    await Logout()
    isLogin.value = false
    userInfo.value = null
    console.log('[UserProfile] 退出登录成功')
  } catch (error) {
    console.error('[UserProfile] 退出登录失败:', error)
  }
}
</script>

<style scoped>
.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  transition: background 0.2s;
}

.user-profile:hover {
  background: rgba(255, 255, 255, 0.08);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-btn {
  font-size: 12px;
  color: #999;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  text-align: left;
  transition: color 0.2s;
}

.logout-btn:hover {
  color: #fff;
}

.login-btn {
  padding: 8px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.login-btn:active {
  transform: translateY(0);
}
</style>

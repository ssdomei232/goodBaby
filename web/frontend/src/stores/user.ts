import { defineStore } from 'pinia'
import { ref } from 'vue'
import { userApi } from '@/api'
import type { UserInfo } from '@/api/types'

export const useUserStore = defineStore('user', () => {
  const user = ref<UserInfo | null>(null)
  const loaded = ref(false)

  /** 拉取当前登录用户，未登录时静默失败 */
  async function fetchUser(): Promise<UserInfo | null> {
    try {
      user.value = await userApi.info()
    } catch {
      user.value = null
    } finally {
      loaded.value = true
    }
    return user.value
  }

  async function logout() {
    try {
      await userApi.logout()
    } finally {
      user.value = null
    }
  }

  function clear() {
    user.value = null
  }

  return { user, loaded, fetchUser, logout, clear }
})

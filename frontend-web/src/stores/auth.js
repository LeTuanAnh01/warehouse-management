import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')

  async function login(credentials) {
    try {
      const response = await authService.login(credentials)
      
      if (response.success) {
        token.value = response.data.access_token
        refreshToken.value = response.data.refresh_token
        user.value = response.data.user

        localStorage.setItem('access_token', response.data.access_token)
        localStorage.setItem('refresh_token', response.data.refresh_token)
        localStorage.setItem('user', JSON.stringify(response.data.user))

        return { success: true, data: response.data }
      }
      
      return { success: false, error: 'Login failed' }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || 'Login failed' 
      }
    }
  }

  async function register(userData) {
    try {
      const response = await authService.register(userData)
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || 'Registration failed' 
      }
    }
  }

  function logout() {
    token.value = ''
    refreshToken.value = ''
    user.value = null

    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  function updateUser(userData) {
    user.value = { ...user.value, ...userData }
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  async function refreshAccessToken() {
    try {
      const response = await authService.refreshToken(refreshToken.value)
      
      if (response.success) {
        token.value = response.data.access_token
        refreshToken.value = response.data.refresh_token
        
        localStorage.setItem('access_token', response.data.access_token)
        localStorage.setItem('refresh_token', response.data.refresh_token)
        
        return true
      }
      
      return false
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    userRole,
    login,
    register,
    logout,
    refreshAccessToken,
    updateUser,
  }
})

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { userService } from '@/services'
import { useAuthStore } from '@/stores/auth'

export const useUserStore = defineStore('user', () => {
  const users = ref([])
  const total = ref(0)
  const loading = ref(false)
  const error = ref(null)

  async function fetchUsers(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await userService.getUsers(page, pageSize)
      users.value = response.data || []
      total.value = response.total || 0
    } catch (err) {
      error.value = err.response?.data?.message || 'Lấy danh sách người dùng thất bại'
    } finally {
      loading.value = false
    }
  }

  async function updateUser(id, userData) {
    loading.value = true
    error.value = null
    try {
      const data = await userService.updateUser(id, userData)

      // Cập nhật trong danh sách nếu có
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        users.value[index] = { ...users.value[index], ...data }
      }

      // Nếu đang update chính mình thì cập nhật authStore luôn
      const authStore = useAuthStore()
      if (authStore.user?.id === id) {
        authStore.user = { ...authStore.user, ...data }  // ← trực tiếp gán, không gọi method
        localStorage.setItem('user', JSON.stringify(authStore.user))
      }

      return data
    } catch (err) {
      error.value = err.response?.data?.message || 'Cập nhật thất bại'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteUser(id) {
    loading.value = true
    error.value = null
    try {
      await userService.deleteUser(id)
      users.value = users.value.filter(u => u.id !== id)
      total.value -= 1
    } catch (err) {
      error.value = err.response?.data?.message || 'Xóa người dùng thất bại'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    users,
    total,
    loading,
    error,
    fetchUsers,
    updateUser,
    deleteUser,
  }
})
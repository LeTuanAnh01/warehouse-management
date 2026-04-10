<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  error.value = ''
  
  if (!form.username || !form.password) {
    error.value = 'Vui lòng nhập đầy đủ thông tin'
    return
  }

  loading.value = true

  try {
    const result = await authStore.login({
      username: form.username,
      password: form.password,
    })

    if (result.success) {
      router.push('/')
    } else {
      error.value = result.error || 'Đăng nhập thất bại'
    }
  } catch (err) {
    error.value = 'Đã xảy ra lỗi. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
    <div class="w-full max-w-md">
      <!-- Logo & Title -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-primary-600 rounded-2xl mb-4">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Quản Lý Kho</h1>
        <p class="text-gray-600">Đăng nhập vào hệ thống</p>
      </div>

      <!-- Login Card -->
      <div class="card">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Error Alert -->
          <div 
            v-if="error" 
            class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg flex items-start gap-3"
          >
            <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            <span class="text-sm">{{ error }}</span>
          </div>

          <!-- Username Field -->
          <div>
            <label for="username" class="label">
              Tên đăng nhập
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              class="input"
              placeholder="Nhập tên đăng nhập"
              :disabled="loading"
              autocomplete="username"
            />
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="label">
              Mật khẩu
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              class="input"
              placeholder="Nhập mật khẩu"
              :disabled="loading"
              autocomplete="current-password"
            />
          </div>

          <!-- Remember & Forgot -->
          <div class="flex items-center justify-between">
            <label class="flex items-center">
              <input type="checkbox" class="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500">
              <span class="ml-2 text-sm text-gray-600">Ghi nhớ đăng nhập</span>
            </label>
            <a href="#" class="text-sm text-primary-600 hover:text-primary-700">
              Quên mật khẩu?
            </a>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="btn btn-primary w-full py-3 text-base"
            :disabled="loading"
          >
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Đang đăng nhập...
            </span>
            <span v-else>Đăng nhập</span>
          </button>

          <!-- Register Link -->
          <div class="text-center pt-4 border-t border-gray-200">
            <p class="text-sm text-gray-600">
              Chưa có tài khoản?
              <router-link to="/register" class="text-primary-600 hover:text-primary-700 font-medium">
                Đăng ký ngay
              </router-link>
            </p>
          </div>
        </form>
      </div>

      <!-- Footer -->
      <p class="text-center text-sm text-gray-500 mt-8">
        © 2026 Warehouse Management System. All rights reserved.
      </p>
    </div>
  </div>
</template>

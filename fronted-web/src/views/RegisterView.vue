<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  fullName: '',
  phone: '',
})

const loading = ref(false)
const error = ref('')
const success = ref(false)

const handleRegister = async () => {
  error.value = ''
  
  // Validation
  if (!form.username || !form.email || !form.password || !form.fullName) {
    error.value = 'Vui lòng nhập đầy đủ thông tin bắt buộc'
    return
  }

  if (form.password.length < 6) {
    error.value = 'Mật khẩu phải có ít nhất 6 ký tự'
    return
  }

  if (form.password !== form.confirmPassword) {
    error.value = 'Mật khẩu xác nhận không khớp'
    return
  }

  loading.value = true

  try {
    const result = await authStore.register({
      username: form.username,
      email: form.email,
      password: form.password,
      full_name: form.fullName,
      phone: form.phone,
    })

    if (result.success) {
      success.value = true
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } else {
      error.value = result.error || 'Đăng ký thất bại'
    }
  } catch (err) {
    error.value = 'Đã xảy ra lỗi. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4 py-8">
    <div class="w-full max-w-md">
      <!-- Logo & Title -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-primary-600 rounded-2xl mb-4">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Tạo Tài Khoản</h1>
        <p class="text-gray-600">Đăng ký để sử dụng hệ thống</p>
      </div>

      <!-- Register Card -->
      <div class="card">
        <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- Success Alert -->
          <div 
            v-if="success" 
            class="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg flex items-start gap-3"
          >
            <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <div class="text-sm">
              <p class="font-medium">Đăng ký thành công!</p>
              <p class="mt-1">Đang chuyển đến trang đăng nhập...</p>
            </div>
          </div>

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

          <!-- Full Name -->
          <div>
            <label for="fullName" class="label">
              Họ và tên <span class="text-red-500">*</span>
            </label>
            <input
              id="fullName"
              v-model="form.fullName"
              type="text"
              class="input"
              placeholder="Nguyễn Văn A"
              :disabled="loading"
              required
            />
          </div>

          <!-- Username -->
          <div>
            <label for="username" class="label">
              Tên đăng nhập <span class="text-red-500">*</span>
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              class="input"
              placeholder="username"
              :disabled="loading"
              autocomplete="username"
              required
            />
            <p class="text-xs text-gray-500 mt-1">Tối thiểu 3 ký tự</p>
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="label">
              Email <span class="text-red-500">*</span>
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="input"
              placeholder="email@example.com"
              :disabled="loading"
              autocomplete="email"
              required
            />
          </div>

          <!-- Phone -->
          <div>
            <label for="phone" class="label">
              Số điện thoại
            </label>
            <input
              id="phone"
              v-model="form.phone"
              type="tel"
              class="input"
              placeholder="0123456789"
              :disabled="loading"
              autocomplete="tel"
            />
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="label">
              Mật khẩu <span class="text-red-500">*</span>
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              class="input"
              placeholder="••••••••"
              :disabled="loading"
              autocomplete="new-password"
              required
            />
            <p class="text-xs text-gray-500 mt-1">Tối thiểu 6 ký tự</p>
          </div>

          <!-- Confirm Password -->
          <div>
            <label for="confirmPassword" class="label">
              Xác nhận mật khẩu <span class="text-red-500">*</span>
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              class="input"
              placeholder="••••••••"
              :disabled="loading"
              autocomplete="new-password"
              required
            />
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="btn btn-primary w-full py-3 text-base"
            :disabled="loading || success"
          >
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Đang xử lý...
            </span>
            <span v-else>Đăng ký</span>
          </button>

          <!-- Login Link -->
          <div class="text-center pt-4 border-t border-gray-200">
            <p class="text-sm text-gray-600">
              Đã có tài khoản?
              <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
                Đăng nhập ngay
              </router-link>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

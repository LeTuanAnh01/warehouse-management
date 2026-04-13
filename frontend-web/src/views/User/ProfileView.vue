<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userService, authService } from '@/services/index' // ← sửa import

const authStore = useAuthStore()

const user = computed(() => authStore.user)

const activeTab = ref('info') // 'info' or 'password'

const profileForm = reactive({
  full_name: '',
  email: '',
  phone: '',
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const loading = ref(false)
const error = ref('')
const success = ref('')

onMounted(() => {
  if (user.value) {
    profileForm.full_name = user.value.full_name || ''
    profileForm.email = user.value.email || ''
    profileForm.phone = user.value.phone || ''
  }
})

const handleUpdateProfile = async () => {
  error.value = ''
  success.value = ''

  if (!profileForm.full_name || !profileForm.email) {
    error.value = 'Vui lòng nhập đầy đủ thông tin'
    return
  }

  loading.value = true

  try {
    const data = await userService.updateUser(user.value.id, {  // ← dùng userService
      full_name: profileForm.full_name,
      email: profileForm.email,
      phone: profileForm.phone,
    })

    if (data) {
      authStore.updateUser(data)
    }

    success.value = 'Cập nhật thông tin thành công'
  } catch (err) {
    error.value = err.response?.message || 'Cập nhật thất bại. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  error.value = ''
  success.value = ''

  if (!passwordForm.currentPassword || !passwordForm.newPassword) {
    error.value = 'Vui lòng nhập đầy đủ thông tin'
    return
  }

  if (passwordForm.newPassword.length < 6) {
    error.value = 'Mật khẩu mới phải có ít nhất 6 ký tự'
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    error.value = 'Mật khẩu xác nhận không khớp'
    return
  }

  loading.value = true

  try {
    await authService.changePassword({
      old_password: passwordForm.currentPassword,  // ← sửa current_password → old_password
      new_password: passwordForm.newPassword,
    })

    success.value = 'Đổi mật khẩu thành công'
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (err) {
    error.value = err.response?.data?.message || 'Đổi mật khẩu thất bại'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Thông tin cá nhân</h1>
      <p class="text-gray-600 mt-1">Quản lý thông tin tài khoản và bảo mật</p>
    </div>

    <!-- Profile Card -->
    <div class="card">
      <div class="flex items-start gap-6 pb-6 border-b border-gray-200">
        <!-- Avatar -->
        <div class="w-24 h-24 bg-primary-100 rounded-full flex items-center justify-center flex-shrink-0">
          <span class="text-primary-700 font-bold text-3xl">
            {{ user?.full_name?.charAt(0).toUpperCase() || 'U' }}
          </span>
        </div>

        <!-- User Info -->
        <div class="flex-1">
          <h2 class="text-xl font-semibold text-gray-900">{{ user?.full_name || 'User' }}</h2>
          <p class="text-gray-600">{{ user?.email || 'email@example.com' }}</p>
          <div class="flex items-center gap-4 mt-3">
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800"
            >
              {{ user?.role || 'VIEWER' }}
            </span>
            <span class="text-sm text-gray-500">
              Tham gia: {{ new Date(user?.created_at || Date.now()).toLocaleDateString('vi-VN') }}
            </span>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="border-b border-gray-200 mt-6">
        <nav class="-mb-px flex gap-8">
          <button
            @click="activeTab = 'info'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === 'info'
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            ]"
          >
            Thông tin cơ bản
          </button>
          <button
            @click="activeTab = 'password'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === 'password'
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            ]"
          >
            Đổi mật khẩu
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- Profile Info Tab -->
        <div v-if="activeTab === 'info'">
          <form @submit.prevent="handleUpdateProfile" class="space-y-6">
            <!-- Alert Messages -->
            <div
              v-if="error"
              class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg flex items-start gap-3"
            >
              <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-sm">{{ error }}</span>
            </div>

            <div
              v-if="success"
              class="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg flex items-start gap-3"
            >
              <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-sm">{{ success }}</span>
            </div>

            <!-- Username (readonly) -->
            <div>
              <label class="label">Tên đăng nhập</label>
              <input
                :value="user?.username"
                type="text"
                class="input bg-gray-50"
                disabled
                readonly
              />
              <p class="text-xs text-gray-500 mt-1">Tên đăng nhập không thể thay đổi</p>
            </div>

            <!-- Full Name -->
            <div>
              <label for="full_name" class="label">
                Họ và tên <span class="text-red-500">*</span>
              </label>
              <input
                id="full_name"
                v-model="profileForm.full_name"
                type="text"
                class="input"
                placeholder="Nguyễn Văn A"
                :disabled="loading"
                required
              />
            </div>

            <!-- Email -->
            <div>
              <label for="email" class="label">
                Email <span class="text-red-500">*</span>
              </label>
              <input
                id="email"
                v-model="profileForm.email"
                type="email"
                class="input"
                placeholder="email@example.com"
                :disabled="loading"
                required
              />
            </div>

            <!-- Phone -->
            <div>
              <label for="phone" class="label">Số điện thoại</label>
              <input
                id="phone"
                v-model="profileForm.phone"
                type="tel"
                class="input"
                placeholder="0123456789"
                :disabled="loading"
              />
            </div>

            <!-- Submit Button -->
            <div class="pt-4">
              <button type="submit" class="btn btn-primary" :disabled="loading">
                <span v-if="loading">Đang cập nhật...</span>
                <span v-else>Lưu thay đổi</span>
              </button>
            </div>
          </form>
        </div>

        <!-- Change Password Tab -->
        <div v-if="activeTab === 'password'">
          <form @submit.prevent="handleChangePassword" class="space-y-6">
            <!-- Alert Messages -->
            <div
              v-if="error"
              class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg flex items-start gap-3"
            >
              <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-sm">{{ error }}</span>
            </div>

            <div
              v-if="success"
              class="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg flex items-start gap-3"
            >
              <svg class="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-sm">{{ success }}</span>
            </div>

            <!-- Current Password -->
            <div>
              <label for="current_password" class="label">
                Mật khẩu hiện tại <span class="text-red-500">*</span>
              </label>
              <input
                id="current_password"
                v-model="passwordForm.currentPassword"
                type="password"
                class="input"
                placeholder="••••••••"
                :disabled="loading"
                autocomplete="current-password"
                required
              />
            </div>

            <!-- New Password -->
            <div>
              <label for="new_password" class="label">
                Mật khẩu mới <span class="text-red-500">*</span>
              </label>
              <input
                id="new_password"
                v-model="passwordForm.newPassword"
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
              <label for="confirm_password" class="label">
                Xác nhận mật khẩu mới <span class="text-red-500">*</span>
              </label>
              <input
                id="confirm_password"
                v-model="passwordForm.confirmPassword"
                type="password"
                class="input"
                placeholder="••••••••"
                :disabled="loading"
                autocomplete="new-password"
                required
              />
            </div>

            <!-- Submit Button -->
            <div class="pt-4">
              <button type="submit" class="btn btn-primary" :disabled="loading">
                <span v-if="loading">Đang xử lý...</span>
                <span v-else>Đổi mật khẩu</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
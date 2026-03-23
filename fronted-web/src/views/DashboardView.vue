<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const stats = ref([
  {
    name: 'Tổng sản phẩmmmmmmm',
    value: '0',
    icon: 'cube',
    color: 'blue',
    loading: true,
  },
  {
    name: 'Tổng kho',
    value: '0',
    icon: 'warehouse',
    color: 'green',
    loading: true,
  },
  {
    name: 'Giao dịch hôm nay',
    value: '0',
    icon: 'transaction',
    color: 'yellow',
    loading: true,
  },
  {
    name: 'Giá trị tồn kho',
    value: '0đ',
    icon: 'money',
    color: 'purple',
    loading: true,
  },
])

const recentActivities = ref([
  { type: 'import', product: 'Sản phẩm A', quantity: 100, time: '2 giờ trước' },
  { type: 'export', product: 'Sản phẩm B', quantity: 50, time: '3 giờ trước' },
  { type: 'create', product: 'Sản phẩm C', quantity: 1, time: '5 giờ trước' },
])

const currentUser = computed(() => authStore.user)

onMounted(async () => {
  // TODO: Fetch real stats from API
  setTimeout(() => {
    stats.value[0].value = '125'
    stats.value[0].loading = false
    stats.value[1].value = '3'
    stats.value[1].loading = false
    stats.value[2].value = '12'
    stats.value[2].loading = false
    stats.value[3].value = '15,450,000đ'
    stats.value[3].loading = false
  }, 1000)
})

const getActivityIcon = (type) => {
  switch (type) {
    case 'import':
      return 'M12 4v16m8-8H4'
    case 'export':
      return 'M17 8l4 4m0 0l-4 4m4-4H3'
    case 'create':
      return 'M12 4v16m8-8H4'
    default:
      return 'M12 4v16m8-8H4'
  }
}

const getActivityColor = (type) => {
  switch (type) {
    case 'import':
      return 'text-green-600 bg-green-50'
    case 'export':
      return 'text-red-600 bg-red-50'
    case 'create':
      return 'text-blue-600 bg-blue-50'
    default:
      return 'text-gray-600 bg-gray-50'
  }
}

const getActivityLabel = (type) => {
  switch (type) {
    case 'import':
      return 'Nhập kho'
    case 'export':
      return 'Xuất kho'
    case 'create':
      return 'Tạo mới'
    default:
      return 'Khác'
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Welcome Section -->
    <div class="bg-gradient-to-r from-primary-600 to-primary-700 rounded-xl p-6 text-white">
      <h1 class="text-2xl font-bold mb-2">
        Xin chào, {{ currentUser?.full_name || 'User' }}! 👋
      </h1>
      <p class="text-primary-100">
        Chào mừng bạn quay lại hệ thống quản lý kho
      </p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="stat in stats"
        :key="stat.name"
        class="card"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 mb-1">{{ stat.name }}</p>
            <p v-if="stat.loading" class="text-2xl font-bold text-gray-400 animate-pulse">
              --
            </p>
            <p v-else class="text-2xl font-bold text-gray-900">
              {{ stat.value }}
            </p>
          </div>
          <div
            :class="[
              'w-12 h-12 rounded-lg flex items-center justify-center',
              stat.color === 'blue' && 'bg-blue-100',
              stat.color === 'green' && 'bg-green-100',
              stat.color === 'yellow' && 'bg-yellow-100',
              stat.color === 'purple' && 'bg-purple-100',
            ]"
          >
            <svg
              :class="[
                'w-6 h-6',
                stat.color === 'blue' && 'text-blue-600',
                stat.color === 'green' && 'text-green-600',
                stat.color === 'yellow' && 'text-yellow-600',
                stat.color === 'purple' && 'text-purple-600',
              ]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                v-if="stat.icon === 'cube'"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
              <path
                v-else-if="stat.icon === 'warehouse'"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
              <path
                v-else-if="stat.icon === 'transaction'"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
              />
              <path
                v-else-if="stat.icon === 'money'"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Two Column Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Activities -->
      <div class="lg:col-span-2">
        <div class="card">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-lg font-semibold text-gray-900">Hoạt động gần đây</h2>
            <router-link to="/transactions" class="text-sm text-primary-600 hover:text-primary-700">
              Xem tất cả →
            </router-link>
          </div>
          
          <div class="space-y-4">
            <div
              v-for="(activity, index) in recentActivities"
              :key="index"
              class="flex items-center gap-4 p-4 bg-gray-50 rounded-lg"
            >
              <div :class="['w-10 h-10 rounded-lg flex items-center justify-center', getActivityColor(activity.type)]">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    :d="getActivityIcon(activity.type)"
                  />
                </svg>
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900">{{ activity.product }}</p>
                <p class="text-xs text-gray-500">
                  {{ getActivityLabel(activity.type) }} - {{ activity.quantity }} sản phẩm
                </p>
              </div>
              <span class="text-xs text-gray-400">{{ activity.time }}</span>
            </div>
          </div>

          <div v-if="recentActivities.length === 0" class="text-center py-8 text-gray-500">
            <svg class="w-12 h-12 mx-auto mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
            <p>Chưa có hoạt động nào</p>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div>
        <div class="card">
          <h2 class="text-lg font-semibold text-gray-900 mb-6">Thao tác nhanh</h2>
          <div class="space-y-3">
            <router-link
              to="/products/create"
              class="flex items-center gap-3 p-3 bg-primary-50 text-primary-700 rounded-lg hover:bg-primary-100 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              <span class="font-medium">Thêm sản phẩm mới</span>
            </router-link>
            
            <button class="w-full flex items-center gap-3 p-3 bg-green-50 text-green-700 rounded-lg hover:bg-green-100 transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
              </svg>
              <span class="font-medium">Nhập kho</span>
            </button>
            
            <button class="w-full flex items-center gap-3 p-3 bg-red-50 text-red-700 rounded-lg hover:bg-red-100 transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
              <span class="font-medium">Xuất kho</span>
            </button>
            
            <router-link
              to="/reports"
              class="flex items-center gap-3 p-3 bg-gray-50 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="font-medium">Xem báo cáo</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

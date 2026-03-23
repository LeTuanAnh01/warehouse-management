<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProductStore } from '@/stores/product'

const router = useRouter()
const productStore = useProductStore()

const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const products = computed(() => productStore.products)
const totalProducts = computed(() => productStore.total)
const totalPages = computed(() => Math.ceil(totalProducts.value / pageSize.value))

const filteredProducts = computed(() => {
  if (!searchQuery.value) return products.value

  const query = searchQuery.value.toLowerCase()
  return products.value.filter(
    (p) =>
      p.name?.toLowerCase().includes(query) ||
      p.sku?.toLowerCase().includes(query) ||
      p.description?.toLowerCase().includes(query)
  )
})

onMounted(async () => {
  await fetchProducts()
})

const fetchProducts = async () => {
  loading.value = true
  try {
    await productStore.fetchProducts({
      page: currentPage.value,
      limit: pageSize.value,
    })
  } catch (error) {
    console.error('Failed to fetch products:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchProducts()
}

const handleDelete = async (id) => {
  if (!confirm('Bạn có chắc chắn muốn xóa sản phẩm này?')) return

  try {
    await productStore.deleteProduct(id)
    await fetchProducts()
  } catch (error) {
    alert('Xóa sản phẩm thất bại')
  }
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
  }).format(value)
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('vi-VN')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Quản lý sản phẩm</h1>
        <p class="text-gray-600 mt-1">Danh sách tất cả sản phẩm trong hệ thống</p>
      </div>
      <router-link to="/products/create" class="btn btn-primary">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Thêm sản phẩm
      </router-link>
    </div>

    <!-- Search & Filter -->
    <div class="card">
      <div class="flex items-center gap-4">
        <div class="flex-1 relative">
          <svg
            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Tìm kiếm sản phẩm theo tên, SKU..."
            class="input pl-10"
          />
        </div>
        <button class="btn btn-secondary">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
            />
          </svg>
          Lọc
        </button>
      </div>
    </div>

    <!-- Products Table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Sản phẩm
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                SKU
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Đơn vị
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Giá
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Trạng thái
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Ngày tạo
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Thao tác
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <!-- Loading State -->
            <tr v-if="loading">
              <td colspan="7" class="px-6 py-12 text-center">
                <div class="flex items-center justify-center gap-3">
                  <svg class="animate-spin h-6 w-6 text-primary-600" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  <span class="text-gray-600">Đang tải...</span>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-else-if="filteredProducts.length === 0">
              <td colspan="7" class="px-6 py-12 text-center">
                <svg
                  class="w-16 h-16 mx-auto mb-4 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                  />
                </svg>
                <p class="text-gray-600 mb-2">Không tìm thấy sản phẩm nào</p>
                <router-link to="/products/create" class="text-primary-600 hover:text-primary-700">
                  Thêm sản phẩm đầu tiên
                </router-link>
              </td>
            </tr>

            <!-- Product Rows -->
            <tr v-for="product in filteredProducts" :key="product.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div>
                  <div class="font-medium text-gray-900">{{ product.name }}</div>
                  <div class="text-sm text-gray-500">{{ product.description || 'N/A' }}</div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-900">
                {{ product.sku || 'N/A' }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900">
                {{ product.unit || 'N/A' }}
              </td>
              <td class="px-6 py-4 text-sm font-medium text-gray-900">
                {{ formatCurrency(product.price || 0) }}
              </td>
              <td class="px-6 py-4">
                <span
                  :class="[
                    'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                    product.is_active
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800',
                  ]"
                >
                  {{ product.is_active ? 'Hoạt động' : 'Ngừng' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">
                {{ formatDate(product.created_at) }}
              </td>
              <td class="px-6 py-4 text-right text-sm font-medium">
                <div class="flex items-center justify-end gap-2">
                  <router-link
                    :to="`/products/${product.id}/edit`"
                    class="text-primary-600 hover:text-primary-900"
                  >
                    Sửa
                  </router-link>
                  <button
                    @click="handleDelete(product.id)"
                    class="text-red-600 hover:text-red-900"
                  >
                    Xóa
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Hiển thị <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span> đến
            <span class="font-medium">{{ Math.min(currentPage * pageSize, totalProducts) }}</span>
            trong tổng số <span class="font-medium">{{ totalProducts }}</span> sản phẩm
          </div>
          <div class="flex items-center gap-2">
            <button
              :disabled="currentPage === 1"
              @click="handlePageChange(currentPage - 1)"
              class="btn btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Trước
            </button>
            <div class="flex gap-1">
              <button
                v-for="page in totalPages"
                :key="page"
                @click="handlePageChange(page)"
                :class="[
                  'px-3 py-1 rounded-md text-sm font-medium',
                  page === currentPage
                    ? 'bg-primary-600 text-white'
                    : 'text-gray-700 hover:bg-gray-100',
                ]"
              >
                {{ page }}
              </button>
            </div>
            <button
              :disabled="currentPage === totalPages"
              @click="handlePageChange(currentPage + 1)"
              class="btn btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Sau
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useProductStore } from '@/stores/product'

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()

const isEditMode = computed(() => !!route.params.id)
const pageTitle = computed(() => (isEditMode.value ? 'Chỉnh sửa sản phẩm' : 'Thêm sản phẩm mới'))

const loading = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  sku: '',
  description: '',
  unit: 'cái',
  price: 0,
  is_active: true,
})

const units = ['cái', 'hộp', 'thùng', 'kg', 'g', 'lít', 'ml', 'm', 'cm']

onMounted(async () => {
  if (isEditMode.value) {
    await loadProduct()
  }
})

const loadProduct = async () => {
  loading.value = true
  try {
    const product = await productStore.getProduct(route.params.id)
    if (product) {
      form.name = product.name
      form.sku = product.sku || ''
      form.description = product.description || ''
      form.unit = product.unit || 'cái'
      form.price = product.price || 0
      form.is_active = product.is_active !== false
    }
  } catch (err) {
    error.value = 'Không thể tải thông tin sản phẩm'
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  error.value = ''

  // Validation
  if (!form.name) {
    error.value = 'Vui lòng nhập tên sản phẩm'
    return
  }

  if (form.price < 0) {
    error.value = 'Giá sản phẩm không được âm'
    return
  }

  loading.value = true

  try {
    if (isEditMode.value) {
      await productStore.updateProduct(route.params.id, form)
    } else {
      await productStore.createProduct(form)
    }
    router.push('/products')
  } catch (err) {
    error.value = err.response?.data?.message || 'Đã xảy ra lỗi. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  router.push('/products')
}
</script>

<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <button @click="handleCancel" class="text-gray-600 hover:text-gray-900">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>
      <div>
        <h1 class="text-2xl font-bold text-gray-900">{{ pageTitle }}</h1>
        <p class="text-gray-600 mt-1">Nhập thông tin sản phẩm</p>
      </div>
    </div>

    <!-- Form -->
    <div class="card">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Error Alert -->
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

        <!-- Product Name -->
        <div>
          <label for="name" class="label">
            Tên sản phẩm <span class="text-red-500">*</span>
          </label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            class="input"
            placeholder="Nhập tên sản phẩm"
            :disabled="loading"
            required
          />
        </div>

        <!-- SKU -->
        <div>
          <label for="sku" class="label">Mã SKU</label>
          <input
            id="sku"
            v-model="form.sku"
            type="text"
            class="input"
            placeholder="SKU-001"
            :disabled="loading"
          />
          <p class="text-xs text-gray-500 mt-1">Mã định danh duy nhất cho sản phẩm</p>
        </div>

        <!-- Description -->
        <div>
          <label for="description" class="label">Mô tả</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="4"
            class="input resize-none"
            placeholder="Nhập mô tả chi tiết về sản phẩm"
            :disabled="loading"
          ></textarea>
        </div>

        <!-- Unit & Price Row -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Unit -->
          <div>
            <label for="unit" class="label">
              Đơn vị tính <span class="text-red-500">*</span>
            </label>
            <select id="unit" v-model="form.unit" class="input" :disabled="loading" required>
              <option v-for="unit in units" :key="unit" :value="unit">
                {{ unit }}
              </option>
            </select>
          </div>

          <!-- Price -->
          <div>
            <label for="price" class="label">
              Giá (VNĐ) <span class="text-red-500">*</span>
            </label>
            <input
              id="price"
              v-model.number="form.price"
              type="number"
              class="input"
              placeholder="0"
              :disabled="loading"
              min="0"
              step="1000"
              required
            />
          </div>
        </div>

        <!-- Active Status -->
        <div class="flex items-center gap-3">
          <input
            id="is_active"
            v-model="form.is_active"
            type="checkbox"
            class="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            :disabled="loading"
          />
          <label for="is_active" class="text-sm text-gray-700">Sản phẩm đang hoạt động</label>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-3 pt-6 border-t border-gray-200">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="flex items-center gap-2">
              <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Đang xử lý...
            </span>
            <span v-else>
              {{ isEditMode ? 'Cập nhật' : 'Tạo mới' }}
            </span>
          </button>
          <button type="button" @click="handleCancel" class="btn btn-secondary" :disabled="loading">
            Hủy
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

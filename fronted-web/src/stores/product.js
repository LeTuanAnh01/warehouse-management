import { defineStore } from 'pinia'
import { ref } from 'vue'
import { productService } from '@/services'

export const useProductStore = defineStore('product', () => {
  const products = ref([])
  const currentProduct = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const pagination = ref({
    page: 1,
    limit: 20,
    total_records: 0,
    total_pages: 0,
  })

  async function fetchProducts(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productService.getProducts(page, pageSize)
      
      if (response.success) {
        products.value = response.data
        pagination.value = response.pagination
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch products'
    } finally {
      loading.value = false
    }
  }

  async function fetchProduct(id) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productService.getProduct(id)
      
      if (response.success) {
        currentProduct.value = response.data
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch product'
    } finally {
      loading.value = false
    }
  }

  async function createProduct(productData) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productService.createProduct(productData)
      
      if (response.success) {
        products.value.unshift(response.data)
        return { success: true, data: response.data }
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create product'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function updateProduct(id, productData) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productService.updateProduct(id, productData)
      
      if (response.success) {
        const index = products.value.findIndex(p => p.id === id)
        if (index !== -1) {
          products.value[index] = response.data
        }
        return { success: true, data: response.data }
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update product'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function deleteProduct(id) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productService.deleteProduct(id)
      
      if (response.success) {
        products.value = products.value.filter(p => p.id !== id)
        return { success: true }
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete product'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  return {
    products,
    currentProduct,
    loading,
    error,
    pagination,
    fetchProducts,
    fetchProduct,
    createProduct,
    updateProduct,
    deleteProduct,
  }
})

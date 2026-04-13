import apiClient from './api'

export const authService = {
  async login(credentials) {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  async register(userData) {
    const response = await apiClient.post('/auth/register', userData)
    return response.data
  },

  async refreshToken(refreshToken) {
    const response = await apiClient.post('/auth/refresh', { refresh_token: refreshToken })
    return response.data
  },

  async logout() {
    const response = await apiClient.post('/auth/logout')
    return response.data
  },

  async changePassword(passwords) {
    const response = await apiClient.post('/auth/change-password', passwords)
    return response.data
  },
}

export const productService = {
  async getProducts(page = 1, pageSize = 20) {
    const response = await apiClient.get('/products', {
      params: { page, page_size: pageSize }
    })
    return response.data
  },

  async getProduct(id) {
    const response = await apiClient.get(`/products/${id}`)
    return response.data
  },

  async createProduct(productData) {
    const response = await apiClient.post('/products', productData)
    return response.data
  },

  async updateProduct(id, productData) {
    const response = await apiClient.put(`/products/${id}`, productData)
    return response.data
  },

  async deleteProduct(id) {
    const response = await apiClient.delete(`/products/${id}`)
    return response.data
  },

  async getProductInventory(id) {
    const response = await apiClient.get(`/products/${id}/inventory`)
    return response.data
  },
}

export const userService = {
  async getUsers(page = 1, pageSize = 20) {
    const response = await apiClient.get('/users', {
      params: { page, page_size: pageSize }
    })
    return response.data
  },

  async getUser(id) {
    const response = await apiClient.get(`/users/${id}`)
    return response.data
  },

  async updateUser(id, userData) {
    const response = await apiClient.put(`/users/${id}`, userData)
    return response.data.data  // ← lấy .data.data vì backend wrap { success, message, data: {...} }
  },

  async deleteUser(id) {
    const response = await apiClient.delete(`/users/${id}`)
    return response.data
  },
}
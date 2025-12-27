import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: add token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor: adapt to backend unified response format
apiClient.interceptors.response.use(
  (response) => {
    // Backend response format: { code: 0, message: "success", data: {...} }
    const result = response.data

    // Check if code is 0 (success)
    if (result.code === 0) {
      return result.data // Return data field
    } else {
      // Error response
      return Promise.reject({
        code: result.code,
        message: result.message,
      })
    }
  },
  (error: AxiosError) => {
    // Handle HTTP errors
    if (error.response) {
      const status = error.response.status

      // 401 Unauthorized - clear local storage
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('user')

        // Only redirect if not already on login page or register page
        const currentPath = window.location.pathname
        if (currentPath !== '/login' && currentPath !== '/register') {
          window.location.href = '/login'
        }
      }

      // Backend error response format: { code: xxx, message: "error", data: null }
      const errorData = error.response.data as { code?: number; message?: string }
      return Promise.reject({
        code: errorData?.code || status,
        message: errorData?.message || error.message,
      })
    }

    // Network error or other errors
    return Promise.reject({
      code: 0,
      message: error.message || 'Network error',
    })
  }
)

export default apiClient

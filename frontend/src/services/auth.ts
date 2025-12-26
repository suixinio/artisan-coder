import apiClient from './api'
import type { LoginRequest, RegisterRequest, AuthResponse } from '@/types/api'
import type { User } from '@/types/models'

export const authService = {
  // Login
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post('/api/auth/login', data) as AuthResponse
    // Save token and user info
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Register
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post('/api/auth/register', data) as AuthResponse
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Logout
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/auth/logout')
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    return await apiClient.get('/api/auth/me') as User
  },

  // Refresh token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post('/api/auth/refresh') as AuthResponse
    localStorage.setItem('token', response.token)
    return response
  },
}

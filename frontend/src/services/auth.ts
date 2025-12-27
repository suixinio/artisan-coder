import apiClient from './api'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import type { User } from '@/types/models'

interface AuthResponse {
  user: User
  token: string
  refreshToken: string
}

export const authService = {
  // Login
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/login', data) as unknown as AuthResponse
    // Save token and user info
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Register
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/register', data) as unknown as AuthResponse
    // Save token and user info
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Logout
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/auth/logout')
    } finally {
      // Always clear local storage regardless of backend success
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
    }
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    return await apiClient.get<User>('/api/auth/me') as unknown as User
  },

  // Refresh token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/refresh') as unknown as AuthResponse
    // Update local storage
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },
}

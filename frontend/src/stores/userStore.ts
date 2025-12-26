import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types/models'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import { authService } from '@/services/auth'
import { mockAuthService } from '@/services/mockApi'

// Use Mock or real API
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

interface UserState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean

  // Actions
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshToken: () => Promise<void>
  initialize: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true })
        try {
          const response = USE_MOCK
            ? await mockAuthService.login(credentials.email, credentials.password)
            : await authService.login(credentials)

          set({
            user: response.user,
            token: response.token,
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      register: async (data: RegisterRequest) => {
        set({ isLoading: true })
        try {
          const response = USE_MOCK
            ? await mockAuthService.register(data.username, data.email, data.password)
            : await authService.register(data)

          set({
            user: response.user,
            token: response.token,
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: async () => {
        set({ isLoading: true })
        try {
          if (USE_MOCK) {
            await mockAuthService.logout()
          } else {
            await authService.logout()
          }
          set({
            user: null,
            token: null,
            isAuthenticated: false,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      refreshToken: async () => {
        const { token } = get()
        if (!token) return

        try {
          const response = await authService.refreshToken()
          set({
            user: response.user,
            token: response.token,
          })
        } catch (error) {
          set({ user: null, token: null, isAuthenticated: false })
          throw error
        }
      },

      initialize: () => {
        const token = localStorage.getItem('token')
        const userStr = localStorage.getItem('user')
        const user = userStr ? JSON.parse(userStr) : null

        if (token && user) {
          set({
            user,
            token,
            isAuthenticated: true,
          })
        }
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({ user: state.user, token: state.token }),
    }
  )
)

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
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  sidebarOpen: boolean

  // Actions
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshAccessToken: () => Promise<void>
  initialize: () => void
  toggleSidebar: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      sidebarOpen: false,

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true })
        try {
          const response = USE_MOCK
            ? await mockAuthService.login(credentials.email, credentials.password)
            : await authService.login(credentials)

          set({
            user: response.user,
            token: response.token,
            refreshToken: response.refreshToken,
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
            refreshToken: response.refreshToken,
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
            refreshToken: null,
            isAuthenticated: false,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      refreshAccessToken: async () => {
        const { refreshToken: currentRefreshToken } = get()
        if (!currentRefreshToken) {
          throw new Error('No refresh token available')
        }

        try {
          const response = await authService.refreshToken()
          set({
            user: response.user,
            token: response.token,
            refreshToken: response.refreshToken,
          })
        } catch (error) {
          // Token refresh failed, clear state
          set({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
          })
          throw error
        }
      },

      initialize: () => {
        const token = localStorage.getItem('token')
        const refreshToken = localStorage.getItem('refreshToken')
        const userStr = localStorage.getItem('user')
        const user = userStr ? JSON.parse(userStr) : null

        if (token && user) {
          set({
            user,
            token,
            refreshToken,
            isAuthenticated: true,
          })
        }
      },

      toggleSidebar: () => {
        set({ sidebarOpen: !get().sidebarOpen })
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        refreshToken: state.refreshToken,
      }),
    }
  )
)

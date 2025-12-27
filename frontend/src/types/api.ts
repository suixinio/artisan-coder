// Authentication related types
export interface LoginRequest {
  email: string
  password: string
  rememberMe?: boolean
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirmPassword: string
}

export interface AuthResponse {
  user: import('./models').User
  token: string
  refreshToken: string
}

// API error response
export interface ApiError {
  code: number
  message: string
}

// Unified API response (backend format)
// Backend returns: { code: 0, message: "success", data: {...} }
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// Paginated response
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

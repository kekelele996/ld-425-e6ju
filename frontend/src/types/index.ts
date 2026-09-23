export * from './enums'
export * from './project'
export * from './design'
export * from './material'
export * from './budget'
export * from './construction'

export interface User {
  id: number
  username: string
  role: string
  name: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

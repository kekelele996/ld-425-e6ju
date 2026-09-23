import { apiGet, apiPost } from '@/utils/request'
import type { LoginResponse, User } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function login(username: string, password: string): Promise<LoginResponse> {
  return apiPost<LoginResponse>(API_PATHS.login, { username, password })
}

export function getMe(): Promise<User> {
  return apiGet<User>(API_PATHS.me)
}

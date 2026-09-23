import axios, { AxiosError } from 'axios'
import type { ApiResponse } from '@/types'

const request = axios.create({
  baseURL: '',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('home_renovation_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('home_renovation_token')
      localStorage.removeItem('home_renovation_user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export async function apiGet<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const response = await request.get<ApiResponse<T>>(url, { params })
  return response.data.data
}

export async function apiPost<T>(url: string, body?: unknown): Promise<T> {
  const response = await request.post<ApiResponse<T>>(url, body)
  return response.data.data
}

export async function apiPut<T>(url: string, body?: unknown): Promise<T> {
  const response = await request.put<ApiResponse<T>>(url, body)
  return response.data.data
}

export async function apiDelete<T>(url: string): Promise<T> {
  const response = await request.delete<ApiResponse<T>>(url)
  return response.data.data
}

export function extractErrorMessage(error: unknown, fallback = '请求失败'): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as ApiResponse | undefined
    return data?.message || error.message || fallback
  }
  return fallback
}

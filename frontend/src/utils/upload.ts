import axios from 'axios'
import type { ApiResponse } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

// 上传文件并返回服务端可访问 URL。
export async function uploadFile(file: File): Promise<string> {
  const formData = new FormData()
  formData.append('file', file)
  const token = localStorage.getItem('home_renovation_token')
  const response = await axios.post<ApiResponse<{ url: string }>>(API_PATHS.upload, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
      Authorization: token ? `Bearer ${token}` : '',
    },
  })
  return response.data.data.url
}

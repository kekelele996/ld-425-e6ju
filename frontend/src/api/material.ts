import { apiDelete, apiGet, apiPost, apiPut } from '@/utils/request'
import type { CreateMaterialRequest, MaterialItem, PageResult } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function listMaterials(params?: Record<string, unknown>): Promise<PageResult<MaterialItem> | MaterialItem[]> {
  return apiGet<PageResult<MaterialItem> | MaterialItem[]>(API_PATHS.materials, params)
}

export function getMaterial(id: number): Promise<MaterialItem> {
  return apiGet<MaterialItem>(`${API_PATHS.materials}/${id}`)
}

export function createMaterial(body: CreateMaterialRequest): Promise<MaterialItem> {
  return apiPost<MaterialItem>(API_PATHS.materials, body)
}

export function updateMaterial(id: number, body: Partial<CreateMaterialRequest>): Promise<MaterialItem> {
  return apiPut<MaterialItem>(`${API_PATHS.materials}/${id}`, body)
}

export function deleteMaterial(id: number): Promise<void> {
  return apiDelete<void>(`${API_PATHS.materials}/${id}`)
}

export function updateMaterialStatus(id: number, status: string): Promise<MaterialItem> {
  return apiPut<MaterialItem>(`${API_PATHS.materials}/${id}/status`, { status })
}

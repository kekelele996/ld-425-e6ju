import { apiDelete, apiGet, apiPost, apiPut } from '@/utils/request'
import type { CreateDesignRequest, DesignPhase, PageResult, ReviewDesignRequest } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function listDesigns(params?: Record<string, unknown>): Promise<PageResult<DesignPhase> | DesignPhase[]> {
  return apiGet<PageResult<DesignPhase> | DesignPhase[]>(API_PATHS.designs, params)
}

export function getDesign(id: number): Promise<DesignPhase> {
  return apiGet<DesignPhase>(`${API_PATHS.designs}/${id}`)
}

export function createDesign(body: CreateDesignRequest): Promise<DesignPhase> {
  return apiPost<DesignPhase>(API_PATHS.designs, body)
}

export function updateDesign(id: number, body: Partial<CreateDesignRequest>): Promise<DesignPhase> {
  return apiPut<DesignPhase>(`${API_PATHS.designs}/${id}`, body)
}

export function deleteDesign(id: number): Promise<void> {
  return apiDelete<void>(`${API_PATHS.designs}/${id}`)
}

export function submitDesign(id: number, body: { description?: string; file_urls?: string[] }): Promise<DesignPhase> {
  return apiPut<DesignPhase>(`${API_PATHS.designs}/${id}/submit`, body)
}

export function reviewDesign(id: number, body: ReviewDesignRequest): Promise<DesignPhase> {
  return apiPut<DesignPhase>(`${API_PATHS.designs}/${id}/review`, body)
}

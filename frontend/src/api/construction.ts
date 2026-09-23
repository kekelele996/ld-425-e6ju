import { apiDelete, apiGet, apiPost, apiPut } from '@/utils/request'
import type { AcceptConstructionRequest, ConstructionNode, CreateConstructionRequest, PageResult } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function listConstructions(params?: Record<string, unknown>): Promise<PageResult<ConstructionNode> | ConstructionNode[]> {
  return apiGet<PageResult<ConstructionNode> | ConstructionNode[]>(API_PATHS.constructions, params)
}

export function getConstruction(id: number): Promise<ConstructionNode> {
  return apiGet<ConstructionNode>(`${API_PATHS.constructions}/${id}`)
}

export function createConstruction(body: CreateConstructionRequest): Promise<ConstructionNode> {
  return apiPost<ConstructionNode>(API_PATHS.constructions, body)
}

export function updateConstruction(id: number, body: Partial<CreateConstructionRequest>): Promise<ConstructionNode> {
  return apiPut<ConstructionNode>(`${API_PATHS.constructions}/${id}`, body)
}

export function deleteConstruction(id: number): Promise<void> {
  return apiDelete<void>(`${API_PATHS.constructions}/${id}`)
}

export function updateConstructionStatus(id: number, status: string): Promise<ConstructionNode> {
  return apiPut<ConstructionNode>(`${API_PATHS.constructions}/${id}/status`, { status })
}

export function acceptConstruction(id: number, body: AcceptConstructionRequest): Promise<ConstructionNode> {
  return apiPut<ConstructionNode>(`${API_PATHS.constructions}/${id}/accept`, body)
}

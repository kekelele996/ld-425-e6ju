import { apiDelete, apiGet, apiPost, apiPut } from '@/utils/request'
import type { CreateProjectRequest, PageResult, RenovationProject, UpdateProjectRequest } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function listProjects(params?: Record<string, unknown>): Promise<PageResult<RenovationProject>> {
  return apiGet<PageResult<RenovationProject>>(API_PATHS.projects, params)
}

export function getProject(id: number): Promise<RenovationProject> {
  return apiGet<RenovationProject>(`${API_PATHS.projects}/${id}`)
}

export function createProject(body: CreateProjectRequest): Promise<RenovationProject> {
  return apiPost<RenovationProject>(API_PATHS.projects, body)
}

export function updateProject(id: number, body: UpdateProjectRequest): Promise<RenovationProject> {
  return apiPut<RenovationProject>(`${API_PATHS.projects}/${id}`, body)
}

export function deleteProject(id: number): Promise<void> {
  return apiDelete<void>(`${API_PATHS.projects}/${id}`)
}

export function updateProjectStatus(id: number, status: string): Promise<RenovationProject> {
  return apiPut<RenovationProject>(`${API_PATHS.projects}/${id}/status`, { status })
}

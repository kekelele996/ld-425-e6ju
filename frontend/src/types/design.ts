import type { PhaseStatus } from './enums'

export interface DesignPhase {
  id: number
  project_id: number
  name: string
  designer_id: number
  status: PhaseStatus
  version: number
  description: string
  file_urls: string[]
  review_comment: string
  reviewer_id: number
  created_at: string
  updated_at: string
}

export interface CreateDesignRequest {
  project_id: number
  name: string
  designer_id?: number
  description?: string
  file_urls?: string[]
}

export interface ReviewDesignRequest {
  approved: boolean
  comment?: string
}

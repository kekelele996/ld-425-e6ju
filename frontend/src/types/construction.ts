import type { ConstructionStatus, AcceptanceStatus } from './enums'

export interface ConstructionNode {
  id: number
  project_id: number
  name: string
  planned_start_date: string | null
  planned_end_date: string | null
  actual_start_date: string | null
  actual_end_date: string | null
  status: ConstructionStatus
  acceptance_status: AcceptanceStatus
  acceptance_photos: string[]
  acceptance_note: string
  created_at: string
  updated_at: string
}

export interface AcceptConstructionRequest {
  accepted: boolean
  photos?: string[]
  note?: string
}

export interface CreateConstructionRequest {
  project_id: number
  name: string
  planned_start_date?: string
  planned_end_date?: string
}

import type { DecorStyle, ProjectStatus, HouseType } from './enums'

export interface RenovationProject {
  id: number
  name: string
  house_type: string
  area: number
  decor_style: DecorStyle
  address: string
  owner_id: number
  designer_id: number
  foreman_id: number
  status: ProjectStatus
  contract_amount: number
  start_date: string | null
  expected_end_date: string | null
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  house_type: HouseType
  area: number
  decor_style: DecorStyle
  address?: string
  owner_id: number
  designer_id?: number
  foreman_id?: number
  status?: ProjectStatus
  contract_amount?: number
  start_date?: string
  expected_end_date?: string
}

export interface UpdateProjectRequest extends Partial<CreateProjectRequest> {}

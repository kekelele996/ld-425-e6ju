import type { PurchaseStatus } from './enums'

export interface MaterialItem {
  id: number
  project_id: number
  name: string
  category: string
  spec: string
  brand: string
  quantity: number
  unit: string
  unit_price: number
  total_price: number
  purchase_status: PurchaseStatus
  supplier: string
  space: string
  created_at: string
  updated_at: string
}

export interface CreateMaterialRequest {
  project_id: number
  name: string
  category: string
  spec?: string
  brand?: string
  quantity: number
  unit: string
  unit_price?: number
  supplier?: string
  space?: string
}

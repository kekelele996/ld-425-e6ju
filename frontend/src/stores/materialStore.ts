import { create } from 'zustand'
import type { MaterialItem } from '@/types'
import { listMaterials } from '@/api/material'

interface MaterialState {
  materials: MaterialItem[]
  loading: boolean
  fetchMaterials: (projectId?: number) => Promise<void>
}

export const useMaterialStore = create<MaterialState>((set) => ({
  materials: [],
  loading: false,
  fetchMaterials: async (projectId) => {
    set({ loading: true })
    try {
      const result = await listMaterials(projectId ? { project_id: projectId } : { page: 1, page_size: 100 })
      if (Array.isArray(result)) {
        set({ materials: result, loading: false })
      } else {
        set({ materials: result.list, loading: false })
      }
    } catch {
      set({ loading: false })
    }
  },
}))

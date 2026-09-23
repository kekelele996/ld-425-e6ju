import { useMemo } from 'react'
import { ProjectStatus } from '@/types/enums'

// 封装项目阶段状态流转逻辑。
const PHASE_ORDER: Record<string, number> = {
  [ProjectStatus.Designing]: 0,
  [ProjectStatus.Quoting]: 1,
  [ProjectStatus.InProgress]: 2,
  [ProjectStatus.Completed]: 3,
  [ProjectStatus.Archived]: 4,
}

const PHASE_LABELS: Record<string, string> = {
  [ProjectStatus.Designing]: '设计阶段',
  [ProjectStatus.Quoting]: '报价阶段',
  [ProjectStatus.InProgress]: '施工进行中',
  [ProjectStatus.Completed]: '已完工',
  [ProjectStatus.Archived]: '已归档',
}

export function useProjectPhase(currentStatus: string) {
  return useMemo(() => {
    const current = PHASE_ORDER[currentStatus] ?? 0
    const nextStatuses = Object.values(ProjectStatus).filter(
      (status) => (PHASE_ORDER[status] ?? 0) === current + 1,
    )
    return {
      current,
      label: PHASE_LABELS[currentStatus] || currentStatus,
      nextStatuses,
      phaseOrder: PHASE_ORDER,
      phaseLabels: PHASE_LABELS,
    }
  }, [currentStatus])
}

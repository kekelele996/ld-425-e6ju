// 共享枚举：与后端 internal/constants/enums.go 保持一致。

export const DecorStyle = {
  Modern: 'Modern',
  Chinese: 'Chinese',
  Nordic: 'Nordic',
  Japanese: 'Japanese',
  Industrial: 'Industrial',
  Mediterranean: 'Mediterranean',
  Minimalist: 'Minimalist',
} as const
export type DecorStyle = (typeof DecorStyle)[keyof typeof DecorStyle]

export const PhaseStatus = {
  NotStarted: 'NotStarted',
  InProgress: 'InProgress',
  Revision: 'Revision',
  Approved: 'Approved',
} as const
export type PhaseStatus = (typeof PhaseStatus)[keyof typeof PhaseStatus]

export const ConstructionPhase = {
  Demolition: 'Demolition',
  Plumbing: 'Plumbing',
  Carpentry: 'Carpentry',
  Tiling: 'Tiling',
  Painting: 'Painting',
  Installation: 'Installation',
  SoftFurnishing: 'SoftFurnishing',
} as const
export type ConstructionPhase = (typeof ConstructionPhase)[keyof typeof ConstructionPhase]

export const PurchaseStatus = {
  NotPurchased: 'NotPurchased',
  Ordered: 'Ordered',
  Delivered: 'Delivered',
  Installed: 'Installed',
} as const
export type PurchaseStatus = (typeof PurchaseStatus)[keyof typeof PurchaseStatus]

export const ProjectStatus = {
  Designing: 'Designing',
  Quoting: 'Quoting',
  InProgress: 'InProgress',
  Completed: 'Completed',
  Archived: 'Archived',
} as const
export type ProjectStatus = (typeof ProjectStatus)[keyof typeof ProjectStatus]

export const ConstructionStatus = {
  Pending: 'Pending',
  InProgress: 'InProgress',
  Completed: 'Completed',
  Delayed: 'Delayed',
} as const
export type ConstructionStatus = (typeof ConstructionStatus)[keyof typeof ConstructionStatus]

export const AcceptanceStatus = {
  Pending: 'Pending',
  Passed: 'Passed',
  Failed: 'Failed',
} as const
export type AcceptanceStatus = (typeof AcceptanceStatus)[keyof typeof AcceptanceStatus]

export const Role = {
  Admin: 'Admin',
  Designer: 'Designer',
  Contractor: 'Contractor',
  Owner: 'Owner',
  ProjectManager: 'ProjectManager',
} as const
export type Role = (typeof Role)[keyof typeof Role]

export const HouseType = ['一室', '两室', '三室', '四室', '别墅', '复式'] as const
export type HouseType = (typeof HouseType)[number]

export const DesignPhaseName = ['方案设计', '效果图', '施工图', '软装方案'] as const
export const MaterialCategory = ['瓷砖', '地板', '油漆', '灯具', '卫浴', '五金', '板材', '其他'] as const
export const MaterialSpace = ['客厅', '卧室', '厨房', '卫生间', '阳台'] as const
export const BudgetCategory = ['Design', 'Material', 'Labor', 'Furniture', 'Appliance', 'Other'] as const
export const ConstructionName = ['拆改', '水电', '木工', '瓦工', '油漆', '安装', '软装'] as const

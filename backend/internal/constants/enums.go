package constants

// DecorStyle 装修风格。
const (
	DecorStyleModern        = "Modern"
	DecorStyleChinese       = "Chinese"
	DecorStyleNordic        = "Nordic"
	DecorStyleJapanese      = "Japanese"
	DecorStyleIndustrial    = "Industrial"
	DecorStyleMediterranean = "Mediterranean"
	DecorStyleMinimalist    = "Minimalist"
)

// PhaseStatus 设计阶段状态。
const (
	PhaseStatusNotStarted = "NotStarted"
	PhaseStatusInProgress = "InProgress"
	PhaseStatusRevision   = "Revision"
	PhaseStatusApproved   = "Approved"
)

// ConstructionPhase 施工节点名称。
const (
	ConstructionPhaseDemolition     = "Demolition"
	ConstructionPhasePlumbing       = "Plumbing"
	ConstructionPhaseCarpentry      = "Carpentry"
	ConstructionPhaseTiling         = "Tiling"
	ConstructionPhasePainting       = "Painting"
	ConstructionPhaseInstallation   = "Installation"
	ConstructionPhaseSoftFurnishing = "SoftFurnishing"
)

// PurchaseStatus 材料采购状态。
const (
	PurchaseStatusNotPurchased = "NotPurchased"
	PurchaseStatusOrdered      = "Ordered"
	PurchaseStatusDelivered    = "Delivered"
	PurchaseStatusInstalled    = "Installed"
)

// ProjectStatus 项目状态。
const (
	ProjectStatusDesigning  = "Designing"
	ProjectStatusQuoting    = "Quoting"
	ProjectStatusInProgress = "InProgress"
	ProjectStatusCompleted  = "Completed"
	ProjectStatusArchived   = "Archived"
)

// ConstructionStatus 施工节点状态。
const (
	ConstructionStatusPending    = "Pending"
	ConstructionStatusInProgress = "InProgress"
	ConstructionStatusCompleted  = "Completed"
	ConstructionStatusDelayed    = "Delayed"
)

// AcceptanceStatus 验收状态。
const (
	AcceptanceStatusPending = "Pending"
	AcceptanceStatusPassed  = "Passed"
	AcceptanceStatusFailed  = "Failed"
)

// Role 系统角色。
const (
	RoleAdmin          = "Admin"
	RoleDesigner       = "Designer"
	RoleContractor     = "Contractor"
	RoleOwner          = "Owner"
	RoleProjectManager = "ProjectManager"
)

var (
	DecorStyles = []string{
		DecorStyleModern, DecorStyleChinese, DecorStyleNordic, DecorStyleJapanese,
		DecorStyleIndustrial, DecorStyleMediterranean, DecorStyleMinimalist,
	}
	PhaseStatuses = []string{
		PhaseStatusNotStarted, PhaseStatusInProgress, PhaseStatusRevision, PhaseStatusApproved,
	}
	ConstructionPhases = []string{
		ConstructionPhaseDemolition, ConstructionPhasePlumbing, ConstructionPhaseCarpentry,
		ConstructionPhaseTiling, ConstructionPhasePainting, ConstructionPhaseInstallation,
		ConstructionPhaseSoftFurnishing,
	}
	PurchaseStatuses = []string{
		PurchaseStatusNotPurchased, PurchaseStatusOrdered, PurchaseStatusDelivered, PurchaseStatusInstalled,
	}
	ProjectStatuses = []string{
		ProjectStatusDesigning, ProjectStatusQuoting, ProjectStatusInProgress, ProjectStatusCompleted, ProjectStatusArchived,
	}
	ConstructionStatuses = []string{
		ConstructionStatusPending, ConstructionStatusInProgress, ConstructionStatusCompleted, ConstructionStatusDelayed,
	}
	AcceptanceStatuses = []string{AcceptanceStatusPending, AcceptanceStatusPassed, AcceptanceStatusFailed}
	Roles              = []string{RoleAdmin, RoleDesigner, RoleContractor, RoleOwner, RoleProjectManager}
)

// Contains 判断字符串是否在集合内。
func Contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

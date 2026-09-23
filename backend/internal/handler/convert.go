package handler

import (
	"encoding/json"

	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/model"
)

func toProjectDTO(project *model.RenovationProject) dto.ProjectDTO {
	return dto.ProjectDTO{
		ID:              project.ID,
		Name:            project.Name,
		HouseType:       project.HouseType,
		Area:            project.Area,
		DecorStyle:      project.DecorStyle,
		Address:         project.Address,
		OwnerID:         project.OwnerID,
		DesignerID:      project.DesignerID,
		ForemanID:       project.ForemanID,
		Status:          project.Status,
		ContractAmount:  project.ContractAmount,
		StartDate:       project.StartDate,
		ExpectedEndDate: project.ExpectedEndDate,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
	}
}

func toDesignDTO(phase *model.DesignPhase) dto.DesignDTO {
	return dto.DesignDTO{
		ID:            phase.ID,
		ProjectID:     phase.ProjectID,
		Name:          phase.Name,
		DesignerID:    phase.DesignerID,
		Status:        phase.Status,
		Version:       phase.Version,
		Description:   phase.Description,
		FileURLs:      parseStringSlice(phase.FileURLs),
		ReviewComment: phase.ReviewComment,
		ReviewerID:    phase.ReviewerID,
		CreatedAt:     phase.CreatedAt,
		UpdatedAt:     phase.UpdatedAt,
	}
}

func toMaterialDTO(item *model.MaterialItem) dto.MaterialDTO {
	return dto.MaterialDTO{
		ID:             item.ID,
		ProjectID:      item.ProjectID,
		Name:           item.Name,
		Category:       item.Category,
		Spec:           item.Spec,
		Brand:          item.Brand,
		Quantity:       item.Quantity,
		Unit:           item.Unit,
		UnitPrice:      item.UnitPrice,
		TotalPrice:     item.TotalPrice,
		PurchaseStatus: item.PurchaseStatus,
		Supplier:       item.Supplier,
		Space:          item.Space,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func toBudgetDTO(item *model.BudgetItem) dto.BudgetDTO {
	return dto.BudgetDTO{
		ID:           item.ID,
		ProjectID:    item.ProjectID,
		Category:     item.Category,
		BudgetAmount: item.BudgetAmount,
		ActualAmount: item.ActualAmount,
		Variance:     item.Variance,
		Remark:       item.Remark,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func toConstructionDTO(node *model.ConstructionNode) dto.ConstructionDTO {
	return dto.ConstructionDTO{
		ID:               node.ID,
		ProjectID:        node.ProjectID,
		Name:             node.Name,
		PlannedStartDate: node.PlannedStartDate,
		PlannedEndDate:   node.PlannedEndDate,
		ActualStartDate:  node.ActualStartDate,
		ActualEndDate:    node.ActualEndDate,
		Status:           node.Status,
		AcceptanceStatus: node.AcceptanceStatus,
		AcceptancePhotos: parseStringSlice(node.AcceptancePhotos),
		AcceptanceNote:   node.AcceptanceNote,
		CreatedAt:        node.CreatedAt,
		UpdatedAt:        node.UpdatedAt,
	}
}

func parseStringSlice(raw string) []string {
	if raw == "" || raw == "null" {
		return []string{}
	}
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return []string{}
	}
	return values
}

func toProjectDTOList(projects []model.RenovationProject) []dto.ProjectDTO {
	out := make([]dto.ProjectDTO, 0, len(projects))
	for _, item := range projects {
		out = append(out, toProjectDTO(&item))
	}
	return out
}

func toDesignDTOList(phases []model.DesignPhase) []dto.DesignDTO {
	out := make([]dto.DesignDTO, 0, len(phases))
	for _, item := range phases {
		out = append(out, toDesignDTO(&item))
	}
	return out
}

func toMaterialDTOList(items []model.MaterialItem) []dto.MaterialDTO {
	out := make([]dto.MaterialDTO, 0, len(items))
	for _, item := range items {
		out = append(out, toMaterialDTO(&item))
	}
	return out
}

func toBudgetDTOList(items []model.BudgetItem) []dto.BudgetDTO {
	out := make([]dto.BudgetDTO, 0, len(items))
	for _, item := range items {
		out = append(out, toBudgetDTO(&item))
	}
	return out
}

func toConstructionDTOList(nodes []model.ConstructionNode) []dto.ConstructionDTO {
	out := make([]dto.ConstructionDTO, 0, len(nodes))
	for _, item := range nodes {
		out = append(out, toConstructionDTO(&item))
	}
	return out
}

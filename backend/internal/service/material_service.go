package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	apperrors "github.com/home-renovation/platform/internal/errors"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
	"github.com/home-renovation/platform/internal/utils"
)

// MaterialService 材料服务接口。
type MaterialService interface {
	Create(req *dto.CreateMaterialRequest) (*model.MaterialItem, error)
	GetByID(id uint) (*model.MaterialItem, error)
	List(projectID uint, category, space string, page, pageSize int) ([]model.MaterialItem, int64, error)
	ListByProjectID(projectID uint) ([]model.MaterialItem, error)
	Update(id uint, req *dto.UpdateMaterialRequest) (*model.MaterialItem, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) (*model.MaterialItem, error)
}

type materialService struct {
	repo   repository.MaterialRepository
	logger *slog.Logger
}

// NewMaterialService 构造材料服务。
func NewMaterialService(repo repository.MaterialRepository, logger *slog.Logger) MaterialService {
	return &materialService{repo: repo, logger: logger}
}

func (s *materialService) Create(req *dto.CreateMaterialRequest) (*model.MaterialItem, error) {
	item := &model.MaterialItem{
		ProjectID:      req.ProjectID,
		Name:           req.Name,
		Category:       req.Category,
		Spec:           req.Spec,
		Brand:          req.Brand,
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		UnitPrice:      req.UnitPrice,
		TotalPrice:     utils.MaterialTotal(req.Quantity, req.UnitPrice),
		PurchaseStatus: constants.PurchaseStatusNotPurchased,
		Supplier:       req.Supplier,
		Space:          req.Space,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("create material item: %w", err)
	}
	return item, nil
}

func (s *materialService) GetByID(id uint) (*model.MaterialItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewNotFound("material item not found")
		}
		return nil, fmt.Errorf("get material item: %w", err)
	}
	return item, nil
}

func (s *materialService) List(projectID uint, category, space string, page, pageSize int) ([]model.MaterialItem, int64, error) {
	items, total, err := s.repo.List(repository.MaterialFilter{ProjectID: projectID, Category: category, Space: space}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list material items: %w", err)
	}
	return items, total, nil
}

func (s *materialService) ListByProjectID(projectID uint) ([]model.MaterialItem, error) {
	items, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list material items by project: %w", err)
	}
	return items, nil
}

func (s *materialService) Update(id uint, req *dto.UpdateMaterialRequest) (*model.MaterialItem, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.Spec != nil {
		item.Spec = *req.Spec
	}
	if req.Brand != nil {
		item.Brand = *req.Brand
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Unit != nil {
		item.Unit = *req.Unit
	}
	if req.UnitPrice != nil {
		item.UnitPrice = *req.UnitPrice
	}
	if req.Supplier != nil {
		item.Supplier = *req.Supplier
	}
	if req.Space != nil {
		item.Space = *req.Space
	}
	item.TotalPrice = utils.MaterialTotal(item.Quantity, item.UnitPrice)
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("update material item: %w", err)
	}
	return item, nil
}

func (s *materialService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete material item: %w", err)
	}
	return nil
}

func (s *materialService) UpdateStatus(id uint, status string) (*model.MaterialItem, error) {
	if !constants.Contains(constants.PurchaseStatuses, status) {
		return nil, apperrors.NewBadRequest("invalid purchase status")
	}
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	order := map[string]int{
		constants.PurchaseStatusNotPurchased: 0,
		constants.PurchaseStatusOrdered:      1,
		constants.PurchaseStatusDelivered:    2,
		constants.PurchaseStatusInstalled:    3,
	}
	if order[status] < order[item.PurchaseStatus] {
		return nil, apperrors.NewConflict("purchase status cannot move backwards")
	}
	item.PurchaseStatus = status
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("update material status: %w", err)
	}
	s.logger.Info("material purchase status changed", "material_id", id, "status", status)
	return item, nil
}

package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/home-renovation/platform/internal/dto"
	apperrors "github.com/home-renovation/platform/internal/errors"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
	"github.com/home-renovation/platform/internal/utils"
)

// BudgetService 预算服务接口。
type BudgetService interface {
	Create(req *dto.CreateBudgetRequest) (*model.BudgetItem, error)
	GetByID(id uint) (*model.BudgetItem, error)
	List(projectID uint, category string, page, pageSize int) ([]model.BudgetItem, int64, error)
	ListByProjectID(projectID uint) ([]model.BudgetItem, error)
	Update(id uint, req *dto.UpdateBudgetRequest) (*model.BudgetItem, error)
	Delete(id uint) error
}

type budgetService struct {
	repo   repository.BudgetRepository
	logger *slog.Logger
}

// NewBudgetService 构造预算服务。
func NewBudgetService(repo repository.BudgetRepository, logger *slog.Logger) BudgetService {
	return &budgetService{repo: repo, logger: logger}
}

func (s *budgetService) Create(req *dto.CreateBudgetRequest) (*model.BudgetItem, error) {
	item := &model.BudgetItem{
		ProjectID:    req.ProjectID,
		Category:     req.Category,
		BudgetAmount: utils.Round2(req.BudgetAmount),
		ActualAmount: utils.Round2(req.ActualAmount),
		Remark:       req.Remark,
	}
	item.Variance = utils.BudgetVariance(item.BudgetAmount, item.ActualAmount)
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("create budget item: %w", err)
	}
	return item, nil
}

func (s *budgetService) GetByID(id uint) (*model.BudgetItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewNotFound("budget item not found")
		}
		return nil, fmt.Errorf("get budget item: %w", err)
	}
	return item, nil
}

func (s *budgetService) List(projectID uint, category string, page, pageSize int) ([]model.BudgetItem, int64, error) {
	items, total, err := s.repo.List(repository.BudgetFilter{ProjectID: projectID, Category: category}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list budget items: %w", err)
	}
	return items, total, nil
}

func (s *budgetService) ListByProjectID(projectID uint) ([]model.BudgetItem, error) {
	items, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list budget items by project: %w", err)
	}
	return items, nil
}

func (s *budgetService) Update(id uint, req *dto.UpdateBudgetRequest) (*model.BudgetItem, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.BudgetAmount != nil {
		item.BudgetAmount = utils.Round2(*req.BudgetAmount)
	}
	if req.ActualAmount != nil {
		item.ActualAmount = utils.Round2(*req.ActualAmount)
	}
	if req.Remark != nil {
		item.Remark = *req.Remark
	}
	item.Variance = utils.BudgetVariance(item.BudgetAmount, item.ActualAmount)
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("update budget item: %w", err)
	}
	s.logger.Info("budget item updated", "budget_id", id, "variance", item.Variance)
	return item, nil
}

func (s *budgetService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete budget item: %w", err)
	}
	return nil
}

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

// materialCostStatuses 计入材料自动花费的采购状态：已到货、已安装。
var materialCostStatuses = []string{
	constants.PurchaseStatusDelivered,
	constants.PurchaseStatusInstalled,
}

// BudgetService 预算服务接口。
type BudgetService interface {
	Create(req *dto.CreateBudgetRequest) (*model.BudgetItem, error)
	GetByID(id uint) (*model.BudgetItem, error)
	GetViewByID(id uint) (*BudgetView, error)
	List(projectID uint, category string, page, pageSize int) ([]BudgetView, int64, error)
	ListByProjectID(projectID uint) ([]BudgetView, error)
	Update(id uint, req *dto.UpdateBudgetRequest) (*model.BudgetItem, error)
	Delete(id uint) error
}

type budgetService struct {
	repo         repository.BudgetRepository
	materialRepo repository.MaterialRepository
	logger       *slog.Logger
}

// NewBudgetService 构造预算服务。
func NewBudgetService(repo repository.BudgetRepository, materialRepo repository.MaterialRepository, logger *slog.Logger) BudgetService {
	return &budgetService{repo: repo, materialRepo: materialRepo, logger: logger}
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

// GetViewByID 获取附带材料自动花费的预算视图。
func (s *budgetService) GetViewByID(id uint) (*BudgetView, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	views, err := s.attachAutoAmount([]model.BudgetItem{*item})
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}

func (s *budgetService) List(projectID uint, category string, page, pageSize int) ([]BudgetView, int64, error) {
	items, total, err := s.repo.List(repository.BudgetFilter{ProjectID: projectID, Category: category}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list budget items: %w", err)
	}
	views, err := s.attachAutoAmount(items)
	if err != nil {
		return nil, 0, err
	}
	return views, total, nil
}

func (s *budgetService) ListByProjectID(projectID uint) ([]BudgetView, error) {
	items, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list budget items by project: %w", err)
	}
	return s.attachAutoAmount(items)
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

// attachAutoAmount 为材料类预算项附加按项目汇总的材料自动花费，
// 并重算实际花费合计与差异；其他类别保持手工金额不变。
// 同一项目存在多条材料预算项时，自动花费只计入第一条，避免合计重复。
func (s *budgetService) attachAutoAmount(items []model.BudgetItem) ([]BudgetView, error) {
	projectIDs := make(map[uint]struct{})
	for _, item := range items {
		if item.Category == constants.BudgetCategoryMaterial {
			projectIDs[item.ProjectID] = struct{}{}
		}
	}

	autoByProject := make(map[uint]float64)
	if len(projectIDs) > 0 {
		ids := make([]uint, 0, len(projectIDs))
		for id := range projectIDs {
			ids = append(ids, id)
		}
		sums, err := s.materialRepo.SumTotalPriceByProjectsAndStatuses(ids, materialCostStatuses)
		if err != nil {
			return nil, fmt.Errorf("sum material auto cost: %w", err)
		}
		autoByProject = sums
	}

	views := make([]BudgetView, 0, len(items))
	materialRowSeen := make(map[uint]struct{})
	for _, item := range items {
		view := BudgetView{BudgetItem: item, ManualActualAmount: item.ActualAmount}
		if item.Category == constants.BudgetCategoryMaterial {
			if _, seen := materialRowSeen[item.ProjectID]; !seen {
				view.AutoActualAmount = utils.Round2(autoByProject[item.ProjectID])
				materialRowSeen[item.ProjectID] = struct{}{}
			}
		}
		total := utils.TotalActualAmount(view.ManualActualAmount, view.AutoActualAmount)
		view.ActualAmount = total
		view.Variance = utils.BudgetVariance(view.BudgetAmount, total)
		views = append(views, view)
	}
	return views, nil
}

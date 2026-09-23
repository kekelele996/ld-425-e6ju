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
	if err := s.recomputeVariance(item); err != nil {
		return nil, err
	}
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("create budget item: %w", err)
	}
	if err := s.enrichAutoAmounts([]*model.BudgetItem{item}); err != nil {
		return nil, err
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
	if err := s.enrichAutoAmounts([]*model.BudgetItem{item}); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *budgetService) List(projectID uint, category string, page, pageSize int) ([]model.BudgetItem, int64, error) {
	items, total, err := s.repo.List(repository.BudgetFilter{ProjectID: projectID, Category: category}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list budget items: %w", err)
	}
	if err := s.enrichAutoAmountList(items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *budgetService) ListByProjectID(projectID uint) ([]model.BudgetItem, error) {
	items, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list budget items by project: %w", err)
	}
	if err := s.enrichAutoAmountList(items); err != nil {
		return nil, err
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
	if err := s.recomputeVariance(item); err != nil {
		return nil, err
	}
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("update budget item: %w", err)
	}
	s.logger.Info("budget item updated", "budget_id", id, "variance", item.Variance)
	if err := s.enrichAutoAmounts([]*model.BudgetItem{item}); err != nil {
		return nil, err
	}
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

// recomputeVariance 依据「材料自动花费 + 手工金额」的合计刷新差异金额，
// 使持久化数据与超支判断口径一致。
func (s *budgetService) recomputeVariance(item *model.BudgetItem) error {
	auto, err := s.materialAutoAmount(item.ProjectID, item.Category)
	if err != nil {
		return err
	}
	item.MaterialAutoAmount = auto
	item.Variance = utils.BudgetVariance(item.BudgetAmount, utils.BudgetTotalActual(item.ActualAmount, auto))
	return nil
}

func (s *budgetService) enrichAutoAmountList(items []model.BudgetItem) error {
	refs := make([]*model.BudgetItem, 0, len(items))
	for i := range items {
		refs = append(refs, &items[i])
	}
	return s.enrichAutoAmounts(refs)
}

// enrichAutoAmounts 按项目实时汇总已到货/已安装材料总价，回填到材料类预算项，
// 材料改价或被移除后金额随查询变化。
func (s *budgetService) enrichAutoAmounts(items []*model.BudgetItem) error {
	projectSet := make(map[uint]struct{})
	for _, item := range items {
		if item.Category == constants.BudgetCategoryMaterial {
			item.MaterialAutoAmount = 0
			projectSet[item.ProjectID] = struct{}{}
		}
	}
	if len(projectSet) == 0 {
		return nil
	}
	projectIDs := make([]uint, 0, len(projectSet))
	for id := range projectSet {
		projectIDs = append(projectIDs, id)
	}
	totals, err := s.materialRepo.SumReceivedTotalByProjectIDs(projectIDs)
	if err != nil {
		return fmt.Errorf("enrich budget material auto amounts: %w", err)
	}
	for _, item := range items {
		if item.Category != constants.BudgetCategoryMaterial {
			continue
		}
		item.MaterialAutoAmount = utils.Round2(totals[item.ProjectID])
	}
	return nil
}

func (s *budgetService) materialAutoAmount(projectID uint, category string) (float64, error) {
	if category != constants.BudgetCategoryMaterial {
		return 0, nil
	}
	totals, err := s.materialRepo.SumReceivedTotalByProjectIDs([]uint{projectID})
	if err != nil {
		return 0, fmt.Errorf("sum material auto amount: %w", err)
	}
	return utils.Round2(totals[projectID]), nil
}

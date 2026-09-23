package service

import (
	"log/slog"
	"math"
	"testing"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// fakeBudgetRepo 预算项测试仓储。
type fakeBudgetRepo struct {
	nextID uint
	items  map[uint]*model.BudgetItem
}

func newFakeBudgetRepo() *fakeBudgetRepo {
	return &fakeBudgetRepo{nextID: 1, items: map[uint]*model.BudgetItem{}}
}

func (f *fakeBudgetRepo) Create(item *model.BudgetItem) error {
	item.ID = f.nextID
	f.nextID++
	f.items[item.ID] = item
	return nil
}

func (f *fakeBudgetRepo) GetByID(id uint) (*model.BudgetItem, error) {
	item, ok := f.items[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return item, nil
}

func (f *fakeBudgetRepo) List(filter repository.BudgetFilter, page, pageSize int) ([]model.BudgetItem, int64, error) {
	var out []model.BudgetItem
	for _, item := range f.items {
		if filter.ProjectID != 0 && item.ProjectID != filter.ProjectID {
			continue
		}
		if filter.Category != "" && item.Category != filter.Category {
			continue
		}
		out = append(out, *item)
	}
	return out, int64(len(out)), nil
}

func (f *fakeBudgetRepo) ListByProjectID(projectID uint) ([]model.BudgetItem, error) {
	var out []model.BudgetItem
	for _, item := range f.items {
		if item.ProjectID == projectID {
			out = append(out, *item)
		}
	}
	return out, nil
}

func (f *fakeBudgetRepo) Update(item *model.BudgetItem) error {
	if _, ok := f.items[item.ID]; !ok {
		return repository.ErrNotFound
	}
	f.items[item.ID] = item
	return nil
}

func (f *fakeBudgetRepo) Delete(id uint) error {
	delete(f.items, id)
	return nil
}

// fakeMaterialRepo 材料测试仓储，可动态调整已到货金额，模拟改价/移除。
type fakeMaterialRepo struct {
	receivedTotal map[uint]float64
}

func newFakeMaterialRepo() *fakeMaterialRepo {
	return &fakeMaterialRepo{receivedTotal: map[uint]float64{}}
}

func (f *fakeMaterialRepo) Create(item *model.MaterialItem) error { return nil }
func (f *fakeMaterialRepo) GetByID(id uint) (*model.MaterialItem, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeMaterialRepo) List(filter repository.MaterialFilter, page, pageSize int) ([]model.MaterialItem, int64, error) {
	return nil, 0, nil
}
func (f *fakeMaterialRepo) ListByProjectID(projectID uint) ([]model.MaterialItem, error) {
	return nil, nil
}
func (f *fakeMaterialRepo) SumReceivedTotalByProjectIDs(projectIDs []uint) (map[uint]float64, error) {
	out := make(map[uint]float64, len(projectIDs))
	for _, id := range projectIDs {
		out[id] = f.receivedTotal[id]
	}
	return out, nil
}
func (f *fakeMaterialRepo) Update(item *model.MaterialItem) error { return nil }
func (f *fakeMaterialRepo) Delete(id uint) error                  { return nil }

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 0.001 }

func TestBudgetService_MaterialAutoAndManualAmounts(t *testing.T) {
	budgetRepo := newFakeBudgetRepo()
	materialRepo := newFakeMaterialRepo()
	materialRepo.receivedTotal[7] = 5000 // 已到货/已安装材料合计
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	item, err := svc.Create(&dto.CreateBudgetRequest{
		ProjectID:    7,
		Category:     constants.BudgetCategoryMaterial,
		BudgetAmount: 6000,
		ActualAmount: 800, // 手工补充金额
	})
	if err != nil {
		t.Fatalf("create material budget: %v", err)
	}
	if !almostEqual(item.MaterialAutoAmount, 5000) {
		t.Fatalf("auto amount: expected 5000, got %.2f", item.MaterialAutoAmount)
	}
	// 差异以合计为准：5000 + 800 - 6000 = -200（未超支）。
	if !almostEqual(item.Variance, -200) {
		t.Fatalf("variance: expected -200, got %.2f", item.Variance)
	}
}

func TestBudgetService_NonMaterialCategoryKeepsManualOnly(t *testing.T) {
	budgetRepo := newFakeBudgetRepo()
	materialRepo := newFakeMaterialRepo()
	materialRepo.receivedTotal[7] = 5000
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	item, err := svc.Create(&dto.CreateBudgetRequest{
		ProjectID:    7,
		Category:     constants.BudgetCategoryLabor,
		BudgetAmount: 3000,
		ActualAmount: 3500,
	})
	if err != nil {
		t.Fatalf("create labor budget: %v", err)
	}
	if item.MaterialAutoAmount != 0 {
		t.Fatalf("non-material category must not carry auto amount, got %.2f", item.MaterialAutoAmount)
	}
	if !almostEqual(item.Variance, 500) {
		t.Fatalf("variance: expected 500, got %.2f", item.Variance)
	}
}

func TestBudgetService_AutoAmountFollowsMaterialChanges(t *testing.T) {
	budgetRepo := newFakeBudgetRepo()
	materialRepo := newFakeMaterialRepo()
	materialRepo.receivedTotal[7] = 5000
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	item, err := svc.Create(&dto.CreateBudgetRequest{
		ProjectID: 7, Category: constants.BudgetCategoryMaterial, BudgetAmount: 10000, ActualAmount: 0,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// 材料改价（或状态推进）后自动金额变化：5000 -> 6200。
	materialRepo.receivedTotal[7] = 6200
	got, err := svc.GetByID(item.ID)
	if err != nil {
		t.Fatalf("get after price change: %v", err)
	}
	if !almostEqual(got.MaterialAutoAmount, 6200) {
		t.Fatalf("after price change: expected 6200, got %.2f", got.MaterialAutoAmount)
	}

	// 材料被移除：6200 -> 2000。
	materialRepo.receivedTotal[7] = 2000
	items, err := svc.ListByProjectID(7)
	if err != nil {
		t.Fatalf("list by project: %v", err)
	}
	if len(items) != 1 || !almostEqual(items[0].MaterialAutoAmount, 2000) {
		t.Fatalf("after removal: expected auto amount 2000, got %+v", items)
	}

	// 超支判断以合计为准：2000 自动 + 手工 9000 - 预算 10000 = 1000。
	if _, err := svc.Update(item.ID, &dto.UpdateBudgetRequest{ActualAmount: floatPtr(9000)}); err != nil {
		t.Fatalf("update manual amount: %v", err)
	}
	got, err = svc.GetByID(item.ID)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if !almostEqual(got.Variance, 1000) {
		t.Fatalf("variance after update: expected 1000, got %.2f", got.Variance)
	}
}

func floatPtr(v float64) *float64 { return &v }

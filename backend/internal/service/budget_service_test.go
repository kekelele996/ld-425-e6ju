package service

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// fakeBudgetRepo 预算服务测试用内存仓储。
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
	copied := *item
	f.items[item.ID] = &copied
	return nil
}

func (f *fakeBudgetRepo) GetByID(id uint) (*model.BudgetItem, error) {
	item, ok := f.items[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copied := *item
	return &copied, nil
}

func (f *fakeBudgetRepo) List(filter repository.BudgetFilter, page, pageSize int) ([]model.BudgetItem, int64, error) {
	items, err := f.ListByProjectID(filter.ProjectID)
	if err != nil {
		return nil, 0, err
	}
	if filter.Category != "" {
		filtered := items[:0]
		for _, item := range items {
			if item.Category == filter.Category {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}
	return items, int64(len(items)), nil
}

func (f *fakeBudgetRepo) ListByProjectID(projectID uint) ([]model.BudgetItem, error) {
	var out []model.BudgetItem
	for _, item := range f.items {
		if projectID == 0 || item.ProjectID == projectID {
			out = append(out, *item)
		}
	}
	return out, nil
}

func (f *fakeBudgetRepo) Update(item *model.BudgetItem) error {
	if _, ok := f.items[item.ID]; !ok {
		return repository.ErrNotFound
	}
	copied := *item
	f.items[item.ID] = &copied
	return nil
}

func (f *fakeBudgetRepo) Delete(id uint) error {
	delete(f.items, id)
	return nil
}

// fakeMaterialCostRepo 仅实现自动花费汇总的材料仓储桩。
type fakeMaterialCostRepo struct {
	sums map[uint]float64
}

func (f *fakeMaterialCostRepo) Create(item *model.MaterialItem) error { return errors.New("not implemented") }
func (f *fakeMaterialCostRepo) GetByID(id uint) (*model.MaterialItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMaterialCostRepo) List(filter repository.MaterialFilter, page, pageSize int) ([]model.MaterialItem, int64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeMaterialCostRepo) ListByProjectID(projectID uint) ([]model.MaterialItem, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMaterialCostRepo) SumTotalPriceByProjectsAndStatuses(projectIDs []uint, statuses []string) (map[uint]float64, error) {
	out := make(map[uint]float64)
	for _, id := range projectIDs {
		if sum, ok := f.sums[id]; ok {
			out[id] = sum
		}
	}
	return out, nil
}
func (f *fakeMaterialCostRepo) Update(item *model.MaterialItem) error { return errors.New("not implemented") }
func (f *fakeMaterialCostRepo) Delete(id uint) error                 { return errors.New("not implemented") }

func TestBudgetService_AutoMaterialCost(t *testing.T) {
	tests := []struct {
		name            string
		category        string
		manual          float64
		auto            float64
		wantAuto        float64
		wantManual      float64
		wantActual      float64
		wantOverBudget  bool
	}{
		{
			name:           "material row combines manual and delivered/installed cost",
			category:       constants.BudgetCategoryMaterial,
			manual:         1000,
			auto:           24318,
			wantAuto:       24318,
			wantManual:     1000,
			wantActual:     25318,
			wantOverBudget: false,
		},
		{
			name:           "material row over budget judged by total actual",
			category:       constants.BudgetCategoryMaterial,
			manual:         60000,
			auto:           70000,
			wantAuto:       70000,
			wantManual:     60000,
			wantActual:     130000,
			wantOverBudget: true,
		},
		{
			name:           "non-material row keeps manual amount only",
			category:       constants.BudgetCategoryLabor,
			manual:         30000,
			auto:           24318,
			wantAuto:       0,
			wantManual:     30000,
			wantActual:     30000,
			wantOverBudget: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			budgetRepo := newFakeBudgetRepo()
			materialRepo := &fakeMaterialCostRepo{sums: map[uint]float64{1: tt.auto}}
			svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

			created, err := svc.Create(&dto.CreateBudgetRequest{
				ProjectID: 1, Category: tt.category, BudgetAmount: 120000, ActualAmount: tt.manual,
			})
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			view, err := svc.GetViewByID(created.ID)
			if err != nil {
				t.Fatalf("get view: %v", err)
			}
			if view.AutoActualAmount != tt.wantAuto {
				t.Fatalf("auto actual: want %.2f, got %.2f", tt.wantAuto, view.AutoActualAmount)
			}
			if view.ManualActualAmount != tt.wantManual {
				t.Fatalf("manual actual: want %.2f, got %.2f", tt.wantManual, view.ManualActualAmount)
			}
			if view.TotalActualAmount() != tt.wantActual {
				t.Fatalf("total actual: want %.2f, got %.2f", tt.wantActual, view.TotalActualAmount())
			}
			if view.ActualAmount != tt.wantActual {
				t.Fatalf("embedded actual amount should equal total: want %.2f, got %.2f", tt.wantActual, view.ActualAmount)
			}
			if (view.Variance > 0) != tt.wantOverBudget {
				t.Fatalf("over budget: want %v, got variance %.2f", tt.wantOverBudget, view.Variance)
			}
		})
	}
}

func TestBudgetService_AutoCostReflectsMaterialChanges(t *testing.T) {
	budgetRepo := newFakeBudgetRepo()
	materialRepo := &fakeMaterialCostRepo{sums: map[uint]float64{1: 20000}}
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	created, err := svc.Create(&dto.CreateBudgetRequest{
		ProjectID: 1, Category: constants.BudgetCategoryMaterial, BudgetAmount: 120000, ActualAmount: 1000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// 材料改价/到货状态变化后，汇总金额随之变化（模拟仓储返回新值）。
	materialRepo.sums[1] = 30000
	view, err := svc.GetViewByID(created.ID)
	if err != nil {
		t.Fatalf("get view: %v", err)
	}
	if view.AutoActualAmount != 30000 || view.TotalActualAmount() != 31000 {
		t.Fatalf("expected auto 30000 / total 31000, got auto %.2f / total %.2f", view.AutoActualAmount, view.TotalActualAmount())
	}

	// 材料被移除后自动花费归零，手工补充金额保留。
	delete(materialRepo.sums, 1)
	view, err = svc.GetViewByID(created.ID)
	if err != nil {
		t.Fatalf("get view after removal: %v", err)
	}
	if view.AutoActualAmount != 0 || view.ManualActualAmount != 1000 || view.TotalActualAmount() != 1000 {
		t.Fatalf("expected auto 0 / manual 1000 / total 1000, got auto %.2f / manual %.2f / total %.2f",
			view.AutoActualAmount, view.ManualActualAmount, view.TotalActualAmount())
	}
}

func TestBudgetService_ListScopesAutoCostByProject(t *testing.T) {	budgetRepo := newFakeBudgetRepo()
	materialRepo := &fakeMaterialCostRepo{sums: map[uint]float64{1: 10000, 2: 5000}}
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	if _, err := svc.Create(&dto.CreateBudgetRequest{ProjectID: 1, Category: constants.BudgetCategoryMaterial}); err != nil {
		t.Fatalf("create material budget p1: %v", err)
	}
	if _, err := svc.Create(&dto.CreateBudgetRequest{ProjectID: 2, Category: constants.BudgetCategoryMaterial}); err != nil {
		t.Fatalf("create material budget p2: %v", err)
	}
	if _, err := svc.Create(&dto.CreateBudgetRequest{ProjectID: 1, Category: constants.BudgetCategoryLabor, ActualAmount: 3000}); err != nil {
		t.Fatalf("create labor budget: %v", err)
	}

	views, err := svc.ListByProjectID(1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, view := range views {
		switch view.Category {
		case constants.BudgetCategoryMaterial:
			if view.AutoActualAmount != 10000 {
				t.Fatalf("project 1 material auto want 10000, got %.2f", view.AutoActualAmount)
			}
		case constants.BudgetCategoryLabor:
			if view.AutoActualAmount != 0 || view.ManualActualAmount != 3000 {
				t.Fatalf("labor row must ignore material auto cost")
			}
		}
	}
}

func TestBudgetService_AutoCostNotDuplicatedAcrossMaterialRows(t *testing.T) {
	budgetRepo := newFakeBudgetRepo()
	materialRepo := &fakeMaterialCostRepo{sums: map[uint]float64{1: 10000}}
	svc := NewBudgetService(budgetRepo, materialRepo, slog.Default())

	if _, err := svc.Create(&dto.CreateBudgetRequest{ProjectID: 1, Category: constants.BudgetCategoryMaterial}); err != nil {
		t.Fatalf("create first material budget: %v", err)
	}
	if _, err := svc.Create(&dto.CreateBudgetRequest{ProjectID: 1, Category: constants.BudgetCategoryMaterial}); err != nil {
		t.Fatalf("create second material budget: %v", err)
	}

	views, err := svc.ListByProjectID(1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var autoSum float64
	for _, view := range views {
		autoSum += view.AutoActualAmount
	}
	if autoSum != 10000 {
		t.Fatalf("project auto cost must be counted once, want 10000, got %.2f", autoSum)
	}
}

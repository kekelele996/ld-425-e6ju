package repository

import (
	"math"
	"testing"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/model"
)

func TestMaterialRepository_SumReceivedTotalByProjectIDs(t *testing.T) {
	db := newTestDB(t)
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repo := NewMaterialRepository(db)

	materials := []model.MaterialItem{
		{ProjectID: 1, Name: "已到货瓷砖", PurchaseStatus: constants.PurchaseStatusDelivered, TotalPrice: 1000},
		{ProjectID: 1, Name: "已安装灯具", PurchaseStatus: constants.PurchaseStatusInstalled, TotalPrice: 2500.55},
		{ProjectID: 1, Name: "未采购板材", PurchaseStatus: constants.PurchaseStatusNotPurchased, TotalPrice: 9999},
		{ProjectID: 1, Name: "已下单油漆", PurchaseStatus: constants.PurchaseStatusOrdered, TotalPrice: 8888},
		{ProjectID: 2, Name: "他项目已到货卫浴", PurchaseStatus: constants.PurchaseStatusDelivered, TotalPrice: 3000},
	}
	for i := range materials {
		if err := repo.Create(&materials[i]); err != nil {
			t.Fatalf("create material: %v", err)
		}
	}

	totals, err := repo.SumReceivedTotalByProjectIDs([]uint{1, 2, 3})
	if err != nil {
		t.Fatalf("sum received totals: %v", err)
	}

	wantProject1 := 3500.55
	if math.Abs(totals[1]-wantProject1) > 0.001 {
		t.Fatalf("project 1: expected %.2f, got %.2f", wantProject1, totals[1])
	}
	if math.Abs(totals[2]-3000) > 0.001 {
		t.Fatalf("project 2: expected 3000.00, got %.2f", totals[2])
	}
	if _, ok := totals[3]; ok {
		t.Fatalf("project 3 without materials should be absent, got %.2f", totals[3])
	}

	// 改价后汇总随之变化。
	materials[0].TotalPrice = 4000
	if err := repo.Update(&materials[0]); err != nil {
		t.Fatalf("update material: %v", err)
	}
	totals, err = repo.SumReceivedTotalByProjectIDs([]uint{1})
	if err != nil {
		t.Fatalf("re-sum: %v", err)
	}
	if math.Abs(totals[1]-6500.55) > 0.001 {
		t.Fatalf("after price change: expected 6500.55, got %.2f", totals[1])
	}

	// 移除材料后汇总随之减少。
	if err := repo.Delete(materials[0].ID); err != nil {
		t.Fatalf("delete material: %v", err)
	}
	totals, err = repo.SumReceivedTotalByProjectIDs([]uint{1})
	if err != nil {
		t.Fatalf("sum after delete: %v", err)
	}
	if math.Abs(totals[1]-2500.55) > 0.001 {
		t.Fatalf("after delete: expected 2500.55, got %.2f", totals[1])
	}

	// 空入参不应查库也不报错。
	empty, err := repo.SumReceivedTotalByProjectIDs(nil)
	if err != nil {
		t.Fatalf("empty input: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty map, got %v", empty)
	}
}

package repository

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

func newMaterialTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestMaterialRepository_SumDeliveredAndInstalled(t *testing.T) {
	db := newMaterialTestDB(t)
	repo := NewMaterialRepository(db)

	items := []model.MaterialItem{
		{ProjectID: 1, Name: "已到货", TotalPrice: 1000, PurchaseStatus: constants.PurchaseStatusDelivered},
		{ProjectID: 1, Name: "已安装", TotalPrice: 2000, PurchaseStatus: constants.PurchaseStatusInstalled},
		{ProjectID: 1, Name: "未采购", TotalPrice: 4000, PurchaseStatus: constants.PurchaseStatusNotPurchased},
		{ProjectID: 1, Name: "已下单", TotalPrice: 8000, PurchaseStatus: constants.PurchaseStatusOrdered},
		{ProjectID: 2, Name: "其他项目已安装", TotalPrice: 500, PurchaseStatus: constants.PurchaseStatusInstalled},
	}
	for i := range items {
		if err := repo.Create(&items[i]); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	tests := []struct {
		name       string
		projectIDs []uint
		want       map[uint]float64
	}{
		{
			name:       "project 1 only counts delivered and installed",
			projectIDs: []uint{1},
			want:       map[uint]float64{1: 3000},
		},
		{
			name:       "multiple projects grouped separately",
			projectIDs: []uint{1, 2},
			want:       map[uint]float64{1: 3000, 2: 500},
		},
		{
			name:       "project without qualifying materials is absent",
			projectIDs: []uint{3},
			want:       map[uint]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.SumTotalPriceByProjectsAndStatuses(tt.projectIDs, []string{
				constants.PurchaseStatusDelivered, constants.PurchaseStatusInstalled,
			})
			if err != nil {
				t.Fatalf("sum: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("want %d groups, got %d (%v)", len(tt.want), len(got), got)
			}
			for projectID, wantSum := range tt.want {
				if got[projectID] != wantSum {
					t.Fatalf("project %d: want %.2f, got %.2f", projectID, wantSum, got[projectID])
				}
			}
		})
	}
}

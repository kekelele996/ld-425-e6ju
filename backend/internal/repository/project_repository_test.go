package repository

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.RenovationProject{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestProjectRepository_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	project := &model.RenovationProject{Name: "测试项目", HouseType: "三室", Area: 120, DecorStyle: "Modern", Status: "Designing"}

	if err := repo.Create(project); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.GetByID(project.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != project.Name {
		t.Fatalf("expected %q, got %q", project.Name, got.Name)
	}
}

func TestProjectRepository_GetByIDNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)

	if _, err := repo.GetByID(9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestProjectRepository_List(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	for i := 0; i < 3; i++ {
		if err := repo.Create(&model.RenovationProject{Name: "项目", HouseType: "一室", Area: 50, DecorStyle: "Nordic", Status: "Designing"}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	projects, total, err := repo.List(ProjectFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(projects) != 2 {
		t.Fatalf("expected page size 2, got %d", len(projects))
	}
}

func TestProjectRepository_UpdateAndDelete(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	project := &model.RenovationProject{Name: "旧名", HouseType: "两室", Area: 80, DecorStyle: "Japanese", Status: "Designing"}
	if err := repo.Create(project); err != nil {
		t.Fatalf("create: %v", err)
	}

	project.Name = "新名"
	if err := repo.Update(project); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := repo.GetByID(project.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "新名" {
		t.Fatalf("expected 新名, got %q", got.Name)
	}

	if err := repo.Delete(project.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.GetByID(project.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

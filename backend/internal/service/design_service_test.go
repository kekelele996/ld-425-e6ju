package service

import (
	"log/slog"
	"testing"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// fakeDesignRepo 测试用内存仓储。
type fakeDesignRepo struct {
	nextID uint
	items  map[uint]*model.DesignPhase
}

func newFakeDesignRepo() *fakeDesignRepo {
	return &fakeDesignRepo{nextID: 1, items: map[uint]*model.DesignPhase{}}
}

func (f *fakeDesignRepo) Create(phase *model.DesignPhase) error {
	phase.ID = f.nextID
	f.nextID++
	f.items[phase.ID] = phase
	return nil
}

func (f *fakeDesignRepo) GetByID(id uint) (*model.DesignPhase, error) {
	item, ok := f.items[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return item, nil
}

func (f *fakeDesignRepo) ListByProjectID(projectID uint) ([]model.DesignPhase, error) {
	var out []model.DesignPhase
	for _, item := range f.items {
		if item.ProjectID == projectID {
			out = append(out, *item)
		}
	}
	return out, nil
}

func (f *fakeDesignRepo) ListAll(page, pageSize int) ([]model.DesignPhase, int64, error) {
	var out []model.DesignPhase
	for _, item := range f.items {
		out = append(out, *item)
	}
	return out, int64(len(out)), nil
}

func (f *fakeDesignRepo) Update(phase *model.DesignPhase) error {
	if _, ok := f.items[phase.ID]; !ok {
		return repository.ErrNotFound
	}
	f.items[phase.ID] = phase
	return nil
}

func (f *fakeDesignRepo) Delete(id uint) error {
	delete(f.items, id)
	return nil
}

func TestDesignService_Review(t *testing.T) {
	tests := []struct {
		name     string
		approved bool
		want     string
	}{
		{name: "approve design", approved: true, want: constants.PhaseStatusApproved},
		{name: "reject design", approved: false, want: constants.PhaseStatusRevision},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeDesignRepo()
			svc := NewDesignService(repo, slog.Default())
			phase, err := svc.Create(&dto.CreateDesignRequest{ProjectID: 1, Name: "方案设计", DesignerID: 2})
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			if _, err := svc.Submit(phase.ID, &dto.SubmitDesignRequest{}); err != nil {
				t.Fatalf("submit: %v", err)
			}
			got, err := svc.Review(phase.ID, 4, &dto.ReviewDesignRequest{Approved: tt.approved, Comment: "审核意见"})
			if err != nil {
				t.Fatalf("review: %v", err)
			}
			if got.Status != tt.want {
				t.Fatalf("expected status %q, got %q", tt.want, got.Status)
			}
		})
	}
}

func TestDesignService_SubmitIncrementsVersion(t *testing.T) {
	repo := newFakeDesignRepo()
	svc := NewDesignService(repo, slog.Default())
	phase, err := svc.Create(&dto.CreateDesignRequest{ProjectID: 1, Name: "效果图"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(phase.ID, &dto.SubmitDesignRequest{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if _, err := svc.Review(phase.ID, 4, &dto.ReviewDesignRequest{Approved: false}); err != nil {
		t.Fatalf("review: %v", err)
	}
	got, err := svc.Submit(phase.ID, &dto.SubmitDesignRequest{Description: "修改后重新提交"})
	if err != nil {
		t.Fatalf("resubmit: %v", err)
	}
	if got.Version != 2 {
		t.Fatalf("expected version 2, got %d", got.Version)
	}
}

func TestDesignService_SubmitApprovedConflict(t *testing.T) {
	repo := newFakeDesignRepo()
	svc := NewDesignService(repo, slog.Default())
	phase, err := svc.Create(&dto.CreateDesignRequest{ProjectID: 1, Name: "施工图"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Submit(phase.ID, &dto.SubmitDesignRequest{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if _, err := svc.Review(phase.ID, 4, &dto.ReviewDesignRequest{Approved: true}); err != nil {
		t.Fatalf("review: %v", err)
	}
	if _, err := svc.Submit(phase.ID, &dto.SubmitDesignRequest{}); err == nil {
		t.Fatal("expected conflict error when resubmitting approved design")
	}
}

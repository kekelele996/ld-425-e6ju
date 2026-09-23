package repository

import (
	"errors"
	"fmt"

	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

// DesignRepository 设计阶段仓储接口。
type DesignRepository interface {
	Create(phase *model.DesignPhase) error
	GetByID(id uint) (*model.DesignPhase, error)
	ListByProjectID(projectID uint) ([]model.DesignPhase, error)
	ListAll(page, pageSize int) ([]model.DesignPhase, int64, error)
	Update(phase *model.DesignPhase) error
	Delete(id uint) error
}

type designRepository struct {
	db *gorm.DB
}

// NewDesignRepository 构造设计阶段仓储。
func NewDesignRepository(db *gorm.DB) DesignRepository {
	return &designRepository{db: db}
}

func (r *designRepository) Create(phase *model.DesignPhase) error {
	if err := r.db.Create(phase).Error; err != nil {
		return fmt.Errorf("create design phase: %w", err)
	}
	return nil
}

func (r *designRepository) GetByID(id uint) (*model.DesignPhase, error) {
	var phase model.DesignPhase
	if err := r.db.First(&phase, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get design phase %d: %w", id, err)
	}
	return &phase, nil
}

func (r *designRepository) ListByProjectID(projectID uint) ([]model.DesignPhase, error) {
	var phases []model.DesignPhase
	if err := r.db.Where("project_id = ?", projectID).Order("id ASC").Find(&phases).Error; err != nil {
		return nil, fmt.Errorf("list design phases by project %d: %w", projectID, err)
	}
	return phases, nil
}

func (r *designRepository) ListAll(page, pageSize int) ([]model.DesignPhase, int64, error) {
	var phases []model.DesignPhase
	var total int64
	if err := r.db.Model(&model.DesignPhase{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count design phases: %w", err)
	}
	offset := (page - 1) * pageSize
	if err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&phases).Error; err != nil {
		return nil, 0, fmt.Errorf("list design phases: %w", err)
	}
	return phases, total, nil
}

func (r *designRepository) Update(phase *model.DesignPhase) error {
	if err := r.db.Save(phase).Error; err != nil {
		return fmt.Errorf("update design phase %d: %w", phase.ID, err)
	}
	return nil
}

func (r *designRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.DesignPhase{}, id).Error; err != nil {
		return fmt.Errorf("delete design phase %d: %w", id, err)
	}
	return nil
}

package repository

import (
	"errors"
	"fmt"

	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

// ErrNotFound 哨兵错误：记录不存在。
var ErrNotFound = errors.New("record not found")

// ProjectRepository 装修项目仓储接口。
type ProjectRepository interface {
	Create(project *model.RenovationProject) error
	GetByID(id uint) (*model.RenovationProject, error)
	List(filter ProjectFilter, page, pageSize int) ([]model.RenovationProject, int64, error)
	ListAll() ([]model.RenovationProject, error)
	Update(project *model.RenovationProject) error
	Delete(id uint) error
}

// ProjectFilter 项目查询过滤条件。
type ProjectFilter struct {
	Status string
	Keyword string
}

type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 构造项目仓储。
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *model.RenovationProject) error {
	if err := r.db.Create(project).Error; err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *projectRepository) GetByID(id uint) (*model.RenovationProject, error) {
	var project model.RenovationProject
	if err := r.db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get project %d: %w", id, err)
	}
	return &project, nil
}

func (r *projectRepository) List(filter ProjectFilter, page, pageSize int) ([]model.RenovationProject, int64, error) {
	var projects []model.RenovationProject
	var total int64
	query := r.db.Model(&model.RenovationProject{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+filter.Keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	return projects, total, nil
}

func (r *projectRepository) ListAll() ([]model.RenovationProject, error) {
	var projects []model.RenovationProject
	if err := r.db.Order("id DESC").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("list all projects: %w", err)
	}
	return projects, nil
}

func (r *projectRepository) Update(project *model.RenovationProject) error {
	if err := r.db.Save(project).Error; err != nil {
		return fmt.Errorf("update project %d: %w", project.ID, err)
	}
	return nil
}

func (r *projectRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.RenovationProject{}, id).Error; err != nil {
		return fmt.Errorf("delete project %d: %w", id, err)
	}
	return nil
}

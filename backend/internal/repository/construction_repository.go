package repository

import (
	"errors"
	"fmt"

	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

// ConstructionRepository 施工节点仓储接口。
type ConstructionRepository interface {
	Create(node *model.ConstructionNode) error
	GetByID(id uint) (*model.ConstructionNode, error)
	List(filter ConstructionFilter, page, pageSize int) ([]model.ConstructionNode, int64, error)
	ListByProjectID(projectID uint) ([]model.ConstructionNode, error)
	Update(node *model.ConstructionNode) error
	Delete(id uint) error
}

// ConstructionFilter 施工节点查询过滤条件。
type ConstructionFilter struct {
	ProjectID uint
	Status    string
}

type constructionRepository struct {
	db *gorm.DB
}

// NewConstructionRepository 构造施工节点仓储。
func NewConstructionRepository(db *gorm.DB) ConstructionRepository {
	return &constructionRepository{db: db}
}

func (r *constructionRepository) Create(node *model.ConstructionNode) error {
	if err := r.db.Create(node).Error; err != nil {
		return fmt.Errorf("create construction node: %w", err)
	}
	return nil
}

func (r *constructionRepository) GetByID(id uint) (*model.ConstructionNode, error) {
	var node model.ConstructionNode
	if err := r.db.First(&node, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get construction node %d: %w", id, err)
	}
	return &node, nil
}

func (r *constructionRepository) List(filter ConstructionFilter, page, pageSize int) ([]model.ConstructionNode, int64, error) {
	var nodes []model.ConstructionNode
	var total int64
	query := r.db.Model(&model.ConstructionNode{})
	if filter.ProjectID != 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count construction nodes: %w", err)
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&nodes).Error; err != nil {
		return nil, 0, fmt.Errorf("list construction nodes: %w", err)
	}
	return nodes, total, nil
}

func (r *constructionRepository) ListByProjectID(projectID uint) ([]model.ConstructionNode, error) {
	var nodes []model.ConstructionNode
	if err := r.db.Where("project_id = ?", projectID).Order("id ASC").Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("list construction nodes by project %d: %w", projectID, err)
	}
	return nodes, nil
}

func (r *constructionRepository) Update(node *model.ConstructionNode) error {
	if err := r.db.Save(node).Error; err != nil {
		return fmt.Errorf("update construction node %d: %w", node.ID, err)
	}
	return nil
}

func (r *constructionRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.ConstructionNode{}, id).Error; err != nil {
		return fmt.Errorf("delete construction node %d: %w", id, err)
	}
	return nil
}

package repository

import (
	"errors"
	"fmt"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/model"
	"gorm.io/gorm"
)

// MaterialRepository 材料清单仓储接口。
type MaterialRepository interface {
	Create(item *model.MaterialItem) error
	GetByID(id uint) (*model.MaterialItem, error)
	List(filter MaterialFilter, page, pageSize int) ([]model.MaterialItem, int64, error)
	ListByProjectID(projectID uint) ([]model.MaterialItem, error)
	// SumReceivedTotalByProjectIDs 按项目汇总已到货（Delivered/Installed）材料的总价。
	SumReceivedTotalByProjectIDs(projectIDs []uint) (map[uint]float64, error)
	Update(item *model.MaterialItem) error
	Delete(id uint) error
}

// MaterialFilter 材料查询过滤条件。
type MaterialFilter struct {
	ProjectID uint
	Category  string
	Space     string
}

type materialRepository struct {
	db *gorm.DB
}

// NewMaterialRepository 构造材料仓储。
func NewMaterialRepository(db *gorm.DB) MaterialRepository {
	return &materialRepository{db: db}
}

func (r *materialRepository) Create(item *model.MaterialItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return fmt.Errorf("create material item: %w", err)
	}
	return nil
}

func (r *materialRepository) GetByID(id uint) (*model.MaterialItem, error) {
	var item model.MaterialItem
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get material item %d: %w", id, err)
	}
	return &item, nil
}

func (r *materialRepository) List(filter MaterialFilter, page, pageSize int) ([]model.MaterialItem, int64, error) {
	var items []model.MaterialItem
	var total int64
	query := r.db.Model(&model.MaterialItem{})
	if filter.ProjectID != 0 {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Space != "" {
		query = query.Where("space = ?", filter.Space)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count material items: %w", err)
	}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list material items: %w", err)
	}
	return items, total, nil
}

func (r *materialRepository) ListByProjectID(projectID uint) ([]model.MaterialItem, error) {
	var items []model.MaterialItem
	if err := r.db.Where("project_id = ?", projectID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list material items by project %d: %w", projectID, err)
	}
	return items, nil
}

// SumReceivedTotalByProjectIDs 按项目汇总已到货（Delivered/Installed）材料的总价，
// 材料改价或被移除后随查询实时反映。
func (r *materialRepository) SumReceivedTotalByProjectIDs(projectIDs []uint) (map[uint]float64, error) {
	result := make(map[uint]float64)
	if len(projectIDs) == 0 {
		return result, nil
	}
	type row struct {
		ProjectID uint
		Total     float64
	}
	var rows []row
	err := r.db.Model(&model.MaterialItem{}).
		Select("project_id AS project_id, COALESCE(SUM(total_price), 0) AS total").
		Where("project_id IN ?", projectIDs).
		Where("purchase_status IN ?", []string{
			constants.PurchaseStatusDelivered,
			constants.PurchaseStatusInstalled,
		}).
		Group("project_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sum received material totals: %w", err)
	}
	for _, r := range rows {
		result[r.ProjectID] = r.Total
	}
	return result, nil
}

func (r *materialRepository) Update(item *model.MaterialItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return fmt.Errorf("update material item %d: %w", item.ID, err)
	}
	return nil
}

func (r *materialRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.MaterialItem{}, id).Error; err != nil {
		return fmt.Errorf("delete material item %d: %w", id, err)
	}
	return nil
}

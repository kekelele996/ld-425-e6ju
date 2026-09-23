package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	apperrors "github.com/home-renovation/platform/internal/errors"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// ConstructionService 施工节点服务接口。
type ConstructionService interface {
	Create(req *dto.CreateConstructionRequest) (*model.ConstructionNode, error)
	GetByID(id uint) (*model.ConstructionNode, error)
	List(projectID uint, status string, page, pageSize int) ([]model.ConstructionNode, int64, error)
	ListByProjectID(projectID uint) ([]model.ConstructionNode, error)
	Update(id uint, req *dto.UpdateConstructionRequest) (*model.ConstructionNode, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) (*model.ConstructionNode, error)
	Accept(id uint, req *dto.AcceptConstructionRequest) (*model.ConstructionNode, error)
}

type constructionService struct {
	repo   repository.ConstructionRepository
	logger *slog.Logger
}

// NewConstructionService 构造施工服务。
func NewConstructionService(repo repository.ConstructionRepository, logger *slog.Logger) ConstructionService {
	return &constructionService{repo: repo, logger: logger}
}

func (s *constructionService) Create(req *dto.CreateConstructionRequest) (*model.ConstructionNode, error) {
	start, err := parseDate(req.PlannedStartDate)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid planned_start_date, expect YYYY-MM-DD")
	}
	end, err := parseDate(req.PlannedEndDate)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid planned_end_date, expect YYYY-MM-DD")
	}
	node := &model.ConstructionNode{
		ProjectID:        req.ProjectID,
		Name:             req.Name,
		PlannedStartDate: start,
		PlannedEndDate:   end,
		Status:           constants.ConstructionStatusPending,
		AcceptanceStatus: constants.AcceptanceStatusPending,
		AcceptancePhotos: "[]",
	}
	if err := s.repo.Create(node); err != nil {
		return nil, fmt.Errorf("create construction node: %w", err)
	}
	return node, nil
}

func (s *constructionService) GetByID(id uint) (*model.ConstructionNode, error) {
	node, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewNotFound("construction node not found")
		}
		return nil, fmt.Errorf("get construction node: %w", err)
	}
	return node, nil
}

func (s *constructionService) List(projectID uint, status string, page, pageSize int) ([]model.ConstructionNode, int64, error) {
	nodes, total, err := s.repo.List(repository.ConstructionFilter{ProjectID: projectID, Status: status}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list construction nodes: %w", err)
	}
	return nodes, total, nil
}

func (s *constructionService) ListByProjectID(projectID uint) ([]model.ConstructionNode, error) {
	nodes, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list construction nodes by project: %w", err)
	}
	return nodes, nil
}

func (s *constructionService) Update(id uint, req *dto.UpdateConstructionRequest) (*model.ConstructionNode, error) {
	node, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		node.Name = *req.Name
	}
	if req.PlannedStartDate != nil {
		t, err := parseDate(req.PlannedStartDate)
		if err != nil {
			return nil, apperrors.NewBadRequest("invalid planned_start_date, expect YYYY-MM-DD")
		}
		node.PlannedStartDate = t
	}
	if req.PlannedEndDate != nil {
		t, err := parseDate(req.PlannedEndDate)
		if err != nil {
			return nil, apperrors.NewBadRequest("invalid planned_end_date, expect YYYY-MM-DD")
		}
		node.PlannedEndDate = t
	}
	if err := s.repo.Update(node); err != nil {
		return nil, fmt.Errorf("update construction node: %w", err)
	}
	return node, nil
}

func (s *constructionService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete construction node: %w", err)
	}
	return nil
}

func (s *constructionService) UpdateStatus(id uint, status string) (*model.ConstructionNode, error) {
	if !constants.Contains(constants.ConstructionStatuses, status) {
		return nil, apperrors.NewBadRequest("invalid construction status")
	}
	node, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	switch status {
	case constants.ConstructionStatusInProgress:
		if node.ActualStartDate == nil {
			node.ActualStartDate = &now
		}
	case constants.ConstructionStatusCompleted:
		if node.ActualStartDate == nil {
			node.ActualStartDate = &now
		}
		if node.ActualEndDate == nil {
			node.ActualEndDate = &now
		}
	case constants.ConstructionStatusDelayed:
		if node.Status == constants.ConstructionStatusCompleted {
			return nil, apperrors.NewConflict("completed node cannot be delayed")
		}
	}
	node.Status = status
	if err := s.repo.Update(node); err != nil {
		return nil, fmt.Errorf("update construction status: %w", err)
	}
	s.logger.Info("construction status changed", "node_id", id, "status", status)
	return node, nil
}

func (s *constructionService) Accept(id uint, req *dto.AcceptConstructionRequest) (*model.ConstructionNode, error) {
	node, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if node.Status != constants.ConstructionStatusCompleted {
		return nil, apperrors.NewConflict("only completed node can be accepted")
	}
	if req.Accepted {
		node.AcceptanceStatus = constants.AcceptanceStatusPassed
	} else {
		node.AcceptanceStatus = constants.AcceptanceStatusFailed
	}
	node.AcceptancePhotos = marshalStrings(req.Photos)
	node.AcceptanceNote = req.Note
	if err := s.repo.Update(node); err != nil {
		return nil, fmt.Errorf("accept construction node: %w", err)
	}
	s.logger.Info("construction node accepted", "node_id", id, "accepted", req.Accepted)
	return node, nil
}

package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	apperrors "github.com/home-renovation/platform/internal/errors"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// DesignService 设计阶段服务接口。
type DesignService interface {
	Create(req *dto.CreateDesignRequest) (*model.DesignPhase, error)
	GetByID(id uint) (*model.DesignPhase, error)
	ListByProjectID(projectID uint) ([]model.DesignPhase, error)
	List(page, pageSize int) ([]model.DesignPhase, int64, error)
	Update(id uint, req *dto.UpdateDesignRequest) (*model.DesignPhase, error)
	Delete(id uint) error
	Submit(id uint, req *dto.SubmitDesignRequest) (*model.DesignPhase, error)
	Review(id uint, reviewerID uint, req *dto.ReviewDesignRequest) (*model.DesignPhase, error)
}

type designService struct {
	repo   repository.DesignRepository
	logger *slog.Logger
}

// NewDesignService 构造设计服务。
func NewDesignService(repo repository.DesignRepository, logger *slog.Logger) DesignService {
	return &designService{repo: repo, logger: logger}
}

func (s *designService) Create(req *dto.CreateDesignRequest) (*model.DesignPhase, error) {
	if !constants.Contains(constants.PhaseStatuses, constants.PhaseStatusNotStarted) {
		return nil, apperrors.NewInternal("invalid phase status constant")
	}
	phase := &model.DesignPhase{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		DesignerID:  req.DesignerID,
		Status:      constants.PhaseStatusNotStarted,
		Version:     1,
		Description: req.Description,
		FileURLs:    marshalStrings(req.FileURLs),
	}
	if err := s.repo.Create(phase); err != nil {
		return nil, fmt.Errorf("create design phase: %w", err)
	}
	return phase, nil
}

func (s *designService) GetByID(id uint) (*model.DesignPhase, error) {
	phase, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewNotFound("design phase not found")
		}
		return nil, fmt.Errorf("get design phase: %w", err)
	}
	return phase, nil
}

func (s *designService) ListByProjectID(projectID uint) ([]model.DesignPhase, error) {
	phases, err := s.repo.ListByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("list design phases by project: %w", err)
	}
	return phases, nil
}

func (s *designService) List(page, pageSize int) ([]model.DesignPhase, int64, error) {
	phases, total, err := s.repo.ListAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list design phases: %w", err)
	}
	return phases, total, nil
}

func (s *designService) Update(id uint, req *dto.UpdateDesignRequest) (*model.DesignPhase, error) {
	phase, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if phase.Status == constants.PhaseStatusApproved {
		return nil, apperrors.NewConflict("approved design phase cannot be edited")
	}
	if req.Name != nil {
		phase.Name = *req.Name
	}
	if req.DesignerID != nil {
		phase.DesignerID = *req.DesignerID
	}
	if req.Description != nil {
		phase.Description = *req.Description
	}
	if req.FileURLs != nil {
		phase.FileURLs = marshalStrings(*req.FileURLs)
	}
	if err := s.repo.Update(phase); err != nil {
		return nil, fmt.Errorf("update design phase: %w", err)
	}
	return phase, nil
}

func (s *designService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete design phase: %w", err)
	}
	return nil
}

func (s *designService) Submit(id uint, req *dto.SubmitDesignRequest) (*model.DesignPhase, error) {
	phase, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if phase.Status == constants.PhaseStatusApproved {
		return nil, apperrors.NewConflict("approved design phase cannot be resubmitted")
	}
	if phase.Status == constants.PhaseStatusRevision {
		phase.Version++
	}
	if req.Description != "" {
		phase.Description = req.Description
	}
	if len(req.FileURLs) > 0 {
		phase.FileURLs = marshalStrings(req.FileURLs)
	}
	phase.Status = constants.PhaseStatusInProgress
	phase.ReviewComment = ""
	if err := s.repo.Update(phase); err != nil {
		return nil, fmt.Errorf("submit design phase: %w", err)
	}
	s.logger.Info("design phase submitted", "design_id", id, "version", phase.Version)
	return phase, nil
}

func (s *designService) Review(id uint, reviewerID uint, req *dto.ReviewDesignRequest) (*model.DesignPhase, error) {
	phase, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if phase.Status != constants.PhaseStatusInProgress {
		return nil, apperrors.NewConflict("only in-progress design phase can be reviewed")
	}
	phase.ReviewerID = reviewerID
	phase.ReviewComment = req.Comment
	if req.Approved {
		phase.Status = constants.PhaseStatusApproved
	} else {
		phase.Status = constants.PhaseStatusRevision
	}
	if err := s.repo.Update(phase); err != nil {
		return nil, fmt.Errorf("review design phase: %w", err)
	}
	s.logger.Info("design phase reviewed", "design_id", id, "approved", req.Approved, "reviewer_id", reviewerID)
	return phase, nil
}

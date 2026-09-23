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

// ProjectService 项目服务接口。
type ProjectService interface {
	Create(req *dto.CreateProjectRequest) (*model.RenovationProject, error)
	GetByID(id uint) (*model.RenovationProject, error)
	List(status, keyword string, page, pageSize int) ([]model.RenovationProject, int64, error)
	Update(id uint, req *dto.UpdateProjectRequest) (*model.RenovationProject, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) (*model.RenovationProject, error)
}

type projectService struct {
	repo   repository.ProjectRepository
	logger *slog.Logger
}

// NewProjectService 构造项目服务。
func NewProjectService(repo repository.ProjectRepository, logger *slog.Logger) ProjectService {
	return &projectService{repo: repo, logger: logger}
}

func (s *projectService) Create(req *dto.CreateProjectRequest) (*model.RenovationProject, error) {
	if !constants.Contains(constants.DecorStyles, req.DecorStyle) {
		return nil, apperrors.NewBadRequest("invalid decor_style")
	}
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid start_date, expect YYYY-MM-DD")
	}
	endDate, err := parseDate(req.ExpectedEndDate)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid expected_end_date, expect YYYY-MM-DD")
	}
	project := &model.RenovationProject{
		Name:            req.Name,
		HouseType:       req.HouseType,
		Area:            req.Area,
		DecorStyle:      req.DecorStyle,
		Address:         req.Address,
		OwnerID:         req.OwnerID,
		DesignerID:      req.DesignerID,
		ForemanID:       req.ForemanID,
		Status:          constants.ProjectStatusDesigning,
		ContractAmount:  req.ContractAmount,
		StartDate:       startDate,
		ExpectedEndDate: endDate,
	}
	if req.Status != "" {
		if !constants.Contains(constants.ProjectStatuses, req.Status) {
			return nil, apperrors.NewBadRequest("invalid project status")
		}
		project.Status = req.Status
	}
	if err := s.repo.Create(project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return project, nil
}

func (s *projectService) GetByID(id uint) (*model.RenovationProject, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperrors.NewNotFound("project not found")
		}
		return nil, fmt.Errorf("get project: %w", err)
	}
	return project, nil
}

func (s *projectService) List(status, keyword string, page, pageSize int) ([]model.RenovationProject, int64, error) {
	projects, total, err := s.repo.List(repository.ProjectFilter{Status: status, Keyword: keyword}, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	return projects, total, nil
}

func (s *projectService) Update(id uint, req *dto.UpdateProjectRequest) (*model.RenovationProject, error) {
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.HouseType != nil {
		project.HouseType = *req.HouseType
	}
	if req.Area != nil {
		project.Area = *req.Area
	}
	if req.DecorStyle != nil {
		if !constants.Contains(constants.DecorStyles, *req.DecorStyle) {
			return nil, apperrors.NewBadRequest("invalid decor_style")
		}
		project.DecorStyle = *req.DecorStyle
	}
	if req.Address != nil {
		project.Address = *req.Address
	}
	if req.OwnerID != nil {
		project.OwnerID = *req.OwnerID
	}
	if req.DesignerID != nil {
		project.DesignerID = *req.DesignerID
	}
	if req.ForemanID != nil {
		project.ForemanID = *req.ForemanID
	}
	if req.ContractAmount != nil {
		project.ContractAmount = *req.ContractAmount
	}
	if req.StartDate != nil {
		t, err := parseDate(req.StartDate)
		if err != nil {
			return nil, apperrors.NewBadRequest("invalid start_date, expect YYYY-MM-DD")
		}
		project.StartDate = t
	}
	if req.ExpectedEndDate != nil {
		t, err := parseDate(req.ExpectedEndDate)
		if err != nil {
			return nil, apperrors.NewBadRequest("invalid expected_end_date, expect YYYY-MM-DD")
		}
		project.ExpectedEndDate = t
	}
	if err := s.repo.Update(project); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	return project, nil
}

func (s *projectService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}

func (s *projectService) UpdateStatus(id uint, status string) (*model.RenovationProject, error) {
	if !constants.Contains(constants.ProjectStatuses, status) {
		return nil, apperrors.NewBadRequest("invalid project status")
	}
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project.Status == constants.ProjectStatusArchived {
		return nil, apperrors.NewConflict("archived project cannot change status")
	}
	project.Status = status
	if err := s.repo.Update(project); err != nil {
		return nil, fmt.Errorf("update project status: %w", err)
	}
	return project, nil
}

package service

import (
	"fmt"
	"log/slog"

	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// AuditService 操作日志服务接口。
type AuditService interface {
	Record(userID uint, username, action, resource string, resourceID uint, detail, ip string) error
	List(page, pageSize int) ([]model.AuditLog, int64, error)
}

type auditService struct {
	repo   repository.AuditLogRepository
	logger *slog.Logger
}

// NewAuditService 构造操作日志服务。
func NewAuditService(repo repository.AuditLogRepository, logger *slog.Logger) AuditService {
	return &auditService{repo: repo, logger: logger}
}

func (s *auditService) Record(userID uint, username, action, resource string, resourceID uint, detail, ip string) error {
	log := &model.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Detail:     detail,
		IP:         ip,
	}
	if err := s.repo.Create(log); err != nil {
		s.logger.Error("record audit log failed", "action", action, "resource", resource, "error", err)
		return fmt.Errorf("record audit log: %w", err)
	}
	return nil
}

func (s *auditService) List(page, pageSize int) ([]model.AuditLog, int64, error) {
	logs, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, total, nil
}

package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/Talan-Application/system-handbook-service/internal/domain"
	"github.com/Talan-Application/system-handbook-service/internal/repository"
)

type SubjectService struct {
	repo   repository.SubjectRepository
	logger *zap.Logger
}

func NewSubjectService(repo repository.SubjectRepository, logger *zap.Logger) *SubjectService {
	return &SubjectService{repo: repo, logger: logger}
}

func (s *SubjectService) Create(ctx context.Context, subject *domain.Subject) (*domain.Subject, error) {
	created, err := s.repo.Create(ctx, subject)
	if err != nil {
		s.logger.Error("failed to create subject", zap.Error(err))
		return nil, err
	}
	return created, nil
}

func (s *SubjectService) GetByID(ctx context.Context, id int64) (*domain.Subject, error) {
	subject, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subject", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return subject, nil
}

func (s *SubjectService) GetAll(ctx context.Context, limit *int, offset *int) ([]domain.Subject, error) {
	subjects, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to get subjects", zap.Error(err))
		return nil, err
	}
	return subjects, nil
}

func (s *SubjectService) Update(ctx context.Context, id int64, subject *domain.Subject) (*domain.Subject, error) {
	updated, err := s.repo.Update(ctx, id, subject)
	if err != nil {
		s.logger.Error("failed to update subject", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return updated, nil
}

func (s *SubjectService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete subject", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}

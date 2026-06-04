package service

import (
	"context"

	"github.com/Talan-Application/system-handbook-service/internal/domain"
)

type ICommonSubjectService interface {
	Create(ctx context.Context, subject *domain.CommonSubject) (*domain.CommonSubject, error)
	GetByID(ctx context.Context, id int64) (*domain.CommonSubject, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]domain.CommonSubject, error)
	GetLookup(ctx context.Context) ([]domain.CommonSubject, error)
	Update(ctx context.Context, id int64, subject *domain.CommonSubject) (*domain.CommonSubject, error)
	Delete(ctx context.Context, id int64) error
}

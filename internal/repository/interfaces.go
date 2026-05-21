package repository

import (
	"context"

	"github.com/Talan-Application/system-handbook-service/internal/domain"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *domain.Subject) (*domain.Subject, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, subject *domain.Subject) (*domain.Subject, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]domain.Subject, error)
	GetByID(ctx context.Context, id int64) (*domain.Subject, error)
}

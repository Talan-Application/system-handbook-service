package service

import (
	"context"

	commonsubjectv1 "github.com/Talan-Application/proto-generation/common_subject/v1"
	"github.com/Talan-Application/system-handbook-service/internal/domain"
)

type ICommonSubjectService interface {
	Create(ctx context.Context, req *commonsubjectv1.CreateCommonSubjectRequest) (*domain.CommonSubject, error)
	GetByID(ctx context.Context, id int64) (*domain.CommonSubject, error)
	GetAll(ctx context.Context, req *commonsubjectv1.GetAllCommonSubjectsRequest) ([]domain.CommonSubject, error)
	GetLookup(ctx context.Context) ([]domain.CommonSubject, error)
	Update(ctx context.Context, id int64, req *commonsubjectv1.UpdateCommonSubjectRequest) (*domain.CommonSubject, error)
	Delete(ctx context.Context, id int64) error
}

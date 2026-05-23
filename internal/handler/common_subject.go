package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonsubjectv1 "github.com/Talan-Application/proto-generation/common_subject/v1"
	"github.com/Talan-Application/system-handbook-service/internal/domain"
	"github.com/Talan-Application/system-handbook-service/internal/service"
)

type CommonSubjectHandler struct {
	commonsubjectv1.UnimplementedCommonSubjectServiceServer
	svc service.ICommonSubjectService
	log *zap.Logger
}

func NewCommonSubjectHandler(svc service.ICommonSubjectService, log *zap.Logger) *CommonSubjectHandler {
	return &CommonSubjectHandler{svc: svc, log: log}
}

func (h *CommonSubjectHandler) CreateCommonSubject(ctx context.Context, req *commonsubjectv1.CreateCommonSubjectRequest) (*commonsubjectv1.CommonSubjectResponse, error) {
	subject := &domain.CommonSubject{
		Translations: req.GetTranslations(),
	}

	created, err := h.svc.Create(ctx, subject)
	if err != nil {
		h.log.Error("CreateCommonSubject", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return toProto(created), nil
}

func (h *CommonSubjectHandler) GetCommonSubject(ctx context.Context, req *commonsubjectv1.GetCommonSubjectRequest) (*commonsubjectv1.CommonSubjectResponse, error) {
	subject, err := h.svc.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrSubjectNotFound) {
			return nil, status.Error(codes.NotFound, "subject not found")
		}
		h.log.Error("GetCommonSubject", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return toProto(subject), nil
}

func (h *CommonSubjectHandler) GetAllCommonSubjects(ctx context.Context, req *commonsubjectv1.GetAllCommonSubjectsRequest) (*commonsubjectv1.GetAllCommonSubjectsResponse, error) {
	var limit, offset *int
	if req.Limit != nil {
		n := int(req.GetLimit())
		limit = &n
	}
	if req.Offset != nil {
		n := int(req.GetOffset())
		offset = &n
	}

	subjects, err := h.svc.GetAll(ctx, limit, offset)
	if err != nil {
		h.log.Error("GetAllCommonSubjects", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := make([]*commonsubjectv1.CommonSubjectResponse, len(subjects))
	for i := range subjects {
		resp[i] = toProto(&subjects[i])
	}

	return &commonsubjectv1.GetAllCommonSubjectsResponse{CommonSubjects: resp}, nil
}

func (h *CommonSubjectHandler) UpdateCommonSubject(ctx context.Context, req *commonsubjectv1.UpdateCommonSubjectRequest) (*commonsubjectv1.CommonSubjectResponse, error) {
	subject := &domain.CommonSubject{
		Translations: req.GetTranslations(),
	}

	updated, err := h.svc.Update(ctx, req.GetId(), subject)
	if err != nil {
		if errors.Is(err, domain.ErrSubjectNotFound) {
			return nil, status.Error(codes.NotFound, "subject not found")
		}
		h.log.Error("UpdateCommonSubject", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return toProto(updated), nil
}

func (h *CommonSubjectHandler) DeleteCommonSubject(ctx context.Context, req *commonsubjectv1.DeleteCommonSubjectRequest) (*commonsubjectv1.DeleteCommonSubjectResponse, error) {
	if err := h.svc.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, domain.ErrSubjectNotFound) {
			return nil, status.Error(codes.NotFound, "subject not found")
		}
		h.log.Error("DeleteCommonSubject", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &commonsubjectv1.DeleteCommonSubjectResponse{Message: "deleted successfully"}, nil
}

func toProto(s *domain.CommonSubject) *commonsubjectv1.CommonSubjectResponse {
	return &commonsubjectv1.CommonSubjectResponse{
		Id:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
	}
}

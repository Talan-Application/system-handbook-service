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
	created, err := h.svc.Create(ctx, req)
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
	subjects, err := h.svc.GetAll(ctx, req)
	if err != nil {
		h.log.Error("GetAllCommonSubjects", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := make([]*commonsubjectv1.CommonSubjectResponse, len(subjects))
	for i := range subjects {
		resp[i] = toProtoWithTranslations(&subjects[i])
	}

	return &commonsubjectv1.GetAllCommonSubjectsResponse{CommonSubjects: resp}, nil
}

func (h *CommonSubjectHandler) GetCommonSubjectsLookup(ctx context.Context, req *commonsubjectv1.GetCommonSubjectsLookupRequest) (*commonsubjectv1.GetCommonSubjectsLookupResponse, error) {
	subjects, err := h.svc.GetLookup(ctx)
	if err != nil {
		h.log.Error("GetCommonSubjectsLookup", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	items := make([]*commonsubjectv1.CommonSubjectLookupItem, len(subjects))
	for i := range subjects {
		items[i] = &commonsubjectv1.CommonSubjectLookupItem{
			Id:   subjects[i].ID,
			Name: subjects[i].Name,
		}
	}

	return &commonsubjectv1.GetCommonSubjectsLookupResponse{CommonSubjects: items}, nil
}

func (h *CommonSubjectHandler) UpdateCommonSubject(ctx context.Context, req *commonsubjectv1.UpdateCommonSubjectRequest) (*commonsubjectv1.CommonSubjectResponse, error) {
	updated, err := h.svc.Update(ctx, req.GetId(), req)
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

func toProtoWithTranslations(s *domain.CommonSubject) *commonsubjectv1.CommonSubjectResponse {
	return &commonsubjectv1.CommonSubjectResponse{
		Id:           s.ID,
		CreatedAt:    s.CreatedAt.Unix(),
		UpdatedAt:    s.UpdatedAt.Unix(),
		Translations: s.Translations,
	}
}

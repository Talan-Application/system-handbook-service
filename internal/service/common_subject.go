package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/Talan-Application/system-handbook-service/internal/domain"
	"github.com/Talan-Application/system-handbook-service/internal/repository"
	"github.com/Talan-Application/system-handbook-service/internal/transport/grpc/ctxkeys"
	translationsvc "github.com/Talan-Application/translation-library/service"
)

type CommonSubjectService struct {
	repo           repository.CommonSubjectRepository
	translationSvc translationsvc.TranslationService
	logger         *zap.Logger
}

func NewCommonSubjectService(translationSvc translationsvc.TranslationService, repo repository.CommonSubjectRepository, logger *zap.Logger) *CommonSubjectService {
	return &CommonSubjectService{repo: repo, translationSvc: translationSvc, logger: logger}
}

func (s *CommonSubjectService) Create(ctx context.Context, subject *domain.CommonSubject) (*domain.CommonSubject, error) {
	key, err := s.translationSvc.GenerateKey(ctx)
	if err != nil {
		s.logger.Error("failed to generate translation key", zap.Error(err))
		return nil, err
	}

	if err := s.translationSvc.CreateTranslations(ctx, key, subject.Translations); err != nil {
		s.logger.Error("failed to create translations", zap.Error(err))
		return nil, err
	}

	subject.NameKey = key
	created, err := s.repo.Create(ctx, subject)
	if err != nil {
		s.logger.Error("failed to create common subject", zap.Error(err))
		return nil, err
	}
	return created, nil
}

func (s *CommonSubjectService) GetByID(ctx context.Context, id int64) (*domain.CommonSubject, error) {
	subject, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get common subject", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	locale := localeFromCtx(ctx)
	translation, err := s.translationSvc.GetByLanguageCodeAndKey(ctx, locale, subject.NameKey)
	if err != nil {
		s.logger.Error("failed to get translation", zap.String("key", subject.NameKey), zap.String("locale", locale), zap.Error(err))
		return nil, err
	}
	subject.Name = translation.TranslationValue

	return subject, nil
}

func (s *CommonSubjectService) GetAll(ctx context.Context, limit *int, offset *int) ([]domain.CommonSubject, error) {
	subjects, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to get common subjects", zap.Error(err))
		return nil, err
	}

	if len(subjects) == 0 {
		return subjects, nil
	}

	keys := make([]string, len(subjects))
	for i, s := range subjects {
		keys[i] = s.NameKey
	}

	resolved, err := s.translationSvc.ResolveBulkAllLocales(ctx, keys)
	if err != nil {
		s.logger.Error("failed to resolve translations", zap.Error(err))
		return nil, err
	}

	for i := range subjects {
		subjects[i].Translations = resolved[subjects[i].NameKey]
	}

	return subjects, nil
}

func (s *CommonSubjectService) GetLookup(ctx context.Context) ([]domain.CommonSubject, error) {
	subjects, err := s.repo.GetAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("failed to get common subjects for lookup", zap.Error(err))
		return nil, err
	}

	if len(subjects) == 0 {
		return subjects, nil
	}

	locale := localeFromCtx(ctx)
	keys := make([]string, len(subjects))
	for i, subj := range subjects {
		keys[i] = subj.NameKey
	}

	resolved, err := s.translationSvc.ResolveBulkAllLocales(ctx, keys)
	if err != nil {
		s.logger.Error("failed to resolve translations for lookup", zap.Error(err))
		return nil, err
	}

	for i := range subjects {
		localeMap := resolved[subjects[i].NameKey]
		if name, ok := localeMap[locale]; ok && name != "" {
			subjects[i].Name = name
		} else if name, ok := localeMap["ru"]; ok {
			subjects[i].Name = name
		}
	}

	return subjects, nil
}

func (s *CommonSubjectService) Update(ctx context.Context, id int64, subject *domain.CommonSubject) (*domain.CommonSubject, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get common subject for update", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	for locale, value := range subject.Translations {
		if err := s.translationSvc.UpdateByKey(ctx, existing.NameKey, locale, value); err != nil {
			s.logger.Error("failed to update translation", zap.String("key", existing.NameKey), zap.String("locale", locale), zap.Error(err))
			return nil, err
		}
	}

	updated, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		s.logger.Error("failed to update common subject", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return updated, nil
}

func (s *CommonSubjectService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get common subject for delete", zap.Int64("id", id), zap.Error(err))
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete common subject", zap.Int64("id", id), zap.Error(err))
		return err
	}

	if err := s.translationSvc.DeleteByKey(ctx, existing.NameKey); err != nil {
		s.logger.Error("failed to delete translations for common subject", zap.String("key", existing.NameKey), zap.Error(err))
		return err
	}

	return nil
}

func localeFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxkeys.LocaleKey).(string); ok && v != "" {
		return v
	}
	return "ru"
}

package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Talan-Application/system-handbook-service/internal/domain"
)

type CommonSubjectRepository struct {
	db *pgxpool.Pool
}

func NewCommonSubjectRepository(db *pgxpool.Pool) *CommonSubjectRepository {
	return &CommonSubjectRepository{db}
}

func (r *CommonSubjectRepository) Create(ctx context.Context, subject *domain.CommonSubject) (*domain.CommonSubject, error) {
	query := `INSERT INTO common_subjects (name_key, created_at, updated_at)
				VALUES ($1, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, subject.NameKey).
		Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create common subject: %w", err)
	}

	return subject, nil
}

func (r *CommonSubjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM common_subjects WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete common subject: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("common subject with id=%d: %w", id, domain.ErrSubjectNotFound)
	}

	return nil
}

func (r *CommonSubjectRepository) Update(ctx context.Context, id int64, subject *domain.CommonSubject) (*domain.CommonSubject, error) {
	query := `UPDATE common_subjects SET name_key = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, subject.NameKey, id).Scan(&subject.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update common subject: %w", domain.ErrSubjectNotFound)
		}
		return nil, fmt.Errorf("update common subject: %w", err)
	}

	return subject, nil
}

func (r *CommonSubjectRepository) GetAll(ctx context.Context, limit *int, offset *int) ([]domain.CommonSubject, error) {
	query := `SELECT id, name_key, is_deleted, deleted_at, created_at, updated_at
				FROM common_subjects ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select common subjects: %w", err)
	}
	defer rows.Close()

	var subjects []domain.CommonSubject
	for rows.Next() {
		var s domain.CommonSubject
		var deletedAt pgtype.Timestamptz

		err := rows.Scan(&s.ID, &s.NameKey, &s.IsDeleted, &deletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan common subject: %w", err)
		}

		if deletedAt.Valid {
			s.DeletedAt = deletedAt.Time
		} else {
			s.DeletedAt = time.Time{}
		}

		subjects = append(subjects, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return subjects, nil
}

func (r *CommonSubjectRepository) GetByID(ctx context.Context, id int64) (*domain.CommonSubject, error) {
	query := `SELECT id, name_key, is_deleted, deleted_at, created_at, updated_at
				FROM common_subjects WHERE id = $1`

	s := &domain.CommonSubject{}
	var deletedAt pgtype.Timestamptz

	err := r.db.QueryRow(ctx, query, id).
		Scan(&s.ID, &s.NameKey, &s.IsDeleted, &deletedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get common subject: %w", domain.ErrSubjectNotFound)
		}
		return nil, fmt.Errorf("get common subject: %w", err)
	}

	if deletedAt.Valid {
		s.DeletedAt = deletedAt.Time
	} else {
		s.DeletedAt = time.Time{}
	}

	return s, nil
}

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

type SubjectRepository struct {
	db *pgxpool.Pool
}

func NewSubjectRepository(db *pgxpool.Pool) *SubjectRepository {
	return &SubjectRepository{db}
}

func (r *SubjectRepository) Create(ctx context.Context, subject *domain.Subject) (*domain.Subject, error) {
	query := `INSERT INTO subjects (name_key, created_at, updated_at)
				VALUES ($1, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, subject.NameKey).
		Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create subject: %w", err)
	}

	return subject, nil
}

func (r *SubjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subjects WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subject: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("subject with id=%d: %w", id, domain.ErrSubjectNotFound)
	}

	return nil
}

func (r *SubjectRepository) Update(ctx context.Context, id int64, subject *domain.Subject) (*domain.Subject, error) {
	query := `UPDATE subjects SET name_key = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, subject.NameKey, id).Scan(&subject.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update subject: %w", domain.ErrSubjectNotFound)
		}
		return nil, fmt.Errorf("update subject: %w", err)
	}

	return subject, nil
}

func (r *SubjectRepository) GetAll(ctx context.Context, limit *int, offset *int) ([]domain.Subject, error) {
	query := `SELECT id, name_key, is_deleted, deleted_at, created_at, updated_at
				FROM subjects ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select subjects: %w", err)
	}
	defer rows.Close()

	var subjects []domain.Subject
	for rows.Next() {
		var s domain.Subject
		var deletedAt pgtype.Timestamptz

		err := rows.Scan(&s.ID, &s.NameKey, &s.IsDeleted, &deletedAt, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
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

func (r *SubjectRepository) GetByID(ctx context.Context, id int64) (*domain.Subject, error) {
	query := `SELECT id, name_key, is_deleted, deleted_at, created_at, updated_at
				FROM subjects WHERE id = $1`

	s := &domain.Subject{}
	var deletedAt pgtype.Timestamptz

	err := r.db.QueryRow(ctx, query, id).
		Scan(&s.ID, &s.NameKey, &s.IsDeleted, &deletedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get subject: %w", domain.ErrSubjectNotFound)
		}
		return nil, fmt.Errorf("get subject: %w", err)
	}

	if deletedAt.Valid {
		s.DeletedAt = deletedAt.Time
	} else {
		s.DeletedAt = time.Time{}
	}

	return s, nil
}

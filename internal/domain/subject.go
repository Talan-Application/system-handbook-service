package domain

import "time"

type Subject struct {
	ID        int64     `db:"id"`
	NameKey   string    `db:"name_key"`
	IsDeleted bool      `db:"is_deleted"`
	DeletedAt time.Time `db:"deleted_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

package domain

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrInternal       = errors.New("internal error")
	ErrSubjectNotFound = errors.New("subject not found")
)

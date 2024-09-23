package db

import (
	"errors"
	"github.com/lib/pq"
)

var ErrUniqueViolation = &pq.Error{
	Code: UniqueViolation,
}

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	EmailConstraint     = "users_email_key"
	UsernameConstraint  = "users_name_key"
	ReadTitlesCode      = "titles:r"
)

func PqErrorCode(err error) string {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code)
	}
	return ""
}

func PqErrorConstraint(err error) string {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Constraint
	}
	return ""
}

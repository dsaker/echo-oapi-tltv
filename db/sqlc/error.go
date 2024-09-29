package db

import (
	"errors"
	"github.com/lib/pq"
)

var ErrUniqueViolation = &pq.Error{
	Code: UniqueViolation,
}

var ErrForeignKeyViolation = &pq.Error{
	Code:    ForeignKeyViolation,
	Message: "insert or update on table \"users_permissions\" violates foreign key constraint \"users_permissions_user_id_fkey\"",
}

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	EmailConstraint     = "users_email_key"
	UsernameConstraint  = "users_name_key"
	ReadTitlesCode      = "titles:r"
	WriteTitlesCode     = "titles:w"
	GlobalAdminCode     = "global:admin"
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

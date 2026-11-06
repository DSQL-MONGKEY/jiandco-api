package database

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503" || pgErr.Code == "23001"
	}

	message := err.Error()
	return strings.Contains(message, "SQLSTATE 23503") || strings.Contains(message, "SQLSTATE 23001")
}

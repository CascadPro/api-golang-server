package core_postgres_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ViolatesForeignKeyErrorCode = "23503"
)

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)

func MapErrors(err error) error {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRows
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == ViolatesForeignKeyErrorCode {
				return fmt.Errorf("%v: %w", pgErr.Message, ErrViolatesForeignKey)
			}
		}
	}

	return fmt.Errorf("%v: %w", err, ErrUnknown)
}

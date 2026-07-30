package core_postgres_pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ViolatesForeignKeyErrorCode       = "23503"
	ViolatesUniqueConstraintErrorCode = "23505"
	ViolatesCheckConstraintErrorCode  = "23514"
)

var (
	ErrNoRows                   = errors.New("no rows")
	ErrViolatesForeignKey       = errors.New("violates foreign key")
	ErrViolatesUniqueConstraint = errors.New("violates unique constraint")
	ErrViolatesCheckConstraint  = errors.New("violates check constraint")
	ErrUnknown                  = errors.New("unknown")
)

func MapErrors(err error) error {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRows
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case ViolatesForeignKeyErrorCode:
				return fmt.Errorf("%v: %w", pgErr.Message, ErrViolatesForeignKey)
			case ViolatesUniqueConstraintErrorCode:
				return fmt.Errorf("%v: %w", pgErr.Message, ErrViolatesUniqueConstraint)
			case ViolatesCheckConstraintErrorCode:
				return fmt.Errorf("%v: %w", pgErr.Message, ErrViolatesCheckConstraint)
			}
		}
	}

	return fmt.Errorf("%v: %w", err, ErrUnknown)
}

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, MapErrors(err)
	}

	return rows, nil
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)

	if err != nil {
		return tag, MapErrors(err)
	}

	return tag, nil
}

func (p *ConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.Pool.Begin(ctx)

	if err != nil {
		return nil, MapErrors(err)
	}

	return tx, nil
}

func (p *ConnectionPool) Close() {
	p.Pool.Close()
}

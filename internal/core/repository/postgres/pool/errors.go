package core_postgres_pool

import (
	"context"
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

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	rows, err := p.Pool.Query(ctx, sql, args...)
	return rows, MapErrors(err)
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	tag, err := p.Pool.Exec(ctx, sql, arguments...)
	return tag, MapErrors(err)
}

func (p *ConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	tx, err := p.Pool.Begin(ctx)
	return tx, MapErrors(err)
}

func (p *ConnectionPool) Close() {
	p.Pool.Close()
}

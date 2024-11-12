package repo

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/internal/service/transactional"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DefaultRepo struct {
	pool *pgxpool.Pool
}

func NewDefaultRepo(pool *pgxpool.Pool) *DefaultRepo {
	return &DefaultRepo{pool: pool}
}

func (r *DefaultRepo) QueryRow(ctx context.Context, destination interface{}, sql string, args ...any) error {
	tx := r.GetTX(ctx)
	if tx == nil {
		return pgxscan.Get(ctx, r.pool, destination, sql, args...)
	}
	return pgxscan.Get(ctx, tx, destination, sql, args...)
}

func (r *DefaultRepo) Query(ctx context.Context, destination interface{}, sql string, args ...any) error {
	tx := r.GetTX(ctx)
	if tx == nil {
		return pgxscan.Select(ctx, r.pool, destination, sql, args...)
	}
	return pgxscan.Select(ctx, tx, destination, sql, args...)
}

func (r *DefaultRepo) Exec(ctx context.Context, sql string, args ...any) error {
	tx := r.GetTX(ctx)
	if tx == nil {
		_, err := r.pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("r.pool.Exec: %w", err)
		}
		return nil
	}
	_, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}
	return nil
}

func (r *DefaultRepo) GetTX(ctx context.Context) pgx.Tx {
	txVal := ctx.Value(transactional.TxKey)
	tx, ok := txVal.(pgx.Tx)
	if ok {
		return tx
	}
	return nil
}

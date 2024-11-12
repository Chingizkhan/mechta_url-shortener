package transactional

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/url-shortener/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TxKey = "tx_key"
)

type Transactional struct {
	pool *pgxpool.Pool
}

func New(pg *postgres.Postgres) *Transactional {
	return &Transactional{pg.Pool}
}

func (t *Transactional) Exec(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("BeginTx: %w", err)
	}
	txCtx := context.WithValue(ctx, TxKey, tx)
	if err = fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("CommitTx: %w", err)
	}

	return nil
}

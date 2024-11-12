package url_repo

import (
	"context"
	"github.com/Chingizkhan/url-shortener/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

const (
	dataSource = "postgresql://app:secret@localhost:5490/app?sslmode=disable"
)

var (
	rep  *Repository
	pool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	pg, err := postgres.New(
		dataSource,
		postgres.MaxPoolSize(20),
	)
	if err != nil {
		log.Fatal("can not connect postgres:", err)
	}
	rep = New(pg)
	pool = pg.Pool

	os.Exit(m.Run())
}

func prependCtx(t *testing.T) (context.Context, pgx.Tx, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	require.NoError(t, err)

	ctx = context.WithValue(ctx, "tx_key", tx)

	return ctx, tx, cancel
}

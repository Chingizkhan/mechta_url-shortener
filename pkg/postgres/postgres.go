package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const (
	_defaultMaxPoolSize     = 1
	_defaultMaxConnLifetime = time.Hour * 20
	_defaultMaxConnIdleTime = time.Minute * 30
	_defaultConnAttempts    = 10
	_defaultConnTimeout     = time.Second
)

type Postgres struct {
	maxPoolSize     int
	maxConnLifetime time.Duration
	maxConnIdleTime time.Duration
	connAttempts    int
	connTimeout     time.Duration
	Pool            *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:     _defaultMaxPoolSize,
		maxConnLifetime: _defaultMaxConnLifetime,
		maxConnIdleTime: _defaultMaxConnIdleTime,
		connAttempts:    _defaultConnAttempts,
		connTimeout:     _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)
	poolConfig.MaxConnLifetime = time.Hour * 20
	poolConfig.MaxConnIdleTime = time.Minute * 30

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			break
		}
		if pg.Pool != nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

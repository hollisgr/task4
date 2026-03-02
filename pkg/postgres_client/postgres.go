package postgres_client

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, maxAttempts int, dsn string) (pool *pgxpool.Pool, err error) {

	poolCfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return pool, fmt.Errorf("new pool: parse config err: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// connect to pool
		pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
		if err != nil {
			return err
		}

		// trying to ping pool
		err = pool.Ping(ctx)
		if err != nil {
			pool.Close()
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}

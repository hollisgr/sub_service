package postgres

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
}

func NewPool(ctx context.Context, maxAttempts int, cfg config.Config) (pool *pgxpool.Pool, err error) {
	l := logger.GetLogger()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Postgresql.Username, cfg.Postgresql.Password, cfg.Postgresql.Host,
		cfg.Postgresql.Port, cfg.Postgresql.Database)
	l.Infoln(fmt.Sprintf("connecting to database: postgresql://%s:{password}@%s:%s/%s",
		cfg.Postgresql.Username, cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.Database))
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--
			continue
		}
		return nil
	}
	return
}

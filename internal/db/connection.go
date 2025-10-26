package db

import (
	"context"
	"fmt"
	"time"

	"dnk.com/hoc-golang/internal/config"
	"dnk.com/hoc-golang/internal/db/sqlc"
	"dnk.com/hoc-golang/internal/utils"
	"dnk.com/hoc-golang/pkg/logger"
	"dnk.com/hoc-golang/pkg/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var DB sqlc.Querier
var DBPool *pgxpool.Pool

func InitDB() error {
	connStr := config.NewConfig().DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %v", err)

	}

	sqlLogger := utils.NewLoggerWithPath("sql.log", "info")

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxZerologTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Millisecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}
	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err = pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %v", err)

	}
	DB = sqlc.New(DBPool)
	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("db ping error: %v", err)

	}
	logger.Log.Info().Msg("Connected Database Postgresql")

	return nil
}

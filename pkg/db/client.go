package db

import (
	"context"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Q *Queries
	C *pgxpool.Pool
)

func InitClient(databaseUrl string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	// to fix supabase prep statement issue
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	// to fix supabase conn isuue
	// config.MaxConns = 10

	// config.Tracer = &LoggingQueryTracer{logger: slog.Default()}
	// config.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	C, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = C.Ping(ctx)
	if err != nil {
		return nil, err
	}

	go listenInterrupts(C)
	// init global db client for later usage across the app
	// if err != nil app won't start, so its safe to use global here
	Q = New(C)
	return C, err
}

func listenInterrupts(conn *pgxpool.Pool) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	conn.Close()
	os.Exit(0)
}

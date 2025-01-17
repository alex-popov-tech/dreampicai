package db

import (
	"context"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Client *Queries

func InitClient(databaseUrl string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	// config.Tracer = &LoggingQueryTracer{logger: slog.Default()}
	// config.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	go listenInterrupts(conn)
	// init global db client for later usage across the app
	// if err != nil app won't start, so its safe to use global here
	Client = New(conn)
	return conn, err
}

func listenInterrupts(conn *pgxpool.Pool) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	conn.Close()
	os.Exit(0)
}

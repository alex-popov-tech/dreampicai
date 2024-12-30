package db

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
)

var Client *Queries

func InitClient(databaseUrl string) (*pgx.Conn, error) {
	ctx := context.Background()
	config, err := pgx.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	config.Tracer = &LoggingQueryTracer{logger: slog.Default()}
	config.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	conn, err := pgx.ConnectConfig(ctx, config)
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

func listenInterrupts(conn *pgx.Conn) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	conn.Close(context.Background())
	os.Exit(0)
}

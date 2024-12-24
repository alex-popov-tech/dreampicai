package db

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

type LoggingQueryTracer struct {
	logger *slog.Logger
}

func (l *LoggingQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.logger.
		Info("query start",
			slog.String("sql", prettyPrintSQL(data.SQL)),
			slog.Any("args", data.Args),
		)
	return ctx
}

func (l *LoggingQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// Failure
	if data.Err != nil {
		l.logger.
			Error("query end",
				slog.String("error", data.Err.Error()),
				slog.String("command_tag", data.CommandTag.String()),
			)
		return
	}

	// Success
	l.logger.
		Info("query end",
			slog.String("command_tag", data.CommandTag.String()),
		)
}

var (
	replaceTabs                      = regexp.MustCompile(`\t+`)
	replaceSpacesBeforeOpeningParens = regexp.MustCompile(`\s+\(`)
	replaceSpacesAfterOpeningParens  = regexp.MustCompile(`\(\s+`)
	replaceSpacesBeforeClosingParens = regexp.MustCompile(`\s+\)`)
	replaceSpacesAfterClosingParens  = regexp.MustCompile(`\)\s+`)
	replaceSpaces                    = regexp.MustCompile(`\s+`)
)

func prettyPrintSQL(sql string) string {
	lines := strings.Split(sql, "\n")

	pretty := strings.Join(lines, " ")
	pretty = replaceTabs.ReplaceAllString(pretty, "")
	pretty = replaceSpacesBeforeOpeningParens.ReplaceAllString(pretty, "(")
	pretty = replaceSpacesAfterOpeningParens.ReplaceAllString(pretty, "(")
	pretty = replaceSpacesAfterClosingParens.ReplaceAllString(pretty, ")")
	pretty = replaceSpacesBeforeClosingParens.ReplaceAllString(pretty, ")")

	// Finally, replace multiple spaces with a single space
	pretty = replaceSpaces.ReplaceAllString(pretty, " ")

	return strings.TrimSpace(pretty)
}

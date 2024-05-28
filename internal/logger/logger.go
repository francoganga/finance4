package logger

import (
	"context"
	"log/slog"
	"regexp"
	"slices"

	sqldblogger "github.com/simukti/sqldb-logger"
)

var re = regexp.MustCompile(`--.*\n`)

func removeSQLCComments(s string) string {

	return re.ReplaceAllString(s, "")

}

type SQLLogger struct {
	prevStmtId string
	logger     *slog.Logger
}

func NewSQLLogger(logger *slog.Logger) *SQLLogger {
	return &SQLLogger{
		logger: logger,
	}
}

func (ml *SQLLogger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {

	msgs := []string{"StmtExecContext", "QueryContext", "StmtQueryContext", "ExecContext"}

	if stmtId, ok := data["stmt_id"].(string); ok {

		if ml.prevStmtId == stmtId {
			return
		}

	}

	if !slices.Contains(msgs, msg) {
		return
	}

	if q, ok := data["query"]; ok {

		qe := removeSQLCComments(q.(string))

		ml.logger.Info("SQL", slog.String("query", string(qe)), slog.Any("args", data["args"]))
	}

	if stmtId, ok := data["stmt_id"].(string); ok {
		ml.prevStmtId = stmtId
	}
}


package orm

import (
    "context"
    "github.com/gohade/hade/framework/contract"
    "gorm.io/gorm/logger"
    "time"
)

type OrmLogger struct {
    logger contract.Log
}

func NewOrmLogger(logger contract.Log) *OrmLogger {
    return &OrmLogger{logger: logger}
}

func (o *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
    return o
}

func (o *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Info(ctx, s, fields)
}

func (o *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Warn(ctx, s, fields)
}

func (o *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
    fields := map[string]interface{}{
        "fields": i,
    }
    o.logger.Error(ctx, s, fields)
}

func (o *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    sql, rows := fc()
    elapsed := time.Since(begin)
    fields := map[string]interface{}{
        "begin": begin,
        "error": err,
        "sql":   sql,
        "rows":  rows,
        "time":  elapsed,
    }

    s := "orm trace sql"
    o.logger.Trace(ctx, s, fields)
}

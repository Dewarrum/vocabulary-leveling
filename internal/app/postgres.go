package app

import (
	"database/sql"
	"errors"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var (
	ErrPostgresUrlIsRequired     = errors.New("POSTGRES_URL is required")
	ErrFailedToConnectToPostgres = errors.New("failed to connect to Postgres")
)

func createPostgresConnection(logger zerolog.Logger) (*sqlx.DB, error) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return nil, ErrPostgresUrlIsRequired
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToConnectToPostgres)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	loggerOptions := []sqldblogger.Option{
		sqldblogger.WithSQLQueryFieldname("sql"),
		sqldblogger.WithWrapResult(false),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithLogArguments(false),
	}
	db = sqldblogger.OpenDriver(dsn, db.Driver(), zerologadapter.New(logger), loggerOptions...)

	return sqlx.NewDb(db, "postgres"), nil
}

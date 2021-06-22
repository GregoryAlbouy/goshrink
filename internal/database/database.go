package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	dbSystem  = "mysql"
	dsnFormat = "%s:%s@(%s:%s)/"
)

type DB struct {
	name string
	sqlx *sqlx.DB
}

type Config struct {
	User     string
	Password string
	Domain   string
	Port     string
	Database string
}

// mustConnect makes the connection to the mysql server.
// It panics in case of error.
func (db *DB) MustConnect(cfg Config) {
	dsn := fmt.Sprintf(
		dsnFormat,
		cfg.User, cfg.Password, cfg.Domain, cfg.Port,
	)
	sqlx := sqlx.MustConnect(dbSystem, dsn)
	sqlx.MustExec("USE " + cfg.Database)
	db.sqlx = sqlx
}

// Close closes the connection to the database.
// It is typically used with `defer` after initialization.
func (db *DB) Close() error {
	return db.sqlx.Close()
}

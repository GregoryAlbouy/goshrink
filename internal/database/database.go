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

// MustInit connects to the database and creates the base database and tables.
// It panics in case of error, which is acceptable here as it is called
// before the server starts.
func (db *DB) MustInit(cfg Config) {
	db.MustConnect(cfg)
	db.mustCreateDatabase(cfg.Database)
	db.mustCreateTables()
}

// Close closes the connection to the database.
// It is typically used with `defer` after initialization.
func (db *DB) Close() error {
	return db.sqlx.Close()
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

// mustCreateDatabase creates the main database if not already created.
// It panics in case of error.
func (db *DB) mustCreateDatabase(name string) {
	db.sqlx.MustExec("CREATE DATABASE IF NOT EXISTS " + name)
	db.sqlx.MustExec("USE " + name)
}

// mustCreateTables creates the tables for users and avatars in the database,
// and a view joining them for easy access.
// It panics in case of error.
func (db *DB) mustCreateTables() {
	db.sqlx.MustExec(userSchema)       // user table
	db.sqlx.MustExec(avatarSchema)     // avatar table
	db.sqlx.MustExec(userAvatarSchema) // user_avatar view
}

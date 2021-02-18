package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

// New returns a new database connection pool.
func New() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", "database.sqlite")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open database")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Unable to ping database")
	}
	return db, nil
}

// StatusCheck pings the database to see if it's reachable..
func StatusCheck(db *sqlx.DB) error {
	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "Unable to ping database")
	}
	return nil
}

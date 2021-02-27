package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

// New returns a new database connection pool.
func New(dbName string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", "database.sqlite")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open database")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Unable to ping database")
	}

	const q = `
	PRAGMA foreign_keys = ON;
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = 'WAL';
	PRAGMA cache_size = -64000;
	`
	_, err = db.Exec(q)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to set pragmas")
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

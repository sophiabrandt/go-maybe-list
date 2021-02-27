// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.SqliteDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed once they have been run in production.
//
// Using constants in a .go file is an easy way to ensure the schema is part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1.1,
		Description: "Create table users, maybes, tag",
		Script: `
-- Create users
CREATE TABLE users (
	user_id       UUID NOT NULL,
	name          TEXT NOT NULL,
	email         TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	active        BOOLEAN NOT NULL DEFAULT TRUE,
	created_at    TIMESTAMP NOT NULL,
	updated_at    TIMESTAMP NOT NULL,
PRIMARY KEY(user_id)
);
-- Create maybes
CREATE TABLE maybes (
	maybe_id       UUID NOT NULL,
	user_id        UUID NOT NULL,
	title          TEXT NOT NULL,
	url            TEXT NOT NULL,
	description    TEXT NOT NULL,
	created_at     TIMESTAMP NOT NULL,
	updated_at     TIMESTAMP NOT NULL,
PRIMARY KEY(maybe_id),
-- One-to-many relationship between users and maybes
FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- Create tags
CREATE TABLE tags (
	tag_id         UUID NOT NULL,
	name           TEXT NOT NULL UNIQUE,
PRIMARY KEY(tag_id)
);
-- Linking table for many-to-many relationship between Tag and Maybe
CREATE TABLE maybetags (
	tag_id         UUID NOT NULL,
	maybe_id       UUID NOT NULL,
	user_id		   UUID NOT NULL,
FOREIGN KEY(tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE,
FOREIGN KEY(maybe_id) REFERENCES maybes(maybe_id) ON DELETE CASCADE,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
UNIQUE(tag_id, maybe_id)
);
`,
	},
}

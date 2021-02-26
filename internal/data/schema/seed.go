package schema

import (
	"github.com/jmoiron/sqlx"
)

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.
const seeds = `
-- Create users, maybes and tags
INSERT INTO users (user_id, name, email, password_hash, active, created_at, updated_at) VALUES
	('bbc79841-7feb-4944-9971-07404558dfdd', 'user1', 'user1@email.com', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', 'true', '2019-01-01 00:00:03.000001+00', '2019-01-01 00:00:03.000001+00'),
	('6ae4a9bf-0bff-40d5-9dbc-ce93819f4208', 'user2', 'user2@email.com', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', 'true', '2019-01-01 00:00:03.000001+00', '2019-01-01 00:00:03.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO tags (tag_id, name) VALUES
	('c4c0b2e4-71a2-4676-bf04-d59667209923', 'books'),
	('82d074f5-9136-45ec-8df2-344528320cce', 'go'),
	('ab6f8437-ef58-4cde-9438-9fa6a9608764', 'watchlist')
	ON CONFLICT DO NOTHING;

INSERT INTO maybes (maybe_id, user_id, title, url, description, created_at, updated_at) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'bbc79841-7feb-4944-9971-07404558dfdd', 'Go Web Programming', 'https://www.manning.com/books/go-web-programming', 'how to build web applications with Go', '2019-01-01 00:00:03.000001+00', '2019-01-01 00:00:03.000001+00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6ae4a9bf-0bff-40d5-9dbc-ce93819f4208', 'video placeholder', 'https://www.youtube.com/watch?v=NpEaa2P7qZI', 'a video placeholder on youtube', '2019-01-01 00:00:03.000001+00', '2019-01-01 00:00:03.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO maybetags(maybe_id, user_id, tag_id) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'bbc79841-7feb-4944-9971-07404558dfdd', 'c4c0b2e4-71a2-4676-bf04-d59667209923'),
	('5cf37266-3473-4006-984f-9325122678b7', 'bbc79841-7feb-4944-9971-07404558dfdd', '82d074f5-9136-45ec-8df2-344528320cce'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '6ae4a9bf-0bff-40d5-9dbc-ce93819f420', 'ab6f8437-ef58-4cde-9438-9fa6a9608764')
	ON CONFLICT DO NOTHING;
`

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM users;
DELETE FROM maybes;
DELETE FROM tags;
`

package maybe

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrInvalidTag occurs when a tag cannot be found in the dabase and cannot be created.
	ErrInvalidTag = errors.New("tag not found")
)

// RepositoryDb defines the repository for the book service.
type RepositoryDb struct {
	Db *sqlx.DB
}

// New returns a pointer to a book repo.
func New(db *sqlx.DB) RepositoryDb {
	return RepositoryDb{Db: db}
}

// Query retrieves all maybes from the database.
func (r RepositoryDb) Query() (Infos, error) {
	const q = `
	SELECT
		m.*,
		u.user_id AS user_id
	FROM maybes as m
	LEFT JOIN
		users AS u ON m.user_id = u.user_id
	ORDER BY m.maybe_id
	`
	var maybes Infos
	if err := r.Db.Select(&maybes, q); err != nil {
		return maybes, errors.Wrap(err, "selecting maybes")
	}
	return maybes, nil
}

// QuerybyID retrieves a book by ID from the database.
func (r RepositoryDb) QueryByID(maybeID string) (Info, error) {
	if _, err := uuid.Parse(maybeID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		m.*,
		u.user_id AS user_id
	FROM maybes as m
	LEFT JOIN
		users AS u ON m.user_id = u.user_id
	WHERE
		m.maybe_id = $1
	`

	var maybe Info
	if err := r.Db.Get(&maybe, q, maybeID); err != nil {
		if err == sql.ErrNoRows {
			return maybe, ErrNotFound
		}
		return maybe, errors.Wrapf(err, "selecting maybe with ID %s", maybeID)
	}
	return maybe, nil
}

// QuerybyTitle retrieves an entry by quering the title from the database.
func (r RepositoryDb) QueryByTitle(title string) (Infos, error) {
	const q = `
	SELECT
		m.*,
		u.user_id AS user_id
	FROM maybes as m
	LEFT JOIN
		users AS u ON m.user_id = u.user_id
	WHERE
		m.title LIKE '%' || $1 || '%'
	`
	var maybes Infos
	if err := r.Db.Select(&maybes, q, title); err != nil {
		if err == sql.ErrNoRows {
			return maybes, ErrNotFound
		}
		return maybes, errors.Wrap(err, "selecting maybes by title")
	}
	return maybes, nil
}

// Create adds a new maybe to the database with pre-filled ID and date fields.
func (r RepositoryDb) Create(nm NewOrUpdateMaybe, userID string) (Info, error) {
	maybe := Info{
		ID:          uuid.New().String(),
		Title:       nm.Title,
		Url:         nm.Url,
		Description: nm.Description,
		Tags:        nil,
		DateCreated: time.Now().UTC().String(),
		DateUpdated: time.Now().UTC().String(),
	}

	if nm.Tags != nil {
		maybe.Tags = nm.Tags

		// For each tag:
		// * find the tag's ID in the database OR create a new tag
		// * create an entry in the linking table.
		for _, tagName := range maybe.Tags {
			row := r.Db.QueryRowx("SELECT tag_id FROM tags where name = $1", tagName)
			var tagID string
			err := row.Scan(&tagID)
			if err != nil {
				// tag does not exist in database, create
				tagID = uuid.New().String()
				const q = `
				INSERT INTO	tags (tag_id, name)
				VALUES ($1, $2)
				`
				if _, err := r.Db.Exec(q, tagID, tagName); err != nil {
					return Info{}, ErrInvalidTag
				}
			}
			const q = `
			INSERT INTO maybetags (maybe_id, tag_id)
			VALUES ($1, $2)
			`
			_, err = r.Db.Exec(q, maybe.ID, tagID)
			if err != nil {
				return Info{}, errors.Wrap(err, "inserting into linking table maybetags")
			}
		}
	}

	const q = `
	INSERT INTO maybes
		(maybe_id, user_id, title, url, description, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := r.Db.Exec(q, maybe.ID, userID, maybe.Title, maybe.Url, maybe.Description, maybe.DateCreated, maybe.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting new maybe")
	}
	return maybe, nil
}

// Update updates an existing maybe.
func (r RepositoryDb) Update(um NewOrUpdateMaybe, maybeID string) error {
	if _, err := uuid.Parse(maybeID); err != nil {
		return ErrInvalidID
	}

	maybe, err := r.QueryByID(maybeID)
	if err != nil {
		return errors.Wrap(err, "updating maybe")
	}

	if um.Title != "" {
		maybe.Title = um.Title
	}

	if um.Url != "" {
		maybe.Url = um.Url
	}

	if um.Description != "" {
		maybe.Description = um.Description
	}

	if um.Tags != nil {
		maybe.Tags = um.Tags

		// For each tag:
		// * find the tag's ID in the database OR create a new tag
		// * create an entry in the linking table.
		for _, tagName := range maybe.Tags {
			row := r.Db.QueryRowx("SELECT tag_id FROM tags where name = $1", tagName)
			var tagID string
			err := row.Scan(&tagID)
			if err != nil {
				// tag does not exist in database, create
				tagID = uuid.New().String()
				const q = `
				INSERT INTO	tags (tag_id, name)
				VALUES ($1, $2)
				`
				if _, err := r.Db.Exec(q, tagID, tagName); err != nil {
					return ErrInvalidTag
				}
			}
			const q = `
			INSERT INTO maybetags (maybe_id, tag_id)
			VALUES ($1, $2)
			`
			_, err = r.Db.Exec(q, maybe.ID, tagID)
			if err != nil {
				return errors.Wrap(err, "inserting into linking table maybetags")
			}
		}
	}

	const q = `
	UPDATE maybes
	SET
		title = $2,
		url = $3,
		description = $4
	WHERE
		maybe_id = $1
	`

	if _, err := r.Db.Exec(q, maybeID, maybe.Title, maybe.Url, maybe.Description); err != nil {
		return errors.Wrap(err, "updating product")
	}
	return nil
}

// Delete removes a maybe with given ID from the database.
func (r RepositoryDb) Delete(maybeID string) error {
	if _, err := uuid.Parse(maybeID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		maybes
	WHERE
		maybe_id = $1
	`

	if _, err := r.Db.Exec(q, maybeID); err != nil {
		return errors.Wrapf(err, "deleting maybe %s", maybeID)
	}

	return nil
}

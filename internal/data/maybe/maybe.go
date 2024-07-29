package maybe

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a specific maybe is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrInvalidTag occurs when a tag cannot be found in the dabase and cannot be created.
	ErrInvalidTag = errors.New("tag not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// MaybeRepository defines the repository for the maybe service.
type MaybeRepository struct {
	Db *sqlx.DB
}

// New returns a pointer to a maybe repo.
func New(db *sqlx.DB) MaybeRepository {
	return MaybeRepository{Db: db}
}

// Query retrieves all maybes from the database for the current user.
func (mr MaybeRepository) Query(userID string) (Infos, error) {
	const q = `
	SELECT
		m.*,
		u.user_id AS user_id
	FROM maybes as m
	LEFT JOIN
		users AS u ON m.user_id = u.user_id
	WHERE
		u.user_id = $1
	ORDER BY
		m.maybe_id
	`
	var maybes Infos
	if err := mr.Db.Select(&maybes, q, userID); err != nil {
		return maybes, errors.Wrap(err, "selecting maybes")
	}
	return maybes, nil
}

// QuerybyID retrieves a book by ID from the database.
func (r MaybeRepository) QueryByID(maybeID string, userID string) (Info, error) {
	if _, err := uuid.Parse(maybeID); err != nil {
		return Info{}, ErrInvalidID
	}

	// Get full details from maybes table
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
	if err := r.Db.Get(&maybe, q, maybeID, userID); err != nil {
		if err == sql.ErrNoRows {
			return maybe, ErrNotFound
		}
		return maybe, errors.Wrapf(err, "selecting maybe with ID %q", maybeID)
	}

	if maybe.UserID != userID {
		return Info{}, ErrForbidden
	}

	// Get all tags for the maybe
	const t = `
	SELECT t.*
	FROM tags AS t
	LEFT JOIN
		maybetags AS mt ON t.tag_id = mt.tag_id
	WHERE
		mt.maybe_id = $1
	`

	var tags []Tag
	if err := r.Db.Select(&tags, t, maybeID); err != nil {
		if err == sql.ErrNoRows {
			tags = nil
		}
		return maybe, errors.Wrapf(err, "selecting tags for maybe with ID %q", maybeID)
	}

	// if tags exist, add them to the model
	if tags != nil {
		maybe.Tags = tags
	}

	return maybe, nil
}

// QueryByTag queries the database for all maybes of a certain tag for the current user.
func (r MaybeRepository) QueryByTag(tagID string, userID string) (Infos, error) {
	var maybes Infos
	if _, err := uuid.Parse(tagID); err != nil {
		return maybes, ErrInvalidTag
	}

	const q = `
	SELECT
		m.*,
		u.user_id AS user_id
	FROM maybes as m
	LEFT JOIN
		users AS u ON m.user_id = u.user_id
	LEFT JOIN
		maybetags as mt ON m.maybe_id = mt.maybe_id
	LEFT JOIN
		tags AS t ON mt.tag_id = t.tag_id
	WHERE
		t.tag_id = $1 AND mt.user_id = $2
	ORDER BY
		m.maybe_id
	`
	if err := r.Db.Select(&maybes, q, tagID, userID); err != nil {
		return maybes, errors.Wrapf(err, "selecting maybes by tag %q", tagID)
	}

	return maybes, nil
}

// Create adds a new maybe to the database with pre-filled ID and date fields.
func (r MaybeRepository) Create(nm NewOrUpdateMaybe, userID string) (Info, error) {
	maybe := Info{
		ID:          uuid.New().String(),
		Title:       nm.Title,
		Url:         nm.Url,
		Description: nm.Description,
		Tags:        nil,
		DateCreated: time.Now().UTC().String(),
		DateUpdated: time.Now().UTC().String(),
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

	// add tags for the new maybe
	if len(nm.Tags) != 0 {
		// For each tag:
		// * find the tag's ID in the database OR create a new tag
		// * create an entry in the linking table.
		for _, tagName := range nm.Tags {
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
				_, err := r.Db.Exec(q, tagID, tagName)
				if err != nil {
					return Info{}, ErrInvalidTag
				}
			}

			const q = `
			INSERT OR IGNORE INTO maybetags (tag_id, maybe_id, user_id)
			VALUES ($1, $2, $3)
			`
			_, err = r.Db.Exec(q, tagID, maybe.ID, userID)
			if err != nil {
				return Info{}, errors.Wrap(err, "inserting into linking table maybetags")
			}
		}
	}

	return maybe, nil
}

// Update updates an existing maybe.
func (mr MaybeRepository) Update(um NewOrUpdateMaybe, maybeID string, userID string) error {
	if _, err := uuid.Parse(maybeID); err != nil {
		return ErrInvalidID
	}

	maybe, err := mr.QueryByID(maybeID, userID)
	if err != nil {
		switch errors.Cause(err) {
		case ErrInvalidID:
			return ErrInvalidID
		case ErrForbidden:
			return ErrForbidden
		case ErrNotFound:
			return ErrNotFound
		default:
			return errors.Wrap(err, "updating maybe")
		}
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

	// update the maybe model
	const q = `
	UPDATE maybes
	SET
		title = $2,
		url = $3,
		description = $4,
		updated_at = $5
	WHERE
		maybe_id = $1
	`
	if _, err := mr.Db.Exec(q, maybeID, maybe.Title, maybe.Url, maybe.Description, time.Now().UTC().String()); err != nil {
		return errors.Wrap(err, "updating product")
	}

	// updating/adding new tags
	if len(um.Tags) != 0 {
		// For each tag:
		// * find the tag's ID in the database OR create a new tag
		// * create an entry in the linking table.
		for _, tagName := range um.Tags {
			row := mr.Db.QueryRowx("SELECT tag_id FROM tags where name = $1", tagName)
			var tagID string
			err := row.Scan(&tagID)
			if err != nil {
				// tag does not exist in database, create
				tagID = uuid.New().String()
				const q = `
				INSERT INTO	tags (tag_id, name)
				VALUES ($1, $2)
				`
				if _, err := mr.Db.Exec(q, tagID, tagName); err != nil {
					return ErrInvalidTag
				}
			}

			// insert new tags into the linking table for the maybe
			const q = `
			INSERT OR IGNORE INTO
				maybetags (tag_id, maybe_id, user_id)
			VALUES
				($1, $2, $3)
			`
			_, err = mr.Db.Exec(q, tagID, maybe.ID, userID)
			if err != nil {
				return errors.Wrap(err, "inserting into linking table maybetags")
			}
		}

		// delete tags that don't exist anymore for the specific maybe
		tagNames := make([]interface{}, len(um.Tags))
		for i, tag := range um.Tags {
			tagNames[i] = tag
		}
		query, args, err := sqlx.In(`
		    DELETE FROM maybetags 
		    WHERE maybe_id = ?
		    AND tag_id NOT IN (
			SELECT tag_id
			FROM tags
			WHERE name IN (?)
		    )
		    `, maybe.ID, tagNames)
		if err != nil {
			return errors.Wrap(err, "preparing query for deleting old tags")
		}
		query = mr.Db.Rebind(query)
		_, err = mr.Db.Exec(query, args...)
		if err != nil {
			return errors.Wrapf(err, "deleting tags from linking table maybetags for: %q", um.Tags)
		}
	} else {
		// if the um.Tags slice contains no values, the user has deleted their tags,
		const t = `
		DELETE FROM
			maybetags
		WHERE
			maybe_id = $1
		`
		_, err = mr.Db.Exec(t, maybe.ID)
		if err != nil {
			return errors.Wrap(err, "deleting linking table maybetags")
		}
	}

	return nil
}

// Delete removes a maybe with given ID and its tags from the database.
func (mr MaybeRepository) Delete(maybeID string) error {
	if _, err := uuid.Parse(maybeID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		maybes
	WHERE
		maybe_id = $1
	`

	if _, err := mr.Db.Exec(q, maybeID); err != nil {
		return errors.Wrapf(err, "deleting maybe %q", maybeID)
	}

	// delete orphaned tags
	const d = `
	DELETE FROM
		tags AS t
	WHERE NOT EXISTS (
		SELECT NULL FROM maybetags AS mt
		WHERE
			t.tag_id = mt.tag_id
	)
	`

	if _, err := mr.Db.Exec(d); err != nil {
		return errors.Wrap(err, "deleting tags from linking table maybetags")
	}

	return nil
}

// QueryTags returns all (unique) tags for a given user.
func (mr MaybeRepository) QueryTags(userID string) (Tags, error) {
	const q = `
	SELECT DISTINCT
		t.*
	FROM
		tags AS t
	LEFT JOIN
		maybetags AS mt ON mt.tag_id = t.tag_id
	WHERE
		mt.user_id = $1
	ORDER BY
		t.tag_id
	`

	var tags Tags
	if err := mr.Db.Select(&tags, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return tags, ErrNotFound
		}
		return tags, errors.Wrap(err, "selecting tags")
	}

	return tags, nil
}

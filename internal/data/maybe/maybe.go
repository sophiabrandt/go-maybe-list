package maybe

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// BookRepositoryDb defines the repository for the book service.
type RepositoryDb struct {
	Db *sqlx.DB
}

// Repo is the interface for the maybe repository.
type Repo interface {
	Query() (Infos, error)
	QueryByID(id string) (Info, error)
	QueryByTitle(title string) (Info, error)
	Create(maybe NewMaybe) (Info, error)
	Delete(id string) error
}

// New returns a pointer to a book repo.
func New(db *sqlx.DB) *RepositoryDb {
	return &RepositoryDb{Db: db}
}

// Query retrieves all maybes from the database.
func (r *RepositoryDb) Query() (Infos, error) {
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
func (r *RepositoryDb) QueryByID(id string) (Info, error) {
	if _, err := uuid.Parse(id); err != nil {
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
	if err := r.Db.Get(&maybe, q, id); err != nil {
		if err == sql.ErrNoRows {
			return maybe, ErrNotFound
		}
		return maybe, errors.Wrapf(err, "selecting maybe with ID %s", id)
	}
	return maybe, nil
}

// QuerybyTitle retrieves an entry by quering the title from the database.
func (r *RepositoryDb) QueryByTitle(title string) (Infos, error) {
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

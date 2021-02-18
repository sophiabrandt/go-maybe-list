package maybe

import (
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

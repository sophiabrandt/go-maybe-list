package user

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrDuplicateEmail occus when the email exists in the database.
	ErrDuplicateEmail = errors.New("email already in use")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// RepositoryDb defines the repository for the book service.
type RepositoryDb struct {
	Db *sqlx.DB
}

// Repo is the interface for the maybe repository.
type Repo interface {
	Create(user NewUser) (Info, error)
	Authenticate(user Info) (string, error)
}

// New returns a pointer to a book repo.
func New(db *sqlx.DB) *RepositoryDb {
	return &RepositoryDb{Db: db}
}

func (r *RepositoryDb) Create(user NewUser) (Info, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return Info{}, errors.Wrap(err, "generating password hash")
	}

	usr := Info{
		ID:           uuid.New().String(),
		Name:         user.Name,
		Email:        user.Email,
		Active:       true,
		PasswordHash: hash,
		DateCreated:  time.Now().UTC().String(),
		DateUpdated:  time.Now().UTC().String(),
	}

	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, active, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)`

	if _, err = r.Db.Exec(q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Active, usr.DateCreated, usr.DateUpdated); err != nil {
		var sqLiteError *sqlite.Error
		if errors.As(err, &sqLiteError) {
			if sqLiteError.Code() == 2067 && strings.Contains(sqLiteError.Error(), "users.email") {
				return Info{}, ErrDuplicateEmail
			}
		}
		return Info{}, errors.Wrap(err, "inserting user")
	}

	return usr, nil
}

func (r *RepositoryDb) Authenticate(email, password string) (string, error) {
	var id string
	var hash []byte
	const q = `
	SELECT
		user_id, password_hash
	FROM
		users
	WHERE 
		email = $1
	AND
		active = TRUE
	`
	row := r.Db.QueryRowx(q, email)
	err := row.Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, ErrAuthenticationFailure
		}
		return id, err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return id, ErrAuthenticationFailure
		} else {
			return id, err
		}
	}

	return id, nil
}

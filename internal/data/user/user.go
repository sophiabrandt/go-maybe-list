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

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// RepositoryDb defines the repository for the user service.
type RepositoryDb struct {
	Db *sqlx.DB
}

// New returns a pointer to a user repo.
func New(db *sqlx.DB) RepositoryDb {
	return RepositoryDb{Db: db}
}

// QueryByID gets the specified user from the database.
func (r RepositoryDb) QueryByID(userID string) (Info, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = $1`

	var usr Info
	if err := r.Db.Get(&usr, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", userID)
	}

	return usr, nil
}

// Create inserts a new user into the database.
func (r RepositoryDb) Create(user NewUser) (Info, error) {
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

// Authenticate queries the database for a user with a matching pasword.
func (r RepositoryDb) Authenticate(email, password string) (string, error) {
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

func (r RepositoryDb) ChangePassword(currentPassword, newPassword, userID string) error {
	var currentPasswordHash []byte
	const p = `
	SELECT password_hash
	FROM users
	WHERE user_id = $1
	`

	if err := r.Db.Get(&currentPasswordHash, p, userID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrapf(err, "selecting hashed password for %q", userID)
	}

	err := bcrypt.CompareHashAndPassword(currentPasswordHash, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrAuthenticationFailure
		} else {
			return err
		}
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	const q = `
	UPDATE
		users
	SET
		password_hash = $2
	WHERE
		user_id = $1
	`

	if _, err := r.Db.Exec(q, userID, newPasswordHash); err != nil {
		return errors.Wrapf(err, "updating password for user %q", userID)
	}

	return nil
}

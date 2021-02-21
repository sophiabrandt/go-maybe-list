package user

// Info is the model for a user.
type Info struct {
	ID           string `db:"user_id"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	Active       bool   `db:"active"`
	DateCreated  string `db:"created_at"`
	DateUpdated  string `db:"updated_at"`
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

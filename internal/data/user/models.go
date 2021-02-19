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
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,secure_password"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

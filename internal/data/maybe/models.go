package maybe

// Info is the model for maybes.
type Info struct {
	ID          string   `db:"maybe_id"`
	UserID      string   `db:"user_id"`
	Title       string   `db:"title"`
	Url         string   `db:"url"`
	Description string   `db:"description"`
	Tags        []string `db:"tags"`
	DateCreated string   `db:"created_at"`
	DateUpdated string   `db:"updated_at"`
}

// Tag is the model for a tag.
type Tag struct {
	ID   string `db:"tag_id"`
	Name string `db:"tag_name"`
}

type Infos []Info

// NewMaybe is the data for creating a new maybe.
// Adding Tags is optional.
type NewMaybe struct {
	Title       string
	Url         string
	Description string
	Tags        []string
}

// NewTag is the data for creating a new tag.
type NewTag struct {
	Name string `db:"name"`
}

// UpdateMaybe defines the information for updating an existing maybe.
// Fields are optional.
type UpdateBook struct {
	Title       *string
	Url         *string
	Description *string
	Tags        []string
}

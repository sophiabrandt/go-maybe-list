package maybe

// Tag is the model for a tag.
type Tag struct {
	ID   string `db:"tag_id"`
	Name string `db:"name"`
}

type Tags []Tag

// Info is the model for maybes.
type Info struct {
	ID          string `db:"maybe_id"`
	UserID      string `db:"user_id"`
	Title       string `db:"title"`
	Url         string `db:"url"`
	Description string `db:"description"`
	Tags        []Tag  `db:"tags"`
	DateCreated string `db:"created_at"`
	DateUpdated string `db:"updated_at"`
}

type Infos []Info

// NewOrUpdateMaybe is the data for creating a new maybe
// or updating an existing maybe.
// Adding Tags is optional.
type NewOrUpdateMaybe struct {
	Title       string
	Url         string
	Description string
	Tags        []string
}

// NewTag is the data for creating a new tag.
type NewTag struct {
	Name string `db:"name"`
}

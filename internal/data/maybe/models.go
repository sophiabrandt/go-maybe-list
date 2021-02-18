package maybe

// Info is the model for maybes.
type Info struct {
	ID          string `db:"maybe_id"`
	UserID      string `db:"user_id"`
	Title       string `db:"title"`
	Url         string `db:"url"`
	Description string `db:"description"`
	Tags        []*Tag `db:"tags"`
	CreatedDate string `db:"created_at"`
	UpdatedDate string `db:"updated_at"`
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
	Title       string `json:"title" validate:"required"`
	Url         string `json:"url" validate:"required,url,max=255"`
	Description string `json:"description" validate:"required,max=255"`
	Tags        []*Tag `json:"tags" validate:"omitempty"`
}

// UpdateBook defines the information for updating an existing product.
// Fields are optional.
type UpdateBook struct {
	Title       *string `json:"title" validate:"omitempty"`
	Url         *string `json:"image_url" validate:"omitempty,url"`
	Description *string `json:"description" validate:"omitempty"`
	Tags        []*Tag  `json:"tags" validate:"omitempty"`
}

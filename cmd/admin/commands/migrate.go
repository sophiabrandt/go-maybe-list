package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/internal/adapter/database"
	"github.com/sophiabrandt/go-maybe-list/internal/data/schema"
)

// Migrate creates the schema in the database.
func Migrate() error {
	db, err := database.New()
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return errors.Wrap(err, "migrate database")
	}

	fmt.Println("migrations complete")
	return nil
}

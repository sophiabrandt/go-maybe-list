package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/cmd/admin/commands"
)

func main() {
	log := log.New(os.Stdout, "ADMIN: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("admin: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	command := flag.String("action", "", "admin command: migrate | seed")
	dbName := flag.String("dbName", "database.sqlite", "database name")
	flag.Parse()

	switch *command {
	case "migrate":
		if err := commands.Migrate(*dbName); err != nil {
			return errors.Wrap(err, "migrating database")
		}
	case "seed":
		if err := commands.Seed(*dbName); err != nil {
			return errors.Wrap(err, "seeding database")
		}
	default:
		fmt.Println("ADMIN: Possible commands:")
		fmt.Println("-action=\"migrate\": create the schema in the database")
		fmt.Println("-action=\"seed\": add data to the database")
	}

	return nil
}

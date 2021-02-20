package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-maybe-list/internal/adapter/database"
	"github.com/sophiabrandt/go-maybe-list/internal/env"
	"github.com/sophiabrandt/go-maybe-list/internal/server"
	"github.com/sophiabrandt/go-maybe-list/internal/web/handlers"
	"github.com/sophiabrandt/go-maybe-list/internal/web/session"
	"github.com/sophiabrandt/go-maybe-list/internal/web/templates"
)

func main() {
	// make a channel to listen for interrupt or terminal signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	log := log.New(os.Stdout, "MAYBE-LIST: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// listen for interrupt signals
	go func() {
		oscall := <-c
		log.Printf("main: system call: %+v", oscall)
		cancel()
	}()

	// run APP
	if err := run(ctx, log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *log.Logger) error {
	log.Println("main : Started : Application initializing")

	addr := flag.String("addr", "0.0.0.0:4000", "Http network address")
	secret := flag.String("secret", "60FNA&6bdH+FnhG306I6MNCY8bv_WjwDcjB", "Secret key")
	flag.Parse()

	// database
	db, err := database.New()
	if err != nil {
		return errors.Wrap(err, "could not start server")
	}
	defer db.Close()

	// initialize global dependencies
	tc, err := templates.NewCache("./ui/html")
	ses := session.New(*secret)

	env := env.New(log, tc, ses)

	router := handlers.New(env, db)

	// create server
	srv := server.New(*addr, router)

	go func() {
		log.Printf("main: APP listening on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("main: %+s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("main: Start shutdown")

	// shutdown server
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("main: Shutdown Failed: %+s", err)
	}

	log.Println("main: APP exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return nil
}

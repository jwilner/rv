package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		migrationsDir = flag.String("source-dir", "migrations/", "the migrations directory")
		targetURL     = flag.String(
			"target-url",
			"",
			"the target database to upgrade; defaults to DATABASE_URL from the environment",
		)
		version = flag.Uint("version", 0, "the target version to migrate to; if not provided, migrates to latest version")
	)

	flag.Parse()

	err := run(
		*migrationsDir,
		*targetURL,
		*version,
	)

	switch {
	case err == nil:
		log.Println("Successfully migrated")
	case errors.Is(err, migrate.ErrNoChange):
		log.Println("Migrations already up-to-date")
	default:
		log.Fatalf("Failed migrating: %v", err)
	}
}

func run(migDir, dbURL string, targetVer uint) error {
	if _, err := ioutil.ReadDir(migDir); err != nil {
		return fmt.Errorf("error accessing %v: %v", migDir, err)
	}

	if dbURL == "" {
		var ok bool
		if dbURL, ok = os.LookupEnv("DATABASE_URL"); !ok {
			return errors.New("no target url provided and DATABASE_URL not present in environment")
		}
	}

	m, err := migrate.New("file://"+migDir, dbURL)
	if err != nil {
		return fmt.Errorf("migrate.New(%v, redacted): %w", migDir, err)
	}

	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		select {
		case <-sigs:
			m.GracefulStop <- true
		case <-ctx.Done():
		}
	}()

	if targetVer == 0 {
		return m.Up()
	}
	return m.Migrate(targetVer)
}

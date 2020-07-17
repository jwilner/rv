package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func main() {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL must be present in the environment and wasn't")
	}
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	binPath, err := exec.LookPath("sqlboiler")
	if err != nil {
		log.Fatalf("unable to find sqlboiler bin in the path")
	}

	driverPath, err := exec.LookPath("sqlboiler-psql")
	if err != nil {
		log.Fatalf("unable to find sqlboiler driver bin in the path")
	}

	cmd := exec.CommandContext(
		context.Background(),
		binPath,
		driverPath,
		"--no-auto-timestamps",
		"--no-back-referencing",
		"--no-hooks",
		"--no-tests",
		"--wipe",
	)
	cmd.Env = []string{
		"PSQL_DBNAME=" + strings.Trim(u.Path, "/"),
		"PSQL_SCHEMA=rv",
		fmt.Sprintf("PSQL_HOST=%s", u.Hostname()),
		fmt.Sprintf("PSQL_USER=%s", u.User.Username()),
		"PSQL_SSLMODE=disable",
	}
	if password, ok := u.User.Password(); ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PSQL_PASS=%s", password))
	}
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("sqlboiler command failed: %v", err)
	}
	log.Println("Successfully generated models")
}

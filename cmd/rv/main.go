package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jwilner/rv/internal/platform"
)

func main() {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	stdoutLog, _ := strconv.ParseBool(os.Getenv("STDOUT_LOG"))

	log.SetFlags(log.LstdFlags)
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if stdoutLog {
		log.SetOutput(os.Stdout)
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	staticDir := os.Getenv("STATIC_DIR") // where to serve static assets from, if at all

	signingKey := os.Getenv("TOKEN_SIGNING_KEY")
	if signingKey == "" {
		log.Fatalln("TOKEN_SIGNING_KEY must be provided")
	}

	tokDurS, ok := os.LookupEnv("TOKEN_DURATION")
	if !ok {
		log.Fatalln("TOKEN_DURATION must be provided")
	}
	length, err := time.ParseDuration(tokDurS)
	if err != nil || length == 0 {
		log.Fatalf("Invalid TOKEN_DURATION %q: %v", tokDurS, err)
	}

	if err := platform.Run(debug, dbURL, ":"+port, staticDir, signingKey, length); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/jwilner/rv/internal/platform"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	staticDir := os.Getenv("STATIC_DIR") // where to serve static assets from, if at all
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	stdoutLog, _ := strconv.ParseBool(os.Getenv("STDOUT_LOG"))

	log.SetFlags(log.LstdFlags)
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if stdoutLog {
		log.SetOutput(os.Stdout)
	}

	if err := platform.Run(debug, dbURL, ":"+port, staticDir); err != nil {
		log.Fatal(err)
	}
}

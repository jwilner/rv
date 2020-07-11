package main

import (
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func run(addr string) error {
	return http.ListenAndServe(addr, route(&handler{&db{nil}}))
}

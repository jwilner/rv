package main

import (
	"context"
	"database/sql"
	"net/url"
	"regexp"
)

type db struct {
	db *sql.DB
}

func (d *db) loadBallotFromPath(ctx context.Context, path string) (*ballot, error) {
	key := parseKey(ballotKeyPathRegex, path)
	if key == "" {
		return nil, sql.ErrNoRows
	}
	return nil, sql.ErrNoRows
}

func (d *db) loadElectionFromPath(ctx context.Context, path string) (*ballot, error) {
	key := parseKey(electionKeyPathRegex, path)
	if key == "" {
		return nil, sql.ErrNoRows
	}
	return nil, sql.ErrNoRows
}

var (
	ballotKeyPathRegex   = regexp.MustCompile(`^/b/([a-zA-Z0-9]+)/?$`)
	electionKeyPathRegex = regexp.MustCompile(`^/e/([a-zA-Z0-9]+)/?$`)
)

func parseKey(reg *regexp.Regexp, path string) string {
	if match := ballotKeyPathRegex.FindStringSubmatch(path); len(match) > 1 {
		return match[1]
	}
	return ""
}

type ballot struct {
	key       string
	completed bool
}

func ballotFromForm(vals url.Values) (*ballot, error) {
	return nil, sql.ErrNoRows
}

type election struct {
	key string
}

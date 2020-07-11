package main

import (
	"database/sql"
	"errors"
	"net/http"
)

type handler struct {
	db *db
}

func (h *handler) getIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("cool index"))
}

func (h *handler) postIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("cool index"))
}

func (h *handler) getBallot(w http.ResponseWriter, r *http.Request) {
	b, err := h.db.loadBallotFromPath(r.Context(), r.URL.Path)
	if handleError(w, r, err) {
		return
	}
	if b.completed {
		// serve results
		return
	}
	_, _ = w.Write([]byte("cool ballot"))
}

func (h *handler) postBallot(w http.ResponseWriter, r *http.Request) {
	b, err := h.db.loadBallotFromPath(r.Context(), r.URL.Path)
	if handleError(w, r, err) {
		return
	}
	if b.completed {
		http.Error(w, "ballot already completed", http.StatusUnprocessableEntity)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("cool ballot"))
}

func (h *handler) getElection(w http.ResponseWriter, r *http.Request) {
	_, err := h.db.loadElectionFromPath(r.Context(), r.URL.Path)
	if handleError(w, r, err) {
		return
	}
	_, _ = w.Write([]byte("cool election"))
}

func (h *handler) postElection(w http.ResponseWriter, r *http.Request) {
	_, err := h.db.loadElectionFromPath(r.Context(), r.URL.Path)
	if handleError(w, r, err) {
		return
	}
	_, _ = w.Write([]byte("cool election"))
}

func handleError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, sql.ErrNoRows):
		serve404(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return true
}

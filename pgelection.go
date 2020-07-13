package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jackc/pgtype"
)

// electionPage is the main page for managing an election
type electionPage struct {
	election
	Form electionForm
}

// electionForm is the election management form
type electionForm struct {
	form
}

func (e *electionForm) validate() bool {
	return true
}

// election is the election database model
type election struct {
	Name      string
	Key       string
	Choices   []string
	CreatedAt pgtype.Timestamp
}

func (h *handler) getElection(w http.ResponseWriter, r *http.Request) {
	_, err := h.txM.loadElection(r.Context())
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "election.html", &electionPage{})
}

func (h *handler) postElection(w http.ResponseWriter, r *http.Request) {
	_, err := h.txM.loadElection(r.Context())
	if handleError(w, r, err) {
		return
	}
	_, _ = w.Write([]byte("cool election"))
}

func (t *txMgr) loadElection(ctx context.Context) (*election, error) {
	key := electionKey(ctx)
	if key == "" {
		return nil, sql.ErrNoRows
	}
	return nil, sql.ErrNoRows
}

package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jackc/pgtype"
)

// ballotPage is the view for the page, backing the template
type ballotPage struct {
	Form ballotForm
}

// ballotForm is the ballot submission form
type ballotForm struct {
	form
}

// ballot is the database model for a ballot
type ballot struct {
	key       string
	name      string
	choices   pgtype.TextArray
	createdAt pgtype.Timestamptz
}

func (h *handler) getBallot(w http.ResponseWriter, r *http.Request) {
	var b *ballot
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		b, err = loadBallot(r.Context(), tx)
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "ballot.html", &ballotPage{})
}

func (h *handler) postBallot(w http.ResponseWriter, r *http.Request) {
	var b *ballot
	err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		b, err = loadBallot(r.Context(), tx)
		return
	})
	if handleError(w, r, err) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "unable to parse body", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("cool ballot"))
}

func loadBallot(ctx context.Context, tx *sql.Tx) (*ballot, error) {
	key := ballotKey(ctx)
	if key == "" {
		return nil, sql.ErrNoRows
	}
	b := ballot{key: key}
	err := tx.QueryRowContext(
		ctx,
		`SELECT name, choices, created_at FROM rv.ballot WHERE key = $1`,
	).Scan(
		&b.name,
		&b.choices,
		&b.createdAt,
	)
	return &b, err
}

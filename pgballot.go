package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jwilner/rv/models"
)

// ballotPage is the view for the page, backing the template
type ballotPage struct {
	Form ballotForm
}

// ballotForm is the ballot submission form
type ballotForm struct {
	form
}

func (h *handler) getBallot(w http.ResponseWriter, r *http.Request) {
	var b *models.Ballot
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
	var b *models.Ballot
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

func loadBallot(ctx context.Context, tx *sql.Tx) (*models.Ballot, error) {
	key := ballotKey(ctx)
	if key == "" {
		return nil, sql.ErrNoRows
	}
	return models.Ballots(models.BallotWhere.Key.EQ(key)).One(ctx, tx)
}

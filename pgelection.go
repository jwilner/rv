package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jwilner/rv/models"
)

// electionPage is the main page for managing an election
type electionPage struct {
	*models.Election
	Form electionForm
}

// electionForm is the election management form
type electionForm struct {
	form
}

func (e *electionForm) validate() bool {
	return true
}

func (h *handler) getElection(w http.ResponseWriter, r *http.Request) {
	var el *models.Election
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		el, err = models.Elections(models.ElectionWhere.Key.EQ(keyParam(ctx))).One(ctx, tx)
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "election.html", &electionPage{Election: el})
}

func (h *handler) postElection(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("cool election"))
}

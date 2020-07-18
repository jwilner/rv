package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jwilner/rv/models"
)

// electionPage is the main page for managing an election
type electionPage struct {
	*models.Election
	Report *report
	Form   electionForm
}

// electionForm is the election management form
type electionForm struct {
	form
}

func (e *electionForm) validate() bool {
	return true
}

func (h *handler) getElection(w http.ResponseWriter, r *http.Request) {
	var (
		el    *models.Election
		votes []*models.Vote
	)
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		if el, err = models.Elections(models.ElectionWhere.Key.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", keyParam(ctx), err)
		}
		if votes, err = models.Votes(models.VoteWhere.ElectionID.EQ(el.ID)).All(ctx, tx); err != nil {
			return fmt.Errorf("models.Votes electionId=%d: %w", el.ID, err)
		}
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "election.html", &electionPage{Election: el, Report: calculateReport(votes)})
}

func (h *handler) postElection(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("cool election"))
}

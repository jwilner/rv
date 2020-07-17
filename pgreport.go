package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jwilner/rv/models"
)

type reportPage struct {
	*models.Election
	Votes []*models.Vote
}

func (h *handler) getReport(w http.ResponseWriter, r *http.Request) {
	var (
		e  *models.Election
		vs []*models.Vote
	)
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		if e, err = models.Elections(models.ElectionWhere.BallotKey.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
			return fmt.Errorf("model.Elections ballotKey=%q: %w", keyParam(ctx), err)
		}
		if vs, err = models.Votes(models.VoteWhere.ElectionID.EQ(e.ID)).All(ctx, tx); err != nil {
			return fmt.Errorf("model.Votes electionId=%v: %w", keyParam(ctx), err)
		}
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "report.html", &reportPage{Election: e, Votes: vs})
}

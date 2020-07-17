package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/models"
)

// votePage is the view for the page, backing the template
type votePage struct {
	Form voteForm

	*models.Election
}

// voteForm is the ballot submission form
type voteForm struct {
	form

	Name    string
	Choices string
}

func (b *voteForm) unmarshal(vals url.Values) {
	b.Name = vals.Get("name")
	b.Choices = vals.Get("choices")
}

func (b *voteForm) validate() bool {
	if len(b.Name) == 0 {
		b.setErrorf("Name", "must provide a name")
	}
	if choices := parseChoices(b.Choices); len(choices) == 0 {
		b.setErrorf("Choices", "must have at least one choice")
	} else {
		counts := make(map[string]int)
		for _, c := range choices {
			counts[c]++
		}

		for _, c := range choices {
			if counts[c] > 1 {
				b.setErrorf("Choices", "%q occurs more %d times -- can only occur once", c, counts[c])
			}
			counts[c] = 0
		}
	}
	return b.checkErrors()
}

func (h *handler) getVote(w http.ResponseWriter, r *http.Request) {
	var e *models.Election
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		e, err = models.Elections(models.ElectionWhere.BallotKey.EQ(keyParam(ctx))).One(ctx, tx)
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "vote.html", &votePage{Election: e})
}

func (h *handler) postVote(w http.ResponseWriter, r *http.Request) {
	if handleError(w, r, r.ParseForm()) {
		return
	}

	var vf voteForm
	vf.unmarshal(r.PostForm)

	if !vf.validate() {
		log.Printf("validation failed: %v", vf.Errors)
		h.tmpls.render(r.Context(), w, "vote.html", &votePage{Form: vf})
		return
	}

	var e *models.Election
	err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		if e, err = models.Elections(models.ElectionWhere.BallotKey.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
			return fmt.Errorf("model.Elections: %w", err)
		}

		v := models.Vote{
			ElectionID: e.ID,
			Name:       vf.Name,
			Choices:    parseChoices(vf.Choices),
			CreatedAt:  time.Now().UTC(),
		}

		return v.Insert(ctx, tx, boil.Infer())
	})
	if handleError(w, r, err) {
		return
	}
	http.Redirect(w, r, "/r/"+e.BallotKey, http.StatusFound)
}

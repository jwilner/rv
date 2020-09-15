package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgtype"

	"github.com/jackc/pgconn"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/models"
)

// votePage is the view for the page, backing the template
type votePage struct {
	Form voteForm

	Election *election
	Now      time.Time
}

func (v *votePage) TimeRemaining() string {
	if !v.Election.CloseScheduled(v.Now) {
		return ""
	}
	return fmt.Sprint(v.Election.CloseTime().Sub(v.Now).Round(time.Second))
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

func (b *voteForm) validate() (normalizedSlice, bool) {
	if len(b.Name) == 0 {
		b.setErrorf("Question", "must provide a name")
	}

	choices := parseChoices(b.Choices)
	validateNonDupe(choices, &b.form)
	return choices, b.checkErrors()
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

	now := time.Now().UTC()

	el := newElection(e)
	if el.Closed(now) {
		http.Redirect(w, r, "/r/"+keyParam(r.Context()), http.StatusFound)
		return
	}

	h.tmpls.render(r.Context(), w, "vote.html", &votePage{Now: now, Election: el})
}

var errValidation = errors.New("invalid against db model")

func (h *handler) postVote(w http.ResponseWriter, r *http.Request) {
	if handleError(w, r, r.ParseForm()) {
		return
	}

	now := time.Now().UTC()

	var vf voteForm
	vf.unmarshal(r.PostForm)

	choices, ok := vf.validate()

	if !ok {
		var e *models.Election
		err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) (err error) {
			if e, err = models.Elections(models.ElectionWhere.BallotKey.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
				return fmt.Errorf("model.Elections: %w", err)
			}
			return nil
		})
		if handleError(w, r, err) {
			return
		}
		h.tmpls.render(r.Context(), w, "vote.html", &votePage{Now: now, Form: vf, Election: newElection(e)})
		return
	}

	var e *election
	err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		var el *models.Election
		if el, err = models.Elections(models.ElectionWhere.BallotKey.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
			return fmt.Errorf("model.Elections: %w", err)
		} else {
			e = newElection(el)
		}
		if e.Close.Status == pgtype.Present && !e.CloseTime().After(time.Now()) {
			vf.setErrorf("Choices", "election has already closed.")
			return errValidation
		}
		if undefined := difference(choices, normalize(e.Choices)); len(undefined) > 0 {
			vf.setErrorf("Choices", "Unknown choices: %v", strings.Join(undefined.raw(), ", "))
			return errValidation
		}
		v := models.Vote{
			ElectionID: e.ID,
			Name:       vf.Name,
			Choices:    choices.raw(),
		}
		_ = v.CreatedAt.Set(now)

		return v.Insert(ctx, tx, boil.Infer())
	})
	if errors.Is(err, errValidation) {
		log.Printf("validation failed: %v", vf.Errors)
		h.tmpls.render(r.Context(), w, "vote.html", &votePage{Form: vf, Election: e, Now: now})
		return
	}
	if isDupe(err, "vote_election_id_name_key") {
		vf.setErrorf("Question", "This name has already been used")
		h.tmpls.render(r.Context(), w, "vote.html", &votePage{Form: vf, Election: e})
		return
	}
	if handleError(w, r, err) {
		return
	}
	http.Redirect(w, r, "/r/"+e.BallotKey, http.StatusFound)
}

func isDupe(err error, constraintName string) bool {
	pgErr := new(pgconn.PgError)
	return errors.As(err, &pgErr) &&
		pgErr.Code == "23505" &&
		pgErr.SchemaName == "rv" &&
		pgErr.ConstraintName == constraintName
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jwilner/rv/models"
)

type reportPage struct {
	Election *election
	Report   *report
	Now      time.Time
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
	h.tmpls.render(
		r.Context(),
		w,
		"report.html",
		&reportPage{Election: newElection(e), Report: calculateReport(vs), Now: time.Now().UTC()},
	)
}

type report struct {
	Winner string
	Steps  []*step
}

type step struct {
	Round      int
	Eliminated []string
	Remaining  []*remainingVote
	Counted    map[string]int
}

type remainingVote struct {
	Name    string
	Choices []string
}

func calculateReport(vs []*models.Vote) *report {
	votes := make([]*vote, 0, len(vs))
	for _, v := range vs {
		votes = append(votes, &vote{v.Name, normalize(v.Choices)})
	}

	var r report
	for {
		var (
			eliminated []string
			counted    = make(map[string]int)
		)
		{
			var copied []*vote
			for _, v := range votes {
				if len(v.choices) == 0 {
					eliminated = append(eliminated, v.name)
					continue
				}
				copied = append(copied, v)
				counted[v.choices[0].normalized]++
			}
			votes = copied
		}
		var (
			min, max       *int
			minVal, maxVal = make([]string, 0), make([]string, 0)
		)
		for v, c := range counted {
			if min == nil || c < *min {
				c := c
				min = &c
				minVal = append(minVal[:0], v) // reset
			} else if c == *min {
				minVal = append(minVal, v)
			}
			if max == nil || c > *max {
				c := c
				max = &c
				maxVal = append(maxVal[:0], v) // reset
			} else if c == *max {
				maxVal = append(maxVal, v)
			}
		}

		if max == nil {
			break
		}

		s := step{
			Round:      len(r.Steps) + 1,
			Eliminated: eliminated,
			Remaining:  make([]*remainingVote, 0, len(vs)),
			Counted:    counted,
		}
		for _, v := range votes {
			s.Remaining = append(s.Remaining, &remainingVote{Name: v.name, Choices: v.choices.raw()})
		}
		r.Steps = append(r.Steps, &s)

		if *max > (len(votes) / 2) {
			r.Winner = maxVal[0]
			break
		}

		// choose min
		least := minVal[0]
		for _, v := range minVal {
			if v < least {
				least = v
			}
		}

		// remove least popular
		for _, v := range votes {
			for i := range v.choices {
				if v.choices[i].normalized == least {
					v.choices = append(v.choices[:i], v.choices[i+1:]...)
					break
				}
			}
		}
	}

	return &r
}

type vote struct {
	name    string
	choices normalizedSlice
}

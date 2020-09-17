package platform

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/jwilner/rv/internal/models"
)

func (h *handler) Report(ctx context.Context, req *rvapi.ReportRequest) (*rvapi.ReportResponse, error) {
	var votes []*models.Vote
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		var el *models.Election
		if el, err = models.Elections(models.ElectionWhere.Key.EQ(req.Key)).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", req.Key, err)
		}
		if votes, err = el.Votes().All(ctx, tx); err != nil {
			return fmt.Errorf("Election.Votes byKey=%v: %w", req.Key, err)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return &rvapi.ReportResponse{Report: calculateReport(votes)}, nil
}

func calculateReport(vs []*models.Vote) *rvapi.Report {
	type vote struct {
		name    string
		choices normalizedSlice
	}

	votes := make([]*vote, 0, len(vs))
	for _, v := range vs {
		votes = append(votes, &vote{v.Name, normalize(v.Choices)})
	}

	var r rvapi.Report
	for {
		var (
			eliminated []string
			counted    = make(map[string]int32)
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
			min, max       *int32
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

		s := rvapi.Round{
			Eliminated: eliminated,
			Remaining:  make([]*rvapi.RemainingVote, 0, len(vs)),
			Counted:    counted,
		}
		for _, v := range votes {
			s.Remaining = append(s.Remaining, &rvapi.RemainingVote{Name: v.name, Choices: v.choices.raw()})
		}
		r.Rounds = append(r.Rounds, &s)

		if *max > int32(len(votes)/2) {
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

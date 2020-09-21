package platform

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/jwilner/rv/internal/models"
)

func (h *handler) Report(ctx context.Context, req *rvapi.ReportRequest) (*rvapi.ReportResponse, error) {
	var votes []*models.Vote
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		var el *models.Election
		if req.Key != "" {
			if el, err = models.Elections(models.ElectionWhere.Key.EQ(req.Key)).One(ctx, tx); err != nil {
				return fmt.Errorf("models.Elections key=%v: %w", req.Key, err)
			}
		} else {
			if el, err = models.Elections(models.ElectionWhere.BallotKey.EQ(req.BallotKey)).One(ctx, tx); err != nil {
				return fmt.Errorf("models.Elections ballotKey=%v: %w", req.BallotKey, err)
			}
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
		choices []string
	}

	votes := make([]*vote, 0, len(vs))
	for _, v := range vs {
		votes = append(votes, &vote{v.Name, v.Choices})
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
				counted[v.choices[0]]++
			}
			votes = copied
		}

		tallies := make(sortableTallies, 0, len(counted))
		for k, v := range counted {
			tallies = append(tallies, &rvapi.Tally{Count: v, Choice: k})
		}
		sort.Sort(sort.Reverse(tallies))

		if len(tallies) == 0 {
			break
		}

		r.Rounds = append(r.Rounds, &rvapi.Round{
			OverallVotes: int32(len(votes)),
			Tallies:      tallies,
		})

		if tallies[0].Count > int32(len(votes)/2) {
			r.Winner = tallies[0].Choice
			break
		}

		least := tallies[len(tallies)-1].Choice

		// remove least popular
		for _, v := range votes {
			for i := range v.choices {
				if v.choices[i] == least {
					v.choices = append(v.choices[:i], v.choices[i+1:]...)
					break
				}
			}
		}
	}

	return &r
}

type sortableTallies []*rvapi.Tally

func (c sortableTallies) Len() int {
	return len(c)
}

func (c sortableTallies) Less(i, j int) bool {
	return c[i].Count < c[j].Count || (c[i].Count == c[j].Count && c[i].Choice < c[j].Choice)
}

func (c sortableTallies) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

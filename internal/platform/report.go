package platform

import (
	"context"
	"database/sql"
	"fmt"
	"math"
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
	return &rvapi.ReportResponse{Report: calculateReport(votes, 1)}, nil
}

func calculateReport(vs []*models.Vote, numWinners int) *rvapi.Report {
	type vote struct {
		name        string
		choices     []string
		coefficient float64
	}

	votes := make([]*vote, 0, len(vs))
	for _, v := range vs {
		if len(v.Choices) > 0 {
			votes = append(votes, &vote{v.Name, v.Choices, 1})
		}
	}

	quota := math.Ceil((float64(len(votes)) + 1) / (float64(numWinners) + 1))

	r := &rvapi.Report{Winners: make([]string, 0, numWinners)}
	for len(r.Winners) < numWinners && len(votes) > 0 {
		counted := make(map[string]float64)
		for _, v := range votes {
			counted[v.choices[0]] += v.coefficient
		}

		round := rvapi.Round{Tallies: make([]*rvapi.Tally, 0, len(counted))}
		for k, v := range counted {
			round.Tallies = append(round.Tallies, &rvapi.Tally{Count: v, Choice: k})
			round.OverallVotes += int32(math.Round(v))
		}
		r.Rounds = append(r.Rounds, &round)

		sort.Sort(sort.Reverse(sortableTallies(round.Tallies)))

		var didElect bool
		for i := 0; i < len(round.Tallies) && round.Tallies[i].Count >= quota && len(r.Winners) < numWinners; i++ {
			t := round.Tallies[i]
			t.Outcome = rvapi.Tally_ELECTED

			didElect = true
			r.Winners = append(r.Winners, t.Choice)

			surplus := t.Count - quota
			transferValue := surplus / t.Count

			var cur int
			for _, v := range votes {
				idx := -1
				for i := range v.choices {
					if t.Choice == v.choices[i] {
						idx = i
						break
					}
				}
				if idx >= 0 {
					v.choices = append(v.choices[:idx], v.choices[idx+1:]...)
				}
				if len(v.choices) == 0 {
					continue
				}
				votes[cur] = v
				cur++
				if idx == 0 {
					v.coefficient *= transferValue
				}
			}
			votes = votes[:cur]
		}
		if didElect {
			continue
		}

		// fewer candidates than remaining positions -- all elected
		if len(round.Tallies)+len(r.Winners) <= numWinners {
			for _, t := range round.Tallies {
				t.Outcome = rvapi.Tally_ELECTED
				r.Winners = append(r.Winners, t.Choice)
			}
			break
		}

		least := round.Tallies[len(round.Tallies)-1]
		least.Outcome = rvapi.Tally_ELIMINATED

		// remove least popular
		var cur int
		for _, v := range votes {
			for i := range v.choices {
				if v.choices[i] == least.Choice {
					v.choices = append(v.choices[:i], v.choices[i+1:]...)
					break
				}
			}
			if len(v.choices) > 0 {
				votes[cur] = v
				cur++
			}
		}
		votes = votes[:cur]
	}

	return r
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

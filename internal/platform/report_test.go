package platform

import (
	"testing"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/stretchr/testify/require"

	"github.com/jwilner/rv/internal/models"
)

func Test_calculateReport(t *testing.T) {
	type (
		v  = []*models.Vote
		s  = []*rvapi.Round
		r  = rvapi.Report
		sa = []string
	)

	round := func(overallVotes int32, vals ...interface{}) *rvapi.Round {
		r := rvapi.Round{OverallVotes: overallVotes}
		for i := 0; i < len(vals); i += 3 {
			r.Tallies = append(r.Tallies, &rvapi.Tally{
				Choice:  vals[i].(string),
				Count:   vals[i+1].(float64),
				Outcome: vals[i+2].(rvapi.Tally_Outcome),
			})
		}
		return &r
	}

	tests := []struct {
		name string
		vs   v
		want *rvapi.Report
	}{
		{"empty", nil, &r{}},
		{"nonNil", make(v, 0), &r{}},
		{
			name: "one vote",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}},
			want: &r{Rounds: s{round(1, "bob", 1., rvapi.Tally_ELECTED)}, Winners: sa{"bob"}}},
		{
			name: "agreeing vote",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bob"}}},
			want: &r{Rounds: s{round(2, "bob", 2., rvapi.Tally_ELECTED)}, Winners: sa{"bob"}},
		},
		{
			name: "disagreeing vote eliminates lexicographically least choice",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}},
			want: &r{
				Rounds: s{
					round(2, "bob", 1., rvapi.Tally_UNKNOWN, "bill", 1., rvapi.Tally_ELIMINATED),
					round(1, "bob", 1., rvapi.Tally_ELECTED),
				},
				Winners: []string{"bob"},
			},
		},
		{
			name: "case insensitive winner",
			vs: v{
				{Name: "jack", Choices: sa{"bob"}},
				{Name: "jill", Choices: sa{"bill"}},
				{Name: "jon", Choices: sa{"Barbs", "bill"}},
			},
			want: &r{
				Rounds: s{
					round(
						3,
						"bob", 1., rvapi.Tally_UNKNOWN,
						"bill", 1., rvapi.Tally_UNKNOWN,
						"Barbs", 1., rvapi.Tally_ELIMINATED,
					),
					round(3, "bill", 2., rvapi.Tally_ELECTED, "bob", 1., rvapi.Tally_UNKNOWN),
				},
				Winners: []string{"bill"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.Winners == nil {
				tt.want.Winners = make([]string, 0)
			}
			require.Equal(t, tt.want, calculateReport(tt.vs, 1))
		})
	}
}

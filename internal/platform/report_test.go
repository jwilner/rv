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
		for i := 0; i < len(vals); i += 2 {
			r.Tallies = append(r.Tallies, &rvapi.Tally{
				Choice: vals[i].(string),
				Count:  int32(vals[i+1].(int)),
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
			want: &r{Winner: "bob", Rounds: s{round(1, "bob", 1)}}},
		{
			name: "agreeing vote",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bob"}}},
			want: &r{Winner: "bob", Rounds: s{round(2, "bob", 2)}},
		},
		{
			name: "disagreeing vote eliminates lexicographically least choice",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}},
			want: &r{
				Winner: "bob",
				Rounds: s{
					round(2, "bob", 1, "bill", 1),
					round(1, "bob", 1),
				},
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
				Winner: "bill",
				Rounds: s{
					round(3, "bob", 1, "bill", 1, "Barbs", 1),
					round(3, "bill", 2, "bob", 1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, calculateReport(tt.vs))
		})
	}
}

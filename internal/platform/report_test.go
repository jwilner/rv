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
		rv = []*rvapi.RemainingVote
		vc = map[string]int32
	)

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
			want: &r{Winner: "bob", Rounds: s{{Remaining: rv{{Name: "jack", Choices: sa{"bob"}}}, Counted: vc{"bob": 1}}}}},
		{
			name: "agreeing vote",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bob"}}},
			want: &r{Winner: "bob", Rounds: s{{Remaining: rv{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bob"}}}, Counted: vc{"bob": 2}}}},
		},
		{
			name: "disagreeing vote eliminates lexicographically least choice",
			vs:   v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}},
			want: &r{
				Winner: "bob",
				Rounds: s{
					{Remaining: rv{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}}, Counted: vc{"bob": 1, "bill": 1}},
					{Eliminated: []string{"jill"}, Remaining: rv{{Name: "jack", Choices: sa{"bob"}}}, Counted: vc{"bob": 1}},
				},
			},
		},
		{
			name: "case insensitive winner",
			vs: v{
				{Name: "jack", Choices: sa{"bob"}},
				{Name: "jill", Choices: sa{"bill"}},
				{Name: "jon", Choices: sa{"Barbs", "Bill"}},
			},
			want: &r{
				Winner: "bill",
				Rounds: s{
					{
						Remaining: rv{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}, {Name: "jon", Choices: sa{"Barbs", "Bill"}}},
						Counted:   vc{"bob": 1, "bill": 1, "barbs": 1}},
					{
						Remaining: rv{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}, {Name: "jon", Choices: sa{"Bill"}}},
						Counted:   vc{"bob": 1, "bill": 2},
					},
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

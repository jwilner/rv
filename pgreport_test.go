package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jwilner/rv/models"
)

func Test_calculateReport(t *testing.T) {
	type (
		v  = []*models.Vote
		s  = []*step
		r  = report
		sa = []string
		rv = []*remainingVote
		vc = map[string]int
	)

	tests := []struct {
		name string
		vs   v
		want *report
	}{
		{"empty", nil, &r{}},
		{"nonNil", make(v, 0), &r{}},
		{
			"one vote",
			v{{Name: "jack", Choices: sa{"bob"}}},
			&r{"bob", s{{1, nil, rv{{"jack", sa{"bob"}}}, vc{"bob": 1}}}}},
		{
			"agreeing vote",
			v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bob"}}},
			&r{"bob", s{{1, nil, rv{{"jack", sa{"bob"}}, {"jill", sa{"bob"}}}, vc{"bob": 2}}}},
		},
		{
			"disagreeing vote eliminates lexicographically least choice",
			v{{Name: "jack", Choices: sa{"bob"}}, {Name: "jill", Choices: sa{"bill"}}},
			&r{
				"bob",
				s{
					{1, nil, rv{{"jack", sa{"bob"}}, {"jill", sa{"bill"}}}, vc{"bob": 1, "bill": 1}},
					{2, []string{"jill"}, rv{{"jack", sa{"bob"}}}, vc{"bob": 1}},
				},
			},
		},
		{
			"case insensitive winner",
			v{
				{Name: "jack", Choices: sa{"bob"}},
				{Name: "jill", Choices: sa{"bill"}},
				{Name: "jon", Choices: sa{"Barbs", "Bill"}},
			},
			&r{
				"bill",
				s{
					{
						1,
						nil,
						rv{{"jack", sa{"bob"}}, {"jill", sa{"bill"}}, {"jon", sa{"Barbs", "Bill"}}},
						vc{"bob": 1, "bill": 1, "barbs": 1}},
					{
						2,
						nil,
						rv{{"jack", sa{"bob"}}, {"jill", sa{"bill"}}, {"jon", sa{"Bill"}}},
						vc{"bob": 1, "bill": 2},
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

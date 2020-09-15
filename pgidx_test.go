package main

import (
	"reflect"
	"testing"
)

func Test_parseChoices(t *testing.T) {

	tests := []struct {
		name string
		s    string
		want normalizedSlice
	}{
		{
			"obno",
			`separated by spaces
separated_by_newline
comma, 
"quotes"`,
			normalize([]string{"separated by spaces", "separated_by_newline", "comma,", "\"quotes\""}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseChoices(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseChoices() = %q, want %q", got.raw(), tt.want.raw())
			}
		})
	}
}

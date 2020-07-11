package main

import "testing"

func Test_parseBallotKey(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty", "/v/", ""},
		{"single segment", "/v/abcdefgh", "abcdefgh"},
		{"single segment", "/v/abcdefgh/ab", ""},
		{"single segment", "/v/abcdefgh", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

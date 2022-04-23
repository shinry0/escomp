package escomp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		testname string
		names    []string
		vals     []string
		want     map[string]string
		wantErr  bool
	}{
		{
			testname: "empty",
			names:    []string{},
			vals:     []string{},
			want:     map[string]string{},
		},
		{
			testname: "normal params",
			names:    []string{"message", "speaker"},
			vals:     []string{"hello", "me"},
			want:     map[string]string{"message": "hello", "speaker": "me"},
		},
		{
			testname: "error: length not equal",
			names:    []string{"hoge"},
			vals:     []string{},
			wantErr:  true,
		},
		{
			testname: "error: length not equal",
			names:    []string{},
			vals:     []string{"hoge"},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			got, err := ParseParams(test.names, test.vals)
			if test.wantErr {
				if err == nil {
					t.Errorf("expected error didn't occur: got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error occurred: %s", err)
				}
				if diff := cmp.Diff(test.want, got); diff != "" {
					t.Errorf("parse results mismatch:\n%s", diff)
				}
			}
		})
	}
}

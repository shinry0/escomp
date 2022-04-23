package escomp

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockESSearcher struct{}

func (m *mockESSearcher) Search(index, query string, size int) (esResult, error) {
	q := `{ "query": { "match": { "title": "TITLE" } } }`
	qn := `{ "query": { "match": { "title": "NO HITS" } } }`

	es := map[string]map[string]esResult{
		"INDEX": {
			q: {
				Hits: ESHits{
					Hits: []ESDoc{
						{"1", 1.0, map[string]any{"title": "TITLE"}},
					},
					Total: ESHitCount{"eq", 10},
				},
			},
			qn: {
				Hits: ESHits{
					Hits:  []ESDoc{},
					Total: ESHitCount{"eq", 0},
				},
			},
		},
	}

	idx, ok := es[index]
	if !ok {
		return esResult{}, errors.New("index not found")
	}
	return idx[query], nil
}

func TestSearch(t *testing.T) {
	tests := []struct {
		testname   string
		searchCase SearchCase
		params     map[string]string
		size       int
		want       SearcherResult
		wantErr    bool
	}{
		{
			testname: "hit",
			searchCase: SearchCase{
				Name:  "alpha",
				Index: "INDEX",
				Query: `{ "query": { "match": { "title": "{{keyword}}" } } }`,
			},
			params: map[string]string{"keyword": "TITLE"},
			size:   8,
			want: SearcherResult{
				Name: "alpha",
				ESHits: ESHits{
					Hits: []ESDoc{
						{"1", 1.0, map[string]any{"title": "TITLE"}},
					},
					Total: ESHitCount{"eq", 10},
				},
			},
		},
		{
			testname: "no hits",
			searchCase: SearchCase{
				Name:  "beta",
				Index: "INDEX",
				Query: `{ "query": { "match": { "title": "{{keyword}}" } } }`,
			},
			params: map[string]string{"keyword": "NO HITS"},
			size:   8,
			want: SearcherResult{
				Name: "beta",
				ESHits: ESHits{
					Hits:  []ESDoc{},
					Total: ESHitCount{"eq", 0},
				},
			},
		},
		{
			testname: "index not found",
			searchCase: SearchCase{
				Name:  "beta",
				Index: "NOTFOUND",
				Query: `{ "query": { "match": { "title": "{{keyword}}" } } }`,
			},
			params:  map[string]string{"keyword": "TITLE"},
			size:    8,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			mock := new(mockESSearcher)
			got, err := NewSearcher(test.searchCase, mock).Search(test.params, test.size)
			if test.wantErr {
				if err == nil {
					t.Errorf("expected error didn't occur: got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error occurred: %s", err)
				}
				if diff := cmp.Diff(test.want, *got); diff != "" {
					t.Errorf("search result mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestEmbedParams(t *testing.T) {
	query := `{ "query": { "match": { "title": "{{keyword}}" } } }`
	params := map[string]string{"keyword": "TITLE"}

	want := `{ "query": { "match": { "title": "TITLE" } } }`
	got, _ := emdedParams(query, params)
	if want != got {
		t.Errorf("embeded query is not equal: want %s, got %s", want, got)
	}
}

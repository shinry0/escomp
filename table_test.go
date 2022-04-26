package escomp

import (
	"testing"

	"github.com/Code-Hex/dd"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/olekukonko/tablewriter"
)

func TestConvert(t *testing.T) {
	conv := TableConverter{
		Fields: []string{"title", "number"},
	}
	results := []SearcherResult{
		{
			Name: "alpha",
			ESHits: ESHits{
				Hits: []ESDoc{
					{
						ID:    "0",
						Score: 3.0,
						Source: map[string]interface{}{
							"title":  "title_0",
							"number": 0,
						},
					},
					{
						ID:    "1",
						Score: 1.0,
						Source: map[string]interface{}{
							"title":  "title_1",
							"number": 1,
						},
					},
					{
						ID:    "2",
						Score: 1.0,
						Source: map[string]interface{}{
							"title":  "title_2",
							"number": 2,
						},
					},
				},
				Total: ESHitCount{
					Relation: "eq",
					Value:    10,
				},
			},
		},
		{
			Name: "beta",
			ESHits: ESHits{
				Hits: []ESDoc{
					{
						ID:    "3",
						Score: 3.0,
						Source: map[string]interface{}{
							"title":  "title_1",
							"number": 1,
						},
					},
					{
						ID:    "4",
						Score: 1.0,
						Source: map[string]interface{}{
							"title":  "title_0",
							"number": 0,
						},
					},
				},
				Total: ESHitCount{
					Relation: "eq",
					Value:    2,
				},
			},
		},
	}

	want := Table{
		Header: []string{"", "alpha", "beta"},
		Content: [][]string{
			{"#1", "title_0 / 0", "title_1 / 1"},
			{"#2", "title_1 / 1", "title_0 / 0"},
			{"#3", "title_2 / 2", ""},
		},
		Footer: []string{"COUNT", "10", "2"},
	}
	got := conv.Convert(results)

	if diff := cmp.Diff(want, *got, cmpopts.IgnoreUnexported(Table{})); diff != "" {
		t.Errorf("table mismatch (-want +got):\n%s", diff)
	}
}

func TestColor(t *testing.T) {
	table := Table{
		Header: []string{"", "alpha", "beta"},
		Content: [][]string{
			{"#1", "title_0", "title_1"},
			{"#2", "title_1", "title_0"},
			{"#3", "title_2", ""},
		},
		Footer: []string{"COUNT", "10", "2"},
	}
	got := table.Color()

	// check for color match
	wantColoredIndices := [][][2]int{
		{{0, 1}, {1, 2}}, // colors[0][1] == colors[1][2]
		{{1, 1}, {0, 2}}, // colors[1][1] == colors[0][2]
	}
	for _, wantIndices := range wantColoredIndices {
		if len(wantIndices) == 0 {
			continue
		}
		idx := wantIndices[0]
		i, j := idx[0], idx[1]
		if i >= len(got.colors) || j >= len(got.colors[i]) {
			t.Errorf("cannot access to colors[%d][%d], got %v", i, j, got.colors)
			continue
		}

		if len(got.colors[i][j]) == 0 {
			t.Errorf("%v of table want to be colored", wantIndices)
		}
		if !haveColors(got.colors, wantIndices[1:], got.colors[i][j]) {
			t.Errorf("%v of table want to have same colors, got\n%s", wantIndices, dd.Dump(got.colors))
		}
	}

	wantNotColoredIndices := [][][2]int{
		{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
	}
	for _, wantIndices := range wantNotColoredIndices {
		wantColors := tablewriter.Color() // no color
		if !haveColors(got.colors, wantIndices, wantColors) {
			t.Errorf("%v of table want to have same colors, got %v", wantIndices, got.colors)
		}
	}
}

func haveColors(colorTable [][]tablewriter.Colors, indices [][2]int, want tablewriter.Colors) bool {
	for _, idx := range indices {
		i, j := idx[0], idx[1]
		if i >= len(colorTable) || j >= len(colorTable[i]) {
			return false
		}
		if !cmp.Equal(colorTable[i][j], want, cmp.Options{
			cmpopts.SortSlices(func(x, y int) bool {
				return x < y
			}),
		}) {
			return false
		}
	}
	return true
}

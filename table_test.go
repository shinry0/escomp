package escomp

import (
	"testing"

	"github.com/Code-Hex/dd"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/olekukonko/tablewriter"
)

var results = []SearcherResult{
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
					ID:    "1",
					Score: 3.0,
					Source: map[string]interface{}{
						"title":  "title_1",
						"number": 1,
					},
				},
				{
					ID:    "0",
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

func TestConvert(t *testing.T) {
	conv := TableConverter{
		Fields: []string{"title", "number"},
		Color:  false,
	}

	want := Table{
		Header: []string{"", "alpha", "beta"},
		Content: [][]string{
			{"#1", "title_0 / 0", "title_1 / 1"},
			{"#2", "title_1 / 1", "title_0 / 0"},
			{"#3", "title_2 / 2", ""},
		},
		Colors: [][]tablewriter.Colors{},
		Footer: []string{"COUNT", "10", "2"},
	}
	got := conv.Convert(results)

	if diff := cmp.Diff(want, *got); diff != "" {
		t.Errorf("table mismatch (-want +got):\n%s", diff)
	}
}

func TestConvert_Color(t *testing.T) {
	conv := TableConverter{
		Fields: []string{"title", "number"},
		Color:  true,
	}

	want := Table{
		Header: []string{"", "alpha", "beta"},
		Content: [][]string{
			{"#1", "title_0 / 0", "title_1 / 1"},
			{"#2", "title_1 / 1", "title_0 / 0"},
			{"#3", "title_2 / 2", ""},
		},
		Colors: [][]tablewriter.Colors{}, // not need to know what colors
		Footer: []string{"COUNT", "10", "2"},
	}
	got := conv.Convert(results)

	if diff := cmp.Diff(want, *got, cmp.Options{
		cmpopts.IgnoreFields(Table{}, "Colors"),
	}); diff != "" {
		t.Errorf("table mismatch (-want +got):\n%s", diff)
	}

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
		if i >= len(got.Colors) || j >= len(got.Colors[i]) {
			t.Errorf("cannot access to (%d, %d), got %v", i, j, got.Colors)
			continue
		}

		wantColors := got.Colors[i][j]
		if !haveColors(wantColors, got.Colors, wantIndices[1:]) {
			t.Errorf("%v of table want to have same colors, got\n%s", wantIndices, dd.Dump(got.Colors))
		}
	}
	wantNotColoredIndices := [][][2]int{
		{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
	}
	for _, wantIndices := range wantNotColoredIndices {
		wantColors := tablewriter.Color() // no color
		if !haveColors(wantColors, got.Colors, wantIndices) {
			t.Errorf("%v of table want to have same colors, got %v", wantIndices, got.Colors)
		}
	}
}

func haveColors(want tablewriter.Colors, colorTable [][]tablewriter.Colors, indices [][2]int) bool {
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

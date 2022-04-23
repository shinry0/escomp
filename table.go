package escomp

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type TableConverter struct {
	Fields []string
	Color  bool
}

func NewTableConverter(fields []string, color bool) *TableConverter {
	return &TableConverter{
		Fields: fields,
		Color:  color,
	}
}

func (c *TableConverter) Convert(results []SearcherResult) *Table {
	t := &Table{
		Header:  makeHeader(results),
		Content: makeContent(results, c.Fields),
		Colors:  [][]tablewriter.Colors{},
		Footer:  makeFooter(results),
	}
	if c.Color {
		t.Colors = makeColors(results)
	}
	return t
}

type Table struct {
	Header  []string
	Content [][]string
	Colors  [][]tablewriter.Colors
	Footer  []string
}

func (t *Table) Render(w io.Writer) {
	writer := tablewriter.NewWriter(w)
	writer.SetRowLine(true)
	writer.SetHeader(t.Header)
	writer.SetFooter(t.Footer)
	writer.SetReflowDuringAutoWrap(false)

	if len(t.Content) == len(t.Colors) {
		for i := range t.Content {
			writer.Rich(t.Content[i], t.Colors[i])
		}
	} else {
		for i := range t.Content {
			writer.Append(t.Content[i])
		}
	}
	writer.Render()
}

func makeHeader(results []SearcherResult) []string {
	header := make([]string, len(results)+1)
	header[0] = "" // header column
	for i, res := range results {
		header[i+1] = res.Name
	}
	return header
}

func makeContent(results []SearcherResult, fields []string) [][]string {
	m, n := maxHitLen(results), len(results)
	content := make([][]string, m)
	for i := 0; i < m; i++ {
		content[i] = make([]string, n+1)
		content[i][0] = fmt.Sprintf("#%d", i+1) // header column
		for j, res := range results {
			if h := res.Hits; len(h) > i {
				content[i][j+1] = makeContentString(h[i], fields)
			}
		}
	}
	return content
}

func makeContentString(doc ESDoc, fields []string) string {
	var strs []string
	for _, f := range fields {
		strs = append(strs, fmt.Sprintf("%v", doc.Source[f]))
	}
	return strings.Join(strs, " / ")
}

func makeColors(results []SearcherResult) [][]tablewriter.Colors {
	// scan search results and mark matching docs
	scanned := make(map[string]bool)
	colorMap, colorID := make(map[string]int), 0
	for _, res := range results {
		for _, doc := range res.Hits {
			if scanned[doc.ID] {
				if _, ok := colorMap[doc.ID]; !ok {
					colorMap[doc.ID] = colorID
					colorID++
				}
				continue
			}
			scanned[doc.ID] = true
		}
	}

	// coloring
	m, n := maxHitLen(results), len(results)
	colors := make([][]tablewriter.Colors, m)
	for i := 0; i < m; i++ {
		colors[i] = make([]tablewriter.Colors, n+1)
		colors[i][0] = tablewriter.Color() // header column
		for j, res := range results {
			if h := res.Hits; len(h) > i {
				doc := h[i]
				if colorID, ok := colorMap[doc.ID]; ok {
					colors[i][j+1] = getColors(colorID)
					continue
				}
			}
		}
	}
	return colors
}

func makeFooter(results []SearcherResult) []string {
	footer := make([]string, len(results)+1)
	footer[0] = "COUNT" // header column
	for i, res := range results {
		switch res.Total.Relation {
		case "eq":
			footer[i+1] = fmt.Sprintf("%d", res.Total.Value)
		case "gte":
			footer[i+1] = fmt.Sprintf("%d+", res.Total.Value)
		}
	}
	return footer
}

func maxHitLen(results []SearcherResult) int {
	var max int
	for _, res := range results {
		if l := len(res.Hits); l > max {
			max = l
		}
	}
	return max
}

func getColors(n int) tablewriter.Colors {
	// ref. https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	// foreground color and background color should be different, so skip 9*k
	m := n + 1 + (n / 8)
	fg, bg := m%8+30, m/8+40
	return tablewriter.Color(4, fg, bg)
}

package escomp

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type TableConverter struct {
	Fields []string
}

func NewTableConverter(fields []string) *TableConverter {
	return &TableConverter{
		Fields: fields,
	}
}

func (c *TableConverter) Convert(results []SearcherResult) *Table {
	t := &Table{
		Header:  c.makeHeader(results),
		Content: c.makeContent(results),
		Footer:  c.makeFooter(results),
	}
	return t
}

func (c *TableConverter) makeHeader(results []SearcherResult) []string {
	header := make([]string, len(results)+1)
	header[0] = "" // header column
	for i, res := range results {
		header[i+1] = res.Name
	}
	return header
}

func (c *TableConverter) makeContent(results []SearcherResult) [][]string {
	m, n := maxHitLen(results), len(results)
	content := make([][]string, m)
	for i := 0; i < m; i++ {
		content[i] = make([]string, n+1)
		content[i][0] = fmt.Sprintf("#%d", i+1) // header column
		for j, res := range results {
			if h := res.Hits; len(h) > i {
				content[i][j+1] = c.makeContentString(h[i])
			}
		}
	}
	return content
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

func (c *TableConverter) makeContentString(doc ESDoc) string {
	var strs []string
	for _, f := range c.Fields {
		strs = append(strs, fmt.Sprintf("%v", doc.Source[f]))
	}
	return strings.Join(strs, " / ")
}

func (c *TableConverter) makeFooter(results []SearcherResult) []string {
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

type Table struct {
	Header  []string
	Content [][]string
	Footer  []string

	colors [][]tablewriter.Colors
}

func (t *Table) Color() *Table {
	// scan table and choose colors
	scanned := make(map[string]bool)
	colorMap := make(map[string]int)
	colorID := 0
	for _, row := range t.Content {
		for _, text := range row {
			if text == "" {
				continue
			}
			if scanned[text] {
				if _, ok := colorMap[text]; !ok {
					colorMap[text] = colorID
					colorID++
				}
				continue
			}
			scanned[text] = true
		}
	}

	// set color
	colors := make([][]tablewriter.Colors, len(t.Content))
	for i, row := range t.Content {
		colors[i] = make([]tablewriter.Colors, len(row))
		for j, text := range row {
			if colorID, ok := colorMap[text]; ok {
				colors[i][j] = getColors(colorID)
			}
		}
	}
	t.colors = colors

	return t
}

func getColors(n int) tablewriter.Colors {
	// ref. https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	// foreground color and background color should be different, so skip 9*k
	m := n + 1 + (n / 8)
	fg, bg := m%8+30, m/8+40
	return tablewriter.Color(4, fg, bg)
}

func (t *Table) Render(w io.Writer) {
	writer := tablewriter.NewWriter(w)
	writer.SetRowLine(true)
	writer.SetHeader(t.Header)
	writer.SetFooter(t.Footer)
	writer.SetReflowDuringAutoWrap(false)

	if len(t.Content) == len(t.colors) {
		for i := range t.Content {
			writer.Rich(t.Content[i], t.colors[i])
		}
	} else {
		for i := range t.Content {
			writer.Append(t.Content[i])
		}
	}
	writer.Render()
}

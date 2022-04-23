package escomp

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadDefFile(t *testing.T) {
	tests := []struct {
		def  string
		want Definition
	}{
		{
			`
params = [ "keyword" ]
fields = [ "text_entry" ]

[esconfig.common]
url = "http://localhost:9200"
username = ""
password = ""

[[search]]
name = "alpha"
es = "common"
index = "shakespeare"
query = "alpha query"

[[search]]
name = "beta"
es = "common"
index = "shakespeare2"
query = "beta query"`,
			Definition{
				Params: []string{"keyword"},
				Fields: []string{"text_entry"},
				ESConfigs: map[string]ESConfig{
					"common": {URL: "http://localhost:9200"},
				},
				SearchCases: []SearchCase{
					{Name: "alpha", ES: "common", Index: "shakespeare", Query: "alpha query"},
					{Name: "beta", ES: "common", Index: "shakespeare2", Query: "beta query"},
				},
			},
		},
	}
	for _, test := range tests {
		r := strings.NewReader(test.def)
		got, _ := LoadDefFile(r)
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("definition mismatch (-want +got):\n%s", diff)
		}
	}
}

package escomp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cbroglie/mustache"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

type Searcher struct {
	SearchCase
	searcher ESSearcher
}

type SearcherResult struct {
	Name string
	ESHits
}

func NewSearcher(sc SearchCase, cli ESSearcher) *Searcher {
	return &Searcher{sc, cli}
}

func (s *Searcher) Search(params map[string]string, size int) (*SearcherResult, error) {
	res := &SearcherResult{Name: s.Name}

	query, err := emdedParams(s.Query, params)
	if err != nil {
		return res, err
	}

	esr, err := s.searcher.Search(s.Index, query, size)
	if err != nil {
		return res, fmt.Errorf("failed to search: %s", err)
	}
	res.ESHits = esr.Hits

	return res, nil
}

func emdedParams(query string, params map[string]string) (string, error) {
	q, err := mustache.Render(query, params)
	if err != nil {
		return "", fmt.Errorf("failed to embed params to query template: %s", err)
	}
	return q, nil
}

type ESSearcher interface {
	Search(index, query string, size int) (esResult, error)
}

type ES7Client struct {
	client *es7.Client
}

func NewES7Client(conf ESConfig) (*ES7Client, error) {
	c := new(ES7Client)
	client, err := es7.NewClient(
		es7.Config{
			Addresses: []string{conf.URL},
			Username:  conf.Username,
			Password:  conf.Password,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ES7 client: %s", err)
	}
	c.client = client
	return c, nil
}

func (c *ES7Client) Search(index, query string, size int) (esResult, error) {
	res, err := c.client.Search(
		c.client.Search.WithContext(context.Background()),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(strings.NewReader(query)),
		c.client.Search.WithSize(size),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return esResult{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e esErrResult
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return esResult{}, err
		}
		return esResult{}, fmt.Errorf(
			"[%s] %s: %s",
			res.Status(), e.Error.Type, e.Error.Reason,
		)
	}

	var esr esResult
	if err := json.NewDecoder(res.Body).Decode(&esr); err != nil {
		return esr, err
	}
	return esr, nil
}

type esResult struct {
	Hits ESHits `json:"hits"`
}

type esErrResult struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
}

type ESHits struct {
	Hits  []ESDoc    `json:"hits"`
	Total ESHitCount `json:"total"`
}

type ESHitCount struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}

type ESDoc struct {
	ID     string         `json:"_id"`
	Score  float64        `json:"_score"`
	Source map[string]any `json:"_source"`
}

package escomp

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type Definition struct {
	Params      []string            `toml:"params"`
	Fields      []string            `toml:"fields"`
	ESConfigs   map[string]ESConfig `toml:"esconfig"`
	SearchCases []SearchCase        `toml:"search"`
}

type SearchCase struct {
	Name  string `toml:"name"`
	ES    string `toml:"es"`
	Index string `toml:"index"`
	Query string `toml:"query"`
}

type ESConfig struct {
	URL      string
	Username string
	Password string
}

func LoadDefFile(r io.Reader) (Definition, error) {
	d := Definition{}
	_, err := toml.NewDecoder(r).Decode(&d)
	if err != nil {
		return d, fmt.Errorf("fail to load the definition file: %w", err)
	}
	return d, nil
}

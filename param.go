package escomp

import (
	"fmt"
)

func ParseParams(names, vals []string) (map[string]string, error) {
	if len(vals) != len(names) {
		return nil, fmt.Errorf("%d param[s] be defined but %d given", len(names), len(vals))
	}

	params := make(map[string]string, len(vals))
	for i := 0; i < len(vals); i++ {
		params[names[i]] = vals[i]
	}
	return params, nil
}

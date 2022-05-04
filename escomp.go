package escomp

import (
	"fmt"
	"os"
)

type Option struct {
	Filename    string
	ParamValues []string
	Size        int
	EnableColor bool
}

func Run(opt Option) error {
	f, err := os.Open(opt.Filename)
	if err != nil {
		return fmt.Errorf("failed to open a file: %w", err)
	}
	defer f.Close()

	def, err := LoadDefFile(f)
	if err != nil {
		return err
	}

	params, err := ParseParams(def.Params, opt.ParamValues)
	if err != nil {
		return err
	}

	results := make([]SearcherResult, len(def.SearchCases))
	for i, sc := range def.SearchCases {
		conf, ok := def.ESConfigs[sc.ES]
		if !ok {
			return fmt.Errorf(`Elasticsearch config "%s" is not exist`, sc.ES)
		}
		cli, err := NewES7Client(conf)
		if err != nil {
			return err
		}
		res, err := NewSearcher(sc, cli).Search(params, def.Fields, opt.Size)
		if err != nil {
			return err
		}
		results[i] = *res
	}

	table := NewTableConverter(def.Fields).Convert(results)
	if opt.EnableColor {
		table.Color()
	}
	table.Render(os.Stdout)

	return nil
}

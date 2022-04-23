package main

import (
	"os"

	"github.com/shinry0/escomp"
	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

var filename string
var size int
var color bool

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "escomp",
		Short: "Elasticsearch Result Comparing Tool",
		Long:  "A CLI tool to compare several search results for Elasticsearch",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := escomp.Run(escomp.Option{
				Filename:    filename,
				ParamValues: args,
				Size:        size,
				EnableColor: color,
			})
			return err
		},
	}
}

func main() {
	rootCmd.Flags().StringVarP(&filename, "file", "f", "", "definition file (required)")
	rootCmd.Flags().IntVarP(&size, "size", "n", 8, "max number of hits")
	rootCmd.Flags().BoolVar(&color, "color", false, "enable coloring")
	rootCmd.MarkFlagRequired("file")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

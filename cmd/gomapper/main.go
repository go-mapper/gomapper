package main

import (
	"fmt"
	"os"

	"github.com/go-mapper/gomapper/cmd/gomapper/generate"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(generate.NewCommand())

	err := root.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "execute: %s\n", err)
		os.Exit(2)
	}
}

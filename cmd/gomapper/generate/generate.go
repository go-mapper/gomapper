package generate

import (
	"fmt"
	"os"

	"github.com/go-mapper/gomapper"
	"github.com/go-mapper/gomapper/internal/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var (
		config string
		output string
	)
	cmd := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loader.Load(config, []string{})
			if err != nil {
				return err
			}

			out, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("out: %w", err)
			}
			defer out.Close()

			gomapper.Output(out, cfg)

			return nil
		},
	}

	cmd.Flags().StringVar(&config, "config", "", "path to the config .go file")
	cmd.Flags().StringVar(&output, "output", "mappings_gen.go", "output file")

	return cmd
}

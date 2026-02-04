package cmd

import (
	"fmt"

	"yj/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("services.yaml")
		if err != nil {
			return err
		}

		for name, yj := range cfg.Services {
			fmt.Printf("- %s (%s)\n", name, yj.Path)
			for script, cmd := range yj.Scripts {
				fmt.Printf("  • %s: %s\n", script, cmd)
			}
			
			if len(yj.Scripts) == 0 {
				fmt.Println("  (no scripts)")
			}
		}

		return nil
	},
}

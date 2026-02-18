package cmd

import (
	"fmt"
	"os"
	"yj/internal/logger"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "yj",
	Short: "Local service manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := logger.Setup(); err != nil {
			fmt.Printf("Warning: failed to setup logger: %v\n", err)
		}

		logger.Debug("CLI command started", "command", cmd.Use)
	},
}

func init() {
	baseCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := baseCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

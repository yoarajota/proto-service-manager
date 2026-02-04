package cmd

import (
	"fmt"
	"os"
	"yj/internal/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yj",
	Short: "Local service manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := logger.Setup(); err != nil {
			fmt.Printf("Warning: failed to setup logger: %v\n", err)
		}
		
		logger.Debug("CLI command started", "command", cmd.Use)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

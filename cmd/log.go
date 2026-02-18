package cmd

import (
	"fmt"
	"os"

	"yj/internal/logger"

	"github.com/spf13/cobra"
)

func init() {
	baseCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"logs"},
	Short:   "Show system logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if logger.LogFilePath == "" {
			if err := logger.Setup(); err != nil {
				return err
			}
		}

		data, err := os.ReadFile(logger.LogFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No logs found.")
				return nil
			}

			return fmt.Errorf("failed to read log file: %w", err)
		}

		fmt.Print(string(data))
		return nil
	},
}

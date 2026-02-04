package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"yj/internal/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop [service]",
	Short: "Stop a running service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		logger.Info("CLI: Stopping service", "service", name)

		home, _ := os.UserHomeDir()
		pidFile := filepath.Join(home, ".yj", "pids", name+".pid")

		data, err := os.ReadFile(pidFile)
		if err != nil {
			logger.Error("Service is not running (pid file not found)", "service", name, "error", err)
			return fmt.Errorf("service %s is not running", name)
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			logger.Error("Failed to read pid from file", "service", name, "error", err)
			return err
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			logger.Error("Failed to find process", "service", name, "pid", pid, "error", err)
			return err
		}

		if err := proc.Signal(syscall.SIGTERM); err != nil {
			logger.Error("Failed to send SIGTERM to process", "service", name, "pid", pid, "error", err)
			return err
		}

		_ = os.Remove(pidFile)

		logger.Info("CLI: Service stopped", "service", name)
		fmt.Printf("Stopped %s\n", name)
		return nil
	},
}

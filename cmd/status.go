package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Show service status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		home, _ := os.UserHomeDir()
		pidFile := filepath.Join(home, ".yj", "pids", name+".pid")

		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Printf("%s: stopped\n", name)
			return nil
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}

		exists, err := process.PidExists(int32(pid))
		if err != nil || !exists {
			fmt.Printf("%s: stopped (stale pid)\n", name)
			_ = os.Remove(pidFile)
			return nil
		}

		fmt.Printf("%s: running (pid %d)\n", name, pid)
		return nil
	},
}

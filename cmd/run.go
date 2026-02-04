package cmd

import (
	"fmt"

	"yj/internal/config"
	"yj/internal/logger"
	"yj/internal/process"
	"yj/internal/service"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [service] [script]",
	Short: "Run a service script",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		scriptName := args[1]
		logger.Info("CLI: Starting service", "service", name, "script", scriptName)

		cfg, err := config.Load("services.yaml")
		if err != nil {
			logger.Error("Failed to load config", "error", err)
			return err
		}

		raw, ok := cfg.Services[name]
		if !ok {
			logger.Error("Service not found", "service", name)
			return fmt.Errorf("service %s not found", name)
		}

		exec, ok := raw.Scripts[scriptName]
		if !ok {
			return fmt.Errorf("script %s not found for service %s", scriptName, name)
		}

		yj := service.Service{
			Name:  name,
			Start: exec,
			Cwd:   raw.Path,
		}

		proc, err := process.Start(yj)
		if err != nil {
			// process.Start already logs the error
			return err
		}

		logger.Info("CLI: Service started successfully", "service", name, "pid", proc.Pid)
		fmt.Printf("Started %s (pid %d)\n", name, proc.Pid)
		return nil
	},
}

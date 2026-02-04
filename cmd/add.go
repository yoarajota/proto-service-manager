package cmd

import (
	"yj/internal/config"
	"yj/internal/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("path", "", "Path to the service directory")
	addCmd.MarkFlagRequired("path")
	addCmd.Flags().String("script", "", "Name of the script to add (e.g., dev)")
	addCmd.Flags().String("exec", "", "Command to execute for the script")
}

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add or update a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path, _ := cmd.Flags().GetString("path")
		script, _ := cmd.Flags().GetString("script")
		exec, _ := cmd.Flags().GetString("exec")

		// Validate input
		if script != "" && exec == "" {
			return config.UpdateService("services.yaml", name, config.ServiceConfig{
				Path: path,
			})
		}
		
		// Load existing to append script or create new
		cfg, err := config.Load("services.yaml")
		var services map[string]config.ServiceConfig
		if err == nil && cfg != nil {
			services = cfg.Services
		} else {
			services = make(map[string]config.ServiceConfig)
		}

		service, ok := services[name]
		if !ok {
			service = config.ServiceConfig{
				Path:    path,
				Scripts: make(map[string]string),
			}
		}
		
		// Update path if provided
		if path != "" {
			service.Path = path
		}

		if script != "" && exec != "" {
			if service.Scripts == nil {
				service.Scripts = make(map[string]string)
			}
			service.Scripts[script] = exec
		}

		err = config.UpdateService("services.yaml", name, service)
		if err != nil {
			logger.Error("Failed to update service config", "service", name, "error", err)
			return err
		}

		logger.Info("CLI: Service config updated", "service", name)
		return nil
	},
}

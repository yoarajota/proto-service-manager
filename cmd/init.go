package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	baseCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("force", false, "Overwrite an existing global config if it exists")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a global services.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		if envPath := os.Getenv("YJ_CONFIG"); envPath != "" && !force {
			return fmt.Errorf("YJ_CONFIG is set (%s); refusing to create global config without --force", envPath)
		}

		localPath, err := filepath.Abs("services.yaml")
		if err != nil {
			return err
		}
		if _, err := os.Stat(localPath); err == nil && !force {
			return fmt.Errorf("local services.yaml exists at %s; refusing to create global config without --force", localPath)
		}

		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}

		configDir := filepath.Join(userConfigDir, "yj")
		globalPath := filepath.Join(configDir, "services.yaml")
		if _, err := os.Stat(globalPath); err == nil && !force {
			return fmt.Errorf("global config already exists at %s; use --force to overwrite", globalPath)
		}

		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}

		content := "# yj services configuration\nservices:\n  example:\n    path: /path/to/service\n    scripts:\n      dev: command\n"
		if err := os.WriteFile(globalPath, []byte(content), 0644); err != nil {
			return err
		}

		fmt.Printf("Created global config at %s\n", globalPath)
		return nil
	},
}

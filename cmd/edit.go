package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"yj/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	baseCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit services configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := "services.yaml", error(nil)
		resolvedPath, err := config.GetConfigPath()
		if err == nil {
			file = resolvedPath
		}

		// 1. Try VS Code first
		if _, err := exec.LookPath("code"); err == nil {
			c := exec.Command("code", file)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		}

		// 2. Fallback based on OS
		var command string
		var commandArgs []string

		switch runtime.GOOS {
		case "windows":
			command = "notepad.exe"
			commandArgs = []string{file}
		case "darwin":
			command = "open"
			commandArgs = []string{"-t", file}
		case "linux":
			if isWSL() {
				absPath, err := filepath.Abs(file)
				if err != nil {
					return err
				}
				winPath, err := toWindowsPath(absPath)
				if err != nil {
					return err
				}
				command = "notepad.exe"
				commandArgs = []string{winPath}
			} else {
				command = "xdg-open"
				commandArgs = []string{file}
			}
		default:
			return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
		}

		c := exec.Command(command, commandArgs...)
		// Set Dir to /mnt/c for WSL to avoid UNC warnings safely, though less critical for direct executable execution
		if isWSL() {
			c.Dir = "/mnt/c"
		}
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")

	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, "microsoft") || strings.Contains(content, "wsl")
}

func toWindowsPath(linuxPath string) (string, error) {
	cmd := exec.Command("wslpath", "-w", linuxPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

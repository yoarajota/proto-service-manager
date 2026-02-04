package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"yj/internal/logger"
	"yj/internal/service"
)

func pidPath(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yj", "pids", name+".pid")
}

func Start(s service.Service) (*os.Process, error) {
	logger.Info("Starting service process", "service", s.Name, "command", s.Start, "cwd", s.Cwd)

	cmd := exec.Command("sh", "-c", s.Start)
	cmd.Dir = s.Cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start process", "service", s.Name, "error", err)
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(pidPath(s.Name)), 0755); err != nil {
		logger.Error("Failed to create pid directory", "service", s.Name, "error", err)
		return nil, fmt.Errorf("failed to create pid directory: %w", err)
	}

	if err := os.WriteFile(pidPath(s.Name), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		logger.Error("Failed to write pid file", "service", s.Name, "error", err)
		return nil, fmt.Errorf("failed to write pid file: %w", err)
	}

	logger.Info("Service process started", "service", s.Name, "pid", cmd.Process.Pid)
	return cmd.Process, nil
}

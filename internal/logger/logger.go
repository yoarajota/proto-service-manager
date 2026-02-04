package logger

import (
	"fmt"

	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

var LogFilePath string

// Setup initializes the logger to write to a file in the temp directory
func Setup() error {
	tempDir := os.TempDir()
	LogFilePath = filepath.Join(tempDir, "yj-cli.log")

	file, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(file, opts)
	Log = slog.New(handler)

	Log.Info("Logger initialized", "path", LogFilePath)
	return nil
}

func Info(msg string, args ...any) {
	if Log != nil {
		Log.Info(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if Log != nil {
		Log.Error(msg, args...)
	}
}

func Debug(msg string, args ...any) {
	if Log != nil {
		Log.Debug(msg, args...)
	}
}

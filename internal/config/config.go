package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services          map[string]ServiceConfig `yaml:"services"`
	PjsonScriptRunner string                   `yaml:"pjson_script_runner"`
}

type ServiceConfig struct {
	Path    string            `yaml:"path"`
	Scripts map[string]string `yaml:"scripts"`
}


func GetConfigPath() (string, error) {
	// 1. Check environment variable
	if envPath := os.Getenv("YJ_CONFIG"); envPath != "" {
		return envPath, nil
	}

	// 2. Check local directory
	localPath, err := filepath.Abs("services.yaml")
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// 3. Check global configuration
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(userConfigDir, "yj")
	globalPath := filepath.Join(configDir, "services.yaml")

	// Create config directory if it doesn't exist (for future writes)
	// We don't create the file itself here, just the directory structure
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return globalPath, nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.PjsonScriptRunner == "" {
		cfg.PjsonScriptRunner = "pnpm"
	}

	for name, service := range cfg.Services {
		// If scripts map is nil, initialize it
		if service.Scripts == nil {
			service.Scripts = make(map[string]string)
		}

		// Check for package.json
		pjsonPath := filepath.Join(service.Path, "package.json")
		if pjsonData, err := os.ReadFile(pjsonPath); err == nil {
			var pjson struct {
				Scripts map[string]string `json:"scripts"`
			}
			if err := json.Unmarshal(pjsonData, &pjson); err == nil {
				// Add package.json scripts if not already present
				for scriptName := range pjson.Scripts {
					if _, exists := service.Scripts[scriptName]; !exists {
						service.Scripts[scriptName] = fmt.Sprintf("%s run %s", cfg.PjsonScriptRunner, scriptName)
					}
				}
			}
		}
		// Update the service back in the map because it's a value receiver
		cfg.Services[name] = service
	}

	return &cfg, nil
}

func UpdateService(path string, name string, yj ServiceConfig) error {
	// Load existing config (or create empty)
	cfg := Config{
		Services: make(map[string]ServiceConfig),
	}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return err
		}
	}

	// Update or insert
	cfg.Services[name] = yj

	// Marshal back to YAML
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	// Write atomically (safe write)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func DeleteService(path string, name string) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}

	if _, ok := cfg.Services[name]; !ok {
		return fmt.Errorf("service %s not found", name)
	}

	delete(cfg.Services, name)

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

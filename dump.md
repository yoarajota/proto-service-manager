# Project Context: `yj` (Local Service Manager)

## 1. Overview

`yj` is a lightweight CLI/TUI tool designed to manage local development services on your machine. It avoids the complexity of Docker or Kubernetes for simple local setups, relying instead on a YAML configuration file (`services.yaml`) and direct process execution.

- **Goal**: Run services locally without daemons or cloud dependencies.
- **Primary Interface**: CLI (Command Line Interface), with TUI (Terminal User Interface) capabilities (powered by Bubble Tea).

## 2. How to Maintain This Document

This document (`context.md`) serves as the long-term memory and context for both AI Agents and human developers. It must remain accurate, concise, and structured.

### 2.1. When to Update

- **New Features**: When adding a new command or module, update **Section 3 (Tech Stack)** if new libs are added, and **Section 6 (Development Patterns)** if a new pattern emerges.
- **Refactoring**: If the project structure changes, update **Section 4**.
- **Data Model Changes**: If `services.yaml` schema changes, update **Section 5**.

### 2.2. Writing Guidelines

- **Audience**: AI Agents (primary) and Humans (secondary).
- **Style**:
  - Use **bold** for key terms.
  - Use code blocks for file paths, commands, and schemas.
  - Be explicit about relationships (e.g., "A Service _has a_ Path").
  - Avoid ambiguity.
- **Structure**:
  - Maintain the numbered sections.
  - Use `##` for top-level concepts and `###` for specifics.

### 2.3. Pattern for Extending Context

To add a new layer of features (e.g., "Remote Services"):

1.  **Define**: Add a new definition in **Section 5**.
2.  **Architect**: Show where it fits in **Section 4**.
3.  **Guide**: Add a "How-To" in **Section 6**.

## 3. Technology Stack

- **Language**: Go (v1.24+)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (referenced in dependencies)
- **Configuration**: YAML (`gopkg.in/yaml.v3`)
- **Process Management**: `os/exec` (via `internal/process`)
- **Logging**: Custom logger (`internal/logger`)

## 4. Project Structure

The project follows a standard Go CLI layout:

```text
.
├── cmd/                # CLI command implementations (cobra)
│   ├── base.go         # Entry point and global flags
│   ├── init.go         # `yj init` implementation
│   ├── list.go         # `yj list` implementation
│   ├── run.go          # `yj run` implementation
│   └── ...
├── internal/           # Private application and library code
│   ├── config/         # Configuration loading and parsing
│   ├── process/        # Process execution logic
│   ├── service/        # Service domain model
│   └── logger/         # Logging utilities
├── docs/               # Documentation
└── services.yaml       # User configuration file (example)
```

## 5. Key Concepts & Definitions

### 5.1. Core Domain Models

- **Service**: A named unit of work (e.g., "api", "web"). Defined in `services.yaml`. A service has:
  - `path`: The working directory for the service.
  - `scripts`: A collection of named commands (e.g., "dev", "build").

### 5.2. Configuration (`services.yaml`)

The source of truth for services. Resolved via `YJ_CONFIG`, local `./services.yaml`, or `~/.config/yj/services.yaml`.

Create the global config with `yj init`.

```yaml
services:
  myapp:
    path: /path/to/app
    scripts:
      dev: go run main.go
```

- **Auto-Discovery**: `yj` automatically scans `package.json` in the service's `path` and adds scripts (e.g., `npm run dev`) if they are missing from YAML.

### 5.3. Data Flow

1. **Resolve**: `config.GetConfigPath()` determines the file location (Env > Local > Global).
2. **Load**: `config.Load(path)` reads the file.
3. **Merge**: It iterates through services, checks for `package.json` in `path`, and injects discovered scripts into the in-memory config.
4. **Execute**: Commands like `run` use the resolved config to spawn processes.

## 6. Development Patterns

### Adding a New Command

1. Create a new file in `cmd/` (e.g., `cmd/restart.go`).
2. Define a `cobra.Command` variable.
3. Implement `RunE` logic.
4. Register it in `init()`: `baseCmd.AddCommand(restartCmd)`.

### Modifying Configuration

1. Update `ServiceConfig` struct in `internal/config/config.go`.
2. Add corresponding struct tags (`yaml:"..."`).
3. Update `services.yaml` to reflect changes.

### Error Handling

- Use `RunE` in Cobra commands to return errors.
- Use `internal/logger` for structured logging (Debug/Info/Error).

## 7. Future Extensibility

- **TUI**: The project includes `bubbletea`. Future features can implement interactive dashboards in `cmd/` using `tea.NewProgram`.
- **Process Supervision**: Current implementation starts processes. Future layers could add restart policies, log streaming, or status monitoring.

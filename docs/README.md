# yj — Local Service Manager

`yj` is a lightweight CLI/TUI tool to manage local services on your machine using a YAML config and an interactive UI.

Runs entirely locally. No daemon. No cloud. No magic.

---

## Requirements

- Go **1.21+**
- macOS / Linux / Windows

---

## Installation

### Option 1 — Install via Go (recommended)

From the project root:

```bash
go install .
```

This builds and installs the `yj` binary into:

```text
~/go/bin/yj
```

### Ensure `yj` is in your PATH

Check:

```bash
which yj
```

If it returns nothing, add Go’s bin directory to your PATH.

#### macOS / Linux (Conditional)

Check which shell you are using:

```bash
echo $SHELL
```

**If using zsh:**

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**If using bash:**

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Now verify:

```bash
yj --help
```

---

### Option 2 — Build manually

```bash
go build -o yj
```

Run locally:

```bash
./yj list
```

Install globally:

```bash
mv yj /usr/local/bin
```

(You may need `sudo`.)

---

## Configuration

`yj` searches for `services.yaml` in the following order:

1. **Environment Variable**: `YJ_CONFIG`
2. **Local Directory**: `./services.yaml`
3. **Global Config**: `~/.config/yj/services.yaml`

Example `services.yaml`:

```yaml
services:
  service-name:
    path: /path/to/service
    scripts:
      script-name: command
```

If no configuration is found, `yj` can create one in the global location.

Create the global config:

```bash
yj init
```

---

## Usage

List services:

```bash
yj list
```

## Uninstall / Remove

### 1️⃣ Remove the binary

If installed via `go install`:

```bash
rm ~/go/bin/yj
```

If installed manually:

```bash
rm /usr/local/bin/yj
```

Verify removal:

```bash
which yj
```

(No output means it’s gone.)

---

### 2️⃣ (Optional) Remove configuration and data

This deletes **all yj data**:

```bash
rm -rf ~/.config/yj
```

---

## Development

Run without installing:

```bash
go run . list
```

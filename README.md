# CFCTL

A modern command-line interface for managing Cloudflare services with advanced cache management capabilities.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)](#supported-platforms)
[![Cloudflare SDK](https://img.shields.io/badge/Cloudflare%20SDK-v6.5.0-F38020?style=flat)](https://github.com/cloudflare/cloudflare-go)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Security](#security)
- [Architecture](#architecture)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview

CFCTL is a production-grade command-line interface designed for efficient management of Cloudflare services. Built with Go 1.24 and the Cloudflare SDK v6, it provides an interactive terminal user interface (TUI) for cache management, domain administration, and secure multi-account credential handling.

### Key Capabilities

- **Secure Credential Management**: Leverages system-native keyring services for encrypted storage of API credentials
- **Multi-Account Support**: Seamlessly manage multiple Cloudflare accounts from a single interface
- **Advanced Cache Purging**: Five distinct purge methods for granular cache control
- **Interactive Terminal UI**: Built with Bubble Tea framework for smooth, responsive user experience
- **Cross-Platform Compatibility**: Native binaries for macOS (Intel/Apple Silicon) and Linux (AMD64/ARM64)

## Features

### Cache Management

CFCTL provides comprehensive cache purging capabilities:

1. **Purge by URL**: Remove specific files from cache by providing exact URLs
2. **Purge by Hostname**: Clear all cached assets for a specific hostname
3. **Purge by Tag**: Remove cache entries matching specific tags (Enterprise feature)
4. **Purge by Prefix**: Clear cache for all URLs matching a path prefix
5. **Purge Everything**: Complete zone cache invalidation with safety confirmations

### Account Management

- Secure storage of API tokens and global API keys via system keyring
- Support for multiple Cloudflare accounts with easy switching
- Account removal with automatic credential cleanup
- Persistent configuration with YAML-based storage

### Domain Operations

- List all zones associated with configured accounts
- Interactive domain selection interface
- Cached domain listings with configurable TTL

### User Interface

- Modern terminal UI with smooth animations
- Keyboard-driven navigation (Vim-style key bindings supported)
- Colored output with optional monochrome mode
- Confirmation prompts for destructive operations
- Real-time operation feedback

## System Requirements

### Supported Platforms

| Platform | Architecture | Minimum Version |
|----------|-------------|-----------------|
| macOS | AMD64 (Intel) | macOS 10.15+ |
| macOS | ARM64 (Apple Silicon) | macOS 11.0+ |
| Linux | AMD64 | Kernel 3.10+ |
| Linux | ARM64 | Kernel 3.10+ |

### Dependencies

**Runtime Requirements:**
- System keyring service:
  - macOS: Keychain Services (built-in)
  - Linux: Secret Service API (GNOME Keyring, KDE Wallet, or compatible)

**Build Requirements** (for compilation from source):
- Go 1.24 or later
- Make (GNU Make 3.81+)
- Git

### Network Requirements

- HTTPS connectivity to Cloudflare API endpoints
- Outbound access to `api.cloudflare.com` (port 443)

## Installation

### Quick Install from Source

```bash
# Clone the repository
git clone https://github.com/siyamsarker/cfctl.git
cd cfctl

# Build and install
make build
sudo ./scripts/install.sh
```

The installation script will:
- Detect your system architecture automatically
- Install the binary to `/usr/local/bin/cfctl`
- Create configuration directory at `~/.config/cfctl`
- Set appropriate executable permissions

### Manual Installation

```bash
# Build for your platform
make build

# Install to system path
sudo cp bin/cfctl /usr/local/bin/

# Verify installation
cfctl --version
```

### Platform-Specific Builds

```bash
# Build for macOS (Intel)
make build-darwin

# Build for Linux (AMD64)
make build-linux

# Build for all platforms
make build-all
```

Compiled binaries will be available in the `bin/` directory:
- `cfctl-darwin-amd64` - macOS Intel
- `cfctl-darwin-arm64` - macOS Apple Silicon
- `cfctl-linux-amd64` - Linux AMD64
- `cfctl-linux-arm64` - Linux ARM64

### Uninstallation

```bash
sudo ./scripts/uninstall.sh
```

This will remove:
- Binary from `/usr/local/bin/cfctl`
- Configuration directory (with user confirmation)
- Stored credentials from system keyring (with user confirmation)

## Quick Start

### Initial Setup

1. **Launch CFCTL**
   ```bash
   cfctl
   ```

2. **Configure Cloudflare Account**
   - Select "Configure Cloudflare Account" from the main menu
   - Choose authentication method (API Token recommended)
   - Enter account details and credentials

3. **Obtain API Credentials**

   **API Token (Recommended)**
   - Navigate to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Click "Create Token"
   - Use "Edit zone DNS" template or create custom with permissions:
     - Zone - Zone - Read
     - Zone - Cache Purge - Purge
   - Copy the generated token

   **Global API Key (Legacy)**
   - Navigate to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Locate "Global API Key" section
   - Click "View" and copy the key

### Basic Operations

**Managing Domains**
```bash
cfctl
# Navigate to: Manage Domains → Select Domain → Choose Operation
```

**Purging Cache**
1. Select domain from the domain list
2. Choose purge method:
   - Purge by URL (for specific files)
   - Purge by Hostname (for entire subdomains)
   - Purge Everything (with double confirmation)
3. Follow interactive prompts
4. Confirm operation

## Usage

### Command-Line Interface

```bash
cfctl [flags]
```

### Available Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--account` | `-a` | string | Use specific Cloudflare account |
| `--config` | `-c` | string | Config file path (default: `~/.config/cfctl/config.yaml`) |
| `--no-color` | | boolean | Disable colored output |
| `--debug` | | boolean | Enable debug mode with verbose logging |
| `--quiet` | `-q` | boolean | Suppress non-error output |
| `--version` | `-v` | boolean | Display version information |
| `--help` | `-h` | boolean | Display help information |

### Common Usage Patterns

**Use specific account**
```bash
cfctl --account production
```

**Custom configuration file**
```bash
cfctl --config /path/to/config.yaml
```

**Debug mode for troubleshooting**
```bash
cfctl --debug
```

**Disable colors for CI/CD environments**
```bash
cfctl --no-color
```

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `Enter` | Select / Confirm |
| `Esc` / `q` | Back / Cancel |
| `Ctrl+C` | Quit application |
| `Tab` | Next field (forms) |
| `Shift+Tab` | Previous field (forms) |

### Running with Elevated Privileges

CFCTL generally runs without elevated privileges. However, if your system requires sudo for keyring access:

```bash
sudo cfctl
```

When run with sudo, CFCTL automatically uses the invoking user's home directory for configuration and credentials. To specify a custom config path:

```bash
sudo cfctl --config /path/to/config.yaml
```

**Note**: The `--version` and `--help` flags never require elevated privileges.

## Configuration

### Configuration File Location

Default: `~/.config/cfctl/config.yaml`

Override via:
- Command-line flag: `--config /path/to/config.yaml`
- Environment variable: `CFCTL_CONFIG=/path/to/config.yaml`

### Configuration Structure

```yaml
version: 1

defaults:
  account: ""           # Default account name (empty = prompt on startup)
  theme: "dark"         # UI theme: dark, light, system
  output: "interactive" # Output mode: interactive, json, table

api:
  timeout: 30           # Request timeout in seconds (default: 30)
  retries: 3            # Number of retry attempts (default: 3)

ui:
  confirmations: true   # Show confirmation prompts (default: true)
  animations: true      # Enable UI animations (default: true)
  colors: true          # Enable colored output (default: true)

cache:
  domains_ttl: 300      # Domain list cache TTL in seconds (default: 300)
  enabled: true         # Enable local caching (default: true)

accounts: []            # Account list (managed by application)
```

### Configuration Options

**defaults**
- `account`: Automatically select this account on startup (omit to show account selector)
- `theme`: Color scheme for terminal UI
- `output`: Future-proofing for non-interactive modes

**api**
- `timeout`: Maximum wait time for API requests (seconds)
- `retries`: Automatic retry attempts for failed requests

**ui**
- `confirmations`: Require user confirmation for destructive operations
- `animations`: Enable/disable UI transition animations
- `colors`: Control colored output (overridden by `--no-color` flag)

**cache**
- `domains_ttl`: How long to cache domain listings before refreshing
- `enabled`: Toggle local caching of API responses

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CFCTL_CONFIG` | Override config file location |
| `NO_COLOR` | Disable colored output (set to any value) |
| `CFCTL_DEBUG` | Enable debug logging (set to any value) |
| `HOME` | User home directory (for config/credential paths) |
| `XDG_CONFIG_HOME` | XDG base directory (overrides `~/.config`) |

## Security

### Credential Storage

CFCTL uses platform-native keyring services for secure credential storage:

| Platform | Service | Implementation |
|----------|---------|----------------|
| macOS | Keychain Services | Apple Keychain Access |
| Linux | Secret Service API | GNOME Keyring / KDE Wallet |

**Security Guarantees:**
- Credentials are **never** stored in plain text
- API tokens are encrypted at rest using OS-level encryption
- Configuration files (`config.yaml`) contain **no** sensitive data
- Keyring access requires user authentication on first use

### Best Practices

1. **Use API Tokens instead of Global API Keys**
   - Tokens provide scoped permissions
   - Can be revoked without affecting other integrations
   - Support fine-grained access control

2. **Apply Principle of Least Privilege**
   
   Required permissions for cache management:
   ```
   Zone - Zone - Read
   Zone - Cache Purge - Purge
   ```

3. **Rotate Credentials Regularly**
   - Regenerate API tokens every 90 days
   - Remove unused accounts with "Remove Account" feature

4. **Use Separate Tokens per Environment**
   - Development: Limited zone access
   - Staging: Staging zones only
   - Production: Production zones with team review

5. **Audit Account Access**
   - Review configured accounts periodically
   - Remove unused accounts to minimize attack surface

### Security Considerations

- CFCTL does **not** transmit credentials to any third parties
- All API communication uses HTTPS (TLS 1.2+)
- Cloudflare API keys are handled according to Cloudflare's security guidelines
- Local cache files contain **only** non-sensitive metadata (zone names, IDs)

## Architecture

### Project Structure

```
cfctl/
├── cmd/
│   └── cfctl/              # Application entry point
│       └── main.go         # CLI initialization, flag parsing
├── internal/
│   ├── api/                # Cloudflare API client
│   │   ├── cache.go        # Cache purge operations
│   │   ├── client.go       # API client initialization
│   │   ├── client_test.go  # Client unit tests
│   │   └── zones.go        # Zone/domain operations
│   ├── config/             # Configuration management
│   │   ├── accounts.go     # Account CRUD operations
│   │   ├── config.go       # Config file handling
│   │   ├── validator.go    # Input validation
│   │   └── validator_test.go
│   ├── handlers/           # Business logic layer
│   ├── ui/                 # Terminal UI components
│   │   ├── welcome.go      # Welcome screen
│   │   ├── menu.go         # Main menu
│   │   ├── account_*.go    # Account management screens
│   │   ├── domain_list.go  # Domain selection
│   │   ├── purge_*.go      # Cache purge interfaces
│   │   ├── settings.go     # Settings screen
│   │   ├── styles.go       # UI styling (Lip Gloss)
│   │   └── help.go         # Help screen
│   └── utils/              # Utility functions
│       ├── helpers.go      # General helpers
│       ├── validator.go    # Validation utilities
│       └── validator_test.go
├── pkg/
│   └── cloudflare/         # Public types and models
├── configs/
│   └── default.yaml        # Default configuration template
├── scripts/
│   ├── install.sh          # Installation script
│   └── uninstall.sh        # Uninstallation script
├── bin/                    # Compiled binaries (generated)
├── Makefile                # Build automation
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── LICENSE                 # MIT License
└── README.md               
```

### Technology Stack

| Component | Library | Version | Purpose |
|-----------|---------|---------|---------|
| Language | Go | 1.24.0 | Core implementation |
| Cloudflare SDK | cloudflare-go | v6.5.0 | API client |
| TUI Framework | Bubble Tea | v1.3.10 | Interactive terminal UI |
| UI Styling | Lip Gloss | v1.1.0 | Terminal styling |
| UI Components | Bubbles | v0.21.0 | Reusable TUI widgets |
| CLI Framework | Cobra | v1.10.2 | Command-line interface |
| Configuration | Viper | v1.21.0 | Config file management |
| Keyring | go-keyring | v0.2.6 | Secure credential storage |
| Testing | testify | v1.11.1 | Test assertions |

### Component Responsibilities

**cmd/cfctl**: Application bootstrap, flag parsing, environment setup

**internal/api**: Cloudflare API integration layer
- Client initialization with authentication
- Zone listing and filtering
- Cache purge operations (all methods)
- Error handling and retry logic

**internal/config**: Configuration and credential management
- YAML configuration loading/saving
- Account CRUD operations
- Keyring integration for secure storage
- Input validation

**internal/ui**: Bubble Tea-based terminal interface
- Screen navigation and state management
- Form rendering and input handling
- Visual styling and theming
- User interaction flows

**internal/utils**: Shared utilities
- URL validation
- String manipulation
- Data formatting

**pkg/cloudflare**: Public types for external use (future extensibility)

### Build System

The Makefile provides comprehensive build automation:

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make build-all` | Build for all platforms |
| `make build-darwin` | Build macOS binaries (Intel + ARM) |
| `make build-linux` | Build Linux binaries (AMD64 + ARM64) |
| `make test` | Run unit tests |
| `make test-coverage` | Generate coverage report |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with gofmt |
| `make tidy` | Tidy Go modules |
| `make clean` | Remove build artifacts |
| `make install` | Build and install locally |
| `make run` | Build and run application |

**Build Flags:**
```bash
-ldflags="-s -w -X main.version=$(VERSION)"
```
- `-s`: Omit symbol table
- `-w`: Omit DWARF debug info
- `-X main.version=...`: Embed version string

## Development

### Setting Up Development Environment

1. **Clone Repository**
   ```bash
   git clone https://github.com/siyamsarker/cfctl.git
   cd cfctl
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build Application**
   ```bash
   make build
   ```

4. **Run Tests**
   ```bash
   make test
   ```

### Development Workflow

**Running Locally**
```bash
make run
```

**Running with Debug Mode**
```bash
./bin/cfctl --debug
```

**Code Formatting**
```bash
make fmt
```

**Linting** (requires [golangci-lint](https://golangci-lint.run/))
```bash
make lint
```

**Test Coverage**
```bash
make test-coverage
# Opens coverage.html in browser
```

### Testing

**Run All Tests**
```bash
go test ./...
```

**Run Specific Package Tests**
```bash
go test ./internal/config -v
```

**Run Tests with Coverage**
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Test Configuration**
- Unit tests: `*_test.go` files
- Test framework: `stretchr/testify`
- Mock objects: Interface-based mocking

### Adding New Features

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Implement Changes**
   - Add code in appropriate `internal/` package
   - Write unit tests
   - Update UI components if needed

3. **Test Changes**
   ```bash
   make test
   make lint
   ```

4. **Update Documentation**
   - Update README.md if user-facing
   - Add inline code documentation
   - Update default.yaml if adding config options

5. **Submit Pull Request**
   - Ensure all tests pass
   - Follow conventional commit messages
   - Reference related issues

### Code Style Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Keep functions focused and testable
- Document exported functions and types
- Write descriptive commit messages

### Debugging

**Enable Debug Logging**
```bash
cfctl --debug
```

**Check Configuration**
```bash
cat ~/.config/cfctl/config.yaml
```

**Verify Credentials**
```bash
# macOS
security find-generic-password -s "cfctl" -a "<account_name>"

# Linux (using secret-tool)
secret-tool lookup service cfctl account "<account_name>"
```

## Troubleshooting

### Installation Issues

**Error: Binary not found**
```bash
# Ensure you've built the application first
make build

# Verify binary location
ls -la bin/
```

**Error: Permission denied during installation**
```bash
# Use sudo for system-wide installation
sudo ./scripts/install.sh
```

**Error: Command not found after installation**
```bash
# Verify installation path
which cfctl

# Check if /usr/local/bin is in PATH
echo $PATH | grep "/usr/local/bin"

# If not, add to your shell profile (~/.bashrc, ~/.zshrc):
export PATH="/usr/local/bin:$PATH"
```

### Runtime Issues

**Error: Unable to load configuration**
```bash
# Check config file exists and is valid YAML
cat ~/.config/cfctl/config.yaml

# Reset to default configuration
rm ~/.config/cfctl/config.yaml
cfctl  # Will recreate with defaults
```

**Error: Keyring access denied**

*macOS:*
```bash
# Grant terminal access to keychain in System Preferences
# Security & Privacy → Privacy → Full Disk Access → Add Terminal
```

*Linux:*
```bash
# Ensure keyring service is running
systemctl --user status gnome-keyring

# Or for KDE
systemctl --user status kwalletd5
```

**Error: API authentication failed**
```bash
# Verify credentials at Cloudflare dashboard
# Regenerate API token if necessary
cfctl
# Navigate to: Configure Account → Update credentials
```

**Error: No domains found**
```bash
# Verify API token has Zone:Read permission
# Check account has zones configured at Cloudflare dashboard
```

### API Issues

**Error: Request timeout**
```yaml
# Increase timeout in ~/.config/cfctl/config.yaml
api:
  timeout: 60  # Increase from default 30 seconds
```

**Error: Rate limit exceeded**
```bash
# Wait for rate limit reset
# Consider reducing operation frequency
# Enterprise plans have higher rate limits
```

**Error: Invalid zone ID**
```bash
# Clear cached zone list
rm -rf ~/.cache/cfctl/  # If cache directory exists
cfctl  # Refresh zone list
```

### UI Issues

**UI rendering incorrectly**
```bash
# Try disabling colors
cfctl --no-color

# Check terminal compatibility
echo $TERM  # Should be xterm-256color or similar
```

**Keyboard navigation not working**
```bash
# Ensure terminal supports ANSI escape sequences
# Try different terminal emulator (iTerm2, GNOME Terminal, etc.)
```

### Building from Source

**Error: Go version mismatch**
```bash
# Check Go version
go version

# Install Go 1.24 or later from https://go.dev/dl/
```

**Error: Module dependency issues**
```bash
# Clean and refresh modules
go clean -modcache
go mod download
go mod tidy
```

**Error: Build failed with linker error**
```bash
# Ensure sufficient disk space
df -h

# Try clean build
make clean
make build
```

### Getting Help

If you encounter issues not covered here:

1. **Check Existing Issues**: [GitHub Issues](https://github.com/siyamsarker/cfctl/issues)
2. **Enable Debug Mode**: Run with `--debug` flag and include output in issue report
3. **Create New Issue**: Provide:
   - Operating system and version
   - Go version (`go version`)
   - CFCTL version (`cfctl --version`)
   - Steps to reproduce
   - Debug output if applicable

## Contributing

Contributions are welcome and appreciated. To contribute:

### Reporting Issues

1. Search existing issues to avoid duplicates
2. Use issue templates when available
3. Provide detailed reproduction steps
4. Include system information (OS, Go version, CFCTL version)

### Submitting Pull Requests

1. **Fork the Repository**
   ```bash
   git clone https://github.com/your-username/cfctl.git
   cd cfctl
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**
   - Write clean, documented code
   - Follow existing code style
   - Add tests for new functionality
   - Update documentation

4. **Test Thoroughly**
   ```bash
   make test
   make lint
   make build
   ```

5. **Commit Changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```
   
   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` New feature
   - `fix:` Bug fix
   - `docs:` Documentation changes
   - `test:` Test additions/changes
   - `refactor:` Code refactoring
   - `chore:` Maintenance tasks

6. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then open a pull request on GitHub

### Development Guidelines

- **Code Quality**: All code must pass `make lint` and `make test`
- **Documentation**: Update README.md for user-facing changes
- **Tests**: Maintain or improve test coverage
- **Commits**: Use clear, descriptive commit messages
- **Dependencies**: Minimize external dependencies

### Review Process

1. Automated CI checks must pass
2. Code review by maintainers
3. Address review feedback
4. Squash commits if requested
5. Merge upon approval

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for full text.

---

**Repository**: [github.com/siyamsarker/cfctl](https://github.com/siyamsarker/cfctl)  
**Documentation**: [Cloudflare API Documentation](https://developers.cloudflare.com/api/)  
**Issues**: [GitHub Issues](https://github.com/siyamsarker/cfctl/issues)  
**Author**: [Siyam Sarker](https://github.com/siyamsarker)

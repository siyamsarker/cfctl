<div align="center">

# CFCTL

**A modern, interactive command-line interface for managing Cloudflare services with focus on cache management.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)](#)
[![Cloudflare SDK](https://img.shields.io/badge/Cloudflare%20SDK-v6-F38020?style=flat)](https://github.com/cloudflare/cloudflare-go)
</div>

## ✨ Features

- 🔐 **Secure Credential Management** - Store API credentials securely in system keyring
- 👥 **Multi-Account Support** - Manage multiple Cloudflare accounts effortlessly
- 🌐 **Domain Management** - List and select domains/zones with ease
- 🗑️ **Advanced Cache Purging** - Five different purge methods:
  - Purge by URL (specific files)
  - Purge by Hostname (all assets for a host)
  - Purge by Tag (Enterprise feature)
  - Purge by Prefix (path-based purging)
  - Purge Everything (entire zone cache)
- 🎨 **Beautiful TUI** - Interactive terminal UI with smooth navigation
- ⚡ **Fast & Lightweight** - Single binary, no dependencies required
- 🔄 **Cross-Platform** - Works on macOS and Linux

## 📦 Installation

### macOS & Linux

#### Quick Install (from source)

```bash
# Clone the repository
git clone https://github.com/siyamsarker/cfctl.git
cd cfctl

# Build and install
make build
sudo ./scripts/install.sh
```

#### Manual Installation

```bash
# Build for your platform
make build

# Move to your PATH
sudo cp bin/cfctl /usr/local/bin/

# Verify installation
cfctl --version
```

## 🚀 Quick Start

### 1. First Run

Launch cfctl for the first time:

```bash
cfctl
```

You'll see the welcome screen. Press Enter to continue.

### 2. Configure Your Cloudflare Account

1. Select **"Configure Cloudflare Account"** from the main menu
2. Choose authentication method:
   - **API Token** (Recommended) - More secure, scoped permissions
   - **Global API Key** - Full account access
3. Enter your credentials:
   - Account name (friendly identifier)
   - Email (for Global API Key only)
   - API Token or Key

### 3. Get Your Cloudflare Credentials

#### API Token (Recommended)

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Click **"Create Token"**
3. Use **"Edit zone DNS"** template or create custom with:
   - **Permissions:**
     - Zone - Zone - Read
     - Zone - Cache Purge - Purge
4. Copy the token

#### Global API Key

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Scroll to **"Global API Key"**
3. Click **"View"** and copy the key

### 4. Start Managing Cache

1. Select **"Manage Domains"**
2. Choose a domain from the list
3. Select your preferred purge method
4. Follow the prompts

## 📖 Usage

### Sudo Usage Notes

- `cfctl --version` and `cfctl --help` work without sudo.
- Interactive mode should be run without sudo whenever possible.
- If your environment requires sudo (for example, restricted keyring access), use:

```bash
sudo cfctl
```

When run with sudo, cfctl automatically uses the invoking user’s home directory for config and credentials. To use a custom config path, pass `--config`:

```bash
sudo cfctl --config /path/to/config.yaml
```

### Interactive Mode (Default)

Simply run:

```bash
cfctl
```

Navigate using:
- **↑/↓ or j/k** - Navigate menus
- **Enter** - Select/Confirm
- **Esc or q** - Back/Cancel
- **Ctrl+C** - Quit application
- **Tab/Shift+Tab** - Navigate form fields

### Command-Line Flags

cfctl supports several flags for enhanced control:

```bash
# Use a specific account
cfctl --account production

# Use a custom config file
cfctl --config ~/.cfctl-work.yaml

# Disable colored output (useful for CI/CD or logging)
cfctl --no-color

# Enable debug mode with verbose logging
cfctl --debug

# Suppress non-error output
cfctl --quiet

# Display version information
cfctl --version

# Display help
cfctl --help
```

**Available Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--account` | `-a` | Use specific Cloudflare account |
| `--config` | `-c` | Config file path (default: `~/.config/cfctl/config.yaml`) |
| `--no-color` | | Disable colored output |
| `--debug` | | Enable debug mode with verbose logging |
| `--quiet` | `-q` | Suppress non-error output |
| `--version` | `-v` | Display version information |
| `--help` | `-h` | Display help information |

### Configuration

Configuration is stored at `~/.config/cfctl/config.yaml`

Default configuration:

```yaml
version: 1

defaults:
  account: ""
  theme: dark
  output: interactive

api:
  timeout: 30
  retries: 3

ui:
  confirmations: true
  animations: true
  colors: true

cache:
  domains_ttl: 300
  enabled: true

accounts: []
```

### Environment Variables

- `CFCTL_CONFIG` - Override config file location
- `NO_COLOR` - Disable colored output (set by `--no-color` flag)
- `CFCTL_DEBUG` - Enable debug mode (set by `--debug` flag)

## 🎯 Use Cases

### Purge Specific Files After Deployment

```
1. Launch cfctl
2. Select your domain
3. Choose "Purge by URL"
4. Enter URLs:
   https://example.com/css/main.css
   https://example.com/js/app.js
5. Confirm
```

### Clear Cache for Entire Subdomain

```
1. Launch cfctl
2. Select your domain
3. Choose "Purge by Hostname"
4. Enter: blog.example.com
5. Confirm
```

### Emergency Full Cache Clear

```
1. Launch cfctl
2. Select your domain
3. Choose "Purge Everything"
4. Confirm twice (safety feature)
```

## 🔒 Security

### Credential Storage

cfctl stores credentials securely using your system's native keyring:

- **macOS**: Keychain Services
- **Linux**: Secret Service API (GNOME Keyring / KDE Wallet)

API tokens/keys are **never** stored in plain text configuration files.

### Best Practices

1. Use **API Tokens** instead of Global API Keys
2. Create tokens with **minimal required permissions**:
   - Zone:Read
   - Cache Purge:Purge
3. Rotate tokens regularly
4. Use different tokens for different environments

## 🛠️ Development

### Technical Stack

- **Language**: Go 1.21+
- **Cloudflare SDK**: [cloudflare-go v6](https://github.com/cloudflare/cloudflare-go)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Keyring**: [go-keyring](https://github.com/zalando/go-keyring)

### Project Structure

```
cfctl/
├── .gitignore          # Git ignore rules
├── LICENSE             # MIT License
├── Makefile            # Build automation and tasks
├── README.md           # Project documentation
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
├── bin/                # Compiled binaries
├── cmd/
│   └── cfctl/          # Application entry point (main.go)
├── configs/            # Default configuration files
├── internal/
│   ├── api/            # Cloudflare API client (v6 SDK)
│   ├── config/         # Configuration management & keyring
│   ├── handlers/       # Business logic handlers
│   ├── ui/             # Terminal UI components (Bubble Tea)
│   └── utils/          # Utility functions & validators
├── pkg/
│   └── cloudflare/     # Public types and models
└── scripts/            # Installation & uninstallation scripts
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platform
make build-darwin
make build-linux

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
make test-coverage

# Run specific package tests
go test ./internal/config -v
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Write tests for new features
- Update documentation as needed
- Format code with `gofmt`
- Run linter before submitting

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cloudflare](https://www.cloudflare.com/) for their excellent API and Go SDK
- [Charm](https://charm.sh/) for Bubble Tea TUI framework and Lip Gloss styling
- [Cobra](https://cobra.dev/) for CLI command framework
- [Viper](https://github.com/spf13/viper) for configuration management
- All open-source contributors and the Go community

## 📞 Support

- **Documentation**: [Cloudflare API Docs](https://developers.cloudflare.com/api/)
- **Issues**: Report bugs and request features via GitHub Issues
- **Discussions**: Join community discussions on GitHub

## 🗺️ Roadmap

### Planned Features

- [ ] **DNS Management** - Add, edit, and delete DNS records
- [ ] **SSL/TLS Configuration** - Manage SSL settings and certificates
- [ ] **Firewall Rules** - Create and manage firewall rules
- [ ] **Page Rules** - Configure page rules for zones
- [ ] **Non-interactive Mode** - CLI flags for automation and scripting
- [ ] **Configuration Presets** - Save and load configuration profiles
- [ ] **Batch Operations** - Perform operations on multiple domains
- [ ] **Cache Analytics** - View cache hit rates and statistics
- [ ] **Export/Import** - Export configurations and import them

### Distribution

- [ ] Docker image for containerized deployment
- [ ] Homebrew formula for macOS users
- [ ] APT/RPM packages for Linux distributions
- [ ] Pre-built binaries for releases

---

<p align="center">
  <strong>Made with ❤️ by <a href="https://github.com/siyamsarker">Siyam Sarker</a></strong>
  <br>
  <sub>© 2026 cfctl. Released under the MIT License.</sub>
</p>

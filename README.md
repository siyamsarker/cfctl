# cfctl - Cloudflare CLI Management Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)](https://github.com/siyamsarker/cfctl)

> A modern, interactive command-line interface for managing Cloudflare services with focus on cache management.

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
./scripts/install.sh
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

### Build from Source

Requirements:
- Go 1.21 or higher

```bash
# Clone and build
git clone https://github.com/siyamsarker/cfctl.git
cd cfctl
make build

# Run directly
./bin/cfctl
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

### Project Structure

```
cfctl/
├── cmd/cfctl/          # Application entry point
├── internal/
│   ├── api/            # Cloudflare API client
│   ├── config/         # Configuration management
│   ├── ui/             # Terminal UI components
│   ├── handlers/       # Business logic handlers
│   └── utils/          # Utility functions
├── pkg/cloudflare/     # Public types and models
├── scripts/            # Installation scripts
└── Makefile            # Build automation
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

- [Cloudflare](https://www.cloudflare.com/) for their excellent API
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal styling

## 📞 Support

- **Documentation**: [Cloudflare API Docs](https://developers.cloudflare.com/api/)
- **Issues**: [GitHub Issues](https://github.com/siyamsarker/cfctl/issues)
- **Discussions**: [GitHub Discussions](https://github.com/siyamsarker/cfctl/discussions)

## 🗺️ Roadmap

- [ ] Add support for more Cloudflare services
- [ ] Non-interactive CLI mode for scripting
- [ ] Configuration presets/profiles
- [ ] Batch operations support
- [ ] Cache analytics and insights
- [ ] Docker image
- [ ] Homebrew formula

## 📊 Project Status

Current Version: **1.0.0**

Status: **Production Ready** ✅

---

**Made with ❤️ by [Siyam Sarker](https://github.com/siyamsarker)**

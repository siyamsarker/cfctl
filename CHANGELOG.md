# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-13

### Added
- Initial release of cfctl
- Interactive TUI with Bubble Tea framework
- Multi-account support with secure credential storage
- Five cache purge methods:
  - Purge by URL (specific files)
  - Purge by Hostname (all assets for a host)
  - Purge by Tag (Enterprise feature)
  - Purge by Prefix (path-based purging)
  - Purge Everything (entire zone cache)
- Domain/zone listing and selection
- Account configuration and management
- Settings management
- Help and documentation system
- Cross-platform support (macOS & Linux)
- Secure credential storage using system keyring
- Beautiful terminal UI with color themes
- Comprehensive error handling
- Installation and uninstallation scripts
- Makefile for build automation
- Complete documentation

### Security
- Credentials stored in system keyring (macOS Keychain / Linux Secret Service)
- Support for both API Tokens and Global API Keys
- No plain-text credential storage
- API token verification before saving

### Technical
- Built with Go 1.21+
- Uses official Cloudflare Go SDK v2
- Bubble Tea for TUI
- Lip Gloss for styling
- Viper for configuration management
- Cobra for CLI framework

[1.0.0]: https://github.com/siyamsarker/cfctl/releases/tag/v1.0.0

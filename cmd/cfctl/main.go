package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"

	// Global flags
	configFile  string
	accountName string
	noColor     bool
	debug       bool
	quiet       bool

	rootCmd = &cobra.Command{
		Use:   "cfctl",
		Short: "A modern CLI for managing Cloudflare services",
		Long: `CFCTL - Cloudflare CLI Management Tool

A modern, interactive command-line interface for managing Cloudflare services.
Features secure credential management, multi-account support, and advanced
cache purging capabilities with a beautiful terminal UI.

Examples:
  # Launch interactive mode (default)
  cfctl

  # Use specific account
  cfctl --account production

  # Use custom config file
  cfctl --config ~/.cfctl.yaml

  # Disable colored output
  cfctl --no-color

Documentation: https://github.com/siyamsarker/cfctl
Report bugs: https://github.com/siyamsarker/cfctl/issues`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			setSudoUserEnv()

			// Handle config file override
			if configFile != "" {
				os.Setenv("CFCTL_CONFIG", configFile)
			}

			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				if !quiet {
					fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				}
				os.Exit(1)
			}

			// Override default account if specified
			if accountName != "" {
				cfg.Defaults.Account = accountName
			}

			// Handle no-color flag (can be extended to disable colors in UI)
			if noColor {
				// Set environment variable that UI can check
				os.Setenv("NO_COLOR", "1")
			}

			// Handle debug mode (can be used for verbose logging)
			if debug {
				os.Setenv("CFCTL_DEBUG", "1")
			}

			// Attempt to resize terminal for better visibility (100 cols x 30 rows)
			// This uses ANSI escape sequence \x1b[8;{rows};{cols}t
			fmt.Print("\x1b[8;30;100t")

			// Launch interactive mode
			p := tea.NewProgram(
				ui.NewWelcomeModel(version, cfg),
				tea.WithAltScreen(),
			)

			if _, err := p.Run(); err != nil {
				if !quiet {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				os.Exit(1)
			}
		},
	}
)

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default: ~/.config/cfctl/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "", "use specific Cloudflare account")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode with verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-error output")

	// Customize help template
	rootCmd.SetHelpTemplate(`{{.Long}}

Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	// Customize version template
	rootCmd.SetVersionTemplate(`cfctl version {{.Version}}

Build information:
  Version:    {{.Version}}
  Go version: go1.21+
  Platform:   darwin/arm64 or linux/amd64

Cloudflare SDK: v6.5.0
`)
}

func setSudoUserEnv() {
	if os.Geteuid() != 0 {
		return
	}

	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser == "" {
		return
	}

	u, err := user.Lookup(sudoUser)
	if err != nil || u.HomeDir == "" {
		return
	}

	if os.Getenv("HOME") == "" || os.Getenv("HOME") == "/var/root" || os.Getenv("HOME") == "/root" {
		_ = os.Setenv("HOME", u.HomeDir)
	}

	if os.Getenv("XDG_CONFIG_HOME") == "" {
		_ = os.Setenv("XDG_CONFIG_HOME", filepath.Join(u.HomeDir, ".config"))
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if !quiet {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

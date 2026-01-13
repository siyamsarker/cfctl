package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	rootCmd = &cobra.Command{
		Use:     "cfctl",
		Short:   "Cloudflare CLI Management Tool",
		Long:    `A modern, interactive CLI for managing Cloudflare services via API`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Launch interactive mode
			p := tea.NewProgram(
				ui.NewWelcomeModel(version, cfg),
				tea.WithAltScreen(),
			)

			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

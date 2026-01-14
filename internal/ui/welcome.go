package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type WelcomeModel struct {
	version string
	config  *config.Config
}

func NewWelcomeModel(version string, cfg *config.Config) WelcomeModel {
	return WelcomeModel{
		version: version,
		config:  cfg,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			// Transition to main menu
			return NewMainMenuModel(m.config), nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	// Modern ASCII logo with gradient effect
	logo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render(`
    ╔══════════════════════════════════════════╗
    ║                                          ║
    ║     ██████╗███████╗ ██████╗████████╗██╗  ║
    ║    ██╔════╝██╔════╝██╔════╝╚══██╔══╝██║  ║
    ║    ██║     █████╗  ██║        ██║   ██║  ║
    ║    ██║     ██╔══╝  ██║        ██║   ██║  ║
    ║    ╚██████╗██║     ╚██████╗   ██║   ███║ ║
    ║     ╚═════╝╚═╝      ╚═════╝   ╚═╝   ╚══╝ ║
    ║                                          ║
    ╚══════════════════════════════════════════╝`)

	// Title and version
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		MarginTop(1).
		Render("Cloudflare CLI Management Tool")

	version := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render("v" + m.version)

	titleSection := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", version)

	// Subtitle with icon
	subtitle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		MarginTop(1).
		MarginBottom(2).
		Render("⚡ A powerful, interactive CLI for managing Cloudflare services")

	// Account status card
	var statusCard string
	if len(m.config.Accounts) > 0 {
		accountCount := lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Render(fmt.Sprintf("%d", len(m.config.Accounts)))

		statusCard = HighlightCardStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				SuccessStyle.Render("Ready"),
				" • Configured accounts: ",
				accountCount,
			),
		)
	} else {
		statusCard = CardStyle.Render(
			WarningStyle.Render("No accounts configured") +
				lipgloss.NewStyle().Foreground(MutedColor).Render("\nConfigure your Cloudflare account to get started"),
		)
	}

	// Feature highlights
	features := lipgloss.NewStyle().
		Foreground(MutedColor).
		MarginTop(2).
		MarginBottom(2).
		Render(
			"🔐 Secure credential management  •  🌐 Multi-account support\n" +
				"🗑️  Advanced cache purging       •  ⚡ Fast & lightweight",
		)

	// Navigation prompt
	prompt := lipgloss.NewStyle().
		MarginTop(2).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().
					Foreground(SuccessColor).
					Bold(true).
					Render("Press Enter"),
				lipgloss.NewStyle().
					Foreground(MutedColor).
					Render(" to continue  •  "),
				lipgloss.NewStyle().
					Foreground(ErrorColor).
					Bold(true).
					Render("q"),
				lipgloss.NewStyle().
					Foreground(MutedColor).
					Render(" to quit"),
			),
		)

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		titleSection,
		subtitle,
		statusCard,
		features,
		prompt,
	)

	// Center everything on screen
	return lipgloss.Place(
		120, 30,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

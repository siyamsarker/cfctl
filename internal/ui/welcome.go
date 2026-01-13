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
	logo := `
╔═══════════════════════════════════════╗
║                                       ║
║   ██████╗███████╗ ██████╗████████╗██╗ ║
║  ██╔════╝██╔════╝██╔════╝╚══██╔══╝██║ ║
║  ██║     █████╗  ██║        ██║   ██║ ║
║  ██║     ██╔══╝  ██║        ██║   ██║ ║
║  ╚██████╗██║     ╚██████╗   ██║   ███║║
║   ╚═════╝╚═╝      ╚═════╝   ╚═╝   ╚══╝║
║                                       ║
╚═══════════════════════════════════════╝
`

	title := TitleStyle.Render("Cloudflare CLI Management Tool v" + m.version)
	subtitle := SubtitleStyle.Render("A powerful, interactive CLI for managing Cloudflare services")

	var accountInfo string
	if len(m.config.Accounts) > 0 {
		accountInfo = SuccessStyle.Render("✓ Accounts configured: ") +
			InfoStyle.Render(fmt.Sprintf("%d", len(m.config.Accounts)))
	} else {
		accountInfo = WarningStyle.Render("⚠ No accounts configured yet")
	}

	prompt := SuccessStyle.Render("\nPress Enter to continue or 'q' to quit...")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		InfoStyle.Render(logo),
		"",
		title,
		"",
		subtitle,
		"",
		"",
		accountInfo,
		prompt,
	)

	return lipgloss.Place(
		80, 20,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

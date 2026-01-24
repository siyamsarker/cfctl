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
	width   int
	height  int
}

func NewWelcomeModel(version string, cfg *config.Config) WelcomeModel {
	return WelcomeModel{
		version: version,
		config:  cfg,
		width:   80,
		height:  24,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
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
	// Responsive logo - smaller for narrow terminals
	var logo string
	if m.width >= 60 {
		logo = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Render(`
   ██████╗███████╗ ██████╗████████╗██╗     
  ██╔════╝██╔════╝██╔════╝╚══██╔══╝██║     
  ██║     █████╗  ██║        ██║   ██║     
  ██║     ██╔══╝  ██║        ██║   ██║     
  ╚██████╗██║     ╚██████╗   ██║   ███████╗
   ╚═════╝╚═╝      ╚═════╝   ╚═╝   ╚══════╝`)
	} else {
		logo = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Render("[ CFCTL ]")
	}

	// Title and version with responsive layout
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		Render("Cloudflare CLI Management Tool")

	version := lipgloss.NewStyle().
		Foreground(MutedColor).
		Background(HighlightColor).
		Padding(0, 1).
		Render("v" + m.version)

	titleSection := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", version)

	// Subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		Render("A powerful, interactive CLI for managing Cloudflare services")

	// Responsive divider
	dividerWidth := min(m.width-10, 60)
	if dividerWidth < 20 {
		dividerWidth = 20
	}
	divider := lipgloss.NewStyle().
		Foreground(BorderColor).
		Render(repeatStr("─", dividerWidth))

	// Account status card - responsive
	var statusCard string
	cardWidth := min(m.width-20, 50)
	if cardWidth < 30 {
		cardWidth = 30
	}

	statusCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Width(cardWidth).
		Align(lipgloss.Center)

	if len(m.config.Accounts) > 0 {
		accountCount := lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Render(fmt.Sprintf("%d", len(m.config.Accounts)))

		statusCard = statusCardStyle.Copy().
			BorderForeground(SuccessColor).
			Render(
				lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("✓ Ready") +
					lipgloss.NewStyle().Foreground(MutedColor).Render(" • Accounts: ") +
					accountCount,
			)
	} else {
		statusCard = statusCardStyle.Copy().
			BorderForeground(WarningColor).
			Render(
				lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("⚠ No accounts configured") +
					"\n" +
					lipgloss.NewStyle().Foreground(MutedColor).Render("Configure your Cloudflare account to get started"),
			)
	}

	// Feature highlights - responsive grid
	var features string
	if m.width >= 70 {
		features = lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(
				"🔐 Secure credential management  •  🌐 Multi-account support\n" +
					"🗑️  Advanced cache purging       •  ⚡ Fast & lightweight",
			)
	} else {
		features = lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(
				"🔐 Secure credentials\n" +
					"🌐 Multi-account\n" +
					"🗑️  Cache purging\n" +
					"⚡ Lightweight",
			)
	}

	// Navigation prompt - modern pill style
	enterKey := lipgloss.NewStyle().
		Background(SuccessColor).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1).
		Render("Enter")

	quitKey := lipgloss.NewStyle().
		Background(BorderColor).
		Foreground(TextColor).
		Padding(0, 1).
		Render("q")

	prompt := lipgloss.JoinHorizontal(
		lipgloss.Center,
		enterKey,
		lipgloss.NewStyle().Foreground(MutedColor).Render(" to continue  •  "),
		quitKey,
		lipgloss.NewStyle().Foreground(MutedColor).Render(" to quit"),
	)

	// Combine all elements with proper spacing
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		titleSection,
		"",
		subtitle,
		"",
		divider,
		"",
		statusCard,
		"",
		features,
		"",
		divider,
		"",
		prompt,
	)

	// Center everything on screen with actual terminal dimensions
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

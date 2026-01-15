package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	returnTo tea.Model
	width    int
	height   int
}

func NewHelpModel(returnTo tea.Model) HelpModel {
	return HelpModel{
		returnTo: returnTo,
		width:    80,
		height:   24,
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			return m.returnTo, nil
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	// Header
	dividerWidth := min(m.width-8, 60)
	if dividerWidth < 30 {
		dividerWidth = 30
	}
	divider := lipgloss.NewStyle().
		Foreground(BorderColor).
		Render(repeatStr("─", dividerWidth))

	title := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("❓ Help & Documentation")

	// Keyboard shortcuts section
	keySection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("⌨️  Keyboard Shortcuts")

	keyStyle := lipgloss.NewStyle().
		Background(BorderColor).
		Foreground(TextColor).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(MutedColor)

	shortcuts := lipgloss.JoinVertical(
		lipgloss.Left,
		keyStyle.Render("↑↓")+descStyle.Render(" Navigate  ")+keyStyle.Render("Enter")+descStyle.Render(" Select"),
		keyStyle.Render("Esc")+descStyle.Render(" Back  ")+keyStyle.Render("q")+descStyle.Render(" Quit  ")+keyStyle.Render("Tab")+descStyle.Render(" Next field"),
		keyStyle.Render("/")+descStyle.Render(" Filter  ")+keyStyle.Render("Ctrl+C")+descStyle.Render(" Force quit"),
	)

	// Features section
	featSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("✨ Features")

	features := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render(
			"• Multi-account management with secure keyring storage\n" +
				"• Domain listing with filtering support\n" +
				"• Advanced cache purging (URL, hostname, tag, prefix)\n" +
				"• Full zone cache purge capability",
		)

	// Auth section
	authSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("🔐 Authentication")

	authInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("API Token")+" "+
			lipgloss.NewStyle().Foreground(MutedColor).Render("(Recommended)"),
		lipgloss.NewStyle().Foreground(MutedColor).Render("  dash.cloudflare.com/profile/api-tokens"),
		"",
		lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("Global API Key")+" "+
			lipgloss.NewStyle().Foreground(MutedColor).Render("(Legacy)"),
		lipgloss.NewStyle().Foreground(MutedColor).Render("  Full account access, use tokens when possible"),
	)

	// Links
	linksSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("🔗 Links")

	links := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render(
			"• API Docs: developers.cloudflare.com/api/\n" +
				"• GitHub: github.com/siyamsarker/cfctl",
		)

	// Footer
	prompt := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().
			Background(AccentColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render("Enter"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" or "),
		lipgloss.NewStyle().
			Background(BorderColor).
			Foreground(TextColor).
			Padding(0, 1).
			Render("Esc"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" to return"),
	)

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		divider,
		"",
		keySection,
		shortcuts,
		"",
		featSection,
		features,
		"",
		authSection,
		authInfo,
		"",
		linksSection,
		links,
		"",
		divider,
		prompt,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

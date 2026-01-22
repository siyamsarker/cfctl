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
	// Responsive sizing
	dividerWidth := min(m.width-8, 62)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header
	title := MakeSectionHeader("❓", "Help & Documentation", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	// Keyboard shortcuts card
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

	shortcutsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		keyStyle.Render("↑↓")+descStyle.Render(" Navigate")+"   "+keyStyle.Render("Enter")+descStyle.Render(" Select"),
		keyStyle.Render("Esc")+descStyle.Render(" Back")+"  "+keyStyle.Render("q")+descStyle.Render(" Quit")+"  "+keyStyle.Render("Tab")+descStyle.Render(" Next field"),
		keyStyle.Render("/")+descStyle.Render(" Filter")+"  "+keyStyle.Render("Ctrl+C")+descStyle.Render(" Force quit"),
	)

	shortcutsCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, keySection, "", shortcutsContent))

	// Features card
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

	featuresCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, featSection, "", features))

	// Auth card
	authSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("🔐 Authentication")

	authInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("API Token")+" "+
			lipgloss.NewStyle().Foreground(MutedColor).Render("(Recommended)"),
		lipgloss.NewStyle().Foreground(MutedColor).Render("dash.cloudflare.com/profile/api-tokens"),
		"",
		lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("Global API Key")+" "+
			lipgloss.NewStyle().Foreground(MutedColor).Render("(Legacy)"),
		lipgloss.NewStyle().Foreground(MutedColor).Render("Full account access, use tokens when possible"),
	)

	authCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, authSection, "", authInfo))

	// Links card
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

	linksCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, linksSection, "", links))

	// Modern footer
	footerHints := []KeyHint{
		{Key: "Enter", Description: "Return", IsAction: true},
		{Key: "Esc", Description: "Back", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	// Combine all sections with left alignment
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		shortcutsCard,
		"",
		featuresCard,
		"",
		authCard,
		"",
		linksCard,
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		footer,
	)

	// Polished container with responsive sizing
	containerWidth := min(m.width-10, 68)
	if containerWidth < 54 {
		containerWidth = 54
	}

	// Adjust padding based on available height
	verticalPadding := 1
	if m.height < 30 {
		verticalPadding = 0
	}

	container := lipgloss.NewStyle().
		Width(containerWidth).
		Padding(verticalPadding, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
	)
}

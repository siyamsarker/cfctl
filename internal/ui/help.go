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
	header := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render("Documentation"),
		SubtitleStyle.Render("System Help & Shortcuts"),
	)

	// Shortcuts Section
	shortcuts := lipgloss.JoinVertical(lipgloss.Left,
		SectionTitleStyle.Render("⌨️  Shortcuts"),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Width(30).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					row("↑/↓", "Navigate"),
					row("Enter", "Select/Confirm"),
				),
			),
			lipgloss.NewStyle().Width(30).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					row("Esc", "Back"),
					row("q", "Quit"),
				),
			),
		),
	)

	// Auth Section
	auth := lipgloss.JoinVertical(lipgloss.Left,
		SectionTitleStyle.Render("🔐 Authentication"),
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Width(12).Render("API Token"),
				lipgloss.NewStyle().Foreground(MutedColor).Render("Recommended. Supports fine-grained permissions."),
			),
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Width(12).Render("Global Key"),
				lipgloss.NewStyle().Foreground(MutedColor).Render("Legacy. Full account access (use with caution)."),
			),
		),
	)

	// Features Section
	features := lipgloss.JoinVertical(lipgloss.Left,
		SectionTitleStyle.Render("✨ Capabilities"),
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Multi-account management with secure keyring"),
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Advanced cache purging (Tag, Prefix, Host)"),
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Domain management and filtering"),
	)

	// Assemble Content
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		MakeDivider(min(m.width-10, 60)),
		shortcuts,
		"",
		auth,
		"",
		features,
		"",
	)

	// Container
	container := ContainerStyle.
		Width(min(m.width-4, 70)).
		Render(content)

	// Footer
	footer := MakeFooter([]KeyHint{
		{Key: "Esc", Description: "Return", IsAction: true},
	})

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			container,
			lipgloss.NewStyle().MarginTop(1).Render(footer),
		),
	)
}

func row(key, desc string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left,
		KeyStyle.Render(key),
		lipgloss.NewStyle().Foreground(SubTextColor).PaddingLeft(1).Render(desc),
	)
}

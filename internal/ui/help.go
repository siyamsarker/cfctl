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
			if menu, ok := m.returnTo.(MainMenuModel); ok {
				menu.applySize(m.width, m.height)
				return menu, nil
			}
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

	// Shortcuts Section - Clean 2-column layout
	shortcutsTitle := SectionTitleStyle.Render("Shortcuts")

	col1 := lipgloss.JoinVertical(lipgloss.Left,
		row("↑/↓", "Navigate"),
		row("Enter", "Select/Confirm"),
		row("/", "Filter"),
	)

	col2 := lipgloss.JoinVertical(lipgloss.Left,
		row("Esc", "Back"),
		row("q", "Quit"),
		row("Tab", "Next Field"),
	)

	shortcuts := lipgloss.JoinVertical(lipgloss.Left,
		shortcutsTitle,
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(35).Render(col1),
			lipgloss.NewStyle().Width(35).Render(col2),
		),
	)

	// Auth Section
	authTitle := SectionTitleStyle.Render("Authentication")
	auth := lipgloss.JoinVertical(lipgloss.Left,
		authTitle,
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Foreground(SuccessColor).Width(15).Render("API Token"),
			lipgloss.NewStyle().Foreground(SubTextColor).Render("Recommended. Supports fine-grained permissions."),
		),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Foreground(WarningColor).Width(15).Render("Global Key"),
			lipgloss.NewStyle().Foreground(SubTextColor).Render("Legacy. Full account access."),
		),
	)

	// Features Section
	featTitle := SectionTitleStyle.Render("Capabilities")
	features := lipgloss.JoinVertical(lipgloss.Left,
		featTitle,
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Multi-account management with secure keyring"),
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Advanced cache purging (Tag, Prefix, Host)"),
		lipgloss.NewStyle().Foreground(SubTextColor).Render("• Domain management and filtering"),
	)

	// Assemble Content
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		MakeDivider(min(m.width-4, 60)),
		shortcuts,
		"",
		auth,
		"",
		features,
	)

	// Container
	container := ContainerStyle.
		Width(min(m.width-2, 70)).
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
	// Fixed width key column for alignment
	k := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Width(10).
		Render(key)

	d := lipgloss.NewStyle().
		Foreground(SubTextColor).
		Render(desc)

	return lipgloss.JoinHorizontal(lipgloss.Left, k, d)
}

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
	title := TitleStyle.Render("Help & Documentation")

	help := `
CFCTL - Cloudflare CLI Management Tool

KEYBOARD SHORTCUTS:
  ↑/↓, j/k     Navigate menus
  Enter        Select/Confirm
  Esc, q       Back/Cancel
  Ctrl+C       Quit application
  Tab          Next field
  Shift+Tab    Previous field

FEATURES:
  • Configure multiple Cloudflare accounts
  • Switch between accounts easily
  • List and manage domains/zones
  • Purge cache by URL, hostname, tag, prefix
  • Purge entire zone cache
  • Secure credential storage

AUTHENTICATION:
  API Token (Recommended):
    Create at: dash.cloudflare.com/profile/api-tokens
    Required permissions: Zone:Read, Cache Purge:Purge

  Global API Key:
    Found at: dash.cloudflare.com/profile/api-tokens
    Less secure, use API tokens when possible

DOCUMENTATION:
  Cloudflare API: developers.cloudflare.com/api/
  Project GitHub: github.com/siyamsarker/cfctl

VERSION: 1.0.0
`

	content := lipgloss.NewStyle().
		Width(70).
		Padding(1, 2).
		Render(help)

	prompt := HelpStyle.Render("Press any key to return...")

	full := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		prompt,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(full),
	)
}

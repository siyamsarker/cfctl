package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type SettingsModel struct {
	config   *config.Config
	returnTo tea.Model
	width    int
	height   int
}

func NewSettingsModel(cfg *config.Config) SettingsModel {
	return SettingsModel{
		config:   cfg,
		returnTo: NewMainMenuModel(cfg),
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m SettingsModel) View() string {
	title := TitleStyle.Render("Settings")

	content := lipgloss.NewStyle().
		Width(60).
		Padding(1, 2).
		Render(
			"Theme: " + m.config.Defaults.Theme + "\n" +
				"Output: " + m.config.Defaults.Output + "\n" +
				"API Timeout: " + lipgloss.NewStyle().Render(fmt.Sprintf("%d seconds", m.config.API.Timeout)) + "\n" +
				"API Retries: " + lipgloss.NewStyle().Render(fmt.Sprintf("%d", m.config.API.Retries)) + "\n" +
				"Confirmations: " + lipgloss.NewStyle().Render(fmt.Sprintf("%t", m.config.UI.Confirmations)) + "\n" +
				"Animations: " + lipgloss.NewStyle().Render(fmt.Sprintf("%t", m.config.UI.Animations)) + "\n" +
				"Colors: " + lipgloss.NewStyle().Render(fmt.Sprintf("%t", m.config.UI.Colors)) + "\n" +
				"Cache Enabled: " + lipgloss.NewStyle().Render(fmt.Sprintf("%t", m.config.Cache.Enabled)) + "\n" +
				"Cache TTL: " + lipgloss.NewStyle().Render(fmt.Sprintf("%d seconds", m.config.Cache.DomainsTTL)),
		)

	prompt := HelpStyle.Render("\nPress any key to return...")

	full := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		prompt,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(full),
	)
}

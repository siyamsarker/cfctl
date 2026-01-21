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
		width:    80,
		height:   24,
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
	// Responsive sizing
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header
	title := MakeSectionHeader("⚙️", " Settings", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	// Settings card
	cardWidth := min(m.width-10, 50)
	if cardWidth < 35 {
		cardWidth = 35
	}

	// Helper for settings rows
	settingRow := func(label, value string, active bool) string {
		labelStyle := lipgloss.NewStyle().Foreground(MutedColor).Width(16)
		valueStyle := lipgloss.NewStyle().Foreground(TextColor)
		if active {
			valueStyle = valueStyle.Foreground(SuccessColor).Bold(true)
		}
		return labelStyle.Render(label) + valueStyle.Render(value)
	}

	boolVal := func(b bool) string {
		if b {
			return "✓ Enabled"
		}
		return "✗ Disabled"
	}

	// General section
	generalSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("General")

	generalSettings := lipgloss.JoinVertical(
		lipgloss.Center,
		settingRow("Theme:", m.config.Defaults.Theme, false),
		settingRow("Output:", m.config.Defaults.Output, false),
	)

	// API section
	apiSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("API")

	apiSettings := lipgloss.JoinVertical(
		lipgloss.Center,
		settingRow("Timeout:", fmt.Sprintf("%ds", m.config.API.Timeout), false),
		settingRow("Retries:", fmt.Sprintf("%d", m.config.API.Retries), false),
	)

	// UI section
	uiSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("UI")

	uiSettings := lipgloss.JoinVertical(
		lipgloss.Center,
		settingRow("Confirmations:", boolVal(m.config.UI.Confirmations), m.config.UI.Confirmations),
		settingRow("Animations:", boolVal(m.config.UI.Animations), m.config.UI.Animations),
		settingRow("Colors:", boolVal(m.config.UI.Colors), m.config.UI.Colors),
	)

	// Cache section
	cacheSection := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("Cache")

	cacheSettings := lipgloss.JoinVertical(
		lipgloss.Center,
		settingRow("Enabled:", boolVal(m.config.Cache.Enabled), m.config.Cache.Enabled),
		settingRow("TTL:", fmt.Sprintf("%ds", m.config.Cache.DomainsTTL), false),
	)

	// Combine settings card
	settingsCard := ProfessionalCardStyle.Copy().
		Width(cardWidth).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				generalSection,
				generalSettings,
				"",
				apiSection,
				apiSettings,
				"",
				uiSection,
				uiSettings,
				"",
				cacheSection,
				cacheSettings,
			),
		)

	// Note
	note := lipgloss.NewStyle().
		Foreground(MutedColor).
		Italic(true).
		Render("Edit ~/.cfctl/config.yaml to modify settings")

	// Modern footer
	footerHints := []KeyHint{
		{Key: "Enter", Description: "Return", IsAction: true},
		{Key: "Esc", Description: "Back", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		settingsCard,
		"",
		note,
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		footer,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

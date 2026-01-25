package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type AccountItem struct {
	account   config.Config
	name      string
	email     string
	isDefault bool
}

func (i AccountItem) Title() string {
	prefix := "  "
	if i.isDefault {
		prefix = "✓ "
	}
	return prefix + i.name
}

func (i AccountItem) Description() string {
	if i.email == "" {
		return "Token authentication"
	}
	return i.email
}

func (i AccountItem) FilterValue() string {
	return i.name
}

type AccountSelectModel struct {
	config *config.Config
	list   list.Model
	width  int
	height int
}

func NewAccountSelectModel(cfg *config.Config) AccountSelectModel {
	items := make([]list.Item, len(cfg.Accounts))
	for i, acc := range cfg.Accounts {
		items[i] = AccountItem{
			name:      acc.Name,
			email:     acc.Email,
			isDefault: acc.Default,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(0, 0, 0, 2)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(AccentColor).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(TextColor).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(MutedColor).
		Padding(0, 0, 0, 2)

	l := list.New(items, delegate, 60, 12)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return AccountSelectModel{
		config: cfg,
		list:   l,
		width:  80,
		height: 24,
	}
}

func (m AccountSelectModel) Init() tea.Cmd {
	return nil
}

func (m AccountSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := min(msg.Width-10, 60)
		listHeight := min(msg.Height-10, 12)
		if listWidth < 40 {
			listWidth = 40
		}
		if listHeight < 6 {
			listHeight = 6
		}
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			menu := NewMainMenuModel(m.config)
			menu.applySize(m.width, m.height)
			return menu, nil
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(AccountItem)
				if err := m.config.SetDefaultAccount(item.name); err == nil {
					menu := NewMainMenuModel(m.config)
					menu.applySize(m.width, m.height)
					return NewMessageModel(
						"Success",
						fmt.Sprintf("Default account set to: %s", item.name),
						SuccessColor,
						menu,
					), nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AccountSelectModel) View() string {
	// Responsive sizing
	dividerWidth := min(m.width-8, 50)
	if dividerWidth < 25 {
		dividerWidth = 25
	}

	// Modern header
	title := MakeSectionHeader("👤", "Select Account", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	// Enhanced current account badge
	var currentBadge string
	for _, acc := range m.config.Accounts {
		if acc.Default {
			currentBadge = lipgloss.JoinHorizontal(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(MutedColor).Render("Current: "),
				SuccessStatusBadge.Render(acc.Name),
			)
			break
		}
	}

	// Modern footer with keyboard hints
	footerHints := []KeyHint{
		{Key: "↑↓", Description: "Navigate", IsAction: false},
		{Key: "Enter", Description: "Select", IsAction: true},
		{Key: "/", Description: "Filter", IsAction: false},
		{Key: "Esc", Description: "Back", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	// Build content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		currentBadge,
		"",
		m.list.View(),
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

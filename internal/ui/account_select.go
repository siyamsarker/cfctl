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
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(PrimaryColor)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Select Cloudflare Account"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle

	return AccountSelectModel{
		config: cfg,
		list:   l,
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
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return NewMainMenuModel(m.config), nil
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(AccountItem)
				if err := m.config.SetDefaultAccount(item.name); err == nil {
					return NewMessageModel(
						"Success",
						fmt.Sprintf("Default account set to: %s", item.name),
						NewMainMenuModel(m.config),
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
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type MenuItem struct {
	title       string
	description string
	action      string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

type MainMenuModel struct {
	list   list.Model
	config *config.Config
	width  int
	height int
}

func NewMainMenuModel(cfg *config.Config) MainMenuModel {
	items := []list.Item{
		MenuItem{title: "Configure Cloudflare Account", description: "Add or manage API credentials", action: "configure"},
		MenuItem{title: "Select Cloudflare Account", description: "Switch between accounts", action: "select"},
		MenuItem{title: "Manage Domains", description: "View and select domains", action: "domains"},
		MenuItem{title: "Settings", description: "Configure application settings", action: "settings"},
		MenuItem{title: "Help", description: "View help and documentation", action: "help"},
		MenuItem{title: "Exit", description: "Exit application", action: "exit"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(PrimaryColor)

	l := list.New(items, delegate, 80, 20)
	l.Title = "CFCTL - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return MainMenuModel{
		list:   l,
		config: cfg,
		width:  80,
		height: 20,
	}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(MenuItem)
			switch selected.action {
			case "configure":
				return NewAccountConfigModel(m.config), nil
			case "select":
				if len(m.config.Accounts) == 0 {
					return NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Cloudflare Account' option.",
						m,
					), nil
				}
				return NewAccountSelectModel(m.config), nil
			case "domains":
				if len(m.config.Accounts) == 0 {
					return NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Cloudflare Account' option.",
						m,
					), nil
				}
				return NewDomainListModel(m.config), nil
			case "settings":
				return NewSettingsModel(m.config), nil
			case "help":
				return NewHelpModel(m), nil
			case "exit":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

// MessageModel displays a simple message with an OK button
type MessageModel struct {
	title    string
	message  string
	returnTo tea.Model
	width    int
	height   int
}

func NewMessageModel(title, message string, returnTo tea.Model) MessageModel {
	return MessageModel{
		title:    title,
		message:  message,
		returnTo: returnTo,
	}
}

func (m MessageModel) Init() tea.Cmd {
	return nil
}

func (m MessageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", "q":
			return m.returnTo, nil
		}
	}
	return m, nil
}

func (m MessageModel) View() string {
	title := TitleStyle.Render(m.title)
	message := lipgloss.NewStyle().
		Width(60).
		Padding(1, 2).
		Render(m.message)
	prompt := HelpStyle.Render("Press Enter to continue...")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		"",
		prompt,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(content),
	)
}

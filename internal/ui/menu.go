package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type MenuItem struct {
	title       string
	description string
	action      string
	icon        string
}

func (i MenuItem) Title() string       { return i.icon + " " + i.title }
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
		MenuItem{title: "Configure Account", description: "Add or manage API credentials", action: "configure", icon: "🔧"},
		MenuItem{title: "Select Account", description: "Switch between configured accounts", action: "select", icon: "👤"},
		MenuItem{title: "Remove Account", description: "Delete a configured account", action: "remove", icon: "🗑️"},
		MenuItem{title: "Manage Domains", description: "View and manage your domains", action: "domains", icon: "🌐"},
		MenuItem{title: "Settings", description: "Configure application preferences", action: "settings", icon: "⚙️"},
		MenuItem{title: "Help", description: "View documentation and shortcuts", action: "help", icon: "❓"},
		MenuItem{title: "Exit", description: "Close the application", action: "exit", icon: "🚪"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle.Copy().
		Foreground(PrimaryColor).
		Padding(0, 0)
	delegate.Styles.SelectedDesc = SelectedMenuItemStyle.Copy().
		Foreground(PrimaryDim).
		Padding(0, 0)
	delegate.Styles.NormalTitle = MenuItemStyle.Copy().Padding(0, 0)
	delegate.Styles.NormalDesc = MenuItemStyle.Copy().Foreground(MutedColor).Padding(0, 0)

	// Cleaner spacing
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 65, 18)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return MainMenuModel{
		list:   l,
		config: cfg,
		width:  80,
		height: 24,
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

		// Responsive sizing
		listWidth := min(msg.Width-10, 70)
		listHeight := min(msg.Height-10, 20)
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(MenuItem)
			switch selected.action {
			case "configure":
				model := NewAccountConfigModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "select":
				if len(m.config.Accounts) == 0 {
					return m.showMessage("No Accounts", "Please configure an account first.", WarningColor)
				}
				model := NewAccountSelectModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "remove":
				if len(m.config.Accounts) == 0 {
					return m.showMessage("No Accounts", "There are no accounts to remove.", WarningColor)
				}
				model := NewAccountRemoveModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "domains":
				if len(m.config.Accounts) == 0 {
					return m.showMessage("No Accounts", "Please configure an account first.", WarningColor)
				}
				domainModel := NewDomainListModel(m.config)
				domainModel.width = m.width
				domainModel.height = m.height
				return domainModel, domainModel.Init()
			case "settings":
				model := NewSettingsModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "help":
				model := NewHelpModel(m)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "exit":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) showMessage(title, desc string, color lipgloss.Color) (tea.Model, tea.Cmd) {
	msgModel := NewMessageModel(title, desc, m)
	msgModel.width = m.width
	msgModel.height = m.height
	return msgModel, nil
}

func (m MainMenuModel) View() string {
	// Header Section
	header := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render("CFCTL"),
		SubtitleStyle.Render("Advanced Cloudflare Controller"),
	)

	// Status Section
	var statusBadge string
	if len(m.config.Accounts) > 0 {
		defaultAcc, err := m.config.GetDefaultAccount()
		accName := "Unknown"
		if err == nil {
			accName = defaultAcc.Name
		}
		statusBadge = lipgloss.JoinHorizontal(lipgloss.Left,
			SuccessBadgeStyle.Render("ACTIVE"),
			lipgloss.NewStyle().Foreground(SubTextColor).PaddingLeft(1).Render(accName),
		)
	} else {
		statusBadge = WarningBadgeStyle.Render("NO ACCOUNT CONFIGURED")
	}

	// Content assembly
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		MakeDivider(min(m.width-10, 60)),
		"",
		statusBadge,
		"",
		m.list.View(),
		"",
	)

	// Container
	container := ContainerStyle.
		Width(min(m.width-4, 70)).
		Render(content)

	// Footer
	footer := MakeFooter([]KeyHint{
		{Key: "↑/↓", Description: "Navigate"},
		{Key: "Enter", Description: "Select", IsAction: true},
		{Key: "q", Description: "Quit"},
	})

	// Full Screen Layout
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			container,
			lipgloss.NewStyle().MarginTop(1).Render(footer),
		),
	)
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
		width:    80,
		height:   24,
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
	// Determine style based on title (heuristic)
	var tint lipgloss.Color = InfoColor
	if strings.Contains(strings.ToLower(m.title), "error") {
		tint = ErrorColor
	} else if strings.Contains(strings.ToLower(m.title), "warning") {
		tint = WarningColor
	} else if strings.Contains(strings.ToLower(m.title), "success") {
		tint = SuccessColor
	}

	title := lipgloss.NewStyle().Foreground(tint).Bold(true).Render(m.title)
	desc := lipgloss.NewStyle().Foreground(TextColor).Render(m.message)

	btn := ActionKeyStyle.Copy().
		Background(tint).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render("OK")

	card := ContainerStyle.
		BorderForeground(tint).
		Padding(1, 3).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			title,
			"",
			desc,
			"",
			btn,
		))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		card,
	)
}

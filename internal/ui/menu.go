package ui

import (
	"fmt"

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

func (i MenuItem) Title() string       { return i.icon + "  " + i.title } // Added extra space for proper icon separation
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
		MenuItem{title: "Configure Account", description: "", action: "configure", icon: "⚙"},
		MenuItem{title: "Select Account", description: "", action: "select", icon: "◉"},
		MenuItem{title: "Remove Account", description: "", action: "remove", icon: "✕"},
		MenuItem{title: "Manage Domains", description: "", action: "domains", icon: "◈"},
		MenuItem{title: "Settings", description: "", action: "settings", icon: "◐"},
		MenuItem{title: "Help", description: "", action: "help", icon: "?"},
		MenuItem{title: "Exit", description: "", action: "exit", icon: "→"},
	}

	delegate := list.NewDefaultDelegate()
	// Clean selection style (left border indicator)
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle.Copy()
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Height(0) // Hide descriptions

	// Normal style
	delegate.Styles.NormalTitle = MenuItemStyle.Copy()
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Height(0) // Hide descriptions

	// Spacing
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	l := list.New(items, delegate, 65, 18)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return MainMenuModel{
		list:   l,
		config: cfg,
		width:  0,
		height: 0,
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

		listWidth := min(msg.Width-4, 70)
		listHeight := min(msg.Height-8, 20)
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
	msgModel := NewMessageModel(title, desc, color, m)
	msgModel.width = m.width
	msgModel.height = m.height
	return msgModel, nil
}

func (m MainMenuModel) View() string {
	// Don't render until we have terminal dimensions
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Header Section - Clean and professional
	header := lipgloss.JoinVertical(lipgloss.Center,
		TitleStyle.Render("CFCTL"),
		SubtitleStyle.Render("Cloudflare Management Console"),
	)

	// Account Status Card - More informative
	var accountInfo string
	if len(m.config.Accounts) > 0 {
		defaultAcc, err := m.config.GetDefaultAccount()
		accName := "Unknown"
		if err == nil {
			accName = defaultAcc.Name
		}

		// Account count badge
		accountCount := lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Render(fmt.Sprintf("%d", len(m.config.Accounts)))

		accountLabel := lipgloss.NewStyle().
			Foreground(MutedColor).
			Render("Active: ")

		accountName := lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true).
			Render(accName)

		accountInfo = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Foreground(SuccessColor).Render("● "),
				accountLabel,
				accountName,
			),
			lipgloss.NewStyle().Foreground(MutedColor).Render(fmt.Sprintf("  %s configured account(s)", accountCount)),
		)
	} else {
		accountInfo = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("⚠ No Account Configured"),
			lipgloss.NewStyle().Foreground(MutedColor).Render("  Select 'Configure Account' to get started"),
		)
	}

	// Account info card with subtle styling
	accountCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Width(min(m.width-8, 62)).
		Align(lipgloss.Left).
		Render(accountInfo)

	// Menu section header
	menuHeader := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		Render("MAIN MENU")

	// Content assembly with better spacing
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		MakeDivider(min(m.width-4, 60)),
		"",
		accountCard,
		"",
		menuHeader,
		"",
		m.list.View(),
	)

	// Container
	container := ContainerStyle.
		Width(min(m.width-2, 70)).
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

// MessageModel methods remain largely the same, just clean styling
type MessageModel struct {
	title    string
	message  string
	color    lipgloss.Color
	returnTo tea.Model
	width    int
	height   int
}

func NewMessageModel(title, message string, color lipgloss.Color, returnTo tea.Model) MessageModel {
	return MessageModel{
		title:    title,
		message:  message,
		color:    color,
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
	tint := m.color

	title := lipgloss.NewStyle().Foreground(tint).Bold(true).Render(m.title)
	desc := lipgloss.NewStyle().Foreground(TextColor).Render(m.message)

	btn := ActionKeyStyle.Copy().Render("OK [Enter]")

	card := ContainerStyle.
		BorderForeground(tint).
		Padding(1, 4).
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

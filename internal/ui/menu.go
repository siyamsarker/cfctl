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
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(0, 0, 0, 3)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(AccentColor).
		Padding(0, 0, 0, 3)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(TextColor).
		Padding(0, 0, 0, 3)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(MutedColor).
		Padding(0, 0, 0, 3)

	// Increased spacing for better visual clarity
	delegate.SetSpacing(1)

	l := list.New(items, delegate, 65, 22)
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

		// Make list responsive - ensure all 6 items fit on one page
		listWidth := min(msg.Width-10, 75)
		listHeight := min(msg.Height-10, 22) // Increased for better spacing
		if listWidth < 50 {
			listWidth = 50
		}
		if listHeight < 18 {
			listHeight = 18 // Minimum height with proper spacing
		}
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
					msgModel := NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Account' option.",
						m,
					)
					msgModel.width = m.width
					msgModel.height = m.height
					return msgModel, nil
				}
				model := NewAccountSelectModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "remove":
				if len(m.config.Accounts) == 0 {
					msgModel := NewMessageModel(
						"No Accounts Configured",
						"There are no accounts to remove.",
						m,
					)
					msgModel.width = m.width
					msgModel.height = m.height
					return msgModel, nil
				}
				model := NewAccountRemoveModel(m.config)
				model.width = m.width
				model.height = m.height
				return model, nil
			case "domains":
				if len(m.config.Accounts) == 0 {
					msgModel := NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Account' option.",
						m,
					)
					msgModel.width = m.width
					msgModel.height = m.height
					return msgModel, nil
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

func (m MainMenuModel) View() string {
	// Responsive sizing
	dividerWidth := min(m.width-8, 60)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header with improved styling
	header := MakeSectionHeader("CFCTL", "", "Main Menu")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	// Enhanced account status badge with better visual hierarchy
	var statusBadge string
	if len(m.config.Accounts) > 0 {
		defaultAcc, err := m.config.GetDefaultAccount()
		accName := "Unknown"
		if err == nil {
			accName = defaultAcc.Name
		}

		statusBadge = lipgloss.JoinHorizontal(
			lipgloss.Left,
			SuccessStatusBadge.Render("✓ Active"),
			lipgloss.NewStyle().
				Foreground(TextColor).
				Padding(0, 1).
				Render(accName),
		)
	} else {
		statusBadge = WarningStatusBadge.Render("⚠ No Account")
	}

	// Modern footer with keyboard hints using helper
	footerHints := []KeyHint{
		{Key: "↑↓", Description: "Navigate", IsAction: false},
		{Key: "Enter", Description: "Select", IsAction: true},
		{Key: "q", Description: "Quit", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	// Build complete view with professional spacing
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		statusBadge,
		"",
		"",
		m.list.View(),
		"",
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		footer,
	)

	// Polished container for a modern card-like layout
	containerWidth := min(m.width-10, 72)
	if containerWidth < 54 {
		containerWidth = 54
	}
	container := lipgloss.NewStyle().
		Width(containerWidth).
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Render(content)

	// Center in terminal
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
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
	// Icon based on content
	var icon string
	var borderColor lipgloss.Color

	switch {
	case strings.Contains(strings.ToLower(m.title), "error"):
		icon = "✗"
		borderColor = ErrorColor
	case strings.Contains(strings.ToLower(m.title), "success"):
		icon = "✓"
		borderColor = SuccessColor
	case strings.Contains(strings.ToLower(m.title), "warning"):
		icon = "⚠"
		borderColor = WarningColor
	default:
		icon = "ℹ"
		borderColor = InfoColor
	}

	// Title with icon
	title := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().
			Foreground(borderColor).
			Bold(true).
			Render(icon+" "),
		lipgloss.NewStyle().
			Foreground(borderColor).
			Bold(true).
			Render(m.title),
	)

	// Responsive message card
	cardWidth := min(m.width-20, 50)
	if cardWidth < 30 {
		cardWidth = 30
	}

	messageCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(cardWidth).
		Render(
			lipgloss.NewStyle().
				Foreground(TextColor).
				Render(m.message),
		)

	// Continue prompt
	prompt := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().
			Background(SuccessColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render("Enter"),
		lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(" to continue"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		messageCard,
		"",
		prompt,
	)

	// Polished container
	containerWidth := min(m.width-10, 58)
	if containerWidth < 38 {
		containerWidth = 38
	}
	container := lipgloss.NewStyle().
		Width(containerWidth).
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
	)
}

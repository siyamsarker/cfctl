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
				return NewAccountConfigModel(m.config), nil
			case "select":
				if len(m.config.Accounts) == 0 {
					return NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Account' option.",
						m,
					), nil
				}
				return NewAccountSelectModel(m.config), nil
			case "domains":
				if len(m.config.Accounts) == 0 {
					return NewMessageModel(
						"No Accounts Configured",
						"Please configure an account first using 'Configure Account' option.",
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
	// Responsive header divider
	dividerWidth := min(m.width-8, 60)
	if dividerWidth < 30 {
		dividerWidth = 30
	}
	divider := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Render(repeatStr("─", dividerWidth))

	// Title with branding
	title := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("CFCTL")

	subtitle := lipgloss.NewStyle().
		Foreground(MutedColor).
		Render(" Main Menu")

	header := lipgloss.JoinHorizontal(lipgloss.Left, title, subtitle)

	// Account status badge
	var statusBadge string
	if len(m.config.Accounts) > 0 {
		defaultAcc, err := m.config.GetDefaultAccount()
		accName := "Unknown"
		if err == nil {
			accName = defaultAcc.Name
		}

		statusBadge = lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().
				Background(SuccessColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				Padding(0, 1).
				Render("✓ Active"),
			lipgloss.NewStyle().
				Foreground(TextColor).
				Padding(0, 1).
				Render(accName),
		)
	} else {
		statusBadge = lipgloss.NewStyle().
			Background(WarningColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render("⚠ No Account")
	}

	// Footer with keyboard shortcuts
	keys := []struct {
		key  string
		desc string
	}{
		{"↑↓", "Navigate"},
		{"Enter", "Select"},
		{"q", "Quit"},
	}

	var keyHints []string
	for _, k := range keys {
		keyHints = append(keyHints,
			lipgloss.NewStyle().
				Background(BorderColor).
				Foreground(TextColor).
				Padding(0, 1).
				Render(k.key)+
				lipgloss.NewStyle().
					Foreground(MutedColor).
					Render(" "+k.desc),
		)
	}

	footer := strings.Join(keyHints, "  ")

	// Build complete view with improved spacing
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.NewStyle().Foreground(MutedColor).Render(divider),
		"",
		statusBadge,
		"",
		"",
		m.list.View(),
		"",
		"",
		lipgloss.NewStyle().Foreground(MutedColor).Render(divider),
		"",
		footer,
	)

	// Container with increased padding for a cleaner look
	container := lipgloss.NewStyle().
		Padding(2, 4).
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
			Foreground(AccentColor).
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
		lipgloss.Center,
		lipgloss.NewStyle().
			Background(AccentColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render("Enter"),
		lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(" to continue"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		messageCard,
		"",
		prompt,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

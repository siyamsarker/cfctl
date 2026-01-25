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

func (m *MainMenuModel) applySize(width, height int) {
	m.width = width
	m.height = height

	listWidth := min(width-4, 70)
	itemsHeight := len(m.list.Items())*2 - 1
	availableHeight := height - 6
	if availableHeight < 5 {
		availableHeight = 5
	}
	listHeight := min(availableHeight, itemsHeight)

	m.list.SetWidth(listWidth)
	m.list.SetHeight(listHeight)
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

	// Keep items compact to avoid pagination
	delegate.SetSpacing(1)
	delegate.ShowDescription = false

	l := list.New(items, delegate, 70, 22)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	// Tighten initial height to item count to avoid extra empty space
	initialItemsHeight := len(items)
	if initialItemsHeight < 5 {
		initialItemsHeight = 5
	}
	l.SetHeight(initialItemsHeight)

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
		m.applySize(msg.Width, msg.Height)
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
	// ASCII Logo - Compact version
	logo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("CFCTL")

	// For wider terminals, use a smaller ASCII art
	if m.width >= 70 {
		logo = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Render(`
   ██████╗███████╗ ██████╗████████╗██╗     
  ██╔════╝██╔════╝██╔════╝╚══██╔══╝██║     
  ██║     █████╗  ██║        ██║   ██║     
  ██║     ██╔══╝  ██║        ██║   ██║     
  ╚██████╗██║     ╚██████╗   ██║   ███████╗
   ╚═════╝╚═╝      ╚═════╝   ╚═╝   ╚══════╝`)
	}

	// Header Section - Clean and professional
	header := lipgloss.JoinVertical(lipgloss.Center,
		logo,
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
		accountCount := fmt.Sprintf("%d", len(m.config.Accounts))
		accountPlural := "account"
		if len(m.config.Accounts) > 1 {
			accountPlural = "accounts"
		}

		statusLine := lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("● "),
			lipgloss.NewStyle().Foreground(TextColor).Bold(true).Render(accName),
		)

		countLine := lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingLeft(2).
			Render(fmt.Sprintf("%s %s configured", accountCount, accountPlural))

		accountInfo = lipgloss.JoinVertical(lipgloss.Left, statusLine, countLine)
	} else {
		warningLine := lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("⚠ "),
			lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Render("No Account Configured"),
		)

		hintLine := lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingLeft(2).
			Render("Select 'Configure Account' to get started")

		accountInfo = lipgloss.JoinVertical(lipgloss.Left, warningLine, hintLine)
	}

	// Account info card with refined styling
	accountCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Width(min(m.width-10, 60)).
		Align(lipgloss.Left).
		Render(accountInfo)

	// Menu section header with refined styling
	menuHeader := lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		PaddingLeft(1).
		Render("MAIN MENU")

	// Content assembly with optimized spacing
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		MakeDivider(min(m.width-6, 60)),
		"",
		accountCard,
		"",
		menuHeader,
		"", // Gap after MAIN MENU
		m.list.View(),
	)

	// Container with optimized padding
	container := ContainerStyle.
		Width(min(m.width-4, 70)).
		Padding(1, 2).
		Render(content)

	// Footer
	footer := MakeFooter([]KeyHint{
		{Key: "↑/↓", Description: "Navigate"},
		{Key: "Enter", Description: "Select", IsAction: true},
		{Key: "q", Description: "Quit"},
	})

	// Full Screen Layout - Centered like welcome screen
	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		container,
		"",
		footer,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
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
			if menu, ok := m.returnTo.(MainMenuModel); ok {
				menu.applySize(m.width, m.height)
				return menu, nil
			}
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

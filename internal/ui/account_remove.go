package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
)

type RemoveAccountItem struct {
	name      string
	email     string
	isDefault bool
}

func (i RemoveAccountItem) Title() string {
	prefix := "  "
	if i.isDefault {
		prefix = "★ "
	}
	return prefix + i.name
}

func (i RemoveAccountItem) Description() string {
	if i.email == "" {
		return "Token authentication"
	}
	return i.email
}

func (i RemoveAccountItem) FilterValue() string {
	return i.name
}

type AccountRemoveModel struct {
	config      *config.Config
	list        list.Model
	width       int
	height      int
	confirmMode bool
	selected    string
	err         error
}

func NewAccountRemoveModel(cfg *config.Config) AccountRemoveModel {
	items := make([]list.Item, len(cfg.Accounts))
	for i, acc := range cfg.Accounts {
		items[i] = RemoveAccountItem{
			name:      acc.Name,
			email:     acc.Email,
			isDefault: acc.Default,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true).
		Padding(0, 0, 0, 2)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(MutedColor).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(TextColor).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(MutedColor).
		Padding(0, 0, 0, 2)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 60, 14)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return AccountRemoveModel{
		config: cfg,
		list:   l,
		width:  80,
		height: 24,
	}
}

func (m AccountRemoveModel) Init() tea.Cmd {
	return nil
}

func (m AccountRemoveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := min(msg.Width-10, 60)
		listHeight := 14
		if listWidth < 40 {
			listWidth = 40
		}
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, nil

	case tea.KeyMsg:
		if m.confirmMode {
			switch msg.String() {
			case "y", "Y":
				// Remove account and credential
				if err := config.DeleteCredential(m.selected); err != nil {
					m.err = err
					m.confirmMode = false
					return m, nil
				}
				if err := m.config.RemoveAccount(m.selected); err != nil {
					m.err = err
					m.confirmMode = false
					return m, nil
				}
				return NewMessageModel(
					"Success",
					fmt.Sprintf("Account '%s' has been removed.", m.selected),
					SuccessColor,
					NewMainMenuModel(m.config),
				), nil
			case "n", "N", "esc":
				m.confirmMode = false
				m.selected = ""
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "esc", "q":
			return NewMainMenuModel(m.config), nil
		case "enter", "d":
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(RemoveAccountItem)
				m.selected = item.name
				m.confirmMode = true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AccountRemoveModel) View() string {
	// Responsive sizing
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 25 {
		dividerWidth = 25
	}

	// Modern header
	title := MakeSectionHeader("🗑️", " Remove Account", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	if m.confirmMode {
		// Enhanced confirmation dialog with warning styling
		confirmCard := WarningCardStyle.Copy().
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Render("⚠ Confirm Deletion"),
					"",
					lipgloss.NewStyle().Foreground(TextColor).Render(fmt.Sprintf("Are you sure you want to remove '%s'?", m.selected)),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Render("This will delete the stored credentials."),
				),
			)

		// Modern footer for confirmation
		footerHints := []KeyHint{
			{Key: "Y", Description: "Yes", IsAction: false},
			{Key: "N", Description: "No", IsAction: true},
		}
		footer := MakeFooter(footerHints)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			confirmCard,
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)

		// Polished container
		containerWidth := min(m.width-10, 58)
		if containerWidth < 48 {
			containerWidth = 48
		}
		container := lipgloss.NewStyle().
			Width(containerWidth).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Render(content)

		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			container,
		)
	}

	// Error display
	var errDisplay string
	if m.err != nil {
		errDisplay = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Render("Error: " + m.err.Error())
	}

	// Instruction
	instruction := lipgloss.NewStyle().
		Foreground(WarningColor).
		Italic(true).
		Render("Select an account to remove")

	// Modern footer
	footerHints := []KeyHint{
		{Key: "↑↓", Description: "Navigate", IsAction: false},
		{Key: "Enter", Description: "Remove", IsAction: true},
		{Key: "Esc", Description: "Back", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		instruction,
		"",
		m.list.View(),
		errDisplay,
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		footer,
	)

	// Polished container
	containerWidth := min(m.width-10, 66)
	if containerWidth < 54 {
		containerWidth = 54
	}
	container := lipgloss.NewStyle().
		Width(containerWidth).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		container,
	)
}

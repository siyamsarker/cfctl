package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/api"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

type DomainItem struct {
	zone cloudflare.Zone
}

func (i DomainItem) Title() string {
	return i.zone.Name
}

func (i DomainItem) Description() string {
	return fmt.Sprintf("Status: %s | Plan: %s", i.zone.Status, i.zone.Plan.Name)
}

func (i DomainItem) FilterValue() string {
	return i.zone.Name
}

type DomainListModel struct {
	config  *config.Config
	list    list.Model
	zones   []cloudflare.Zone
	loading bool
	err     error
	width   int
	height  int
}

type zonesLoadedMsg struct {
	zones []cloudflare.Zone
	err   error
}

func NewDomainListModel(cfg *config.Config) DomainListModel {
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

	l := list.New([]list.Item{}, delegate, 60, 12)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return DomainListModel{
		config:  cfg,
		list:    l,
		loading: true,
		width:   80,
		height:  24,
	}
}

func (m DomainListModel) loadZones() tea.Msg {
	// Get default account
	account, err := m.config.GetDefaultAccount()
	if err != nil {
		return zonesLoadedMsg{err: err}
	}

	// Get credential
	credential, err := config.GetCredential(account.Name)
	if err != nil {
		return zonesLoadedMsg{err: fmt.Errorf("failed to get credentials: %w", err)}
	}

	// Create API client
	var cfg api.ClientConfig
	if account.AuthType == "token" {
		cfg = api.ClientConfig{
			APIToken: credential,
			Timeout:  m.config.API.Timeout,
			Retries:  m.config.API.Retries,
		}
	} else {
		cfg = api.ClientConfig{
			APIKey:  credential,
			Email:   account.Email,
			Timeout: m.config.API.Timeout,
			Retries: m.config.API.Retries,
		}
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return zonesLoadedMsg{err: err}
	}

	// List zones
	ctx := context.Background()
	zones, err := client.ListZones(ctx)
	if err != nil {
		return zonesLoadedMsg{err: err}
	}

	return zonesLoadedMsg{zones: zones}
}

func (m DomainListModel) Init() tea.Cmd {
	return m.loadZones
}

func (m DomainListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := min(msg.Width-10, 65)
		listHeight := min(msg.Height-12, 12)
		if listWidth < 40 {
			listWidth = 40
		}
		if listHeight < 6 {
			listHeight = 6
		}
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, nil

	case zonesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.zones = msg.zones
		items := make([]list.Item, len(msg.zones))
		for i, zone := range msg.zones {
			items[i] = DomainItem{zone: zone}
		}
		m.list.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			if msg.String() == "esc" || msg.String() == "q" {
				return NewMainMenuModel(m.config), nil
			}
			return m, nil
		}

		switch msg.String() {
		case "esc", "q":
			return NewMainMenuModel(m.config), nil
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(DomainItem)
				return NewPurgeMenuModel(m.config, item.zone), nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m DomainListModel) View() string {
	// Header
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 25 {
		dividerWidth = 25
	}
	divider := lipgloss.NewStyle().
		Foreground(BorderColor).
		Render(repeatStr("─", dividerWidth))

	if m.loading {
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Width(min(m.width-10, 45)).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render("◐ Loading Domains..."),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Render("Fetching zones from Cloudflare API"),
				),
			)

		keys := lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().
				Background(BorderColor).
				Foreground(TextColor).
				Padding(0, 1).
				Render("Esc"),
			lipgloss.NewStyle().Foreground(MutedColor).Render(" Cancel"),
		)

		content := lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("🌐 Domains"),
			divider,
			"",
			loadingCard,
			"",
			divider,
			keys,
		)

		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	if m.err != nil {
		errorCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorColor).
			Padding(1, 2).
			Width(min(m.width-10, 50)).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Render("✗ Error Loading Domains"),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Width(min(m.width-20, 45)).Render(m.err.Error()),
				),
			)

		keys := lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().
				Background(BorderColor).
				Foreground(TextColor).
				Padding(0, 1).
				Render("Esc"),
			lipgloss.NewStyle().Foreground(MutedColor).Render(" Return to menu"),
		)

		content := lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("🌐 Domains"),
			divider,
			"",
			errorCard,
			"",
			divider,
			keys,
		)

		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Normal view with domain list
	title := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("🌐 Domains")

	// Account and count badge
	var infoBadge string
	account, _ := m.config.GetDefaultAccount()
	if account != nil {
		infoBadge = lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().
				Background(AccentColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				Padding(0, 1).
				Render(fmt.Sprintf("%d zones", len(m.zones))),
			lipgloss.NewStyle().Foreground(MutedColor).Render("  "),
			lipgloss.NewStyle().Foreground(MutedColor).Render("Account: "),
			lipgloss.NewStyle().Foreground(TextColor).Bold(true).Render(account.Name),
		)
	}

	// Footer
	keys := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().
			Background(BorderColor).
			Foreground(TextColor).
			Padding(0, 1).
			Render("↑↓"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" Navigate  "),
		lipgloss.NewStyle().
			Background(SuccessColor).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Render("Enter"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" Purge  "),
		lipgloss.NewStyle().
			Background(BorderColor).
			Foreground(TextColor).
			Padding(0, 1).
			Render("/"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" Filter  "),
		lipgloss.NewStyle().
			Background(BorderColor).
			Foreground(TextColor).
			Padding(0, 1).
			Render("Esc"),
		lipgloss.NewStyle().Foreground(MutedColor).Render(" Back"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		divider,
		"",
		infoBadge,
		"",
		m.list.View(),
		"",
		divider,
		keys,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

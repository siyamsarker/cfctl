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
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(PrimaryColor)

	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "Domains"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = TitleStyle

	return DomainListModel{
		config:  cfg,
		list:    l,
		loading: true,
		width:   80,
		height:  20,
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
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
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

		// Update title with count
		account, _ := m.config.GetDefaultAccount()
		if account != nil {
			m.list.Title = fmt.Sprintf("Domains (%d) - Account: %s", len(msg.zones), account.Name)
		}
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
	if m.loading {
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			TitleStyle.Render("Loading Domains..."),
			"",
			InfoStyle.Render("⠋ Fetching zones from Cloudflare..."),
			"",
			HelpStyle.Render("Press Esc to cancel"),
		)
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			BorderStyle.Render(content),
		)
	}

	if m.err != nil {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			ErrorStyle.Render("Error Loading Domains"),
			"",
			lipgloss.NewStyle().Width(60).Render(m.err.Error()),
			"",
			HelpStyle.Render("Press any key to return to menu"),
		)
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			BorderStyle.Render(content),
		)
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

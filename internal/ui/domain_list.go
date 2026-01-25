package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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
	spinner spinner.Model
	zones   []cloudflare.Zone
	loading bool
	err     error
	width   int
	height  int
	status  string
}

type zonesLoadedMsg struct {
	zones []cloudflare.Zone
	err   error
}

type zonesTimeoutMsg struct{}

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

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = SpinnerStyle

	l := list.New([]list.Item{}, delegate, 60, 12)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return DomainListModel{
		config:  cfg,
		list:    l,
		spinner: sp,
		loading: true,
		status:  "Authenticating...",
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

	// Verify token first (update status)
	// We can't update visual status here directly without a cmd, but we can assume success if we proceed
	// or we properly chain msgs. For now, we just proceed to ListZones which is the main blocking call.

	// List zones
	ctx := context.Background()
	zones, err := client.ListZones(ctx)
	if err != nil {
		return zonesLoadedMsg{err: err}
	}

	return zonesLoadedMsg{zones: zones}
}

func (m DomainListModel) Init() tea.Cmd {
	return tea.Batch(m.loadZones, m.zonesTimeout(), m.spinner.Tick)
}

func (m DomainListModel) zonesTimeout() tea.Cmd {
	return tea.Tick(20*time.Second, func(time.Time) tea.Msg {
		return zonesTimeoutMsg{}
	})
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
		if len(m.zones) == 0 {
			m.err = fmt.Errorf("no zones found. Ensure your API token has Zone.Zone.Read permission and access to at least one zone")
			return m, nil
		}
		items := make([]list.Item, len(msg.zones))
		for i, zone := range msg.zones {
			items[i] = DomainItem{zone: zone}
		}
		m.list.SetItems(items)
		return m, nil

	case zonesTimeoutMsg:
		if m.loading {
			m.loading = false
			m.err = fmt.Errorf("timeout fetching zones. Check network connectivity and ensure your API token has Zone.Zone.Read permission")
			return m, nil
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			if msg.String() == "esc" || msg.String() == "q" {
				model := NewMainMenuModel(m.config)
				model.applySize(m.width, m.height)
				return model, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "esc", "q":
			model := NewMainMenuModel(m.config)
			model.applySize(m.width, m.height)
			return model, nil
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(DomainItem)
				model := NewPurgeMenuModel(m.config, item.zone)
				model.width = m.width
				model.height = m.height
				return model, nil
			}
		}
	default:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m DomainListModel) View() string {
	// Responsive sizing
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 25 {
		dividerWidth = 25
	}

	// Modern header
	title := MakeSectionHeader("🌐", "Domains", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	if m.loading {
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render(fmt.Sprintf("%s Loading Domains...", m.spinner.View())),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Render("Fetching zones from Cloudflare API"),
				),
			)

		footerHints := []KeyHint{
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			loadingCard,
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

	if m.err != nil {
		errorCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorColor).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Render("✗ Error Loading Domains"),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Render(m.err.Error()),
				),
			)

		footerHints := []KeyHint{
			{Key: "Esc", Description: "Return to menu", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			errorCard,
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)

		// Polished container
		containerWidth := min(m.width-10, 62)
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

	// Account and count badge
	var infoBadge string
	account, _ := m.config.GetDefaultAccount()
	if account != nil {
		infoBadge = lipgloss.JoinHorizontal(
			lipgloss.Left,
			InfoStatusBadge.Render(fmt.Sprintf("%d zones", len(m.zones))),
			lipgloss.NewStyle().Foreground(MutedColor).Render("  "),
			lipgloss.NewStyle().Foreground(MutedColor).Render("Account: "),
			lipgloss.NewStyle().Foreground(TextColor).Bold(true).Render(account.Name),
		)
	}

	// Modern footer
	footerHints := []KeyHint{
		{Key: "↑↓", Description: "Navigate", IsAction: false},
		{Key: "Enter", Description: "Purge", IsAction: true},
		{Key: "/", Description: "Filter", IsAction: false},
		{Key: "Esc", Description: "Back", IsAction: false},
	}
	footer := MakeFooter(footerHints)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		"",
		infoBadge,
		"",
		m.list.View(),
		"",
		lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
		footer,
	)

	// Polished container
	containerWidth := min(m.width-10, 68)
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

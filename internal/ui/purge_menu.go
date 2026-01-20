package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

type PurgeMenuItem struct {
	title       string
	description string
	purgeType   string
	icon        string
}

func (i PurgeMenuItem) Title() string       { return i.icon + " " + i.title }
func (i PurgeMenuItem) Description() string { return i.description }
func (i PurgeMenuItem) FilterValue() string { return i.title }

type PurgeMenuModel struct {
	config *config.Config
	zone   cloudflare.Zone
	list   list.Model
	width  int
	height int
}

func NewPurgeMenuModel(cfg *config.Config, zone cloudflare.Zone) PurgeMenuModel {
	items := []list.Item{
		PurgeMenuItem{
			title:       "Purge by URL",
			description: "Purge specific URLs (exact match)",
			purgeType:   "url",
			icon:        "🔗",
		},
		PurgeMenuItem{
			title:       "Purge by Hostname",
			description: "Purge all assets for a hostname",
			purgeType:   "hostname",
			icon:        "🌐",
		},
		PurgeMenuItem{
			title:       "Purge by Tag",
			description: "Purge assets with Cache-Tag headers",
			purgeType:   "tag",
			icon:        "🏷️",
		},
		PurgeMenuItem{
			title:       "Purge by Prefix",
			description: "Purge assets under a path/prefix",
			purgeType:   "prefix",
			icon:        "📁",
		},
		PurgeMenuItem{
			title:       "Purge Everything",
			description: "Clear entire cache (use with caution)",
			purgeType:   "everything",
			icon:        "🗑️",
		},
		PurgeMenuItem{
			title:       "Back",
			description: "Return to domain list",
			purgeType:   "back",
			icon:        "←",
		},
	}

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
	// Compact spacing - no extra space between items
	delegate.SetSpacing(0)

	// Height needs to accommodate 6 items * 2 lines each = 12 lines minimum
	l := list.New(items, delegate, 60, 18)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return PurgeMenuModel{
		config: cfg,
		zone:   zone,
		list:   l,
		width:  80,
		height: 24,
	}
}

func (m PurgeMenuModel) Init() tea.Cmd {
	return nil
}

func (m PurgeMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Make list height accommodate all 6 items (2 lines each = 12 lines)
		listWidth := min(msg.Width-10, 60)
		listHeight := 18 // Fixed height to show all items
		if listWidth < 40 {
			listWidth = 40
		}
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			domainModel := NewDomainListModel(m.config)
			return domainModel, domainModel.Init()
		case "enter":
			selected := m.list.SelectedItem().(PurgeMenuItem)
			switch selected.purgeType {
			case "url":
				return NewPurgeByURLModel(m.config, m.zone), nil
			case "hostname":
				return NewPurgeByHostnameModel(m.config, m.zone), nil
			case "tag":
				return NewPurgeByTagModel(m.config, m.zone), nil
			case "prefix":
				return NewPurgeByPrefixModel(m.config, m.zone), nil
			case "everything":
				return NewPurgeEverythingModel(m.config, m.zone), nil
			case "back":
				domainModel := NewDomainListModel(m.config)
				return domainModel, domainModel.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PurgeMenuModel) View() string {
	// Header
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 30 {
		dividerWidth = 30
	}
	divider := lipgloss.NewStyle().
		Foreground(BorderColor).
		Render(repeatStr("─", dividerWidth))

	title := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("🗑️  Cache Purge")

	// Zone badge
	zoneBadge := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(MutedColor).Render("Zone: "),
		lipgloss.NewStyle().
			Background(AccentColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render(m.zone.Name),
	)

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
		lipgloss.NewStyle().Foreground(MutedColor).Render(" Select  "),
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
		zoneBadge,
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

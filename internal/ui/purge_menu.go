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
}

func (i PurgeMenuItem) Title() string       { return i.title }
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
		},
		PurgeMenuItem{
			title:       "Purge by Hostname",
			description: "Purge all assets for a hostname",
			purgeType:   "hostname",
		},
		PurgeMenuItem{
			title:       "Purge by Tag",
			description: "Purge assets with Cache-Tag headers (Enterprise)",
			purgeType:   "tag",
		},
		PurgeMenuItem{
			title:       "Purge by Prefix",
			description: "Purge assets under a path/prefix",
			purgeType:   "prefix",
		},
		PurgeMenuItem{
			title:       "Purge Everything",
			description: "Clear entire cache (use with caution)",
			purgeType:   "everything",
		},
		PurgeMenuItem{
			title:       "Back to Domain Selection",
			description: "Return to domain list",
			purgeType:   "back",
		},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedMenuItemStyle
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(PrimaryColor)

	l := list.New(items, delegate, 80, 20)
	l.Title = "Cache Purge - " + zone.Name
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return PurgeMenuModel{
		config: cfg,
		zone:   zone,
		list:   l,
		width:  80,
		height: 20,
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
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return NewDomainListModel(m.config), nil
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
				return NewDomainListModel(m.config), nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PurgeMenuModel) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

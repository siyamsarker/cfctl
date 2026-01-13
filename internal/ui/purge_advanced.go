package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/api"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/internal/utils"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

// PurgeByTagModel
type PurgeByTagModel struct {
	config   *config.Config
	zone     cloudflare.Zone
	textarea textarea.Model
	err      error
	success  bool
	purging  bool
	width    int
	height   int
}

func NewPurgeByTagModel(cfg *config.Config, zone cloudflare.Zone) PurgeByTagModel {
	ta := textarea.New()
	ta.Placeholder = "Enter cache tags, one per line or comma-separated\nExample:\nheader-image\nfooter-content"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(70)
	ta.SetHeight(10)

	return PurgeByTagModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
	}
}

func (m PurgeByTagModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m PurgeByTagModel) executePurge() tea.Msg {
	input := m.textarea.Value()
	lines := strings.Split(input, "\n")
	var tags []string
	for _, line := range lines {
		parsed := utils.ParseCommaSeparated(line)
		tags = append(tags, parsed...)
	}

	if err := utils.ValidateTags(tags); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	account, err := m.config.GetDefaultAccount()
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	credential, err := config.GetCredential(account.Name)
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

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
		return purgeResultMsg{success: false, err: err}
	}

	ctx := context.Background()
	req := cloudflare.PurgeRequest{
		Tags: tags,
	}

	if err := client.PurgeCache(ctx, m.zone.ID, req); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	return purgeResultMsg{success: true}
}

func (m PurgeByTagModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case purgeResultMsg:
		m.purging = false
		if msg.success {
			m.success = true
			m.err = nil
		} else {
			m.err = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		if m.success {
			return NewPurgeMenuModel(m.config, m.zone), nil
		}

		switch msg.String() {
		case "esc":
			if !m.purging {
				return NewPurgeMenuModel(m.config, m.zone), nil
			}
		case "ctrl+s":
			if !m.purging && m.textarea.Value() != "" {
				m.purging = true
				m.err = nil
				return m, m.executePurge
			}
		}
	}

	if !m.purging && !m.success {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m PurgeByTagModel) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Purge by Tag - %s", m.zone.Name))
	warning := WarningStyle.Render("Note: Cache tag purging requires Enterprise plan")

	var content string
	if m.purging {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			InfoStyle.Render("⠋ Purging cache..."),
		)
	} else if m.success {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			SuccessStyle.Render("✓ Cache purged successfully!"),
			"",
			HelpStyle.Render("Press any key to return..."),
		)
	} else {
		instructions := MutedStyle.Render("Enter cache tags to purge (one per line or comma-separated)")

		var errorMsg string
		if m.err != nil {
			errorMsg = "\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		}

		help := HelpStyle.Render("\nCtrl+S: Submit | Esc: Cancel")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			warning,
			"",
			instructions,
			"",
			m.textarea.View(),
			errorMsg,
			help,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(content),
	)
}

// PurgeByPrefixModel
type PurgeByPrefixModel struct {
	config   *config.Config
	zone     cloudflare.Zone
	textarea textarea.Model
	err      error
	success  bool
	purging  bool
	width    int
	height   int
}

func NewPurgeByPrefixModel(cfg *config.Config, zone cloudflare.Zone) PurgeByPrefixModel {
	ta := textarea.New()
	ta.Placeholder = "Enter URL prefixes, one per line or comma-separated\nExample:\nhttps://example.com/images/\nhttps://example.com/static/"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(70)
	ta.SetHeight(10)

	return PurgeByPrefixModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
	}
}

func (m PurgeByPrefixModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m PurgeByPrefixModel) executePurge() tea.Msg {
	input := m.textarea.Value()
	lines := strings.Split(input, "\n")
	var prefixes []string
	for _, line := range lines {
		parsed := utils.ParseCommaSeparated(line)
		prefixes = append(prefixes, parsed...)
	}

	if err := utils.ValidatePrefixes(prefixes); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	account, err := m.config.GetDefaultAccount()
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	credential, err := config.GetCredential(account.Name)
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

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
		return purgeResultMsg{success: false, err: err}
	}

	ctx := context.Background()
	req := cloudflare.PurgeRequest{
		Prefixes: prefixes,
	}

	if err := client.PurgeCache(ctx, m.zone.ID, req); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	return purgeResultMsg{success: true}
}

func (m PurgeByPrefixModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case purgeResultMsg:
		m.purging = false
		if msg.success {
			m.success = true
			m.err = nil
		} else {
			m.err = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		if m.success {
			return NewPurgeMenuModel(m.config, m.zone), nil
		}

		switch msg.String() {
		case "esc":
			if !m.purging {
				return NewPurgeMenuModel(m.config, m.zone), nil
			}
		case "ctrl+s":
			if !m.purging && m.textarea.Value() != "" {
				m.purging = true
				m.err = nil
				return m, m.executePurge
			}
		}
	}

	if !m.purging && !m.success {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m PurgeByPrefixModel) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Purge by Prefix - %s", m.zone.Name))

	var content string
	if m.purging {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			InfoStyle.Render("⠋ Purging cache..."),
		)
	} else if m.success {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			SuccessStyle.Render("✓ Cache purged successfully!"),
			"",
			HelpStyle.Render("Press any key to return..."),
		)
	} else {
		instructions := MutedStyle.Render("Enter URL prefixes to purge (one per line or comma-separated)")

		var errorMsg string
		if m.err != nil {
			errorMsg = "\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		}

		help := HelpStyle.Render("\nCtrl+S: Submit | Esc: Cancel")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			instructions,
			"",
			m.textarea.View(),
			errorMsg,
			help,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(content),
	)
}

// PurgeEverythingModel
type PurgeEverythingModel struct {
	config  *config.Config
	zone    cloudflare.Zone
	input   textinput.Model
	step    int // 0: first confirm, 1: type domain name, 2: purging, 3: done
	err     error
	success bool
	width   int
	height  int
}

func NewPurgeEverythingModel(cfg *config.Config, zone cloudflare.Zone) PurgeEverythingModel {
	ti := textinput.New()
	ti.Placeholder = "Type domain name to confirm"
	ti.Focus()
	ti.Width = 50

	return PurgeEverythingModel{
		config: cfg,
		zone:   zone,
		input:  ti,
		step:   0,
	}
}

func (m PurgeEverythingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PurgeEverythingModel) executePurge() tea.Msg {
	account, err := m.config.GetDefaultAccount()
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	credential, err := config.GetCredential(account.Name)
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

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
		return purgeResultMsg{success: false, err: err}
	}

	ctx := context.Background()
	req := cloudflare.PurgeRequest{
		PurgeEverything: true,
	}

	if err := client.PurgeCache(ctx, m.zone.ID, req); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	return purgeResultMsg{success: true}
}

func (m PurgeEverythingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case purgeResultMsg:
		if msg.success {
			m.success = true
			m.err = nil
			m.step = 3
		} else {
			m.err = msg.err
			m.step = 1
		}
		return m, nil

	case tea.KeyMsg:
		if m.step == 3 {
			return NewPurgeMenuModel(m.config, m.zone), nil
		}

		switch msg.String() {
		case "esc":
			if m.step != 2 {
				return NewPurgeMenuModel(m.config, m.zone), nil
			}
		case "enter":
			if m.step == 0 {
				m.step = 1
				return m, textinput.Blink
			} else if m.step == 1 {
				if m.input.Value() == m.zone.Name {
					m.step = 2
					return m, m.executePurge
				} else {
					m.err = fmt.Errorf("domain name doesn't match")
					return m, nil
				}
			}
		case "n":
			if m.step == 0 {
				return NewPurgeMenuModel(m.config, m.zone), nil
			}
		case "y":
			if m.step == 0 {
				m.step = 1
				return m, textinput.Blink
			}
		}
	}

	if m.step == 1 {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m PurgeEverythingModel) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Purge Everything - %s", m.zone.Name))
	warning := WarningStyle.Render("⚠️  WARNING: This will clear ALL cached content for this domain!")

	var content string
	switch m.step {
	case 0:
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			warning,
			"",
			"This action will:",
			"• Clear all cached files",
			"• May impact website performance temporarily",
			"• Cannot be undone",
			"",
			"Are you sure you want to continue? (y/n)",
			"",
			HelpStyle.Render("Y: Yes | N/Esc: No"),
		)

	case 1:
		var errorMsg string
		if m.err != nil {
			errorMsg = ErrorStyle.Render(fmt.Sprintf("\n%v\n", m.err))
		}

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			warning,
			"",
			fmt.Sprintf("Type the domain name '%s' to confirm:", m.zone.Name),
			"",
			m.input.View(),
			errorMsg,
			"",
			HelpStyle.Render("Enter: Confirm | Esc: Cancel"),
		)

	case 2:
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			InfoStyle.Render("⠋ Purging entire cache..."),
			"",
			MutedStyle.Render("This may take a moment..."),
		)

	case 3:
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			SuccessStyle.Render("✓ All cache purged successfully!"),
			"",
			MutedStyle.Render("Cache will rebuild as visitors access your site."),
			"",
			HelpStyle.Render("Press any key to return..."),
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(content),
	)
}

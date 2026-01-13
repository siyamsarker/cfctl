package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/api"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/internal/utils"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

type PurgeByURLModel struct {
	config   *config.Config
	zone     cloudflare.Zone
	textarea textarea.Model
	err      error
	success  bool
	purging  bool
	width    int
	height   int
}

type purgeResultMsg struct {
	success bool
	err     error
}

func NewPurgeByURLModel(cfg *config.Config, zone cloudflare.Zone) PurgeByURLModel {
	ta := textarea.New()
	ta.Placeholder = "Enter URLs, one per line or comma-separated\nExample:\nhttps://example.com/style.css\nhttps://example.com/script.js"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(70)
	ta.SetHeight(10)

	return PurgeByURLModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
	}
}

func (m PurgeByURLModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m PurgeByURLModel) executePurge() tea.Msg {
	// Parse URLs
	input := m.textarea.Value()
	lines := strings.Split(input, "\n")
	var urls []string
	for _, line := range lines {
		parsed := utils.ParseCommaSeparated(line)
		urls = append(urls, parsed...)
	}

	// Validate URLs
	if err := utils.ValidateURLs(urls); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	// Get account and credentials
	account, err := m.config.GetDefaultAccount()
	if err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	credential, err := config.GetCredential(account.Name)
	if err != nil {
		return purgeResultMsg{success: false, err: err}
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
		return purgeResultMsg{success: false, err: err}
	}

	// Execute purge
	ctx := context.Background()
	req := cloudflare.PurgeRequest{
		Files: urls,
	}

	if err := client.PurgeCache(ctx, m.zone.ID, req); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	return purgeResultMsg{success: true}
}

func (m PurgeByURLModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m PurgeByURLModel) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Purge by URL - %s", m.zone.Name))

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
		instructions := MutedStyle.Render("Enter URLs to purge (one per line or comma-separated)")

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

// Similar models for other purge types
type PurgeByHostnameModel struct {
	config   *config.Config
	zone     cloudflare.Zone
	textarea textarea.Model
	err      error
	success  bool
	purging  bool
	width    int
	height   int
}

func NewPurgeByHostnameModel(cfg *config.Config, zone cloudflare.Zone) PurgeByHostnameModel {
	ta := textarea.New()
	ta.Placeholder = "Enter hostnames, one per line or comma-separated\nExample:\nwww.example.com\napi.example.com"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(70)
	ta.SetHeight(10)

	return PurgeByHostnameModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
	}
}

func (m PurgeByHostnameModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m PurgeByHostnameModel) executePurge() tea.Msg {
	input := m.textarea.Value()
	lines := strings.Split(input, "\n")
	var hostnames []string
	for _, line := range lines {
		parsed := utils.ParseCommaSeparated(line)
		hostnames = append(hostnames, parsed...)
	}

	if err := utils.ValidateHostnames(hostnames); err != nil {
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
		Hosts: hostnames,
	}

	if err := client.PurgeCache(ctx, m.zone.ID, req); err != nil {
		return purgeResultMsg{success: false, err: err}
	}

	return purgeResultMsg{success: true}
}

func (m PurgeByHostnameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m PurgeByHostnameModel) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Purge by Hostname - %s", m.zone.Name))

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
		instructions := MutedStyle.Render("Enter hostnames to purge (one per line or comma-separated)")

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

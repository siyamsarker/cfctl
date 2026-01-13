package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/siyamsarker/cfctl/internal/api"
	"github.com/siyamsarker/cfctl/internal/config"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

type AccountConfigModel struct {
	config     *config.Config
	inputs     []textinput.Model
	focusIndex int
	authType   string // "token" or "key"
	step       int    // 0: auth type, 1: inputs, 2: verifying, 3: done
	err        error
	verified   bool
	width      int
	height     int
}

func NewAccountConfigModel(cfg *config.Config) AccountConfigModel {
	m := AccountConfigModel{
		config:   cfg,
		authType: "token",
		step:     0,
	}
	m.initInputs()
	return m
}

func (m *AccountConfigModel) initInputs() {
	m.inputs = make([]textinput.Model, 3)

	// Account name
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "My Account"
	m.inputs[0].Focus()
	m.inputs[0].CharLimit = 50
	m.inputs[0].Width = 50
	m.inputs[0].Prompt = "Account Name: "

	// Email (for global key) or token
	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "user@example.com"
	m.inputs[1].CharLimit = 100
	m.inputs[1].Width = 50

	// API Key or Token
	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "Your API token or key"
	m.inputs[2].CharLimit = 200
	m.inputs[2].Width = 50
	m.inputs[2].EchoMode = textinput.EchoPassword
	m.inputs[2].EchoCharacter = '•'

	m.updateInputLabels()
}

func (m *AccountConfigModel) updateInputLabels() {
	if m.authType == "token" {
		m.inputs[1].Placeholder = "Not needed for token auth"
		m.inputs[1].Prompt = "Email (optional): "
		m.inputs[2].Prompt = "API Token: "
		m.inputs[2].Placeholder = "Your API token"
	} else {
		m.inputs[1].Placeholder = "user@example.com"
		m.inputs[1].Prompt = "Email: "
		m.inputs[2].Prompt = "API Key: "
		m.inputs[2].Placeholder = "Your global API key"
	}
}

func (m AccountConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

type verifyMsg struct {
	success bool
	err     error
}

func (m AccountConfigModel) verifyCredentials() tea.Msg {
	accountName := m.inputs[0].Value()
	email := m.inputs[1].Value()
	credential := m.inputs[2].Value()

	// Validate inputs
	if err := config.ValidateAccountName(accountName); err != nil {
		return verifyMsg{success: false, err: err}
	}

	var cfg api.ClientConfig
	if m.authType == "token" {
		if err := config.ValidateAPIToken(credential); err != nil {
			return verifyMsg{success: false, err: err}
		}
		cfg = api.ClientConfig{
			APIToken: credential,
			Timeout:  m.config.API.Timeout,
			Retries:  m.config.API.Retries,
		}
	} else {
		if err := config.ValidateEmail(email); err != nil {
			return verifyMsg{success: false, err: err}
		}
		if err := config.ValidateAPIKey(credential); err != nil {
			return verifyMsg{success: false, err: err}
		}
		cfg = api.ClientConfig{
			APIKey:  credential,
			Email:   email,
			Timeout: m.config.API.Timeout,
			Retries: m.config.API.Retries,
		}
	}

	// Create API client and verify
	client, err := api.NewClient(cfg)
	if err != nil {
		return verifyMsg{success: false, err: err}
	}

	ctx := context.Background()
	if err := client.VerifyToken(ctx); err != nil {
		return verifyMsg{success: false, err: fmt.Errorf("credential verification failed: %w", err)}
	}

	// Store credential in keyring
	if err := config.StoreCredential(accountName, credential); err != nil {
		return verifyMsg{success: false, err: fmt.Errorf("failed to store credential: %w", err)}
	}

	// Add account to config
	account := cloudflare.Account{
		Name:     accountName,
		Email:    email,
		AuthType: m.authType,
		Default:  len(m.config.Accounts) == 0,
	}

	if err := m.config.AddAccount(account); err != nil {
		return verifyMsg{success: false, err: err}
	}

	return verifyMsg{success: true, err: nil}
}

func (m AccountConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case verifyMsg:
		if msg.success {
			m.verified = true
			m.step = 3
		} else {
			m.err = msg.err
			m.step = 1
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == 0 || m.step == 1 {
				return NewMainMenuModel(m.config), nil
			}
			return m, nil

		case "enter":
			if m.step == 0 {
				// Move to input step
				m.step = 1
				m.inputs[0].Focus()
				return m, textinput.Blink
			} else if m.step == 1 {
				if m.focusIndex == len(m.inputs)-1 {
					// Submit form
					m.step = 2
					m.err = nil
					return m, m.verifyCredentials
				}
				// Move to next input
				m.focusIndex++
				return m, m.updateFocus()
			} else if m.step == 3 {
				// Done, return to menu
				return NewMainMenuModel(m.config), nil
			}

		case "tab", "shift+tab", "up", "down":
			if m.step == 0 {
				// Toggle auth type
				if m.authType == "token" {
					m.authType = "key"
				} else {
					m.authType = "token"
				}
				m.updateInputLabels()
				return m, nil
			} else if m.step == 1 {
				s := msg.String()
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs)-1 {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) - 1
				}

				return m, m.updateFocus()
			}
		}
	}

	// Handle character input
	if m.step == 1 {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	return m, nil
}

func (m *AccountConfigModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *AccountConfigModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m AccountConfigModel) View() string {
	var content string

	title := TitleStyle.Render("Configure Cloudflare Account")

	switch m.step {
	case 0:
		// Auth type selection
		var options strings.Builder
		options.WriteString("\nSelect authentication method:\n\n")

		if m.authType == "token" {
			options.WriteString(SelectedMenuItemStyle.Render("▸ API Token (Recommended)") + "\n")
			options.WriteString(MenuItemStyle.Render("  Global API Key") + "\n")
		} else {
			options.WriteString(MenuItemStyle.Render("  API Token (Recommended)") + "\n")
			options.WriteString(SelectedMenuItemStyle.Render("▸ Global API Key") + "\n")
		}

		options.WriteString("\n")
		options.WriteString(HelpStyle.Render("Use arrow keys to select, Enter to continue, Esc to cancel"))

		content = lipgloss.JoinVertical(lipgloss.Left, title, options.String())

	case 1:
		// Input form
		var inputs strings.Builder
		inputs.WriteString("\n")

		authTypeLabel := "API Token"
		if m.authType == "key" {
			authTypeLabel = "Global API Key"
		}
		inputs.WriteString(InfoStyle.Render(fmt.Sprintf("Authentication: %s\n\n", authTypeLabel)))

		for i, input := range m.inputs {
			inputs.WriteString(input.View())
			inputs.WriteString("\n")

			// Skip email field for token auth
			if m.authType == "token" && i == 1 {
				continue
			}
		}

		if m.err != nil {
			inputs.WriteString("\n")
			inputs.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		}

		inputs.WriteString("\n")
		inputs.WriteString(HelpStyle.Render("Tab/Shift+Tab: navigate • Enter: submit • Esc: cancel"))

		content = lipgloss.JoinVertical(lipgloss.Left, title, inputs.String())

	case 2:
		// Verifying
		spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"\n",
			InfoStyle.Render(string(spinner[0])+" Verifying credentials..."),
		)

	case 3:
		// Success
		var success strings.Builder
		success.WriteString("\n")
		success.WriteString(SuccessStyle.Render("✓ Account configured successfully!"))
		success.WriteString("\n\n")
		success.WriteString(fmt.Sprintf("Account: %s\n", m.inputs[0].Value()))
		if m.authType == "key" {
			success.WriteString(fmt.Sprintf("Email: %s\n", m.inputs[1].Value()))
		}
		success.WriteString(fmt.Sprintf("Auth Type: %s\n", m.authType))
		success.WriteString("\n")
		success.WriteString(HelpStyle.Render("Press Enter to return to main menu"))

		content = lipgloss.JoinVertical(lipgloss.Left, title, success.String())
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BorderStyle.Render(content),
	)
}

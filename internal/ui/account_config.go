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
		width:    80,
		height:   24,
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
			switch m.step {
			case 0:
				// Move to input step
				m.step = 1
				m.inputs[0].Focus()
				return m, textinput.Blink
			case 1:
				if m.focusIndex == len(m.inputs)-1 {
					// Submit form
					m.step = 2
					m.err = nil
					return m, m.verifyCredentials
				}
				// Move to next input
				m.focusIndex++
				return m, m.updateFocus()
			case 3:
				// Done, return to menu
				return NewMainMenuModel(m.config), nil
			}

		case "tab", "shift+tab", "up", "down":
			switch m.step {
			case 0:
				// Toggle auth type
				if m.authType == "token" {
					m.authType = "key"
				} else {
					m.authType = "token"
				}
				m.updateInputLabels()
				return m, nil
			case 1:
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

	// Responsive card width
	cardWidth := min(m.width-10, 55)
	if cardWidth < 35 {
		cardWidth = 35
	}

	// Modern header with responsive divider
	dividerWidth := min(m.width-8, 50)
	if dividerWidth < 25 {
		dividerWidth = 25
	}
	divider := lipgloss.NewStyle().
		Foreground(BorderColor).
		Render(repeatStr("─", dividerWidth))

	title := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("🔐 Configure Account")

	switch m.step {
	case 0:
		// Modern auth type selection with cards
		description := lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Render("Choose your authentication method")

		// Token card
		tokenSelected := m.authType == "token"
		tokenCard := m.buildAuthCard(
			"API Token",
			"✓ Recommended",
			"Scoped permissions • Easy to rotate • Secure",
			tokenSelected,
			cardWidth,
		)

		// Key card
		keySelected := m.authType == "key"
		keyCard := m.buildAuthCard(
			"Global API Key",
			"⚠ Full access",
			"Requires email • Complete control • Legacy",
			keySelected,
			cardWidth,
		)

		// Footer
		footerHints := []KeyHint{
			{Key: "↑↓", Description: "Select", IsAction: false},
			{Key: "Enter", Description: "Continue", IsAction: true},
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			description,
			"",
			tokenCard,
			"",
			keyCard,
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)

	case 1:
		// Modern input form
		authLabel := "API Token"
		if m.authType == "key" {
			authLabel = "Global API Key"
		}

		authBadge := lipgloss.NewStyle().
			Background(AccentColor).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1).
			Render("🔑 " + authLabel)

		// Build input fields
		var inputFields []string

		for i, input := range m.inputs {
			// Skip email for token auth
			if m.authType == "token" && i == 1 {
				continue
			}

			// Label
			label := lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true).
				Render(strings.TrimSuffix(input.Prompt, ": "))

			// Input styling
			inputStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Padding(0, 1).
				Width(cardWidth - 4)

			if m.focusIndex == i {
				inputStyle = inputStyle.BorderForeground(PrimaryColor)
			}

			field := lipgloss.JoinVertical(
				lipgloss.Left,
				label,
				inputStyle.Render(input.View()),
			)
			inputFields = append(inputFields, field)
		}

		// Error display
		if m.err != nil {
			errBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ErrorColor).
				Foreground(ErrorColor).
				Padding(0, 1).
				Render("✗ " + m.err.Error())
			inputFields = append(inputFields, errBox)
		}

		// Footer
		footerHints := []KeyHint{
			{Key: "Tab", Description: "Next", IsAction: false},
			{Key: "Enter", Description: "Submit", IsAction: true},
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			authBadge,
			"",
			lipgloss.JoinVertical(lipgloss.Left, inputFields...),
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)

	case 2:
		// Loading state with animation hint
		loadingIcon := lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Render("◐")

		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Width(cardWidth).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(
						lipgloss.Left,
						loadingIcon,
						lipgloss.NewStyle().Foreground(TextColor).Bold(true).Render(" Verifying credentials..."),
					),
					"",
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Connecting to Cloudflare API"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Validating authentication"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Storing secure credentials"),
				),
			)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			loadingCard,
		)

	case 3:
		// Success state
		successHeader := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().
				Background(SuccessColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				Padding(0, 1).
				Render("✓ Success"),
		)

		// Details card
		var details []string
		details = append(details,
			lipgloss.NewStyle().Foreground(MutedColor).Render("Name: ")+
				lipgloss.NewStyle().Foreground(TextColor).Bold(true).Render(m.inputs[0].Value()),
		)
		if m.authType == "key" && m.inputs[1].Value() != "" {
			details = append(details,
				lipgloss.NewStyle().Foreground(MutedColor).Render("Email: ")+
					lipgloss.NewStyle().Foreground(TextColor).Render(m.inputs[1].Value()),
			)
		}
		details = append(details,
			lipgloss.NewStyle().Foreground(MutedColor).Render("Auth: ")+
				lipgloss.NewStyle().
					Background(AccentColor).
					Foreground(lipgloss.Color("#000000")).
					Padding(0, 1).
					Render(m.authType),
		)

		detailsCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2).
			Render(lipgloss.JoinVertical(lipgloss.Left, details...))

		prompt := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().
				Background(SuccessColor).
				Foreground(lipgloss.Color("#000000")).
				Bold(true).
				Padding(0, 1).
				Render("Enter"),
			lipgloss.NewStyle().Foreground(MutedColor).Render(" to continue"),
		)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			successHeader,
			"",
			detailsCard,
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			prompt,
		)
	}

	// Polished container with responsive sizing
	containerWidth := min(m.width-10, 62)
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

// Helper to build auth selection cards
func (m AccountConfigModel) buildAuthCard(title, badge, description string, selected bool, width int) string {
	borderColor := BorderColor
	titleStyle := lipgloss.NewStyle().Foreground(AccentColor).Bold(true)

	if selected {
		borderColor = PrimaryColor
		// Add indicator inside the title
		title = "▸ " + title
		titleStyle = titleStyle.Foreground(PrimaryColor)
	}

	badgeColor := SuccessColor
	if strings.Contains(badge, "⚠") {
		badgeColor = WarningColor
	}

	cardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		lipgloss.NewStyle().Foreground(badgeColor).Render(badge),
		"",
		lipgloss.NewStyle().Foreground(MutedColor).Render(description),
	)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Render(cardContent)

	return card
}

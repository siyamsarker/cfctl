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
	ta.Placeholder = "Enter cache tags, one per line\nExample: header-image, footer-content"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(50)
	ta.SetHeight(8)

	return PurgeByTagModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
		width:    80,
		height:   24,
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
			model := NewPurgeMenuModel(m.config, m.zone)
			model.width = m.width
			model.height = m.height
			return model, nil
		}

		switch msg.String() {
		case "esc":
			if !m.purging {
				model := NewPurgeMenuModel(m.config, m.zone)
				model.width = m.width
				model.height = m.height
				return model, nil
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
	// Responsive sizing
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header
	title := MakeSectionHeader("🏷️", "Purge by Tag", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	zoneBadge := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(MutedColor).Render("Zone: "),
		InfoStatusBadge.Render(m.zone.Name),
	)

	warning := lipgloss.NewStyle().
		Foreground(WarningColor).
		Render("⚠ Enterprise plan required")

	var content string
	if m.purging {
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Render(
				lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render("◐ Purging cache..."),
			)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			loadingCard,
		)
	} else if m.success {
		successCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2).
			Render(
				lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("✓ Cache purged successfully!"),
			)

		prompt := lipgloss.NewStyle().Foreground(MutedColor).Render("Press any key to continue")

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			successCard,
			"",
			prompt,
		)
	} else {
		instructions := lipgloss.NewStyle().
			Foreground(MutedColor).
			Render("Enter cache tags (one per line or comma-separated)")

		// Resize textarea
		taWidth := min(m.width-15, 50)
		if taWidth < 30 {
			taWidth = 30
		}
		m.textarea.SetWidth(taWidth)

		var errorMsg string
		if m.err != nil {
			errorMsg = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ErrorColor).
				Foreground(ErrorColor).
				Padding(0, 1).
				Render("✗ " + m.err.Error())
		}

		// Modern footer
		footerHints := []KeyHint{
			{Key: "Ctrl+S", Description: "Submit", IsAction: true},
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			warning,
			"",
			instructions,
			"",
			m.textarea.View(),
			"",
			errorMsg,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
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
	ta.Placeholder = "Enter URL prefixes, one per line\nExample: https://example.com/images/"
	ta.Focus()
	ta.CharLimit = 0
	ta.SetWidth(50)
	ta.SetHeight(8)

	return PurgeByPrefixModel{
		config:   cfg,
		zone:     zone,
		textarea: ta,
		width:    80,
		height:   24,
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
			model := NewPurgeMenuModel(m.config, m.zone)
			model.width = m.width
			model.height = m.height
			return model, nil
		}

		switch msg.String() {
		case "esc":
			if !m.purging {
				model := NewPurgeMenuModel(m.config, m.zone)
				model.width = m.width
				model.height = m.height
				return model, nil
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
	// Header
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header
	title := MakeSectionHeader("📁", "Purge by Prefix", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	zoneBadge := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(MutedColor).Render("Zone: "),
		InfoStatusBadge.Render(m.zone.Name),
	)

	var content string
	if m.purging {
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Render(
				lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render("◐ Purging cache..."),
			)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			loadingCard,
		)
	} else if m.success {
		successCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2).
			Render(
				lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("✓ Cache purged successfully!"),
			)

		prompt := lipgloss.NewStyle().Foreground(MutedColor).Render("Press any key to continue")

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			successCard,
			"",
			prompt,
		)
	} else {
		instructions := lipgloss.NewStyle().
			Foreground(MutedColor).
			Render("Enter URL prefixes (one per line or comma-separated)")

		// Resize textarea
		taWidth := min(m.width-15, 50)
		if taWidth < 30 {
			taWidth = 30
		}
		m.textarea.SetWidth(taWidth)

		var errorMsg string
		if m.err != nil {
			errorMsg = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ErrorColor).
				Foreground(ErrorColor).
				Padding(0, 1).
				Render("✗ " + m.err.Error())
		}

		// Modern footer
		footerHints := []KeyHint{
			{Key: "Ctrl+S", Description: "Submit", IsAction: true},
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			instructions,
			"",
			m.textarea.View(),
			"",
			errorMsg,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
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
	ti.Width = 40

	return PurgeEverythingModel{
		config: cfg,
		zone:   zone,
		input:  ti,
		step:   0,
		width:  80,
		height: 24,
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
			model := NewPurgeMenuModel(m.config, m.zone)
			model.width = m.width
			model.height = m.height
			return model, nil
		}

		switch msg.String() {
		case "esc":
			if m.step != 2 {
				model := NewPurgeMenuModel(m.config, m.zone)
				model.width = m.width
				model.height = m.height
				return model, nil
			}
		case "enter":
			switch m.step {
			case 0:
				m.step = 1
				return m, textinput.Blink
			case 1:
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
				model := NewPurgeMenuModel(m.config, m.zone)
				model.width = m.width
				model.height = m.height
				return model, nil
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
	// Header
	dividerWidth := min(m.width-8, 55)
	if dividerWidth < 30 {
		dividerWidth = 30
	}

	// Modern header
	title := MakeSectionHeader("🗑️", "Purge Everything", "")
	divider := MakeDivider(dividerWidth, PrimaryColor)

	zoneBadge := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(MutedColor).Render("Zone: "),
		InfoStatusBadge.Render(m.zone.Name),
	)

	warning := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ErrorColor).
		Foreground(ErrorColor).
		Bold(true).
		Padding(0, 2).
		Render("⚠ WARNING: This clears ALL cached content!")

	var content string
	switch m.step {
	case 0:
		infoCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(WarningColor).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(TextColor).Render("This action will:"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Clear all cached files"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Impact performance temporarily"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("• Cannot be undone"),
				),
			)

		// Modern footer
		footerHints := []KeyHint{
			{Key: "Y", Description: "Continue", IsAction: true},
			{Key: "N", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			warning,
			"",
			infoCard,
			"",
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			footer,
		)

	case 1:
		confirmPrompt := lipgloss.NewStyle().
			Foreground(MutedColor).
			Render(fmt.Sprintf("Type '%s' to confirm:", m.zone.Name))

		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1).
			Render(m.input.View())

		var errorMsg string
		if m.err != nil {
			errorMsg = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Render("✗ " + m.err.Error())
		}

		// Modern footer
		footerHints := []KeyHint{
			{Key: "Enter", Description: "Confirm", IsAction: true},
			{Key: "Esc", Description: "Cancel", IsAction: false},
		}
		footer := MakeFooter(footerHints)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			warning,
			"",
			confirmPrompt,
			"",
			inputBox,
			"",
			errorMsg,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			footer,
		)

	case 2:
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					lipgloss.NewStyle().Foreground(AccentColor).Bold(true).Render("◐ Purging entire cache..."),
					lipgloss.NewStyle().Foreground(MutedColor).Render("This may take a moment"),
				),
			)

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			loadingCard,
		)

	case 3:
		successCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render("✓ Everything purged successfully!"),
					lipgloss.NewStyle().Foreground(MutedColor).Render("Cache will rebuild as visitors access your site"),
				),
			)

		prompt := lipgloss.NewStyle().Foreground(MutedColor).Render("Press any key to continue")

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().Foreground(BorderColor).Render(divider),
			"",
			zoneBadge,
			"",
			successCard,
			"",
			prompt,
		)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

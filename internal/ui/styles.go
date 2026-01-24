package ui

import "github.com/charmbracelet/lipgloss"

// Professional Color Palette
var (
	// Brand Colors - Using a sophisticated Cloudflare-inspired orange but deeper
	PrimaryColor = lipgloss.Color("#F48120") // Cloudflare Orange
	PrimaryDim   = lipgloss.Color("#C7610F") // Darker Orange

	// UI Colors - Dark, slate-based theme for professionalism
	BackgroundColor = lipgloss.Color("#0F172A") // Slate 900
	SurfaceColor    = lipgloss.Color("#1E293B") // Slate 800
	BorderColor     = lipgloss.Color("#334155") // Slate 700

	// Text Colors
	TextColor    = lipgloss.Color("#F8FAFC") // Slate 50
	SubTextColor = lipgloss.Color("#94A3B8") // Slate 400
	MutedColor   = lipgloss.Color("#64748B") // Slate 500

	// Status Colors - Muted, professional tones rather than neon
	SuccessColor = lipgloss.Color("#10B981") // Emerald 500
	ErrorColor   = lipgloss.Color("#EF4444") // Red 500
	WarningColor = lipgloss.Color("#F59E0B") // Amber 500
	InfoColor    = lipgloss.Color("#3B82F6") // Blue 500

	// Accents
	AccentColor = lipgloss.Color("#38BDF8") // Sky 400

	// Legacy / Compatibility Colors
	HighlightColor = lipgloss.Color("#334155") // Slate 700
	SecondaryColor = lipgloss.Color("#64748B") // Slate 500

)

// Global Layout Styles
var (
	// Base application style
	AppStyle = lipgloss.NewStyle().
			Margin(1, 2)

	// Standard container for content
	ContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Background(SurfaceColor)

	// Focused container
	ActiveContainerStyle = ContainerStyle.Copy().
				BorderForeground(PrimaryColor)

	// Helpers for legacy support / variants
	HighlightCardStyle = ContainerStyle.Copy().
				BorderForeground(AccentColor)

	WarningCardStyle = ContainerStyle.Copy().
				BorderForeground(WarningColor)
)

// Typography Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Italic(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true).
			PaddingBottom(1)

	SectionTitleStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true).
				MarginTop(1).
				MarginBottom(0)
)

// Component Styles
var (
	// Menu Items
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			PaddingLeft(2).
			PaddingRight(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(PrimaryColor).
				Bold(true).
				PaddingLeft(2).
				PaddingRight(2)

	// Input Fields
	InputStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	FocusedInputStyle = InputStyle.Copy().
				BorderForeground(PrimaryColor)

	// Badges & Status Indicators
	BadgeStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BorderColor).
			Padding(0, 1).
			Bold(true)

	SuccessBadgeStyle = BadgeStyle.Copy().
				Background(SuccessColor).
				Foreground(lipgloss.Color("#FFFFFF"))

	WarningBadgeStyle = BadgeStyle.Copy().
				Background(WarningColor).
				Foreground(lipgloss.Color("#000000"))

	ErrorBadgeStyle = BadgeStyle.Copy().
			Background(ErrorColor).
			Foreground(lipgloss.Color("#FFFFFF"))

	// Legacy badge support
	SuccessStatusBadge = SuccessBadgeStyle
	WarningStatusBadge = WarningBadgeStyle
	InfoStatusBadge    = BadgeStyle.Copy().Background(InfoColor).Foreground(lipgloss.Color("#FFFFFF"))

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().Foreground(AccentColor)
)

// Helper functions defined with legacy signature support

// MakeDivider creates a visual divider. The color argument is kept for compatibility but overridden or used as a fallback if needed.
func MakeDivider(width int, color ...lipgloss.Color) string {
	c := BorderColor
	if len(color) > 0 {
		c = color[0]
	}
	return lipgloss.NewStyle().
		Foreground(c).
		Render(repeatStr("─", width))
}

// KeyHint represents a keyboard shortcut hint
type KeyHint struct {
	Key         string
	Description string
	IsAction    bool
}

// Footer Styles
var (
	KeyStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Background(BorderColor).
			Padding(0, 1).
			Bold(true)

	ActionKeyStyle = KeyStyle.Copy().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(PrimaryDim)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingRight(2)
)

func MakeFooter(hints []KeyHint) string {
	var parts []string
	for _, hint := range hints {
		kStyle := KeyStyle
		if hint.IsAction {
			kStyle = ActionKeyStyle
		}
		parts = append(parts,
			kStyle.Render(hint.Key),
			KeyDescStyle.Render(" "+hint.Description),
		)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func repeatStr(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var result string
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// Legacy Compat Helpers

func MakeSectionHeader(icon, title, subtitle string) string {
	titleStyled := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render(icon + " " + title)

	if subtitle != "" {
		subtitleStyled := lipgloss.NewStyle().
			Foreground(SubTextColor).
			Render(" " + subtitle)
		return lipgloss.JoinHorizontal(lipgloss.Left, titleStyled, subtitleStyled)
	}
	return titleStyled
}

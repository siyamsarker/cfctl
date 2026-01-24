package ui

import "github.com/charmbracelet/lipgloss"

// Refined Professional Palette
var (
	// Brand Colors
	PrimaryColor = lipgloss.Color("#F48120") // Cloudflare Orange
	PrimaryDim   = lipgloss.Color("#C7610F") // Darker Orange

	// UI Colors - Minimalist Dark Theme
	// We avoid setting backgrounds explicitly to allow terminal transparency/theme to work naturally
	BorderColor  = lipgloss.Color("#475569") // Slate 600
	SubtleBorder = lipgloss.Color("#334155") // Slate 700

	// Text Colors
	TextColor    = lipgloss.Color("#F8FAFC") // Slate 50
	SubTextColor = lipgloss.Color("#94A3B8") // Slate 400
	MutedColor   = lipgloss.Color("#64748B") // Slate 500

	// Status Colors
	SuccessColor = lipgloss.Color("#34D399") // Emerald 400
	ErrorColor   = lipgloss.Color("#F87171") // Red 400
	WarningColor = lipgloss.Color("#FBBF24") // Amber 400
	InfoColor    = lipgloss.Color("#60A5FA") // Blue 400

	// Accents
	AccentColor = lipgloss.Color("#38BDF8") // Sky 400

	// Legacy / Compatibility Colors
	HighlightColor = lipgloss.Color("#334155")
	SecondaryColor = lipgloss.Color("#64748B")
)

// Global Layout Styles
var (
	// Base application style
	AppStyle = lipgloss.NewStyle().
			Margin(0, 1)

	// Clean container - No background, just subtle border
	ContainerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SubtleBorder).
			Padding(1, 2)

	// Focused container
	ActiveContainerStyle = ContainerStyle.Copy().
				BorderForeground(PrimaryColor)

	// Helpers for legacy support
	HighlightCardStyle = ContainerStyle.Copy()
	WarningCardStyle   = ContainerStyle.Copy().BorderForeground(WarningColor)
)

// Typography Styles
var (
	// Elegant Title
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
				MarginBottom(1) // Increased margin for spacing
)

// Component Styles
var (
	// Menu Items - Clean, no block background
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			PaddingLeft(2).
			PaddingRight(1)

	// Selected Item - Indicator bar instead of block
	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, false, true). // Left border only
				BorderForeground(PrimaryColor).
				PaddingLeft(1).
				PaddingRight(1)

	// Input Fields
	InputStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SubtleBorder).
			Padding(0, 1)

	FocusedInputStyle = InputStyle.Copy().
				BorderForeground(PrimaryColor)

	// Badges - Outline/Ghost style for elegance
	BadgeStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SubtleBorder).
			Padding(0, 1)

	SuccessBadgeStyle = BadgeStyle.Copy().
				Foreground(SuccessColor).
				BorderForeground(SuccessColor)

	WarningBadgeStyle = BadgeStyle.Copy().
				Foreground(WarningColor).
				BorderForeground(WarningColor)

	ErrorBadgeStyle = BadgeStyle.Copy().
			Foreground(ErrorColor).
			BorderForeground(ErrorColor)

	// Legacy badge support
	SuccessStatusBadge = SuccessBadgeStyle
	WarningStatusBadge = WarningBadgeStyle
	InfoStatusBadge    = BadgeStyle.Copy().Foreground(InfoColor).BorderForeground(InfoColor)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().Foreground(AccentColor)
)

// Helper functions defined with legacy signature support

func MakeDivider(width int, color ...lipgloss.Color) string {
	c := SubtleBorder
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

// Footer Styles - Clean text
var (
	KeyStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Bold(true)

	ActionKeyStyle = KeyStyle.Copy().
			Foreground(PrimaryColor)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor)
)

func MakeFooter(hints []KeyHint) string {
	var parts []string
	for i, hint := range hints {
		kStyle := KeyStyle
		if hint.IsAction {
			kStyle = ActionKeyStyle
		}

		// Separator
		sep := ""
		if i > 0 {
			sep = lipgloss.NewStyle().Foreground(SubtleBorder).Render("  │  ")
		}

		parts = append(parts,
			sep,
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

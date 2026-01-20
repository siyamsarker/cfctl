package ui

import "github.com/charmbracelet/lipgloss"

// Modern Color Palette - Cloudflare inspired
var (
	// Primary Colors
	PrimaryColor   = lipgloss.Color("#F38020") // Cloudflare Orange
	AccentColor    = lipgloss.Color("#00A8E8") // Bright Blue
	SecondaryColor = lipgloss.Color("#6366F1") // Indigo

	// Status Colors
	SuccessColor = lipgloss.Color("#10B981") // Green
	ErrorColor   = lipgloss.Color("#EF4444") // Red
	WarningColor = lipgloss.Color("#F59E0B") // Amber
	InfoColor    = lipgloss.Color("#3B82F6") // Blue

	// UI Colors
	TextColor       = lipgloss.Color("#F3F4F6") // Light Gray
	MutedColor      = lipgloss.Color("#9CA3AF") // Gray
	BorderColor     = lipgloss.Color("#4B5563") // Dark Gray
	BackgroundColor = lipgloss.Color("#1F2937") // Very Dark Gray
	HighlightColor  = lipgloss.Color("#374151") // Medium Dark Gray
)

// Modern Typography Styles
var (
	// Headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginBottom(1)

	SectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(AccentColor).
				MarginTop(1).
				MarginBottom(1)

	// Menu Styles
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			Foreground(TextColor)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(2).
				Foreground(PrimaryColor).
				Background(HighlightColor).
				Bold(true)

	// Status Styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	// Border Styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(2, 3).
			MarginTop(1)

	ActiveBorderStyle = BorderStyle.Copy().
				BorderForeground(PrimaryColor)

	// Content Styles
	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginTop(2).
			PaddingLeft(1)

	// Input Styles
	InputLabelStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			MarginBottom(0)

	InputStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(HighlightColor).
			Padding(0, 1).
			MarginBottom(1)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(HighlightColor).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(AccentColor).
				Padding(0, 1).
				MarginBottom(1)

	BlurredInputStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Background(BackgroundColor).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(BorderColor).
				Padding(0, 1).
				MarginBottom(1)

	// Badge Styles
	BadgeStyle = lipgloss.NewStyle().
			Background(HighlightColor).
			Foreground(AccentColor).
			Padding(0, 1).
			Bold(true)

	SuccessBadgeStyle = BadgeStyle.Copy().
				Background(SuccessColor).
				Foreground(lipgloss.Color("#000000"))

	// Card Styles
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			MarginBottom(1)

	HighlightCardStyle = CardStyle.Copy().
				BorderForeground(AccentColor).
				Background(HighlightColor)

	// Spinner Style
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)
)

// Helper functions for consistent spacing
var (
	SpacerStyle = lipgloss.NewStyle().
			Height(1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(BorderColor).
			Bold(true)
)

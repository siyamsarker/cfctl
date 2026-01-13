package ui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	PrimaryColor   = lipgloss.Color("#00D4FF")
	SecondaryColor = lipgloss.Color("#FF6B00")
	SuccessColor   = lipgloss.Color("#00FF88")
	ErrorColor     = lipgloss.Color("#FF0044")
	WarningColor   = lipgloss.Color("#FFAA00")
	MutedColor     = lipgloss.Color("#666666")
	TextColor      = lipgloss.Color("#FFFFFF")
)

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			PaddingTop(1).
			PaddingBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)

	SelectedMenuItemStyle = MenuItemStyle.Copy().
				Foreground(PrimaryColor).
				Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	InfoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			PaddingTop(1)

	InputStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 1)

	FocusedInputStyle = InputStyle.Copy().
				BorderForeground(SuccessColor)

	BlurredInputStyle = InputStyle.Copy().
				BorderForeground(MutedColor)
)

package entities

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func TokyoNightTheme() *huh.Theme {
	t := huh.ThemeBase()

	var (
		teal       = lipgloss.Color("#7eccb8")
		magenta    = lipgloss.Color("#bb9af7")
		cyan       = lipgloss.Color("#7dcfff")
		orange     = lipgloss.Color("#ff9e64")
		comment    = lipgloss.Color("#565f89")
		foreground = lipgloss.Color("#c0caf5")
		bgAccent   = lipgloss.Color("#24283b")
	)

	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(1).BorderStyle(lipgloss.ThickBorder()).BorderLeft(true).BorderForeground(cyan)
	t.Focused.Title = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Foreground(comment)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(teal)
	t.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(foreground)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(comment)
	t.Focused.SelectSelector = lipgloss.NewStyle().Foreground(cyan)
	t.Focused.Option = lipgloss.NewStyle().Foreground(foreground)
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(cyan)
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(magenta)
	t.Focused.FocusedButton = lipgloss.NewStyle().Foreground(bgAccent).Background(cyan).Padding(0, 3).Bold(true)
	t.Focused.BlurredButton = lipgloss.NewStyle().Foreground(foreground).Background(bgAccent).Padding(0, 3)

	t.Blurred.Base = lipgloss.NewStyle().PaddingLeft(2)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(comment)
	t.Blurred.Description = lipgloss.NewStyle().Foreground(comment)
	t.Blurred.TextInput.Prompt = lipgloss.NewStyle().Foreground(comment)
	t.Blurred.TextInput.Text = lipgloss.NewStyle().Foreground(comment)

	t.Help.ShortKey = lipgloss.NewStyle().Foreground(magenta)
	t.Help.ShortDesc = lipgloss.NewStyle().Foreground(comment)
	t.Focused.ErrorMessage = lipgloss.NewStyle().Foreground(orange).Bold(true)

	return t
}

func MonokaiRetroTheme() *huh.Theme {
	t := huh.ThemeBase()

	var (
		hotPink    = lipgloss.Color("#F92672") // Primary Focused Accent
		limeGreen  = lipgloss.Color("#A6E22E") // Prompts & Validations
		brightBlue = lipgloss.Color("#66D9EF") // Selections
		starkWhite = lipgloss.Color("#F8F8F2") // Crisp Input Text
		charcoal   = lipgloss.Color("#75715E") // Muted / Blurred States
		bgButton   = lipgloss.Color("#272822")
		orange     = lipgloss.Color("#ff9e64")
	)

	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(1).BorderStyle(lipgloss.ThickBorder()).BorderLeft(true).BorderForeground(hotPink)
	t.Focused.Title = lipgloss.NewStyle().Foreground(hotPink).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Foreground(charcoal)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(limeGreen)
	t.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(limeGreen).SetString("█ ")
	t.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(starkWhite)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(charcoal)
	t.Focused.SelectSelector = lipgloss.NewStyle().Foreground(brightBlue)
	t.Focused.Option = lipgloss.NewStyle().Foreground(starkWhite)
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(brightBlue)
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(limeGreen)
	t.Focused.FocusedButton = lipgloss.NewStyle().Foreground(bgButton).Background(hotPink).Padding(0, 3).Bold(true)
	t.Focused.BlurredButton = lipgloss.NewStyle().Foreground(starkWhite).Background(bgButton).Padding(0, 3)

	t.Blurred.Base = lipgloss.NewStyle().PaddingLeft(2)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(charcoal)
	t.Blurred.Description = lipgloss.NewStyle().Foreground(charcoal)
	t.Blurred.TextInput.Prompt = lipgloss.NewStyle().Foreground(charcoal).SetString("░ ")
	t.Blurred.TextInput.Text = lipgloss.NewStyle().Foreground(charcoal)

	t.Help.ShortKey = lipgloss.NewStyle().Foreground(brightBlue)
	t.Help.ShortDesc = lipgloss.NewStyle().Foreground(charcoal)
	t.Focused.ErrorMessage = lipgloss.NewStyle().Foreground(orange)

	return t
}

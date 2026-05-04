package misc

import "charm.land/lipgloss/v2"

// PrintLogo is just a pretty print of the Kaeru ascii
// when the app starts.
func PrintLogo() string {
	ascii := `
▄▄
██ ▄█▀  ▀▀█▄ ▄█▀█▄ ████▄ ██ ██
████   ▄█▀██ ██▄█▀ ██ ▀▀ ██ ██
██ ▀█▄ ▀█▄██ ▀█▄▄▄ ██    ▀██▀█
`

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff6e"))

	return style.Render(ascii)
}

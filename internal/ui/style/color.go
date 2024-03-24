package style

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/momarques/kibe/internal/logging"
)

func ColorizeTabKeys(row, col int) lipgloss.Style {
	switch {
	case col == 0:
		return lipgloss.NewStyle().
			Background(lipgloss.NoColor{}).
			Foreground(GetColor(ThemeConfig.Tab.ActiveTabContentKeys)).
			Padding(0, 1).
			Bold(true)
	}
	return lipgloss.NewStyle()
}

func GetColor(c string) lipgloss.TerminalColor {
	logging.Log.Info("color from theme -> ", c)
	colorPattern, _ := regexp.Compile(`^#[0-9a-fA-F]{6}$`)
	if colorPattern.MatchString(c) {
		return lipgloss.Color(c)
	} else {
		return lipgloss.NoColor{}
	}
}

// func colorizeInt(string) string {
// 	return ""
// }

// func colorizeIP(string) string {
// 	return ""
// }

// func colorizeResourceType(string) string {
// 	return ""
// }

// func colorizeStatus(string) string {
// 	return ""
// }

// func colorizePort(string) string {
// 	return ""
// }

// func colorizeBool(string) string {
// 	return ""
// }
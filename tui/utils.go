package tui

import "strings"

func cleanTab(s string) string {
	// Remove ANSI reset code and then resetForeground
	return strings.ReplaceAll(s, "\x1b[0m", "") + "\x1b[39m"
}

func (t *Tui) cleanPostName(name string) string {
	return strings.TrimSuffix(strings.ReplaceAll(name, "_", " "), " series")
}
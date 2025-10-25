package tui

import "strings"

// Use to preserve consistent background colour
func cleanTab(s string) string {
	// Remove ANSI reset code and then resetForeground
	return strings.ReplaceAll(s, "\x1b[0m", "") + "\x1b[39m"
}

func cleanPostName(name string) string {
	return strings.TrimSuffix(strings.ReplaceAll(name, "_", " "), " series")
}
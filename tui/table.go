package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/client"
)

func getPostName(p client.Post) string {
	seen := make(map[string]bool)
	var postNames []string
	for _, name := range p.Names {
		if !seen[name.Name] {
			cleanedName := cleanPostName(name.Name)
			postNames = append(postNames, cleanedName)
			seen[name.Name] = true
		}
	}
	postName := strings.Join(postNames, " • ")
	return cleanTab(postName)
}

func createTable(posts []client.Post, terminalWidth, terminalHeight int) table.Model {
	firstColumnWidth := int(float64(terminalWidth) * 0.01)
	columnWidth := int(float64(terminalWidth) * 0.07)
	titleWidth := terminalWidth - ((columnWidth * 2) + firstColumnWidth) - 10
	columns := []table.Column{
		{Title: "No", Width: firstColumnWidth},
		{Title: "Title", Width: titleWidth},
		{Title: "ID", Width: columnWidth},
		{Title: "Score", Width: columnWidth},
	}

	var rows []table.Row
	for i, post := range posts {
		postName := getPostName(post)
		rows = append(rows, table.Row{
			fmt.Sprintf("%v", i+1),
			postName,
			fmt.Sprintf("%v", post.ID),
			fmt.Sprintf("%v", post.Score),
		})
	}

	style := table.DefaultStyles()
	
	style.Header = style.Header.Background(bg).Foreground(highlight).
		Border(paneBorder, false, false, true, false).
		BorderForeground(highlight).BorderBackground(bg)
	
	style.Selected = style.Selected.
			Background(lipgloss.Color("#57534e")).
			Foreground(lipgloss.Color(highlight.Light))
	
		t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(terminalHeight),
		table.WithStyles(style),
	)

	return t
}

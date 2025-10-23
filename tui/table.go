package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/client"
)

var (
	columns = []table.Column{
		{Title: "Post No", Width: 10},
		{Title: "Title", Width: 30},
		{Title: "ID", Width: 10},
		{Title: "Score", Width: 15},
	}
	rows = []table.Row{}

	t = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	tableStyle = lipgloss.NewStyle()
)

func (t *Tui) getPostNames(p client.Post) string {
	seen := make(map[string]bool)
	var postNames []string
	for _, name := range p.Names {
		if !seen[name.Name] {
			cleanedName := t.cleanPostName(name.Name)
			postNames = append(postNames, cleanedName)
			seen[name.Name] = true
		}
	}
	postName := strings.Join(postNames, " • ")
	return titleStyle.Render(postName)
}
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
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
	return cleanTab(titleStyle.Render(postName))
}

func createTable(posts []client.Post, terminalWidth, terminalHeight int) table.Model {
	tableWidth := terminalWidth - 2
	columnWidth := int(float64(tableWidth) * 0.1)
	columns := []table.Column{
			{Title: "Post No", Width: columnWidth},
			{Title: "Title", Width: tableWidth - (columnWidth * 3) - 10},
			{Title: "ID", Width: columnWidth},
			{Title: "Score", Width: columnWidth},
		}
		
		var rows []table.Row
		for i, post := range posts{
			postName := getPostName(post)
			rows = append(rows, table.Row{
				fmt.Sprintf("%v", i+1),
				postName,
				fmt.Sprintf("%v", post.ID),
				fmt.Sprintf("%v", post.ID),
			})
		}
		
		style := table.DefaultStyles()
		style.Header = style.Header.Background(bg)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(terminalHeight-2),
			table.WithStyles(style),
		)
		
		return t
}
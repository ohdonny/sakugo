package tui

import (
	_"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/model"
)

type Tui struct {
	model    *model.Model
	tabIndex int
	spinner  spinner.Model
	postsTable  table.Model
}

func NewTui(m *model.Model) *Tui {
	return &Tui{
		model:    m,
		tabIndex: 0,
	}
}

func (t *Tui) Init() tea.Cmd {
	t.model.SetPosts()
	t.postsTable = createTable(t.model.Posts, t.model.TerminalWidth - 2, t.model.TerminalHeight - 5)
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.model.TerminalHeight = msg.Height
		t.model.TerminalWidth = msg.Width
		if t.model.TerminalWidth < 70 {
			t.model.TerminalWidth = 70
		}
		if t.model.TerminalHeight < 20 {
			t.model.TerminalHeight = 20
		}
		
		t.postsTable = createTable(t.model.Posts, t.model.TerminalWidth - 2, t.model.TerminalHeight - 5)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit
		}
	}
	return t, nil
}

func (t *Tui) View() string {
	switch t.model.ViewState {
	case model.PostsView:
		return t.renderPosts()
	}
	return ""
}

func (t *Tui) renderSearchBar() string {
	searchText := map[borderPosition]string{
		TopLeftBorder:  "Search",
		TopRightBorder: "Press ? for help",
	}
	searchBar := pane.Width(t.model.TerminalWidth - 2).Height(1).Background(bg).Render("")
	return borderize(searchBar, searchText)
}

func (t *Tui) renderPosts() string {
	tableView := t.postsTable.View()
	centeredContent := lipgloss.NewStyle().Width(t.model.TerminalWidth).
	AlignHorizontal(lipgloss.Top).Render(tableView)
		
	posts := t.renderPage(centeredContent)
	return posts
}

func (t *Tui) renderPage(content string) string {
	searchBar := t.renderSearchBar()
	pane := pane.Width(t.model.TerminalWidth - 2).
		Height(t.model.TerminalHeight - 5).Background(bg).Render(content)

	fullPane := lipgloss.JoinVertical(lipgloss.Left, searchBar, pane)
	layout := lipgloss.Place(
		t.model.TerminalWidth, t.model.TerminalHeight, lipgloss.Center, lipgloss.Bottom, fullPane)
	return layout
}
package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type menuitem struct {
	title, desc string
}

func (i menuitem) Title() string       { return i.title }
func (i menuitem) Description() string { return i.desc }
func (i menuitem) FilterValue() string { return i.title }

type clientMenuView struct {
	list       list.Model
	focusIndex int
	stateFn    func(state int) *clientView
}

func newClientMenuView(stateFn func(state int) *clientView) clientMenuView {
	items := []list.Item{
		menuitem{title: "download", desc: "download file from server"},
		menuitem{title: "upload", desc: "upload file to server"},
	}

	m := clientMenuView{
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
		stateFn: stateFn,
	}
	m.list.Title = "My Fave Things"
	return m
}

func (c clientMenuView) Init() tea.Cmd {
	return textinput.Blink
}

func (c clientMenuView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return c, tea.Quit
		case "down":
			c.focusIndex++
		case "up":
			c.focusIndex--
		case "enter":
			home := c.stateFn(2)
			return home.Update(msg)
		}
		if msg.String() == "ctrl+c" {
			return c, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		c.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c clientMenuView) View() string {
	return docStyle.Render(c.list.View())
}

package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type clientView struct {
	state    int
	menu     clientMenuView
	download clientDownloadView
}

func NewClientView() clientView {
	m := clientView{}
	m.menu = newClientMenuView(m.setState)
	m.download = newClientDownloadView(m.setState)
	return m
}

func (c *clientView) setState(state int) *clientView {
	c.state = state
	return c
}

func (c clientView) Init() tea.Cmd {
	return textinput.Blink
}

func (c clientView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch c.state {
	case 0:
		return c.menu.Update(msg)
	case 2:
		return c.download.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c clientView) View() string {
	switch c.state {
	case 0:
		return c.menu.View()
	case 2:
		return c.download.View()
	}

	return "fssync"
}

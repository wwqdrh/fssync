package server

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var exit = make(chan struct{})

type setstate func(int) *serverView

type serverView struct {
	state int

	basic basicView
}

func NewServerView() serverView {
	m := serverView{}
	m.basic = newBasicView(m.setState)
	return m
}

func (c *serverView) setState(state int) *serverView {
	c.state = state
	return c
}

func (c serverView) Init() tea.Cmd {
	return textinput.Blink
}

func (c serverView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if c.state == 0 {
		return c.basic.Update(msg)
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

func (c serverView) View() string {
	if c.state == 0 {
		return c.basic.View()
	}

	return "fssync"
}

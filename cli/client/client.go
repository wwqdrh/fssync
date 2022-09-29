package client

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type setstate func(state int) *clientView

type clientView struct {
	state        int
	menu         clientMenuView
	upload       uploadView
	download     downloadView
	downloadList downloadListView
}

func NewClientView() clientView {
	m := clientView{}
	m.menu = newClientMenuView(m.setState)
	m.upload = newUploadView(m.setState)
	m.download = newDownloadView(m.setState)
	m.downloadList = newDownloadListView(m.setState)
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
	case 1:
		return c.download.Update(msg)
	case 2:
		// c.downloadList.UpdateList()
		return c.downloadList.Update(msg)
	case 3:
		return c.upload.Update(msg)
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
	case 1:
		return c.download.View()
	case 2:
		return c.downloadList.View()
	case 3:
		return c.upload.View()
	}

	return "fssync"
}

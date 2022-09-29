package client

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwqdrh/fssync/client"
)

type fileitem struct {
	title, desc string
}

func (i fileitem) Title() string       { return i.title }
func (i fileitem) Description() string { return i.desc }
func (i fileitem) FilterValue() string { return i.title }

type downloadListView struct {
	list       list.Model
	focusIndex int
	stateFn    func(state int) *clientView
}

func newDownloadListView(stateFn func(state int) *clientView) downloadListView {
	items := []list.Item{
		fileitem{title: "here"},
		fileitem{title: "here2"},
	}

	m := downloadListView{
		list:    list.New(items, list.NewDefaultDelegate(), 0, 2),
		stateFn: stateFn,
	}
	m.list.Title = "downloadlist"
	return m
}

func (c downloadListView) Init() tea.Cmd {
	return textinput.Blink
}

func (c downloadListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			fmt.Println("下载")
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

func (c downloadListView) View() string {
	var b strings.Builder
	b.WriteString(c.list.View())
	// c.UpdateList()
	return b.String()
}

func (c *downloadListView) UpdateList() {
	files, err := client.DownloadList()
	if err != nil {
		fmt.Println(err)
	}
	items := make([]list.Item, 0, len(files))
	for _, item := range files {
		items = append(items, fileitem{title: item})
	}
	c.list = list.New(items, list.NewDefaultDelegate(), 0, 10)
	c.list.Title = "downloadlist"
}

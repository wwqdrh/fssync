package client

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwqdrh/fssync/client"
)

// var docStyle = lipgloss.NewStyle().Margin(1, 2)
var docStyle = lipgloss.NewStyle()

type menuitem struct {
	title, desc string
}

func (i menuitem) Title() string       { return i.title }
func (i menuitem) Description() string { return i.desc }
func (i menuitem) FilterValue() string { return i.title }

type clientMenuView struct {
	inputs      []field
	inputsIndex int

	list       list.Model
	focusIndex int
	stateFn    func(state int) *clientView
}

func newClientMenuView(stateFn func(state int) *clientView) clientMenuView {
	inputs := make([]field, 1)
	var t field
	for i := range inputs {
		t = field{
			field: textinput.New(),
		}
		t.field.CursorStyle = cursorStyle
		t.field.CharLimit = 32

		switch i {
		case 0:
			t.field.Placeholder = " 服务端地址: http://localhost:1080"
			t.field.Focus()
			t.id = "host"
			t.defaultvalue = "http://localhost:1080"
		}
		inputs[i] = t
	}

	items := []list.Item{
		menuitem{title: "download", desc: "输入文件名进行下载"},
		menuitem{title: "downloadlist", desc: "从服务端提供的下载列表中进行选择"},
		menuitem{title: "upload", desc: "上传的文件名"},
	}

	m := clientMenuView{
		inputs: inputs,

		list:    list.New(items, list.NewDefaultDelegate(), 0, 3),
		stateFn: stateFn,
	}
	m.list.Title = "menu"
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
			client.ClientDownloadFlag.DownloadUrl = c.GetField("host")
			home := c.stateFn(c.focusIndex + 1)
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
	return c, tea.Batch(cmd, c.updateInputs(msg))
}

func (c *clientMenuView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	for i := range c.inputs {
		c.inputs[i].field, cmds[i] = c.inputs[i].field.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c clientMenuView) View() string {
	var b strings.Builder
	for i := range c.inputs {
		b.WriteString(c.inputs[i].field.View())
		if i < len(c.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	b.WriteString(c.list.View())
	return b.String()
}

func (c *clientMenuView) GetField(id string) string {
	for _, item := range c.inputs {
		if item.id == id {
			if v := item.field.Value(); v != "" {
				return v
			}
			return item.defaultvalue
		}
	}
	return ""
}

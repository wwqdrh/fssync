package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwqdrh/fssync/server"
	"github.com/wwqdrh/gokit/logger"
)

// 默认页面 传入参数启动服务
var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type basicView struct {
	focusIndex int
	inputs     []field
	stateFn    setstate
}

type field struct {
	id           string
	field        textinput.Model
	defaultvalue string
}

func newBasicView(stateFn setstate) basicView {
	v := basicView{
		inputs:  make([]field, 3),
		stateFn: stateFn,
	}
	var t field
	for i := range v.inputs {
		t = field{field: textinput.New()}
		t.field.CharLimit = 32

		switch i {
		case 0:
			t.field.Placeholder = " 端口号: (:1080)"
			t.field.Focus()
			t.id = "port"
			t.defaultvalue = ":1080"
		case 1:
			t.field.Placeholder = " 上传文件保存路径: (./stores)"
			t.id = "stores"
			t.defaultvalue = "./stores"
		case 2:
			t.field.Placeholder = " 提供下载功能的文件夹: ('')"
			t.id = "download"
			t.defaultvalue = ""
		}
		v.inputs[i] = t
	}

	return v
}

func (c basicView) Init() tea.Cmd {
	return textinput.Blink
}

func (c basicView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			close(exit)
			return c, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && c.focusIndex == len(c.inputs) {
				go c.StartServer()
				return c, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			if c.focusIndex > len(c.inputs) {
				c.focusIndex = 0
			} else if c.focusIndex < 0 {
				c.focusIndex = len(c.inputs)
			}

			cmds := make([]tea.Cmd, len(c.inputs))
			for i := 0; i <= len(c.inputs)-1; i++ {
				if i == c.focusIndex {
					// Set focused state
					cmds[i] = c.inputs[i].field.Focus()
					continue
				}
				// Remove focused state
				c.inputs[i].field.Blur()
			}

			return c, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *basicView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i].field, cmds[i] = c.inputs[i].field.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c basicView) View() string {
	var b strings.Builder

	for i := range c.inputs {
		b.WriteString(c.inputs[i].field.View())
		if i < len(c.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if c.focusIndex == len(c.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	return b.String()
}

func (c *basicView) GetField(id string) string {
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

func (c *basicView) StartServer() {
	server.ServerFlag.Port = c.GetField("port")
	server.ServerFlag.Store = c.GetField("stores")
	server.ServerFlag.ExtraPath = c.GetField("download")

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		if err := server.Start(ctx); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	}()
	<-exit
	cancel()
}

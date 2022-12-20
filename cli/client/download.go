package client

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/gokit/logger"
)

type downloadView struct {
	focusIndex int
	inputs     []field
	cursorMode textinput.CursorMode
	stateFn    setstate
}

func newDownloadView(fn setstate) downloadView {
	v := downloadView{
		inputs:  make([]field, 3),
		stateFn: fn,
	}

	var t field
	for i := range v.inputs {
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
		case 1:
			t.field.Placeholder = " 下载文件: ./example.txt"
			t.id = "download"
		}
		v.inputs[i] = t
	}

	return v
}

func (c downloadView) Init() tea.Cmd {
	return textinput.Blink
}

func (c downloadView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return c, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			c.cursorMode++
			if c.cursorMode > textinput.CursorHide {
				c.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(c.inputs))
			for i := range c.inputs {
				cmds[i] = c.inputs[i].field.SetCursorMode(c.cursorMode)
			}
			return c, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && c.focusIndex == len(c.inputs) {
				c.StartServer()
				return c, tea.Quit
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
					c.inputs[i].field.PromptStyle = focusedStyle
					c.inputs[i].field.TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				c.inputs[i].field.Blur()
				c.inputs[i].field.PromptStyle = noStyle
				c.inputs[i].field.TextStyle = noStyle
			}

			return c, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *downloadView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i].field, cmds[i] = c.inputs[i].field.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c downloadView) View() string {
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

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(c.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (c *downloadView) GetField(id string) string {
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

func (c *downloadView) StartServer() {
	client.ClientDownloadFlag.DownloadUrl = c.GetField("host")
	client.ClientDownloadFlag.FileName = c.GetField("download")

	if err := client.DownloadStart(); err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
}

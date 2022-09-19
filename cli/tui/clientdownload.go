package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type clientDownloadView struct {
	focusIndex int
	inputs     []textinput.Model
	defaultVal []string
	cursorMode textinput.CursorMode
	stateFn    func(state int) *clientView
}

func newClientDownloadView(stateFn func(state int) *clientView) clientDownloadView {
	c := clientDownloadView{
		inputs: make([]textinput.Model, 3),
		defaultVal: []string{
			"http://localhost:1080",
			"./example.txt",
			".",
		},
		stateFn: stateFn,
	}

	var t textinput.Model
	for i := range c.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = " 服务端地址: http://localhost:1080"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = " 上传的文件: ./example.txt"
		case 2:
			t.Placeholder = " 文件分片信息保存地址: ."
		}
		c.inputs[i] = t
	}
	return c
}

func (c clientDownloadView) Init() tea.Cmd {
	return textinput.Blink
}

func (c clientDownloadView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				cmds[i] = c.inputs[i].SetCursorMode(c.cursorMode)
			}
			return c, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && c.focusIndex == len(c.inputs) {
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
					cmds[i] = c.inputs[i].Focus()
					c.inputs[i].PromptStyle = focusedStyle
					c.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				c.inputs[i].Blur()
				c.inputs[i].PromptStyle = noStyle
				c.inputs[i].TextStyle = noStyle
			}

			return c, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := c.updateInputs(msg)

	return c, cmd
}

func (c *clientDownloadView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(c.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (c clientDownloadView) View() string {
	var b strings.Builder

	for i := range c.inputs {
		b.WriteString(c.inputs[i].View())
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

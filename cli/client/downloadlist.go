package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwqdrh/fssync/client"
)

type downloadListView struct {
	choices  []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	intool   bool
	selected map[int]struct{} // which to-do items are selected
	stateFn  func(state int) *clientView
}

func newDownloadListView(stateFn func(state int) *clientView) downloadListView {
	m := downloadListView{
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
		stateFn:  stateFn,
	}
	return m
}

func (c downloadListView) Init() tea.Cmd {
	return nil
}

func (c downloadListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "down":
			if c.cursor < len(c.choices)-1 {
				c.cursor++
			} else {
				c.intool = true
			}
		case "up":
			if c.cursor > 0 {
				c.cursor--
				c.intool = false
			}
		case "enter":
			if c.intool {
				for _, err := range client.DownloadStartByList(c.selectedList()...) {
					fmt.Println(err.Error())
				}
				fmt.Println("done...")
			} else {
				_, ok := c.selected[c.cursor]
				if ok {
					delete(c.selected, c.cursor)
				} else {
					c.selected[c.cursor] = struct{}{}
				}
			}
		}
		if msg.String() == "ctrl+c" {
			return c, tea.Quit
		}
	}

	return c, nil
}

func (c downloadListView) View() string {
	// The header
	s := "What you want download?\n\n"

	// Iterate over our choices
	for i, choice := range c.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if !c.intool && c.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := c.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// button
	if c.intool {
		s += fmt.Sprintf("\n%s\n", focusedButton)
	} else {
		s += fmt.Sprintf("\n%s\n", blurredButton)
	}

	// The footer
	s += "\nPress q or quit to quit.\n"

	// Send the UI for rendering
	return s
}

func (c *downloadListView) UpdateList() {
	files, err := client.DownloadList()
	if err != nil {
		fmt.Println(err)
	}
	c.choices = files
}

func (c *downloadListView) selectedList() []string {
	res := []string{}
	for id := range c.selected {
		res = append(res, c.choices[id])
	}
	return res
}

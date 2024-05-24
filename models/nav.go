package models

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var breadcrumbs = []string{"../", "./"}

type NavModel struct {
	state *State

	currentDir string
	list       []string
	listSize   int
	cursor     int
}

func (m NavModel) String() string {
	return "NavModel"
}

func NewNavModel(state *State) NavModel {
	model := NavModel{
		state:      state,
		currentDir: "",
		listSize:   2,
		cursor:     0,
	}
	model.state.SetCurrentModel(model.String())
	return model
}

func (m NavModel) Init() tea.Cmd {
	return m.ListObjects
}

func (m NavModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.listSize-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				// Go up a directory
			} else if m.cursor == 1 {
				// Go to the current directory
			}
		}

	case listObjectsMsg:
		if msg.err != nil {
			return m, tea.Quit
		} else {
			m.list = msg.list
			m.listSize = len(m.list) + len(breadcrumbs)
			m.cursor = 0
		}
	}

	return m, nil
}

func (m NavModel) View() string {
	s := "File system:\n\n"
	
	for i := range m.listSize {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var path string
		if i < len(breadcrumbs) {
			path = breadcrumbs[i]
		} else {
			path = m.list[i-len(breadcrumbs)]
		}

		s += fmt.Sprintf("%s %s\n", cursor, path)
	}

	s += "\nPress q to quit.\n"
	return s
}

type listObjectsMsg struct {
	list []string
	err  error
}

func (m NavModel) ListObjects() tea.Msg {
	list, err := m.state.store.ListObjects(m.currentDir)
	if err != nil {
		return listObjectsMsg{err: err}
	}
	log.Println(list)
	return listObjectsMsg{list: list}
}

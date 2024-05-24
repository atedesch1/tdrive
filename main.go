package main

import (
	"fmt"
	"log"

	"github.com/atedesch1/odrive/pkg/store"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	objs  []string
	store store.Store
	err   error
}

type getObjectsMsg struct{
	objects []string
	err error
}

func (m model) getObjects() tea.Msg {
	objs, err := m.store.ListObjects()
	if err != nil {
		return getObjectsMsg{nil, err}
	}
	return getObjectsMsg{objs, nil}
}

func (m model) Init() tea.Cmd {
	return m.getObjects
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case getObjectsMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.objs = msg.objects
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	s := "Objects in bucket:\n"
	for _, obj := range m.objs {
		s += fmt.Sprintf("- %s\n", obj)
	}
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	}
	return s
}

func main() {
	store, err := store.NewGCSStorage("odrive_atedeschi")
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(model{store: store})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

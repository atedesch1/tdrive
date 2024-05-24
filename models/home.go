package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var homeChoices = []string{"Navigate", "Configure", "Quit"}

type HomeModel struct {
	cursor int
}

func (m HomeModel) String() string {
	return "HomeModel"
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(homeChoices)-1 {
				m.cursor++
			}
		case "enter":
			switch homeChoices[m.cursor] {
			case "Navigate":
				return NavModel{}, nil
			case "Configure":
				return ConfigModel{}, nil
			case "Quit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m HomeModel) View() string {
	s := "Choose an option:\n\n"

	for i, choice := range homeChoices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

package models

import tea "github.com/charmbracelet/bubbletea"

type NavModel struct {
}

func (m NavModel) String() string {
	return "NavModel"
}

func (m NavModel) Init() tea.Cmd {
	return nil
}

func (m NavModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m NavModel) View() string {
	return "NavModel"
}

package models

import tea "github.com/charmbracelet/bubbletea"

type ConfigModel struct {
}

func (m ConfigModel) String() string {
	return "ConfigModel"
}

func (m ConfigModel) Init() tea.Cmd {
	return nil
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfigModel) View() string {
	return "ConfigModel"
}


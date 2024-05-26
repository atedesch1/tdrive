package models

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	UpDirectoryIndex      = 0
	CurrentDirectoryIndex = 1
)

var breadcrumbs = []string{"../", "./"}

type NavModel struct {
	state *State

	currentPath string
	list        []string
	listSize    int
	cursor      int
}

func (m NavModel) String() string {
	return "NavModel"
}

func NewNavModel(state *State) NavModel {
	model := NavModel{
		state:       state,
		currentPath: "",
		listSize:    2,
		cursor:      0,
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
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "q":
			model := NewHomeModel(m.state)
			return model, model.Init()
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < m.listSize-1 {
				m.cursor++
			}
		case "-":
			// Go up a directory
			if m.currentPath != "" {
				m = m.GoUpADirectory()
				return m, m.ListObjects
			}
		case "~":
			// Go to root directory
			if m.currentPath != "" {
				m.currentPath = ""
				return m, m.ListObjects
			}
		case "enter":
			switch m.cursor {
			case UpDirectoryIndex:
				// Go up a directory (../ chosen)
				if m.currentPath != "" {
					m = m.GoUpADirectory()
					return m, m.ListObjects
				}
			case CurrentDirectoryIndex:
				// Go to the current directory (./ chosen)
			default:
				// Open the selected directory or file
				path := m.list[m.cursor-len(breadcrumbs)]
				isDirectory := path[len(path)-1] == '/'
				if isDirectory {
					// Go to the selected directory
					m.currentPath = path
					return m, m.ListObjects
				} else {
					// Download the selected file
					cmd := m.GetDownloadObjectCmd(path)
					return m, cmd
				}
			}
		}

	case listObjectsMsg:
		if msg.err != nil {
			log.Panic(msg.err)
			return m, tea.Quit
		} else {
			m.list = msg.list
			m.listSize = len(m.list) + len(breadcrumbs)
			m.cursor = 0
		}

	case downloadObjectMsg:
		if msg.err != nil {
			log.Panic(msg.err)
		}
	}

	return m, nil
}

func (m NavModel) View() string {
	s := fmt.Sprintf("Storage Provider: %s\n", m.state.config.StorageProvider)
	s += fmt.Sprintf("Current path: /%s\n\n", m.currentPath)

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

	s += "\nPress ESC or Q to back home.\n"
	return s
}

func (m NavModel) GoUpADirectory() NavModel {
	m.currentPath = strings.TrimSuffix(m.currentPath, "/")
	index := strings.LastIndex(m.currentPath, "/")
	if index == -1 {
		m.currentPath = ""
	} else {
		m.currentPath = m.currentPath[:index]
	}
	return m
}

type listObjectsMsg struct {
	list []string
	err  error
}

func (m NavModel) ListObjects() tea.Msg {
	list, err := m.state.store.ListObjects(m.currentPath)
	if err != nil {
		return listObjectsMsg{err: err}
	}
	return listObjectsMsg{list: list}
}

type downloadObjectMsg struct {
	err error
}

func (m NavModel) GetDownloadObjectCmd(objectPath string) func() tea.Msg {
	return func() tea.Msg {
		err := m.state.store.DownloadObject(objectPath, m.state.config.DownloadPath)
		if err != nil {
			return downloadObjectMsg{err: err}
		}
		return downloadObjectMsg{}
	}
}

package models

import (
	"github.com/atedesch1/odrive/pkg/store"
	tea "github.com/charmbracelet/bubbletea"
)

type State struct {
	currentModel string
	store        store.Store
}

func NewState(store store.Store) *State {
	return &State{
		store: store,
	}
}

func (s *State) InitialModel() tea.Model {
	return NewHomeModel(s)
}

func (s *State) SetCurrentModel(model string) {
	s.currentModel = model
}

func (s State) GetCurrentModel() string {
	return s.currentModel
}

package models

import (
	"log"

	"github.com/atedesch1/tdrive/pkg/config"
	"github.com/atedesch1/tdrive/pkg/store"
	tea "github.com/charmbracelet/bubbletea"
)

type State struct {
	currentModel string
	config       *config.Config
	store        store.Store
}

func NewState(cfg *config.Config) *State {
	if cfg == nil {
		// Start in configuration mode
		return &State{
			currentModel: ConfigModel{}.String(),
		}
	}

	store, err := store.NewStore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &State{
		currentModel: HomeModel{}.String(),
		config:       cfg,
		store:        store,
	}
}

func (s *State) InitialModel() tea.Model {
	switch s.currentModel {
	case ConfigModel{}.String():
		return NewConfigModel(s)
	case HomeModel{}.String():
		return NewHomeModel(s)
	default:
		log.Fatalf("initial model cannot be model %s", s.currentModel)
		return nil
	}
}

func (s *State) SetCurrentModel(model string) {
	s.currentModel = model
}

func (s State) GetCurrentModel() string {
	return s.currentModel
}

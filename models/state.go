package models

type State struct {
	currentModel string
}

func (s State) GetCurrentModel() string {
	return s.currentModel
}

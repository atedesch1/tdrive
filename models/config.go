package models

import (
	"fmt"
	"log"

	"github.com/atedesch1/tdrive/pkg/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfigStep string

const (
	StorageProviderStep ConfigStep = "Storage Provider"
	StorageConfigStep   ConfigStep = "Storage Config"
	DownloadPathStep    ConfigStep = "Download Path"
)

type ConfigModel struct {
	state *State

	step                     ConfigStep
	storageProviderStepModel storageProviderStepModel
	storageConfigStepModel   storageConfigStepModel
	downloadPathStepModel    downloadPathStepModel
}

func (m ConfigModel) String() string {
	return "ConfigModel"
}

func NewConfigModel(state *State) ConfigModel {
	config := &config.Config{}
	state.config = config
	model := ConfigModel{
		state: state,
		step:  StorageProviderStep,
	}
	model.state.SetCurrentModel(model.String())
	return model
}

type nextStepMsg struct {
	step ConfigStep
}

func NextStep(step ConfigStep) tea.Cmd {
	return func() tea.Msg {
		return nextStepMsg{step}
	}
}

type configuredMsg struct {
	config *config.Config
}

func Configured(config *config.Config) tea.Cmd {
	return func() tea.Msg {
		return configuredMsg{config}
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return NextStep(StorageProviderStep)
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nextStepMsg:
		switch msg.step {
		case StorageProviderStep:
			m.storageProviderStepModel = NewStorageProviderStepModel(m.state)
			m.step = StorageProviderStep
			return m, m.storageProviderStepModel.Init()
		case StorageConfigStep:
			m.storageConfigStepModel = NewStorageConfigStepModel(m.state)
			m.step = StorageConfigStep
			return m, m.storageConfigStepModel.Init()
		case DownloadPathStep:
			m.downloadPathStepModel = NewDownloadPathStepModel(m.state)
			m.step = DownloadPathStep
			return m, m.downloadPathStepModel.Init()
		}
	case configuredMsg:
		err := config.WriteConfig(msg.config)
		if err != nil {
			log.Fatalf("failed to write config: %v", err)
		}
		state := NewState(msg.config)
		model := NewNavModel(state)
		return model, model.Init()
	}

	switch m.step {
	case StorageProviderStep:
		model, cmd := m.storageProviderStepModel.Update(msg)
		m.storageProviderStepModel = model.(storageProviderStepModel)
		return m, cmd
	case StorageConfigStep:
		model, cmd := m.storageConfigStepModel.Update(msg)
		m.storageConfigStepModel = model.(storageConfigStepModel)
		return m, cmd
	case DownloadPathStep:
		model, cmd := m.downloadPathStepModel.Update(msg)
		m.downloadPathStepModel = model.(downloadPathStepModel)
		return m, cmd
	default:
		log.Fatalf("unsupported step %s", m.step)
	}

	return m, nil
}

func (m ConfigModel) View() string {
	switch m.step {
	case StorageProviderStep:
		return m.storageProviderStepModel.View()
	case StorageConfigStep:
		return m.storageConfigStepModel.View()
	case DownloadPathStep:
		return m.downloadPathStepModel.View()
	default:
		return "ConfigModel"
	}
}

type storageProviderStepModel struct {
	state *State

	cursor int
}

func NewStorageProviderStepModel(state *State) storageProviderStepModel {
	return storageProviderStepModel{
		state:  state,
		cursor: 0,
	}
}

func (m storageProviderStepModel) Init() tea.Cmd {
	return nil
}

func (m storageProviderStepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(config.StorageProviders)-1 {
				m.cursor++
			}
		case "enter":
			switch config.StorageProviders[m.cursor] {
			case config.GCSStorageProvider:
				m.state.config.StorageProvider = config.GCSStorageProvider
				m.state.config.StorageConfig = config.StorageConfig{
					GCSStorageConfig: &config.GCSStorageConfig{},
				}
				return m, NextStep(StorageConfigStep)
			default:
				log.Fatalf("unsupported storage provider %s", config.StorageProviders[m.cursor])
			}
		}
	}

	return m, nil
}

func (m storageProviderStepModel) View() string {
	s := "Choose your storage provider:\n\n"

	for i, choice := range config.StorageProviders {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		var provider string
		switch choice {
		case config.GCSStorageProvider:
			provider = "Google Cloud Storage"
		default:
			log.Fatalf("unsupported storage provider %s", choice)
		}

		s += fmt.Sprintf("%s %s\n", cursor, provider)
	}

	s += "\nPress q to quit.\n"
	return s
}

type storageConfigStepModel struct {
	state           *State
	bucketNameInput textinput.Model
}

func NewStorageConfigStepModel(state *State) storageConfigStepModel {
	ti := textinput.New()
	ti.Placeholder = "Bucket Name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return storageConfigStepModel{
		state:           state,
		bucketNameInput: ti,
	}
}

func (m storageConfigStepModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m storageConfigStepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state.config.StorageProvider {
	case config.GCSStorageProvider:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyEnter:
				m.state.config.StorageConfig.GCSStorageConfig.BucketName = m.bucketNameInput.Value()
				return m, NextStep(DownloadPathStep)
			default:
				m.bucketNameInput, cmd = m.bucketNameInput.Update(msg)
			}
		}
	default:
		log.Fatalf("unsupported storage provider %s", m.state.config.StorageProvider)
	}

	return m, cmd
}

func (m storageConfigStepModel) View() string {
	s := "Provide the storage bucket's name:\n\n"
	s += m.bucketNameInput.View()
	s += "\n(esc to quit)\n"
	return s
}

type downloadPathStepModel struct {
	state             *State
	downloadPathInput textinput.Model
}

func NewDownloadPathStepModel(state *State) downloadPathStepModel {
	ti := textinput.New()
	ti.Placeholder = "Download Path"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return downloadPathStepModel{
		state:             state,
		downloadPathInput: ti,
	}
}

func (m downloadPathStepModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m downloadPathStepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.state.config.DownloadPath = m.downloadPathInput.Value()
			return m, Configured(m.state.config)
		default:
			m.downloadPathInput, cmd = m.downloadPathInput.Update(msg)
		}
	}

	return m, cmd
}

func (m downloadPathStepModel) View() string {
	s := "Provide a local download target path:\n\n"
	s += m.downloadPathInput.View()
	s += "\n(esc to quit)\n"
	return s
}

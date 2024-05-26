package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atedesch1/odrive/models"
	"github.com/atedesch1/odrive/pkg/config"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	configDir := config.GetConfigPath()
	cfg, err := config.LoadConfig(configDir)
	if err != nil {
		// TODO: handle errors
		log.Println(err)
	}

	state := models.NewState(cfg)

	p := tea.NewProgram(state.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

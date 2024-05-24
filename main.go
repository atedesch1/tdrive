package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atedesch1/odrive/models"
	"github.com/atedesch1/odrive/pkg/store"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// TODO:
	// Check for configuration file
	// If it doesn't exist, begin program on Configure screen
	// If it does exist, read configuration file
	// If configuration file is invalid, begin program on Configure screen
	// If configuration file is valid, set up clients using configuration

	store, err := store.NewGCSStorage("odrive_atedeschi")
	if err != nil {
		log.Fatal(err)
	}

	state := models.NewState(store)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(state.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

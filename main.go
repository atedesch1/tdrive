package main

import (
	"log"

	"github.com/atedesch1/odrive/models"
	tea "github.com/charmbracelet/bubbletea"
)
func main() {
	p := tea.NewProgram(models.HomeModel{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

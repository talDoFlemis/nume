package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/taldoflemis/nume/internal/tui/models"
)

func main() {
	// Start with the welcome screen
	theme := models.ThemeCatppuccin()
	// m := models.NewWelcomeModel(theme)
	m := models.NewMainModel(theme)
	
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

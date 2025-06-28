package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taldoflemis/nume/internal/tui/models"
)

func main() {
	// Start with the welcome screen
	renderer := lipgloss.DefaultRenderer()

	theme := models.ThemeCatppuccin(renderer)

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	m := models.NewWelcomeModel(theme, "TERM", renderer.ColorProfile().Name(), currentUser.Username)
	// m := models.NewMainModel(theme)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

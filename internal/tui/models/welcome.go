package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	text      string
	textIndex int
	finished  bool
	size      tea.WindowSizeMsg
	*Theme
}

type tickMsg time.Time

func NewWelcomeModel(theme *Theme) WelcomeModel {
	return WelcomeModel{
		text:      "nume",
		textIndex: 0,
		finished:  false,
		size: tea.WindowSizeMsg{
			Width:  MinimalWidth,
			Height: MinimalHeight,
		},
		Theme: theme,
	}
}

func (WelcomeModel) Init() tea.Cmd {
	return tick()
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Allow skipping animation with any key
		if !m.finished {
			return m.skipToMain(), nil
		}
	case tea.WindowSizeMsg:
		m.size = msg

	case tickMsg:
		if m.textIndex < len(m.text) {
			m.textIndex++
			return m, tick()
		} else if !m.finished {
			m.finished = true
			return m, tea.Tick(time.Millisecond*TransitionDelay, func(_ time.Time) tea.Msg {
				return transitionMsg{}
			})
		}

	case transitionMsg:
		// Transition to main view
		return m.skipToMain(), nil
	}

	return m, nil
}

func (m WelcomeModel) View() string {
	if m.size.Width < MinimalWidth || m.size.Height < MinimalHeight {
		return lipgloss.Place(
			m.size.Width, m.size.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(m.Theme.Focused.Base.GetBorderBottomForeground()).
				Width(m.size.Width-ComponentPadding).
				Height(m.size.Height-ComponentPadding).
				Padding(ComponentPadding).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(m.Theme.Focused.Base.GetBorderBottomForeground()).
				Border(lipgloss.NormalBorder()).
				Render(fmt.Sprintf(
					"Please resize your terminal to at least %dx%d for optimal experience.",
					MinimalWidth, MinimalHeight,
				)),
		)
	}

	activeStyle := m.Focused

	// Show animated text
	displayText := m.text[:m.textIndex]

	// Add blinking cursor if not finished
	if !m.finished && m.textIndex < len(m.text) {
		displayText += "█"
	}

	flexBox := lipgloss.JoinVertical(
		lipgloss.Center,
		"Welcome to",
		activeStyle.NoteTitle.Render(strings.ToUpper(displayText)),
		"Press any key to continue...",
		fmt.Sprintf("Terminal Size: %d columns x %d rows\n", m.size.Width, m.size.Height),
	)

	content := lipgloss.NewStyle().
		Padding(ComponentPadding).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.Theme.Focused.Base.GetBorderBottomForeground()).
		Border(lipgloss.NormalBorder()).Render(flexBox)

	return lipgloss.Place(
		m.size.Width, m.size.Height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m WelcomeModel) skipToMain() tea.Model {
	model := NewMainModel(m.Theme)
	model.size.Height = m.size.Height
	model.size.Width = m.size.Width
	return model
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*AnimationDelay, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type transitionMsg struct{}

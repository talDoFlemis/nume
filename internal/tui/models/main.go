package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	DerivativeTab Tab = 0
	IntegralTab   Tab = 1
	EigenTab      Tab = 2
)

type MainModel struct {
	tabs      []string
	activeTab Tab
	models    map[Tab]NumeModel
	size      *tea.WindowSizeMsg
	keys      help.KeyMap
	help      help.Model
	*Theme
}

type NumeTabContent interface {
	GetHelpKeys() help.KeyMap
}

type NumeModel interface {
	tea.Model
	NumeTabContent
}

func NewMainModel(theme *Theme) MainModel {
	derivateModel := NewDerivativeModel(theme)
	integralModel := NewIntegralModel()
	eigenModel := NewEigenModel(theme)

	models := make(map[Tab]NumeModel)

	models[DerivativeTab] = derivateModel
	models[IntegralTab] = integralModel
	models[EigenTab] = eigenModel

	return MainModel{
		tabs:      []string{"d Derivatives", "i Integrals", "e Eigen"},
		activeTab: DerivativeTab,
		models:    models,
		size: &tea.WindowSizeMsg{
			Width:  0,
			Height: 0,
		},
		keys:  derivateModel.GetHelpKeys(),
		help:  help.New(),
		Theme: theme,
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.models[m.activeTab].Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = &msg
		// Set help width for responsive design
		m.help.Width = msg.Width

		// Pass window size to child models
		var cmds []tea.Cmd

		for modelTab, model := range m.models {
			var newModel tea.Model
			var cmd tea.Cmd
			newModel, cmd = model.Update(msg)

			if sameModel, ok := newModel.(NumeModel); ok {
				m.models[modelTab] = sameModel
			}

			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case "d":
			if m.activeTab != DerivativeTab {
				m.activeTab = DerivativeTab
				m.keys = m.models[m.activeTab].GetHelpKeys()
			}
			return m, nil
		case "i":
			if m.activeTab != IntegralTab {
				m.activeTab = IntegralTab
				m.keys = m.models[m.activeTab].GetHelpKeys()
			}
			return m, nil
		case "e":
			if m.activeTab != EigenTab {
				m.activeTab = EigenTab
				m.keys = m.models[m.activeTab].GetHelpKeys()
			}
			return m, nil
		}
	}

	// Delegate to active tab's model
	var cmd tea.Cmd

	model := m.models[m.activeTab]

	var newModel tea.Model
	newModel, cmd = model.Update(msg)

	if sameModel, ok := newModel.(NumeModel); ok {
		m.models[m.activeTab] = sameModel
	}

	return m, cmd
}

func (m MainModel) View() string {
	if m.size.Width < MinimalWidth || m.size.Height < MinimalHeight {
		return lipgloss.Place(
			m.size.Width, m.size.Height,
			lipgloss.Center, lipgloss.Center,
			m.Renderer.NewStyle().
				Foreground(m.Focused.Base.GetBorderBottomForeground()).
				Width(m.size.Width-ComponentPadding).
				Height(m.size.Height-ComponentPadding).
				Padding(ComponentPadding).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(m.Focused.Base.GetBorderBottomForeground()).
				Border(lipgloss.NormalBorder()).
				Render(fmt.Sprintf(
					"Please resize your terminal to at least %dx%d for optimal experience.",
					MinimalWidth, MinimalHeight,
				)),
		)
	}

	// Render tabs
	var tabsRender []string
	for i, tab := range m.tabs {
		style := m.Blurred.BlurredButton
		isActive := i == int(m.activeTab)
		if isActive {
			style = m.Focused.FocusedButton
		}

		tabsRender = append(tabsRender, style.Render(tab))
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, tabsRender...)

	// Header with instructions
	header := m.Renderer.NewStyle().
		Bold(true).
		Foreground(m.Focused.Title.GetForeground()).
		Render("NUME - Numerical Methods Calculator")

	// Use the help view directly
	helpView := m.help.View(m.keys)

	// Style the help view
	styledHelp := m.Renderer.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(m.Focused.Base.GetBorderBottomForeground()).
		Render(helpView)

	// Content area
	content := m.models[m.activeTab].View()

	// Layout
	flexBox := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		tabsRow,
		"",
		m.Renderer.NewStyle().
			BorderTop(false).
			Padding(1).
			Render(content),
		"",
		styledHelp,
	)

	return lipgloss.Place(
		m.size.Width, m.size.Height,
		lipgloss.Center, lipgloss.Center,
		flexBox,
	)
}

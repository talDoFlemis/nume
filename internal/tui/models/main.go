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
)

type GlobalState struct {
}

var globalState = &GlobalState{
}

type MainModel struct {
	tabs            []string
	activeTab       Tab
	derivativeModel *DerivativeModel
	integralModel   *IntegralModel
	size            *tea.WindowSizeMsg
	keys            help.KeyMap
	help            help.Model
	*Theme
}

type NumeTabContent interface {
	GetHelpKeys() help.KeyMap
}

type instruction struct {
	shortcut    string
	description string
}

type tabItem struct {
	shortcut string
	name     string
}

func NewMainModel(theme *Theme) MainModel {
	derivateModel := NewDerivativeModel()
	integralModel := NewIntegralModel()

	return MainModel{
		tabs:            []string{"d Derivatives", "i Integrals"},
		activeTab:       DerivativeTab,
		derivativeModel: derivateModel,
		integralModel:   integralModel,
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
	return m.derivativeModel.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = &msg
		// Set help width for responsive design
		m.help.Width = msg.Width

		// Pass window size to child models
		var cmds []tea.Cmd
		if m.derivativeModel != nil {
			var newModel tea.Model
			var cmd tea.Cmd
			newModel, cmd = m.derivativeModel.Update(msg)
			if derivModel, ok := newModel.(*DerivativeModel); ok {
				m.derivativeModel = derivModel
			}
			cmds = append(cmds, cmd)
		}
		if m.integralModel != nil {
			var newModel tea.Model
			var cmd tea.Cmd
			newModel, cmd = m.integralModel.Update(msg)
			if intModel, ok := newModel.(*IntegralModel); ok {
				m.integralModel = intModel
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
				m.keys = m.derivativeModel.GetHelpKeys()
			}
			return m, nil
		case "i":
			if m.activeTab != IntegralTab {
				m.activeTab = IntegralTab
			}
			return m, nil
		}
	}

	// Delegate to active tab's model
	var cmd tea.Cmd
	switch m.activeTab {
	case DerivativeTab:
		var newModel tea.Model
		newModel, cmd = m.derivativeModel.Update(msg)
		if derivModel, ok := newModel.(*DerivativeModel); ok {
			m.derivativeModel = derivModel
		}
	case IntegralTab:
		var newModel tea.Model
		newModel, cmd = m.integralModel.Update(msg)
		if intModel, ok := newModel.(*IntegralModel); ok {
			m.integralModel = intModel
		}
	}

	return m, cmd
}

func (m MainModel) View() string {
	if m.size.Width < MINIMAL_WIDTH || m.size.Height < MINIMAL_HEIGHT {
		return lipgloss.Place(
			m.size.Width, m.size.Height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(m.Theme.Focused.Base.GetBorderBottomForeground()).
				Width(m.size.Width-2).
				Height(m.size.Height-2).
				Padding(2).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(m.Theme.Focused.Base.GetBorderBottomForeground()).
				Border(lipgloss.NormalBorder()).
				Render(fmt.Sprintf(
					"Please resize your terminal to at least %dx%d for optimal experience.",
					MINIMAL_WIDTH, MINIMAL_HEIGHT,
				)),
		)
	}

	// Tab styles
	activeTabStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#282828")).
		Padding(0, 2)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(0, 2)

	// Render tabs
	var tabsRender []string
	for i, tab := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == int(m.activeTab)
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		tabsRender = append(tabsRender, style.Render(tab))
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, tabsRender...)

	// Header with instructions
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Focused.Title.GetForeground()).
		Render("NUME - Numerical Methods Calculator")

	// Use the help view directly
	helpView := m.help.View(m.keys)
	
	// Style the help view
	styledHelp := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(m.Theme.Focused.Base.GetBorderBottomForeground()).
		Render(helpView)

	// Content area
	var content string
	// if m.activeTab == 0 {
	// 	content = m.derivativeModel.View()
	// } else {
	// 	content = m.integralModel.View()
	// }

	// Layout
	flexBox := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		tabsRow,
		"",
		lipgloss.NewStyle().
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

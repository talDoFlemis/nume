package models

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/taldoflemis/nume/internal/expressions"
	"github.com/taldoflemis/nume/internal/usecases"
)

type DerivativeModel struct {
	size *tea.WindowSizeMsg
	// Current focus section (0-4)
	focusedSection int

	// Section 1: Method Selection
	methodOptions  []string
	selectedMethod int

	// Section 2: Function Selection
	functionOptions  []string
	selectedFunction int

	// Section 3: Derivative Settings (Grid layout)
	derivativeOrder int // 1, 2, or 3
	philosophy      int // 0: forward, 1: backward, 2: central

	// Section 4: Error and Delta Settings (Grid layout)
	errorDegree int // 1-4
	delta       float64
	deltaInput  string

	// Section 5: Algorithm Parameters
	epsilon       float64
	maxSteps      int
	epsilonInput  string
	maxStepsInput string

	// Calculation results
	result          string
	showExplanation bool
	explanation     string
	functionExpr    expressions.SingleVariableExpr

	// Styling
	renderer *glamour.TermRenderer
}

// keyMap defines the keybindings for the main model
type derivativeKeyMap struct {
	Quit      key.Binding
	Help      key.Binding
	TabD      key.Binding
	TabI      key.Binding
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Enter     key.Binding
	Space     key.Binding
	Calculate key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k derivativeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k derivativeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TabD, k.TabI, k.Help},                // first column - navigation
		{k.Up, k.Down, k.Left, k.Right},         // second column - movement
		{k.Enter, k.Space, k.Calculate, k.Quit}, // third column - actions
	}
}

var derivativeKeys = derivativeKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	TabD: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "derivatives tab"),
	),
	TabI: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "integrals tab"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/confirm"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle option"),
	),
	Calculate: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "calculate"),
	),
}

// GetHelpKeys implements NumeTabContent.
func (m *DerivativeModel) GetHelpKeys() help.KeyMap {
	return derivativeKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		TabD: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "derivatives tab"),
		),
		TabI: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "integrals tab"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/confirm"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle option"),
		),
		Calculate: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "calculate"),
		),
	}
}

var _ (NumeTabContent) = (*DerivativeModel)(nil)

func NewDerivativeModel() *DerivativeModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(70),
	)

	return &DerivativeModel{
		focusedSection: 0,
		methodOptions: []string{
			"Newton Interpolation",
			"Lagrange Interpolation",
		},
		selectedMethod: 0,
		functionOptions: []string{
			"Polynomial: f(x) = x³ + 2x² - x + 1",
			"Exponential: f(x) = e^x",
			"Trigonometric: f(x) = sin(x)",
			"Logarithmic: f(x) = ln(x)",
			"Rational: f(x) = 1/x",
			"Composite: f(x) = sin(x²)",
		},
		selectedFunction: 0,
		derivativeOrder:  1,
		philosophy:       2, // central
		errorDegree:      2,
		delta:            0.001,
		deltaInput:       "0.001",
		epsilon:          1e-6,
		maxSteps:         100,
		epsilonInput:     "1e-6",
		maxStepsInput:    "100",
		renderer:         renderer,
	}
}

func (m *DerivativeModel) Init() tea.Cmd {
	return nil
}

func (m *DerivativeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusedSection = (m.focusedSection + 1) % 5
			return m, nil
		case "shift+tab":
			m.focusedSection = (m.focusedSection - 1 + 5) % 5
			return m, nil
		case "up", "k":
			return m.handleUp(), nil
		case "down", "j":
			return m.handleDown(), nil
		case "left", "h":
			return m.handleLeft(), nil
		case "right", "l":
			return m.handleRight(), nil
		case "enter":
			return m.handleEnter(), nil
		case "f", "F":
			m.generateResult()
			return m, nil
		case "e":
			m.showExplanation = !m.showExplanation
			if m.showExplanation && m.explanation == "" {
				m.generateExplanation()
			}
			return m, nil
		case "r":
			return NewDerivativeModel(), nil
		}

		// Handle input for delta, epsilon, and maxSteps
		if m.focusedSection == 3 || m.focusedSection == 4 {
			return m.handleTextInput(msg.String()), nil
		}
	}

	return m, nil
}

func (m *DerivativeModel) handleUp() *DerivativeModel {
	switch m.focusedSection {
	case 0: // Method selection
		if m.selectedMethod > 0 {
			m.selectedMethod--
		}
	case 1: // Function selection
		if m.selectedFunction > 0 {
			m.selectedFunction--
		}
	case 2: // Derivative order or philosophy
		// Focus on derivative order (left side)
		if m.derivativeOrder > 1 {
			m.derivativeOrder--
		}
	case 3: // Error degree
		if m.errorDegree > 1 {
			m.errorDegree--
		}
	case 4: // Max steps
		if steps, err := strconv.Atoi(m.maxStepsInput); err == nil && steps > 10 {
			m.maxSteps = steps - 10
			m.maxStepsInput = strconv.Itoa(m.maxSteps)
		}
	}
	return m
}

func (m *DerivativeModel) handleDown() *DerivativeModel {
	switch m.focusedSection {
	case 0: // Method selection
		if m.selectedMethod < len(m.methodOptions)-1 {
			m.selectedMethod++
		}
	case 1: // Function selection
		if m.selectedFunction < len(m.functionOptions)-1 {
			m.selectedFunction++
		}
	case 2: // Derivative order
		if m.derivativeOrder < 3 {
			m.derivativeOrder++
		}
	case 3: // Error degree
		if m.errorDegree < 4 {
			m.errorDegree++
		}
	case 4: // Max steps
		if steps, err := strconv.Atoi(m.maxStepsInput); err == nil {
			m.maxSteps = steps + 10
			m.maxStepsInput = strconv.Itoa(m.maxSteps)
		}
	}
	return m
}

func (m *DerivativeModel) handleLeft() *DerivativeModel {
	switch m.focusedSection {
	case 2: // Philosophy (right side of grid)
		if m.philosophy > 0 {
			m.philosophy--
		}
	}
	return m
}

func (m *DerivativeModel) handleRight() *DerivativeModel {
	switch m.focusedSection {
	case 2: // Philosophy (right side of grid)
		if m.philosophy < 2 {
			m.philosophy++
		}
	}
	return m
}

func (m *DerivativeModel) handleEnter() *DerivativeModel {
	// Generate result when enter is pressed
	m.generateResult()
	return m
}

func (m *DerivativeModel) handleTextInput(input string) *DerivativeModel {
	switch m.focusedSection {
	case 3: // Delta input
		switch input {
		case "backspace":
			if len(m.deltaInput) > 0 {
				m.deltaInput = m.deltaInput[:len(m.deltaInput)-1]
			}
		default:
			if len(input) == 1 &&
				(input[0] >= '0' && input[0] <= '9' || input[0] == '.' || input[0] == 'e' || input[0] == '-') {
				m.deltaInput += input
				if val, err := strconv.ParseFloat(m.deltaInput, 64); err == nil {
					m.delta = val
				}
			}
		}
	case 4: // Epsilon input
		switch input {
		case "backspace":
			if len(m.epsilonInput) > 0 {
				m.epsilonInput = m.epsilonInput[:len(m.epsilonInput)-1]
			}
		default:
			if len(input) == 1 &&
				(input[0] >= '0' && input[0] <= '9' || input[0] == '.' || input[0] == 'e' || input[0] == '-') {
				m.epsilonInput += input
				if val, err := strconv.ParseFloat(m.epsilonInput, 64); err == nil {
					m.epsilon = val
				}
			}
		}
	}
	return m
}

func (m *DerivativeModel) View() string {
	var sections []string

	// Section 1: Method Selection
	sections = append(sections, m.renderMethodSelection())

	// Section 2: Function Selection
	sections = append(sections, m.renderFunctionSelection())

	// Section 3: Derivative Settings (Grid)
	sections = append(sections, m.renderDerivativeSettings())

	// Section 4: Error and Delta Settings (Grid)
	sections = append(sections, m.renderErrorDeltaSettings())

	// Section 5: Algorithm Parameters
	sections = append(sections, m.renderAlgorithmParameters())

	// Results section
	if m.result != "" {
		sections = append(sections, m.renderResults())
	}

	// Join all sections
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Center the entire content
	return lipgloss.Place(
		80, 24,
		lipgloss.Center, lipgloss.Top,
		content,
	)
}

func (m *DerivativeModel) renderMethodSelection() string {
	style := m.getSectionStyle(0)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("1. Method Selection")
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Choose the interpolation method for derivative calculation")

	options := ""
	for i, method := range m.methodOptions {
		prefix := "  "
		if i == m.selectedMethod {
			prefix = "▶ "
		}
		options += fmt.Sprintf("%s%s\n", prefix, method)
	}

	return style.Render(title + "\n" + desc + "\n\n" + options)
}

func (m *DerivativeModel) renderFunctionSelection() string {
	style := m.getSectionStyle(1)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("2. Function Selection")
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Select the mathematical function to differentiate")

	options := ""
	for i, function := range m.functionOptions {
		prefix := "  "
		if i == m.selectedFunction {
			prefix = "▶ "
		}
		options += fmt.Sprintf("%s%s\n", prefix, function)
	}

	return style.Render(title + "\n" + desc + "\n\n" + options)
}

func (m *DerivativeModel) renderDerivativeSettings() string {
	style := m.getSectionStyle(2)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("3. Derivative Configuration")
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Configure derivative order and difference method")

	// Left side: Derivative Order
	orderOptions := []string{"First (f')", "Second (f'')", "Third (f''')"}
	orderSection := "Derivative Order:\n"
	for i, order := range orderOptions {
		prefix := "  "
		if i+1 == m.derivativeOrder {
			prefix = "▶ "
		}
		orderSection += fmt.Sprintf("%s%s\n", prefix, order)
	}

	// Right side: Philosophy
	philosophyOptions := []string{"Forward", "Backward", "Central"}
	philosophySection := "Difference Method:\n"
	for i, phil := range philosophyOptions {
		prefix := "  "
		if i == m.philosophy {
			prefix = "▶ "
		}
		philosophySection += fmt.Sprintf("%s%s\n", prefix, phil)
	}

	// Grid layout
	grid := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(30).Render(orderSection),
		lipgloss.NewStyle().Width(30).Render(philosophySection),
	)

	return style.Render(title + "\n" + desc + "\n\n" + grid)
}

func (m *DerivativeModel) renderErrorDeltaSettings() string {
	style := m.getSectionStyle(3)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("4. Numerical Parameters")
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Set error degree and step size (delta)")

	// Left side: Error Degree
	errorSection := "Error Degree:\n"
	for i := 1; i <= 4; i++ {
		prefix := "  "
		if i == m.errorDegree {
			prefix = "▶ "
		}
		errorSection += fmt.Sprintf("%sO(h^%d)\n", prefix, i)
	}

	// Right side: Delta
	deltaSection := fmt.Sprintf("Delta (h): %s\n", m.deltaInput)
	if m.focusedSection == 3 {
		deltaSection += "Type to edit..."
	}

	// Grid layout
	grid := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(30).Render(errorSection),
		lipgloss.NewStyle().Width(30).Render(deltaSection),
	)

	return style.Render(title + "\n" + desc + "\n\n" + grid)
}

func (m *DerivativeModel) renderAlgorithmParameters() string {
	style := m.getSectionStyle(4)
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("5. Algorithm Parameters")
	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("Configure convergence criteria and iteration limits")

	content := fmt.Sprintf(
		"Epsilon (ε): %s\nMax Steps: %s\n\nPress F to calculate, E for explanation, R to reset",
		m.epsilonInput,
		m.maxStepsInput,
	)

	if m.focusedSection == 4 {
		content += "\nType to edit epsilon or use arrows for max steps"
	}

	return style.Render(title + "\n" + desc + "\n\n" + content)
}

func (m *DerivativeModel) renderResults() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1).
		Width(70)

	content := "RESULTS:\n\n" + m.result

	if m.showExplanation && m.explanation != "" {
		content += "\n\nEXPLANATION:\n"
		if rendered, err := m.renderer.Render(m.explanation); err == nil {
			content += rendered
		} else {
			content += m.explanation
		}
	}

	return style.Render(content)
}

func (m *DerivativeModel) getSectionStyle(sectionIndex int) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(1).
		Width(70).
		Align(lipgloss.Left)

	if m.focusedSection == sectionIndex {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))
	} else {
		style = style.
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#444444"))
	}

	return style
}

func (m *DerivativeModel) generateResult() {
	m.setupFunctionExpression()

	// Choose strategy based on philosophy
	var strategy usecases.DifferenceStrategy
	switch m.philosophy {
	case 0: // forward
		strategy = &usecases.ForwardDifferenceStrategy{}
	case 1: // backward
		strategy = &usecases.BackwardDifferenceStrategy{}
	case 2: // central
		strategy = &usecases.CentralDifferenceStrategy{}
	default:
		strategy = &usecases.CentralDifferenceStrategy{}
	}

	ctx := context.Background()

	// Calculate derivative based on order
	var derivativeExpr expressions.SingleVariableExpr
	var err error

	switch m.derivativeOrder {
	case 1:
		derivativeExpr, err = strategy.Derivative(ctx, m.functionExpr, m.delta)
	case 2:
		derivativeExpr, err = strategy.DoubleDerivative(ctx, m.functionExpr, m.delta)
	case 3:
		// For third derivative, apply derivative twice
		firstDeriv, err1 := strategy.Derivative(ctx, m.functionExpr, m.delta)
		if err1 != nil {
			err = err1
			break
		}
		secondDeriv, err2 := strategy.Derivative(ctx, firstDeriv, m.delta)
		if err2 != nil {
			err = err2
			break
		}
		derivativeExpr, err = strategy.Derivative(ctx, secondDeriv, m.delta)
	}

	if err != nil {
		m.result = fmt.Sprintf("Error calculating derivative: %v", err)
		return
	}

	// Evaluate at appropriate test points
	testPoint := m.getTestPoint()
	derivativeValue := derivativeExpr(testPoint)

	methodName := m.methodOptions[m.selectedMethod]
	functionName := strings.Split(m.functionOptions[m.selectedFunction], ":")[0]
	philosophyName := []string{"Forward", "Backward", "Central"}[m.philosophy]

	m.result = fmt.Sprintf(`Method: %s
Function: %s
Derivative Order: %s
Philosophy: %s difference
Error Degree: O(h^%d)
Delta (h): %.6f
Epsilon (ε): %.2e
Max Steps: %d

Test point: x = %.1f
Derivative value: %.6f`,
		methodName,
		functionName,
		m.getDerivativeOrderText(),
		philosophyName,
		m.errorDegree,
		m.delta,
		m.epsilon,
		m.maxSteps,
		testPoint,
		derivativeValue)
}

func (m *DerivativeModel) getTestPoint() float64 {
	switch m.selectedFunction {
	case 3: // Logarithmic
		return 2.0 // Avoid x=1 where ln(x)=0 and ensure x>0
	case 4: // Rational
		return 2.0 // Avoid x=0 for 1/x
	default:
		return 1.0
	}
}

func (m *DerivativeModel) getDerivativeOrderText() string {
	switch m.derivativeOrder {
	case 1:
		return "First derivative (f'(x))"
	case 2:
		return "Second derivative (f''(x))"
	case 3:
		return "Third derivative (f'''(x))"
	default:
		return "Unknown"
	}
}

func (m *DerivativeModel) setupFunctionExpression() {
	// Define function expressions based on selected function
	switch m.selectedFunction {
	case 0: // Polynomial
		m.functionExpr = func(x float64) float64 {
			return x*x*x + 2*x*x - x + 1
		}
	case 1: // Exponential
		m.functionExpr = func(x float64) float64 {
			return math.Exp(x)
		}
	case 2: // Trigonometric
		m.functionExpr = func(x float64) float64 {
			return math.Sin(x)
		}
	case 3: // Logarithmic
		m.functionExpr = func(x float64) float64 {
			if x <= 0 {
				return math.NaN()
			}
			return math.Log(x)
		}
	case 4: // Rational
		m.functionExpr = func(x float64) float64 {
			if x == 0 {
				return math.NaN()
			}
			return 1.0 / x
		}
	case 5: // Composite
		m.functionExpr = func(x float64) float64 {
			return math.Sin(x * x)
		}
	default:
		m.functionExpr = func(x float64) float64 {
			return x*x*x + 2*x*x - x + 1
		}
	}
}

func (m *DerivativeModel) generateExplanation() {
	philosophyName := []string{"forward", "backward", "central"}[m.philosophy]
	filename := fmt.Sprintf("%s_difference.md", philosophyName)
	explanationPath := filepath.Join("internal", "tui", "views", "explanations", filename)

	if content, err := ioutil.ReadFile(explanationPath); err == nil {
		m.explanation = string(content)
	} else {
		// Fallback explanation
		m.explanation = fmt.Sprintf(`# %s Difference Method

## Overview
The %s difference method for numerical differentiation.

## Configuration
- **Method**: %s
- **Function**: %s  
- **Order**: %s
- **Error**: O(h^%d)

## Parameters
- **Delta**: %.6f
- **Epsilon**: %.2e
- **Max Steps**: %d
`,
			strings.Title(philosophyName),
			philosophyName,
			m.methodOptions[m.selectedMethod],
			strings.Split(m.functionOptions[m.selectedFunction], ":")[0],
			m.getDerivativeOrderText(),
			m.errorDegree,
			m.delta,
			m.epsilon,
			m.maxSteps)
	}
}

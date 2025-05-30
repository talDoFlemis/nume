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
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/taldoflemis/nume/internal/expressions"
	"github.com/taldoflemis/nume/internal/usecases"
)

type DerivativeModel struct {
	size *tea.WindowSizeMsg
	// Current focus section (0-5)
	focusedSection int

	// Section 1: Function Selection
	functionOptions  []string
	selectedFunction int

	// Section 2: Error Order (for polynomial functions)
	polynomialOrder int // 1-4 (linear to 4th degree)

	// Section 3: Derivative Order
	derivativeOrder int // 1, 2, or 3

	// Section 4: Philosophy (difference method)
	philosophy int // 0: forward, 1: backward, 2: central

	// Section 5: Inputs (Delta and Test Point inputs)
	deltaInput     textinput.Model
	testPointInput textinput.Model
	delta          float64
	testPoint      float64

	// Calculation results
	result          string
	showExplanation bool
	explanation     string
	functionExpr    expressions.SingleVariableExpr

	// Styling
	renderer *glamour.TermRenderer
	*Theme
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
	Explain   key.Binding
	Reset     key.Binding
	Section1  key.Binding
	Section2  key.Binding
	Section3  key.Binding
	Section4  key.Binding
	Section5  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k derivativeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k derivativeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TabD, k.TabI, k.Help},                           // first column - navigation
		{k.Up, k.Down, k.Left, k.Right},                    // second column - movement
		{k.Section1, k.Section2, k.Section3, k.Section4, k.Section5}, // third column - sections
		{k.Enter, k.Calculate, k.Explain, k.Reset, k.Quit}, // fourth column - actions
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
	Explain: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "toggle explanation"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	Section1: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "function selection"),
	),
	Section2: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "error order"),
	),
	Section3: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "derivative order"),
	),
	Section4: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "philosophy"),
	),
	Section5: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "inputs"),
	),
}

// GetHelpKeys implements NumeTabContent.
func (m *DerivativeModel) GetHelpKeys() help.KeyMap {
	return derivativeKeys
}

var _ (NumeTabContent) = (*DerivativeModel)(nil)

func NewDerivativeModel(theme *Theme) *DerivativeModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(70),
	)

	// Create delta input
	deltaInput := textinput.New()
	deltaInput.Placeholder = "0.001"
	deltaInput.CharLimit = 20
	deltaInput.SetValue("0.001")

	// Create test point input
	testPointInput := textinput.New()
	testPointInput.Placeholder = "1.0"
	testPointInput.CharLimit = 20
	testPointInput.SetValue("1.0")

	return &DerivativeModel{
		focusedSection: 0,
		functionOptions: []string{
			"Polynomial: customizable degree",
			"Exponential: f(x) = e^x",
			"Trigonometric: f(x) = sin(x)",
			"Logarithmic: f(x) = ln(x)",
			"Rational: f(x) = 1/x",
			"Composite: f(x) = sin(x²)",
		},
		selectedFunction: 0,
		polynomialOrder:  3, // default to cubic
		derivativeOrder:  1,
		philosophy:       2, // central
		deltaInput:       deltaInput,
		testPointInput:   testPointInput,
		delta:            0.001,
		testPoint:        1.0,
		renderer:         renderer,
		Theme:            theme,
	}
}

func (m *DerivativeModel) Init() tea.Cmd {
	return nil
}

func (m *DerivativeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
		case "c":
			m.generateResult()
			return m, nil
		case "x":
			m.showExplanation = !m.showExplanation
			if m.showExplanation && m.explanation == "" {
				m.generateExplanation()
			}
			return m, nil
		case "r":
			return NewDerivativeModel(m.Theme), nil
		case "1":
			m.focusedSection = 0 // Function Selection
			return m, nil
		case "2":
			m.focusedSection = 1 // Error Order
			return m, nil
		case "3":
			m.focusedSection = 2 // Derivative Order
			return m, nil
		case "4":
			m.focusedSection = 3 // Philosophy
			return m, nil
		case "5":
			m.focusedSection = 4 // Inputs
			return m, nil
		}

		// Handle input for text inputs
		if m.focusedSection == 4 {
			var cmd tea.Cmd
			m.deltaInput, cmd = m.deltaInput.Update(msg)
			if val, err := strconv.ParseFloat(m.deltaInput.Value(), 64); err == nil {
				m.delta = val
			}
			cmds = append(cmds, cmd)

			m.testPointInput, cmd = m.testPointInput.Update(msg)
			if val, err := strconv.ParseFloat(m.testPointInput.Value(), 64); err == nil {
				m.testPoint = val
			}
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *DerivativeModel) handleUp() *DerivativeModel {
	switch m.focusedSection {
	case 0: // Function selection
		if m.selectedFunction > 0 {
			m.selectedFunction--
		}
	case 1: // Error order
		if m.polynomialOrder > 1 {
			m.polynomialOrder--
		}
	case 2: // Derivative order
		if m.derivativeOrder > 1 {
			m.derivativeOrder--
		}
	case 3: // Philosophy
		if m.philosophy > 0 {
			m.philosophy--
		}
	case 4: // Inputs - focus delta input
		m.deltaInput.Focus()
		m.testPointInput.Blur()
	}
	return m
}

func (m *DerivativeModel) handleDown() *DerivativeModel {
	switch m.focusedSection {
	case 0: // Function selection
		if m.selectedFunction < len(m.functionOptions)-1 {
			m.selectedFunction++
		}
	case 1: // Error order
		if m.polynomialOrder < 4 {
			m.polynomialOrder++
		}
	case 2: // Derivative order
		if m.derivativeOrder < 3 {
			m.derivativeOrder++
		}
	case 3: // Philosophy
		if m.philosophy < 2 {
			m.philosophy++
		}
	case 4: // Inputs - focus test point input
		m.deltaInput.Blur()
		m.testPointInput.Focus()
	}
	return m
}

func (m *DerivativeModel) handleLeft() *DerivativeModel {
	switch m.focusedSection {
	case 4: // Inputs - focus delta input
		m.deltaInput.Focus()
		m.testPointInput.Blur()
	}
	return m
}

func (m *DerivativeModel) handleRight() *DerivativeModel {
	switch m.focusedSection {
	case 4: // Inputs - focus test point input
		m.deltaInput.Blur()
		m.testPointInput.Focus()
	}
	return m
}

func (m *DerivativeModel) handleEnter() *DerivativeModel {
	// Generate result when enter is pressed
	m.generateResult()
	return m
}

func (m *DerivativeModel) View() string {
	// Create two-column layout: left side navigation, right side content
	leftWidth := 40
	rightWidth := 60

	// Left side - Section navigation
	leftContent := m.renderSectionNavigation()

	// Right side - Markdown content based on focused section
	rightContent := m.renderSectionContent()

	// Join horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(leftWidth).Render(leftContent),
		lipgloss.NewStyle().Width(rightWidth).Render(rightContent),
	)

	// Results section at the bottom if available
	if m.result != "" {
		resultsSection := m.renderResults()
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", resultsSection)
	}

	return content
}

func (m *DerivativeModel) renderSectionNavigation() string {
	var sections []string

	// Section names with tilde formatting
	sectionNames := []string{
		"Function Selection",
		"Error Order",
		"Derivative Order",
		"Philosophy",
		"Inputs",
	}

	for i, name := range sectionNames {
		var style lipgloss.Style
		if i == m.focusedSection {
			// Use focused title color from theme
			style = lipgloss.NewStyle().
				Foreground(m.Theme.Focused.Title.GetForeground()).
				Bold(true)
		} else {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666"))
		}

		// Format with tildes
		formattedName := fmt.Sprintf("~ %s ~", name)
		sections = append(sections, style.Render(formattedName))

		// Add content based on section
		switch i {
		case 0: // Function Selection
			for j, function := range m.functionOptions {
				prefix := "  "
				if j == m.selectedFunction {
					prefix = "▶ "
				}
				functionName := strings.Split(function, ":")[0]
				sections = append(sections, fmt.Sprintf("%s%s", prefix, functionName))
			}
		case 1: // Error Order
			orderNames := []string{"Linear", "Quadratic", "Cubic", "Quartic"}
			for j, orderName := range orderNames {
				prefix := "  "
				if j+1 == m.polynomialOrder {
					prefix = "▶ "
				}
				sections = append(sections, fmt.Sprintf("%s%s (degree %d)", prefix, orderName, j+1))
			}
		case 2: // Derivative Order
			orderOptions := []string{"First", "Second", "Third"}
			for j, order := range orderOptions {
				prefix := "  "
				if j+1 == m.derivativeOrder {
					prefix = "▶ "
				}
				sections = append(sections, fmt.Sprintf("%s%s", prefix, order))
			}
		case 3: // Philosophy
			philosophyOptions := []string{"Forward", "Backward", "Central"}
			for j, phil := range philosophyOptions {
				prefix := "  "
				if j == m.philosophy {
					prefix = "▶ "
				}
				sections = append(sections, fmt.Sprintf("%s%s", prefix, phil))
			}
		case 4: // Inputs
			sections = append(sections, fmt.Sprintf("  Delta: %s", m.deltaInput.View()))
			sections = append(sections, fmt.Sprintf("  Test Point: %s", m.testPointInput.View()))
		}
		sections = append(sections, "") // Add spacing
	}

	return strings.Join(sections, "\n")
}

func (m *DerivativeModel) renderSectionContent() string {
	var content string

	switch m.focusedSection {
	case 0: // Function Selection
		content = `# Function Selection (Press 1)

Choose the mathematical function for derivative calculation:

## Available Functions

- **Polynomial**: Customizable polynomial functions (linear to quartic)
- **Exponential**: f(x) = e^x
- **Trigonometric**: f(x) = sin(x)
- **Logarithmic**: f(x) = ln(x)
- **Rational**: f(x) = 1/x
- **Composite**: f(x) = sin(x²)

Use ↑/↓ arrows to select a function type.
`
	case 1: // Error Order
		content = `# Error Order (Press 2)

Choose the degree of the polynomial function:

## Available Orders

- **Linear (degree 1)**: f(x) = 2x + 1
- **Quadratic (degree 2)**: f(x) = x² + 2x + 1  
- **Cubic (degree 3)**: f(x) = x³ + 2x² - x + 1
- **Quartic (degree 4)**: f(x) = x⁴ + x³ + 2x² - x + 1

Use ↑/↓ arrows to select the polynomial degree.

**Note**: This setting only applies to polynomial functions.`
	case 2: // Derivative Order
		content = `# Derivative Order (Press 3)

Select the order of derivative to calculate:

## Available Orders

- **First derivative**: f'(x) - Rate of change
- **Second derivative**: f''(x) - Concavity and acceleration  
- **Third derivative**: f'''(x) - Rate of change of acceleration

Use ↑/↓ arrows to select the derivative order.

## Mathematical Notation
- First: f'(x) or df/dx
- Second: f''(x) or d²f/dx²
- Third: f'''(x) or d³f/dx³`
	case 3: // Philosophy
		content = `# Philosophy (Press 4)

Choose the finite difference method for numerical differentiation:

## Available Methods

- **Forward Difference**: Uses f(x+h) - f(x)
  - Good for left boundary points
  - First-order accurate: O(h)

- **Backward Difference**: Uses f(x) - f(x-h)  
  - Good for right boundary points
  - First-order accurate: O(h)

- **Central Difference**: Uses f(x+h) - f(x-h)
  - Most accurate for interior points
  - Second-order accurate: O(h²)

Use ↑/↓ arrows to select the difference method.

**Recommended**: Central difference for most applications.`
	case 4: // Inputs
		content = `# Inputs (Press 5)

Configure the numerical calculation parameters:

## Delta (h)
The step size for finite difference calculation.
- Smaller values: More accurate but prone to numerical errors
- Larger values: Less accurate but more stable
- Typical range: 1e-6 to 1e-2
- **Default**: 0.001

## Test Point
The x-coordinate where the derivative is evaluated.
- Choose based on your function's domain
- Avoid singularities (e.g., x=0 for 1/x)
- **Default**: 1.0

Use ←/→ arrows to switch between input fields.
Press **Enter** or **c** to calculate the derivative.`
	}

	// Render with glamour
	if rendered, err := m.renderer.Render(content); err == nil {
		return rendered
	}
	return content
}

func (m *DerivativeModel) renderResults() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1).
		Width(100)

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

	// Evaluate at test point
	derivativeValue := derivativeExpr(m.testPoint)

	functionName := strings.Split(m.functionOptions[m.selectedFunction], ":")[0]
	philosophyName := []string{"Forward", "Backward", "Central"}[m.philosophy]

	// Add polynomial order info if polynomial is selected
	functionDescription := functionName
	if m.selectedFunction == 0 { // Polynomial
		orderNames := []string{"Linear", "Quadratic", "Cubic", "Quartic"}
		if m.polynomialOrder >= 1 && m.polynomialOrder <= 4 {
			functionDescription = fmt.Sprintf("%s (%s - degree %d)", functionName, orderNames[m.polynomialOrder-1], m.polynomialOrder)
		}
	}

	m.result = fmt.Sprintf(`Function: %s
Derivative Order: %s
Philosophy: %s difference
Delta (h): %.6f

Test point: x = %.1f
Derivative value: %.6f`,
		functionDescription,
		m.getDerivativeOrderText(),
		philosophyName,
		m.delta,
		m.testPoint,
		derivativeValue)
}

func (m *DerivativeModel) getTestPoint() float64 {
	return m.testPoint
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
	case 0: // Polynomial - use the selected polynomial order
		switch m.polynomialOrder {
		case 1: // Linear: f(x) = 2x + 1
			m.functionExpr = func(x float64) float64 {
				return 2*x + 1
			}
		case 2: // Quadratic: f(x) = x² + 2x + 1
			m.functionExpr = func(x float64) float64 {
				return x*x + 2*x + 1
			}
		case 3: // Cubic: f(x) = x³ + 2x² - x + 1
			m.functionExpr = func(x float64) float64 {
				return x*x*x + 2*x*x - x + 1
			}
		case 4: // Quartic: f(x) = x⁴ + x³ + 2x² - x + 1
			m.functionExpr = func(x float64) float64 {
				return x*x*x*x + x*x*x + 2*x*x - x + 1
			}
		default: // Default to cubic
			m.functionExpr = func(x float64) float64 {
				return x*x*x + 2*x*x - x + 1
			}
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
- **Function**: %s  
- **Order**: %s
- **Delta**: %.6f

## Parameters
- **Test Point**: %.1f
`,
			strings.Title(philosophyName),
			philosophyName,
			strings.Split(m.functionOptions[m.selectedFunction], ":")[0],
			m.getDerivativeOrderText(),
			m.delta,
			m.testPoint)
	}
}

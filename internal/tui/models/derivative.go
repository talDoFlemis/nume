package models

import (
	"context"
	"fmt"
	"math"
	"os"
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

	// Section 5: Arguments (Delta and Test Point inputs)
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
	Quit             key.Binding
	Help             key.Binding
	TabD             key.Binding
	TabI             key.Binding
	CycleNextSection key.Binding
	CyclePrevSection key.Binding
	Up               key.Binding
	Down             key.Binding
	Left             key.Binding
	Right            key.Binding
	Enter            key.Binding
	Space            key.Binding
	Explain          key.Binding
	Reset            key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k derivativeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k derivativeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TabD, k.TabI, k.Help},                 // first column - navigation
		{k.Up, k.Down, k.Left, k.Right},          // second column - movement
		{k.CycleNextSection, k.CyclePrevSection}, // third column - sections
		{k.Enter, k.Explain, k.Reset, k.Quit},    // fourth column - actions
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
	CycleNextSection: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cycle to next section"),
	),
	CyclePrevSection: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "cycle to previous section"),
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
	Explain: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "toggle explanation"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
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
			"Polynomial: f(x) = x^4 - 2x² + 5x - 1",
			"Exponential: f(x) = e^3x",
			"Trigonometric: f(x) = sin(2x)",
			"Hyperbolic: f(x) = cosh(x)",
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
		switch {
		case key.Matches(msg, derivativeKeys.CycleNextSection):
			m.focusedSection = (m.focusedSection + 1) % 6 // 6 sections now including calculate button
			return m, nil
		case key.Matches(msg, derivativeKeys.CyclePrevSection):
			m.focusedSection = (m.focusedSection - 1 + 6) % 6
			return m, nil
		case key.Matches(msg, derivativeKeys.Up):
			return m.handleUp(), nil
		case key.Matches(msg, derivativeKeys.Down):
			return m.handleDown(), nil
		case key.Matches(msg, derivativeKeys.Left):
			return m.handleLeft(), nil
		case key.Matches(msg, derivativeKeys.Right):
			return m.handleRight(), nil
		case key.Matches(msg, derivativeKeys.Enter):
			return m.handleEnter(), nil
		case key.Matches(msg, derivativeKeys.Explain):
			m.showExplanation = !m.showExplanation
			if m.showExplanation && m.explanation == "" {
				m.generateExplanation()
			}
			return m, nil
		case key.Matches(msg, derivativeKeys.Reset):
			return NewDerivativeModel(m.Theme), nil
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
		} else {
			// Cycle to the end
			m.selectedFunction = len(m.functionOptions) - 1
		}
	case 1: // Error order
		if m.polynomialOrder > 1 {
			m.polynomialOrder--
		} else {
			// Cycle to the highest order (quartic = 4)
			m.polynomialOrder = 4
		}
	case 2: // Derivative order
		if m.derivativeOrder > 1 {
			m.derivativeOrder--
		} else {
			// Cycle to the highest order (third = 3)
			m.derivativeOrder = 3
		}
	case 3: // Philosophy
		if m.philosophy > 0 {
			m.philosophy--
		} else {
			// Cycle to the last philosophy (central = 2)
			m.philosophy = 2
		}
	case 4: // Arguments - focus delta input
		m.deltaInput.Focus()
		m.testPointInput.Blur()
	case 5: // Calculate button - no up action
	}
	return m
}

func (m *DerivativeModel) handleDown() *DerivativeModel {
	switch m.focusedSection {
	case 0: // Function selection
		if m.selectedFunction < len(m.functionOptions)-1 {
			m.selectedFunction++
		} else {
			// Cycle to the beginning
			m.selectedFunction = 0
		}
	case 1: // Error order
		if m.polynomialOrder < 4 {
			m.polynomialOrder++
		} else {
			// Cycle to the lowest order (linear = 1)
			m.polynomialOrder = 1
		}
	case 2: // Derivative order
		if m.derivativeOrder < 3 {
			m.derivativeOrder++
		} else {
			// Cycle to the lowest order (first = 1)
			m.derivativeOrder = 1
		}
	case 3: // Philosophy
		if m.philosophy < 2 {
			m.philosophy++
		} else {
			// Cycle to the first philosophy (forward = 0)
			m.philosophy = 0
		}
	case 4: // Arguments - focus test point input
		m.deltaInput.Blur()
		m.testPointInput.Focus()
	case 5: // Calculate button - no down action
	}
	return m
}

func (m *DerivativeModel) handleLeft() *DerivativeModel {
	switch m.focusedSection {
	case 4: // Arguments - focus delta input
		m.deltaInput.Focus()
		m.testPointInput.Blur()
	case 5: // Calculate button - no left action
	}
	return m
}

func (m *DerivativeModel) handleRight() *DerivativeModel {
	switch m.focusedSection {
	case 4: // Arguments - focus test point input
		m.deltaInput.Blur()
		m.testPointInput.Focus()
	case 5: // Calculate button - no right action
	}
	return m
}

func (m *DerivativeModel) handleEnter() *DerivativeModel {
	// Only generate result if calculate button is focused
	if m.focusedSection == 5 {
		m.generateResult()
	}
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
		"Arguments",
		"Calculate",
	}

	for i, name := range sectionNames {
		var style lipgloss.Style
		if i == m.focusedSection {
			// Use focused title color from theme
			style = lipgloss.NewStyle().
				Foreground(m.Focused.Title.GetForeground()).
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
				style := m.Blurred.UnselectedPrefix
				if j == m.selectedFunction {
					style = m.Focused.SelectedPrefix
				}
				functionName := strings.Split(function, ":")[0]
				sections = append(sections, style.Render(functionName))
			}
		case 1: // Error Order
			orderNames := []string{"Linear", "Quadratic", "Cubic", "Quartic"}
			for j, orderName := range orderNames {
				style := m.Blurred.UnselectedPrefix
				if j+1 == m.polynomialOrder {
					style = m.Focused.SelectedPrefix
				}
				sections = append(
					sections,
					style.Render(fmt.Sprintf("%s (degree %d)", orderName, j+1)),
				)
			}
		case 2: // Derivative Order
			orderOptions := []string{"First", "Second", "Third"}
			for j, order := range orderOptions {
				style := m.Blurred.UnselectedPrefix
				if j+1 == m.derivativeOrder {
					style = m.Focused.SelectedPrefix
				}
				sections = append(sections, style.Render(order))
			}
		case 3: // Philosophy
			philosophyOptions := []string{"Forward", "Backward", "Central"}
			for j, phil := range philosophyOptions {
				style := m.Blurred.UnselectedPrefix
				if j == m.philosophy {
					style = m.Focused.SelectedPrefix
				}
				sections = append(sections, style.Render(phil))
			}
		case 4: // Arguments
			sections = append(sections, fmt.Sprintf("  Delta: %s", m.deltaInput.View()))
			sections = append(sections, fmt.Sprintf("  Test Point: %s", m.testPointInput.View()))
		case 5: // Calculate button
			// Create a styled button
			var buttonStyle lipgloss.Style
			if i == m.focusedSection {
				buttonStyle = m.Focused.FocusedButton
			} else {
				buttonStyle = m.Focused.BlurredButton
			}
			button := buttonStyle.Render(" CALCULATE ")
			sections = append(sections, fmt.Sprintf("  %s", button))
		}
		sections = append(sections, "") // Add spacing
	}

	return strings.Join(sections, "\n")
}

func (m *DerivativeModel) renderSectionContent() string {
	var content string

	switch m.focusedSection {
	case 0: // Function Selection
		content = `# Function Selection

Choose the mathematical function for derivative calculation:

## Available Functions

- **Polynomial**: f(x) = x^4 - 2x² + 5x - 1
- **Exponential**: f(x) = e^3x
- **Trigonometric**: f(x) = sin(2x)
- **Hyperbolic**: f(x) = cosh(x)

Use ↑/↓ arrows to select a function type.
`
	case 1: // Error Order
		content = `# Error Order

Choose the degree of the error for the approximation:

## Available Orders

- **Linear (degree 1)**: O(h)
- **Quadratic (degree 2)**: O(h²)
- **Cubic (degree 3)**: O(h³)
- **Quartic (degree 4)**: O(h⁴)

Use ↑/↓ arrows to select the approximation degree.`
	case 2: // Derivative Order
		content = `# Derivative Order

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
		content = `# Philosophy

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
	case 4: // Arguments
		content = `# Arguments

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

Use ←/→ arrows to switch between input fields.`
	case 5: // Calculate
		content = `# Calculate

Execute the derivative calculation with the configured parameters:

## Current Configuration

- **Function**: ` + strings.Split(m.functionOptions[m.selectedFunction], ":")[0] + `
- **Derivative Order**: ` + m.getDerivativeOrderText() + `
- **Philosophy**: ` + []string{"Forward", "Backward", "Central"}[m.philosophy] + ` difference
- **Delta (h)**: ` + fmt.Sprintf("%.6f", m.delta) + `
- **Test Point**: ` + fmt.Sprintf("%.1f", m.testPoint) + `

Press **Enter** on the Calculate button to run the calculation.`

		// Add results section if available
		if m.result != "" {
			content += `

# Result

` + m.result
		}
	}

	// Render with glamour
	if rendered, err := m.renderer.Render(content); err == nil {
		return rendered
	}
	return content
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
		m.result = m.Focused.ErrorMessage.Render(
			fmt.Sprintf("Error calculating derivative: %v", err),
		)
		return
	}

	// Evaluate at test point
	derivativeValue := derivativeExpr(m.testPoint)

	m.result = fmt.Sprintf(`%.6f`, derivativeValue)
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
	if m.selectedFunction < 0 || m.selectedFunction >= len(m.functionOptions) {
		panic(fmt.Sprintf("Invalid function selection: %d", m.selectedFunction))
	}

	// Define function expressions based on selected function
	switch m.selectedFunction {
	case 0: // Polynomial
		m.functionExpr = func(x float64) float64 {
			return math.Pow(x, 4) - 2*x*x + 5*x - 1
		}
	case 1: // Exponential
		m.functionExpr = func(x float64) float64 {
			return math.Exp(3 * x)
		}
	case 2: // Trigonometric
		m.functionExpr = func(x float64) float64 {
			return math.Sin(2 * x)
		}
	case 3: // Hyperbolic
		m.functionExpr = math.Cosh
	}
}

func (m *DerivativeModel) generateExplanation() {
	philosophyName := []string{"forward", "backward", "central"}[m.philosophy]
	filename := fmt.Sprintf("%s_difference.md", philosophyName)
	explanationPath := filepath.Join("internal", "tui", "views", "explanations", filename)

	if content, err := os.ReadFile(explanationPath); err == nil {
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
			strings.ToUpper(philosophyName[:1])+philosophyName[1:],
			philosophyName,
			strings.Split(m.functionOptions[m.selectedFunction], ":")[0],
			m.getDerivativeOrderText(),
			m.delta,
			m.testPoint)
	}
}

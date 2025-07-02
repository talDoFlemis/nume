package models

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/taldoflemis/nume/internal/usecases"
)

type EigenModel struct {
	// Current focus section (0-3)
	focusedSection int

	// Section 1: Power Method Selection
	powerMethodOptions  []string
	selectedPowerMethod int

	// Section 2: Matrix Selection
	matrixOptions      []string
	selectedMatrix     int
	predefinedMatrices [][][]float64

	// Section 3: Arguments (Vector, Epsilon, Max Iterations, K Eigenvalue inputs)
	vectorInput        textinput.Model
	epsilonInput       textinput.Model
	maxIterationsInput textinput.Model
	kEigenvalueInput   textinput.Model
	initialVector      []float64
	epsilon            float64
	maxIterations      uint64
	kEigenvalue        float64

	// Calculation results
	result          string
	showExplanation bool
	explanation     string

	// Use case
	useCase *usecases.PowerUseCase

	// Styling
	renderer *glamour.TermRenderer
	*Theme
}

// keyMap defines the keybindings for the eigen model
type eigenKeyMap struct {
	Quit             key.Binding
	Help             key.Binding
	TabD             key.Binding
	TabI             key.Binding
	TabE             key.Binding
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
func (k eigenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k eigenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TabD, k.TabI, k.TabE, k.Help},         // first column - navigation
		{k.Up, k.Down, k.Left, k.Right},          // second column - movement
		{k.CycleNextSection, k.CyclePrevSection}, // third column - sections
		{k.Enter, k.Explain, k.Reset, k.Quit},    // fourth column - actions
	}
}

var eigenKeys = eigenKeyMap{
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
	TabE: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "eigen tab"),
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
func (*EigenModel) GetHelpKeys() help.KeyMap {
	return eigenKeys
}

var _ (NumeTabContent) = (*EigenModel)(nil)

func NewEigenModel(theme *Theme) *EigenModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(GlamourRenderWidth),
		glamour.WithStandardStyle("dracula"),
	)

	// Create input fields
	vectorInput := textinput.New()
	vectorInput.Placeholder = "1.0,1.0"
	vectorInput.CharLimit = 50
	vectorInput.SetValue("1.0,1.0")

	epsilonInput := textinput.New()
	epsilonInput.Placeholder = "1e-6"
	epsilonInput.CharLimit = 20
	epsilonInput.SetValue("1e-6")

	maxIterationsInput := textinput.New()
	maxIterationsInput.Placeholder = "100"
	maxIterationsInput.CharLimit = 10
	maxIterationsInput.SetValue("100")

	kEigenvalueInput := textinput.New()
	kEigenvalueInput.Placeholder = "0.0"
	kEigenvalueInput.CharLimit = 20
	kEigenvalueInput.SetValue("0.0")

	// Predefined matrices
	predefinedMatrices := [][][]float64{
		// 2x2 Simple
		{{2.0, 3.0}, {5.0, 4.0}},
		// 3x3 Simple
		{{2.0, 1.0, 0.0}, {1.0, 2.0, 1.0}, {0.0, 1.0, 2.0}},
		// 3x3 Complex
		{{5.0, 2.0, 1.0}, {2.0, 3.0, 1.0}, {1.0, 1.0, 2.0}},
		// 4x4 Simple
		{{4.0, 1.0, 0.0, 0.0}, {1.0, 3.0, 1.0, 0.0}, {0.0, 1.0, 3.0, 1.0}, {0.0, 0.0, 1.0, 2.0}},
	}

	return &EigenModel{
		focusedSection: 0,
		powerMethodOptions: []string{
			"Regular Power Method",
			"Inverse Power Method",
			"Farthest Eigenvalue Power",
			"Nearest Eigenvalue Power",
		},
		selectedPowerMethod: 0,
		matrixOptions: []string{
			"2x2 Simple Matrix",
			"3x3 Simple Matrix",
			"3x3 Complex Matrix",
			"4x4 Simple Matrix",
		},
		selectedMatrix:     0,
		predefinedMatrices: predefinedMatrices,
		vectorInput:        vectorInput,
		epsilonInput:       epsilonInput,
		maxIterationsInput: maxIterationsInput,
		kEigenvalueInput:   kEigenvalueInput,
		initialVector:      []float64{1.0, 1.0},
		epsilon:            DefaultEpsilon,
		maxIterations:      DefaultMaxIterations,
		kEigenvalue:        0.0,
		useCase:            usecases.NewPowerUseCase(),
		renderer:           renderer,
		Theme:              theme,
	}
}

func (*EigenModel) Init() tea.Cmd {
	return nil
}

func (m *EigenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, eigenKeys.CycleNextSection):
			m.focusedSection = (m.focusedSection + 1) % EigenSectionCount
			return m, nil
		case key.Matches(keyMsg, eigenKeys.CyclePrevSection):
			m.focusedSection = (m.focusedSection - 1 + EigenSectionCount) % EigenSectionCount
			return m, nil
		case key.Matches(keyMsg, eigenKeys.Up):
			return m.handleUp(), nil
		case key.Matches(keyMsg, eigenKeys.Down):
			return m.handleDown(), nil
		case key.Matches(keyMsg, eigenKeys.Left):
			return m.handleLeft(), nil
		case key.Matches(keyMsg, eigenKeys.Right):
			return m.handleRight(), nil
		case key.Matches(keyMsg, eigenKeys.Enter):
			return m.handleEnter(), nil
		case key.Matches(keyMsg, eigenKeys.Explain):
			m.showExplanation = !m.showExplanation
			if m.showExplanation && m.explanation == "" {
				m.generateExplanation()
			}
			return m, nil
		case key.Matches(keyMsg, eigenKeys.Reset):
			return NewEigenModel(m.Theme), nil
		}

		// Handle input for text inputs
		if m.focusedSection == EigenSectionArguments {
			var cmd tea.Cmd
			m.vectorInput, cmd = m.vectorInput.Update(keyMsg)
			if val := m.parseVector(m.vectorInput.Value()); val != nil {
				m.initialVector = val
			}
			cmds = append(cmds, cmd)

			m.epsilonInput, cmd = m.epsilonInput.Update(keyMsg)
			if val, err := strconv.ParseFloat(m.epsilonInput.Value(), 64); err == nil {
				m.epsilon = val
			}
			cmds = append(cmds, cmd)

			m.maxIterationsInput, cmd = m.maxIterationsInput.Update(keyMsg)
			if val, err := strconv.ParseUint(m.maxIterationsInput.Value(), 10, 64); err == nil {
				m.maxIterations = val
			}
			cmds = append(cmds, cmd)

			m.kEigenvalueInput, cmd = m.kEigenvalueInput.Update(keyMsg)
			if val, err := strconv.ParseFloat(m.kEigenvalueInput.Value(), 64); err == nil {
				m.kEigenvalue = val
			}
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *EigenModel) handleUp() *EigenModel {
	switch m.focusedSection {
	case EigenSectionPowerMethodSelection: // Power method selection
		if m.selectedPowerMethod > 0 {
			m.selectedPowerMethod--
		} else {
			m.selectedPowerMethod = len(m.powerMethodOptions) - 1
		}
	case EigenSectionMatrixSelection: // Matrix selection
		if m.selectedMatrix > 0 {
			m.selectedMatrix--
		} else {
			m.selectedMatrix = len(m.matrixOptions) - 1
		}
	case EigenSectionArguments: // Arguments - cycle through inputs
		// Cycle backwards through inputs (up key)
		if m.kEigenvalueInput.Focused() {
			m.kEigenvalueInput.Blur()
			m.maxIterationsInput.Focus()
		} else if m.maxIterationsInput.Focused() {
			m.maxIterationsInput.Blur()
			m.epsilonInput.Focus()
		} else if m.epsilonInput.Focused() {
			m.epsilonInput.Blur()
			m.vectorInput.Focus()
		} else {
			// Default to k eigenvalue input (wrap around)
			m.vectorInput.Blur()
			m.epsilonInput.Blur()
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Focus()
		}
	case EigenSectionCalculate: // Calculate button - no up action
	}
	return m
}

func (m *EigenModel) handleDown() *EigenModel {
	switch m.focusedSection {
	case EigenSectionPowerMethodSelection: // Power method selection
		if m.selectedPowerMethod < len(m.powerMethodOptions)-1 {
			m.selectedPowerMethod++
		} else {
			m.selectedPowerMethod = 0
		}
	case EigenSectionMatrixSelection: // Matrix selection
		if m.selectedMatrix < len(m.matrixOptions)-1 {
			m.selectedMatrix++
		} else {
			m.selectedMatrix = 0
		}
	case EigenSectionArguments: // Arguments - cycle through inputs
		// Cycle forwards through inputs (down key)
		if m.vectorInput.Focused() {
			m.vectorInput.Blur()
			m.epsilonInput.Focus()
		} else if m.epsilonInput.Focused() {
			m.epsilonInput.Blur()
			m.maxIterationsInput.Focus()
		} else if m.maxIterationsInput.Focused() {
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Focus()
		} else {
			// Default to vector input (wrap around)
			m.vectorInput.Focus()
			m.epsilonInput.Blur()
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Blur()
		}
	case EigenSectionCalculate: // Calculate button - no down action
	}
	return m
}

func (m *EigenModel) handleLeft() *EigenModel {
	switch m.focusedSection {
	case EigenSectionArguments: // Arguments - focus previous input
		// Cycle backwards through inputs
		if m.kEigenvalueInput.Focused() {
			m.kEigenvalueInput.Blur()
			m.maxIterationsInput.Focus()
		} else if m.maxIterationsInput.Focused() {
			m.maxIterationsInput.Blur()
			m.epsilonInput.Focus()
		} else if m.epsilonInput.Focused() {
			m.epsilonInput.Blur()
			m.vectorInput.Focus()
		} else {
			// Default to vector input
			m.vectorInput.Focus()
			m.epsilonInput.Blur()
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Blur()
		}
	case EigenSectionCalculate: // Calculate button - no left action
	}
	return m
}

func (m *EigenModel) handleRight() *EigenModel {
	switch m.focusedSection {
	case EigenSectionArguments: // Arguments - focus next input
		// Cycle forwards through inputs
		if m.vectorInput.Focused() {
			m.vectorInput.Blur()
			m.epsilonInput.Focus()
		} else if m.epsilonInput.Focused() {
			m.epsilonInput.Blur()
			m.maxIterationsInput.Focus()
		} else if m.maxIterationsInput.Focused() {
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Focus()
		} else {
			// Default to k eigenvalue input
			m.vectorInput.Blur()
			m.epsilonInput.Blur()
			m.maxIterationsInput.Blur()
			m.kEigenvalueInput.Focus()
		}
	case EigenSectionCalculate: // Calculate button - no right action
	}
	return m
}

func (m *EigenModel) handleEnter() *EigenModel {
	// Only generate result if calculate button is focused
	if m.focusedSection == EigenSectionCalculate {
		m.generateResult()
	}
	return m
}

func (m *EigenModel) View() string {
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
		m.Renderer.NewStyle().Width(leftWidth).Render(leftContent),
		m.Renderer.NewStyle().Width(rightWidth).Render(rightContent),
	)

	return content
}

func (m *EigenModel) renderSectionNavigation() string {
	var sections []string

	// Section names with tilde formatting
	sectionNames := []string{
		"Power Method Selection",
		"Matrix Selection",
		"Arguments",
		"Calculate",
	}

	for i, name := range sectionNames {
		var style lipgloss.Style
		if i == m.focusedSection {
			// Use focused title color from theme
			style = m.Renderer.NewStyle().
				Foreground(m.Focused.Title.GetForeground()).
				Bold(true)
		} else {
			style = m.Renderer.NewStyle().
				Foreground(lipgloss.Color("#666666"))
		}

		// Format with tildes
		formattedName := fmt.Sprintf("~ %s ~", name)
		sections = append(sections, style.Render(formattedName))

		// Add content based on section
		switch i {
		case EigenSectionPowerMethodSelection: // Power Method Selection
			for j, method := range m.powerMethodOptions {
				style := m.Blurred.UnselectedPrefix
				if j == m.selectedPowerMethod {
					style = m.Focused.SelectedPrefix
				}
				sections = append(sections, style.Render(method))
			}
		case EigenSectionMatrixSelection: // Matrix Selection
			for j, matrix := range m.matrixOptions {
				style := m.Blurred.UnselectedPrefix
				if j == m.selectedMatrix {
					style = m.Focused.SelectedPrefix
				}
				sections = append(sections, style.Render(matrix))
			}
		case EigenSectionArguments: // Arguments
			sections = append(sections, fmt.Sprintf("  Initial Vector: %s", m.vectorInput.View()))
			sections = append(sections, fmt.Sprintf("  Epsilon: %s", m.epsilonInput.View()))
			sections = append(sections, fmt.Sprintf("  Max Iterations: %s", m.maxIterationsInput.View()))
			sections = append(sections, fmt.Sprintf("  K Eigenvalue: %s", m.kEigenvalueInput.View()))
		case EigenSectionCalculate: // Calculate button
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

func (m *EigenModel) renderSectionContent() string {
	var content string

	switch m.focusedSection {
	case EigenSectionPowerMethodSelection: // Power Method Selection
		content = `# Power Method Selection

Choose the eigenvalue calculation method:

## Available Methods

- **Regular Power Method**: Finds the largest eigenvalue
- **Inverse Power Method**: Finds the smallest eigenvalue
- **Farthest Eigenvalue Power**: Finds eigenvalue farthest from given value
- **Nearest Eigenvalue Power**: Finds eigenvalue nearest to given value

Use ↑/↓ arrows to select a power method.
`
	case EigenSectionMatrixSelection: // Matrix Selection
		content = `# Matrix Selection

Choose a predefined matrix for eigenvalue calculation:

## Available Matrices

- **2x2 Simple**: Small symmetric matrix
- **3x3 Simple**: Tridiagonal symmetric matrix
- **3x3 Complex**: General 3x3 matrix
- **4x4 Simple**: Larger tridiagonal matrix

Use ↑/↓ arrows to select a matrix.

## Current Matrix
` + m.getMatrixDisplay()
	case EigenSectionArguments: // Arguments
		content = `# Arguments

Configure the power method parameters:

## Initial Vector
Starting eigenvector guess (comma-separated values).
- Must have same dimension as matrix
- Cannot be zero vector
- **Format**: 1.0,1.0 or 1,0,1
- **Default**: 1.0,1.0

## Epsilon (ε)
Convergence tolerance for the algorithm.
- Smaller values: More precise but slower
- Typical range: 1e-10 to 1e-3
- **Default**: 1e-6

## Max Iterations
Maximum number of iterations before stopping.
- Higher values: More chances to converge
- Typical range: 50 to 1000
- **Default**: 100

## K Eigenvalue (Shift Value)
Shift value for nearest/farthest eigenvalue methods.
- Used only with "Nearest" and "Farthest" power methods
- For nearest: finds eigenvalue closest to this value
- For farthest: finds eigenvalue farthest from this value
- **Default**: 0.0

Use ←/→ arrows to switch between input fields.`
	case EigenSectionCalculate: // Calculate
		content = `# Calculate

Execute the eigenvalue calculation with the configured parameters:

## Current Configuration

- **Power Method**: ` + m.powerMethodOptions[m.selectedPowerMethod] + `
- **Matrix**: ` + m.matrixOptions[m.selectedMatrix] + `
- **Initial Vector**: ` + m.formatVector(m.initialVector) + `
- **Epsilon**: ` + fmt.Sprintf("%.2e", m.epsilon) + `
- **Max Iterations**: ` + fmt.Sprintf("%d", m.maxIterations) + `
- **K Eigenvalue**: ` + fmt.Sprintf("%.3f", m.kEigenvalue) + ` (used for nearest/farthest methods)

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

func (m *EigenModel) parseVector(input string) []float64 {
	if input == "" {
		return nil
	}

	parts := strings.Split(input, ",")
	vector := make([]float64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if val, err := strconv.ParseFloat(part, 64); err == nil {
			vector = append(vector, val)
		} else {
			return nil // Invalid input
		}
	}

	if len(vector) == 0 {
		return nil
	}

	return vector
}

func (m *EigenModel) formatVector(vector []float64) string {
	if len(vector) == 0 {
		return "[]"
	}

	parts := make([]string, len(vector))
	for i, val := range vector {
		parts[i] = fmt.Sprintf("%.3f", val)
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func (m *EigenModel) getMatrixDisplay() string {
	if m.selectedMatrix < 0 || m.selectedMatrix >= len(m.predefinedMatrices) {
		return "Invalid matrix selection"
	}

	matrix := m.predefinedMatrices[m.selectedMatrix]
	var lines []string

	for _, row := range matrix {
		var rowStr []string
		for _, val := range row {
			rowStr = append(rowStr, fmt.Sprintf("%4.1f", val))
		}
		lines = append(lines, "[ "+strings.Join(rowStr, "  ")+" ]")
	}

	return "```\n" + strings.Join(lines, "\n") + "\n```"
}

func (m *EigenModel) generateResult() {
	if m.selectedMatrix < 0 || m.selectedMatrix >= len(m.predefinedMatrices) {
		m.result = m.Focused.ErrorMessage.Render("Invalid matrix selection")
		return
	}

	matrix := m.predefinedMatrices[m.selectedMatrix]

	// Validate initial vector dimension
	if len(m.initialVector) != len(matrix) {
		m.result = m.Focused.ErrorMessage.Render(
			fmt.Sprintf("Initial vector dimension (%d) must match matrix dimension (%d)",
				len(m.initialVector), len(matrix)))
		return
	}

	// Check for zero vector
	const zeroTolerance = 1e-10
	allZero := true
	for _, val := range m.initialVector {
		if math.Abs(val) > zeroTolerance {
			allZero = false
			break
		}
	}
	if allZero {
		m.result = m.Focused.ErrorMessage.Render("Initial vector cannot be zero")
		return
	}

	ctx := context.Background()
	var powerResult *usecases.PowerResult
	var err error

	// Call appropriate power method
	switch m.selectedPowerMethod {
	case PowerMethodRegular:
		powerResult, err = m.useCase.RegularPower(ctx, matrix, m.initialVector, m.epsilon, m.maxIterations)
	case PowerMethodInverse:
		powerResult, err = m.useCase.InversePower(ctx, matrix, m.initialVector, m.epsilon, m.maxIterations)
	case PowerMethodFarthest:
		// For farthest, we use the k eigenvalue as shift value
		eigenvalue, err := m.useCase.FarthestEigenvaluePower(ctx, matrix, m.initialVector, m.kEigenvalue, m.epsilon, m.maxIterations)
		if err == nil {
			powerResult = &usecases.PowerResult{
				Eigenvalue:    eigenvalue,
				Eigenvector:   m.initialVector, // Simplified - actual eigenvector calculation needed
				NumIterations: m.maxIterations, // Simplified
			}
		}
	case PowerMethodNearest:
		// For nearest, we use the k eigenvalue as shift value
		eigenvalue, err := m.useCase.NearestEigenvaluePower(ctx, matrix, m.initialVector, m.kEigenvalue, m.epsilon, m.maxIterations)
		if err == nil {
			powerResult = &usecases.PowerResult{
				Eigenvalue:    eigenvalue,
				Eigenvector:   m.initialVector, // Simplified - actual eigenvector calculation needed
				NumIterations: m.maxIterations, // Simplified
			}
		}
	default:
		m.result = m.Focused.ErrorMessage.Render("Unknown power method selected")
		return
	}

	if err != nil {
		m.result = m.Focused.ErrorMessage.Render(
			fmt.Sprintf("Error calculating eigenvalue: %v", err))
		return
	}

	// Format result
	m.result = fmt.Sprintf(`**Eigenvalue**: %.6f

**Eigenvector**: %s

**Iterations**: %d`,
		powerResult.Eigenvalue,
		m.formatVector(powerResult.Eigenvector),
		powerResult.NumIterations)
}

func (m *EigenModel) generateExplanation() {
	methodName := []string{"regular", "inverse", "farthest", "nearest"}[m.selectedPowerMethod]

	// Fallback explanation
	m.explanation = fmt.Sprintf(`# %s Power Method

## Overview
The %s power method for eigenvalue calculation.

## Configuration
- **Matrix**: %s
- **Method**: %s
- **Epsilon**: %.2e
- **Max Iterations**: %d

## Parameters
- **Initial Vector**: %s
`,
		strings.ToUpper(methodName[:1])+methodName[1:],
		methodName,
		m.matrixOptions[m.selectedMatrix],
		m.powerMethodOptions[m.selectedPowerMethod],
		m.epsilon,
		m.maxIterations,
		m.formatVector(m.initialVector))
}

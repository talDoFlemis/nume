# Nume

Numerical Methods Calculator - A beautiful terminal-based interface for computing derivatives and integrals using various numerical methods.

## Features

### Current Features
- **Interactive TUI** - Beautiful terminal user interface built with Charm Bubble Tea
- **Derivative Calculations** - Compute numerical derivatives using:
  - Forward difference method
  - Backward difference method  
  - Central difference method
  - First, second, and third derivatives
  - Error analysis from linear O(h) to quartic O(h⁴)
- **Built-in Functions** - Pre-defined mathematical functions:
  - Polynomial: f(x) = x³ + 2x² - x + 1
  - Exponential: f(x) = eˣ
  - Trigonometric: f(x) = sin(x)
  - Logarithmic: f(x) = ln(x)
  - Rational: f(x) = 1/x
  - Composite: f(x) = sin(x²)
- **Mathematical Explanations** - Detailed markdown explanations for each numerical method

### Coming Soon
- Integral calculations using various numerical integration methods
- More function types and custom function input
- Advanced error analysis and visualization

## Quick Start

### Running the TUI Application
```bash
# Build and run the TUI
go build -o nume-tui ./cmd/tui && ./nume-tui

# Or using Task (if available)
task tui
```

### Navigation
- **Tab/Shift+Tab**: Switch between Derivatives and Integrals tabs
- **Arrow keys**: Navigate through options
- **Enter**: Select option
- **Numbers (1-4)**: Quick selection for derivative order, philosophy, and error degree
- **F**: Calculate derivative result
- **E**: Toggle mathematical explanation
- **R**: Reset to start over
- **Backspace**: Go back to previous step
- **Q/Ctrl+C**: Quit application

## Architecture

This project uses a clean architecture approach with:
- **TUI Layer**: Bubble Tea models and views for terminal interface
- **Use Cases**: Business logic for numerical calculations
- **Expressions**: Mathematical function representations
- **Strategies**: Different numerical method implementations

## Development

Originally built as a web application, this project has been refactored to use a modern terminal user interface for better developer experience and portability.

Numerical Methods project for the numerical methods 2 class

package models

// Shared constants for the TUI models
const (
	MinimalWidth  = 80
	MinimalHeight = 24
)

// Numerical constants
const (
	// Glamour rendering width
	GlamourRenderWidth = 70

	// Default numerical values
	DefaultPolynomialOrder = 3
	DefaultPhilosophy      = 2 // central difference
	DefaultDelta           = 0.001
	DefaultTestPoint       = 1.0

	// UI layout constants
	SectionCount       = 6
	MaxPolynomialOrder = 4
	MaxDerivativeOrder = 3
	MaxPhilosophyIndex = 2

	// Animation timing
	AnimationDelay  = 200  // milliseconds
	TransitionDelay = 3000 // milliseconds

	// Component padding
	ComponentPadding = 2

	// Function constants used in mathematical expressions
	PolynomialPower     = 4
	ExponentialMultiple = 3
	TrigMultiple        = 2
)

// Section indices
const (
	SectionFunctionSelection = 0
	SectionErrorOrder        = 1
	SectionDerivativeOrder   = 2
	SectionPhilosophy        = 3
	SectionArguments         = 4
	SectionCalculate         = 5
)

// Philosophy indices
const (
	PhilosophyForward  = 0
	PhilosophyBackward = 1
	PhilosophyCentral  = 2
)

// Philosophy case values for switch statements
const (
	PhilosophyBackwardCase = 2
	PhilosophyCentralCase  = 3
)

// Derivative order indices
const (
	DerivativeOrderFirst  = 1
	DerivativeOrderSecond = 2
	DerivativeOrderThird  = 3
)

// Eigen section indices
const (
	EigenSectionPowerMethodSelection = 0
	EigenSectionMatrixSelection      = 1
	EigenSectionArguments            = 2
	EigenSectionCalculate            = 3
)

// Power method indices
const (
	PowerMethodRegular  = 0
	PowerMethodInverse  = 1
	PowerMethodFarthest = 2
	PowerMethodNearest  = 3
)

// Matrix selection indices
const (
	Matrix2x2Simple = 0
	Matrix3x3Simple = 1
	Matrix3x3Complex = 2
	Matrix4x4Simple = 3
)

// Default eigen values
const (
	DefaultEpsilon      = 1e-6
	DefaultMaxIterations = 100
)

// Eigen section count
const (
	EigenSectionCount = 4
)

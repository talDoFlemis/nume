package newtoncotes

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type closedFormulasTestCase struct {
	leftInterval  float64
	rightInterval float64
	expectedValue float64
	tolerance     float64
	simpleExpr    expressions.SingleVariableExpr
}

// Testing trapezoidal rule separately because it has the shittiest approximation
func TestTrapezoidalRule(t *testing.T) {
	// Arrange
	t.Parallel()

	testCases := []closedFormulasTestCase{
		{
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedValue: 1,
			tolerance:     0.3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			leftInterval:  0,
			rightInterval: math.Pi * 2,
			expectedValue: 0,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
	}

	strategy := &TrapezoidalRule{}

	for _, testCase := range testCases {
		testName := fmt.Sprintf("TrapezoidalRule from %.2f to %.2f",
			testCase.leftInterval, testCase.rightInterval)
		t.Run(testName, func(t *testing.T) {
			// Act
			partitionArea, err := strategy.Integrate(
				t.Context(),
				testCase.simpleExpr,
				testCase.leftInterval,
				testCase.rightInterval,
			)

			// Assert
			assert.NoError(t, err, "Expected no error during integration")
			assert.InDelta(t, testCase.expectedValue, partitionArea, testCase.tolerance,
				"Expected partition area to be within tolerance")
		})
	}
}

func TestClosedFormulas(t *testing.T) {
	// Arrange
	t.Parallel()

	closedFormulas := []NewtonCotesStrategy{
		&SimpsonsOneThirdRule{},
		&SimpsonsThreeEighthsRule{},
	}

	testCases := []closedFormulasTestCase{
		{
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedValue: 1,
			tolerance:     10e-1,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			leftInterval:  0,
			rightInterval: math.Pi * 2,
			expectedValue: 0,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			leftInterval:  0,
			rightInterval: 1,
			expectedValue: 0.5,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return x
			},
		},
		// New test cases without expected values
		{
			leftInterval:  0,
			rightInterval: 1,
			expectedValue: 1 / 3.0,
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return x * x // x^2
			},
		},
		{
			leftInterval:  1,
			rightInterval: 2,
			expectedValue: 15 / 4.0,
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return x * x * x // x^3
			},
		},
		{
			leftInterval:  0,
			rightInterval: 1,
			expectedValue: math.E - 1, // e^x - 1
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return math.Exp(x) // e^x
			},
		},
		{
			leftInterval:  1,
			rightInterval: 2,
			expectedValue: math.Log(2), // ln(x)
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return 1 / x // 1/x
			},
		},
		{
			leftInterval:  0,
			rightInterval: math.Pi / 4,
			expectedValue: 1 / (math.Sqrt(2)), // 1/sqrt(2)
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return math.Cos(x) // cos(x)
			},
		},
		{
			leftInterval:  0,
			rightInterval: 2,
			expectedValue: 4 * math.Sqrt(2) / 3.0,
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return math.Sqrt(x) // sqrt(x)
			},
		},
		{
			leftInterval:  -1,
			rightInterval: 1,
			expectedValue: 10 / 3.0,
			tolerance:     0.1,
			simpleExpr: func(x float64) float64 {
				return x*x*x + 2*x*x - x + 1 // x^3 + 2x^2 - x + 1
			},
		},
	}

	for _, testCase := range testCases {
		for _, strategy := range closedFormulas {
			testName := fmt.Sprintf("%s from %.2f to %.2f",
				strategy.Description(), testCase.leftInterval, testCase.rightInterval)
			t.Run(testName, func(t *testing.T) {
				partitionArea, err := strategy.Integrate(
					t.Context(),
					testCase.simpleExpr,
					testCase.leftInterval,
					testCase.rightInterval,
				)

				assert.NoError(t, err, "Expected no error during integration")
				assert.InDelta(t, testCase.expectedValue, partitionArea, testCase.tolerance,
					"Expected partition area to be within tolerance")
			})
		}
	}
}

package newtoncotes

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenFormulas(t *testing.T) {
	// Arrange
	t.Parallel()

	openFormulas := []NewtonCotesStrategy{
		&OpenTrapezoidalRule{},
		&MilneRule{},
		&ThirdDegreeOpenNewtonCotesStrategy{},
	}

	testCases := []formulasTestCase{
		// Test case 1: f(x) = 1/sqrt(x) on [0, 1]
		// Singularity at x=0 (endpoint). True integral value = 2.0.
		// Open methods avoid x=0 directly.
		{
			name:          "1/√x",
			leftInterval:  0.0,
			rightInterval: 1.0,
			expectedValue: 2.0,
			tolerance:     0.7, // A larger tolerance may be needed for approximation near singularities
			simpleExpr: func(x float64) float64 {
				if x <= 0 { // Function is singular at 0 and undefined for negative
					return math.Inf(1)
				}
				return 1.0 / math.Sqrt(x)
			},
		},
		// Test case 2: f(x) = ln(x) on [0, 1]
		// Singularity at x=0 (endpoint). True integral value = -1.0.
		// Open methods avoid x=0 directly.
		{
			name:          "ln(x)",
			leftInterval:  0.0,
			rightInterval: 1.0,
			expectedValue: -1.0,
			tolerance:     0.3, // A larger tolerance may be needed for approximation near singularities
			simpleExpr: func(x float64) float64 {
				if x <= 0 { // log(0) is undefined, log(negative) is undefined
					return math.Inf(-1)
				}
				return math.Log(x)
			},
		},
	}

	for _, testCase := range testCases {
		for _, strategy := range openFormulas {
			testName := fmt.Sprintf("%s - %s from %.2f to %.2f",
				strategy.Description(), testCase.name,
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
}

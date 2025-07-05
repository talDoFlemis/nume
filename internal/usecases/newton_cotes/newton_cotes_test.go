package newtoncotes

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type newtonCotesTestCase struct {
	name               string
	simpleExpr         expressions.SingleVariableExpr
	leftInterval       float64
	rightInterval      float64
	amountOfPartitions uint64
	tolerance          float64
	expectedValue      float64
}

func TestNewtonCotes(t *testing.T) {
	// Arrange
	t.Parallel()

	strategies := []NewtonCotesStrategy{
		&OpenTrapezoidalRule{},
		&MilneRule{},
		&ThirdDegreeOpenNewtonCotesStrategy{},
		&TrapezoidalRule{},
		&SimpsonsOneThirdRule{},
		&SimpsonsThreeEighthsRule{},
	}

	testCases := []newtonCotesTestCase{
		{
			name:          "sin(x)",
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedValue: 1,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
			amountOfPartitions: 1000,
		},
		{
			name:          "sin(x)",
			leftInterval:  0,
			rightInterval: math.Pi * 2,
			expectedValue: 0,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
			amountOfPartitions: 1000,
		},
		{
			name:          "x",
			leftInterval:  0,
			rightInterval: 1,
			expectedValue: 0.5,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return x
			},
			amountOfPartitions: 1000,
		},
	}

	for _, strategy := range strategies {
		for _, testCase := range testCases {
			testName := fmt.Sprintf("%s - %s from %.2f to %.2f using %d partitions",
				strategy.Description(), testCase.name,
				testCase.leftInterval, testCase.rightInterval, testCase.amountOfPartitions)

			t.Run(testName, func(t *testing.T) {
				// Act
				useCase := NewNewtonCotesUseCase(strategy)

				actualArea, err := useCase.Calculate(
					t.Context(),
					testCase.simpleExpr,
					testCase.leftInterval,
					testCase.rightInterval,
					testCase.amountOfPartitions,
				)

				// Assert
				assert.NoError(t, err, "Expected no error during integration")
				assert.InDelta(t, testCase.expectedValue, actualArea, testCase.tolerance)
			})
		}
	}
}

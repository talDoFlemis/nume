package usecases

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name          string
	inputFunc     SingleVariableExpr
	variable      float64
	delta         float64
	expectedValue float64
	tolerance     float64
}

func TestDoubleDerivatives(t *testing.T) {
	t.Parallel()

	strategies := map[string]DifferenceStrategy{
		"DoubleForward":  &ForwardDifferenceStrategy{},
		"DoubleBackward": &BackwardDifferenceStrategy{},
		"DoubleCentral":  &CentralDifferenceStrategy{},
	}

	tests := []testCase{
		{
			name: "d²(x²)/dx² at x=2 should be 2",
			inputFunc: func(x float64) float64 {
				return x * x
			},
			variable:      2.0,
			delta:         1.0,
			expectedValue: 2.0, // second derivative of x² is 2
			tolerance:     0.001,
		},
		{
			name: "d²(x³)/dx² at x=2 should be 12x",
			inputFunc: func(x float64) float64 {
				return math.Pow(x, 3)
			},
			variable:      2.0,
			delta:         0.0001,
			expectedValue: 12.0, // second derivative of x³ is 6x
			tolerance:     0.001,
		},
		{
			name: "d²(sin(x))/dx² at x=0 should be -sin(0)=-0",
			inputFunc: func(x float64) float64 {
				return math.Sin(x)
			},
			variable:      0.0,
			delta:         0.0001,
			expectedValue: 0.0, // second derivative of sin(x) is -sin(x)
			tolerance:     0.001,
		},
		{
			name: "d²(e^x)/dx² at x=1 should be e",
			inputFunc: func(x float64) float64 {
				return math.Exp(x)
			},
			variable:      1.0,
			delta:         0.0001,
			expectedValue: math.E, // second derivative of e^x is e^x
			tolerance:     0.001,
		},
	}

	ctx := context.Background()

	for strategyName, strategy := range strategies {
		strategyName := strategyName
		strategy := strategy

		t.Run(strategyName, func(t *testing.T) {
			for _, tt := range tests {
				tt := tt
				t.Run(fmt.Sprintf("%s_%s", strategyName, tt.name), func(t *testing.T) {
					// t.Parallel()
					firstDerivative, err := strategy.Derivate(ctx, tt.inputFunc, tt.delta)
					assert.NoError(t, err)

					secondDerivative, err := strategy.Derivate(ctx, firstDerivative, tt.delta)
					assert.NoError(t, err)

					directSecondDerivative, err := strategy.DoubleDerivate(
						ctx,
						tt.inputFunc,
						tt.delta,
					)
					assert.NoError(t, err)

					// Act
					result := secondDerivative(tt.variable)
					resultDirect := directSecondDerivative(tt.variable)

					// Assert
					assert.InDelta(t, tt.expectedValue, result, tt.tolerance,
						"Strategy: %s, Test: %s, Expected: %v, Got: %v",
						strategyName, tt.name, tt.expectedValue, result)
					assert.InDelta(t, tt.expectedValue, resultDirect, tt.tolerance,
						"Strategy: %s, Test: %s, Expected: %v, Got: %v on direct derivative",
						strategyName, tt.name, tt.expectedValue, result)
				})
			}
		})
	}
}

func TestDifferenceStrategies(t *testing.T) {
	// Arrange
	t.Parallel()

	strategies := map[string]DifferenceStrategy{
		"Forward":  &ForwardDifferenceStrategy{},
		"Backward": &BackwardDifferenceStrategy{},
		"Central":  &CentralDifferenceStrategy{},
	}

	tests := []testCase{
		{
			name: "x^2 derivative at x=2",
			inputFunc: func(x float64) float64 {
				return x * x
			},
			variable:      2.0,
			delta:         0.0001,
			expectedValue: 4.0, // derivative of x^2 is 2x, at x=2 it's 4
			tolerance:     0.001,
		},
		{
			name: "sin(x) derivative at x=0",
			inputFunc: func(x float64) float64 {
				return math.Sin(x)
			},
			variable:      0.0,
			delta:         0.0001,
			expectedValue: 1.0, // derivative of sin(x) is cos(x), at x=0 it's 1
			tolerance:     0.001,
		},
		{
			name: "e^x derivative at x=1",
			inputFunc: func(x float64) float64 {
				return math.Exp(x)
			},
			variable:      1.0,
			delta:         0.0001,
			expectedValue: math.E, // derivative of e^x is e^x, at x=1 it's e
			tolerance:     0.001,
		},
		{
			name: "x^3 derivative at x=2",
			inputFunc: func(x float64) float64 {
				return math.Pow(x, 3)
			},
			variable:      2.0,
			delta:         0.0001,
			expectedValue: 12.0, // derivative of x^3 is 3x^2, at x=2 it's 12
			tolerance:     0.001,
		},
	}

	ctx := context.Background()

	for strategyName, strategy := range strategies {
		strategyName := strategyName // capture range variable
		strategy := strategy         // capture range variable

		t.Run(strategyName, func(t *testing.T) {
			for _, tt := range tests {
				tt := tt // capture range variable
				t.Run(fmt.Sprintf("%s_%s", strategyName, tt.name), func(t *testing.T) {
					t.Parallel()

					derivative, err := strategy.Derivate(ctx, tt.inputFunc, tt.delta)
					assert.NoError(t, err)

					// Act
					result := derivative(tt.variable)

					// Assert
					assert.InDelta(t, tt.expectedValue, result, tt.tolerance,
						"Strategy: %s, Test: %s, Expected: %v, Got: %v",
						strategyName, tt.name, tt.expectedValue, result)
				})
			}
		})
	}
}

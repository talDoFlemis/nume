package gaussianquadratures

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type gaussLaguerreTestCase struct {
	name         string
	expr         expressions.SingleVariableExpr
	tolerance    float64
	expectedArea float64
}

func TestGaussLaguerre(t *testing.T) {
	// Arrange
	t.Parallel()

	// Test different orders of Gauss-Laguerre quadrature (orders 3 and 4 only)
	strategies := []GaussianQuadrature{}

	// Create strategies for orders 3 and 4
	for order := 3; order <= 4; order++ {
		strategy, err := NewGaussLaguerre(order)
		assert.NoError(t, err, "Should create Gauss-Laguerre strategy without error")
		strategies = append(strategies, strategy)
	}

	testCases := []gaussLaguerreTestCase{
		// Polynomials multiplied by weight function - Gauss-Laguerre integrates f(x)*e^(-x) from 0 to +∞
		{
			name:         "1 (constant)",
			expectedArea: 1.0, // ∫₀^∞ e^(-x) dx = 1
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return 1.0
			},
		},
		{
			name:         "x (linear)",
			expectedArea: 1.0, // ∫₀^∞ x*e^(-x) dx = 1
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:         "x² (quadratic)",
			expectedArea: 2.0, // ∫₀^∞ x²*e^(-x) dx = 2
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:         "x³ (cubic)",
			expectedArea: 6.0, // ∫₀^∞ x³*e^(-x) dx = 6
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:         "x⁴ (quartic)",
			expectedArea: 24.0, // ∫₀^∞ x⁴*e^(-x) dx = 24
			tolerance:    1e-8,
			expr: func(x float64) float64 {
				return x * x * x * x
			},
		},
		{
			name:         "x⁵ (quintic)",
			expectedArea: 120.0, // ∫₀^∞ x⁵*e^(-x) dx = 120
			tolerance:    1e-6,
			expr: func(x float64) float64 {
				return x * x * x * x * x
			},
		},
		// Test with exponential functions
		{
			name:         "e^(-x) (exponential)",
			expectedArea: 0.5, // ∫₀^∞ e^(-x)*e^(-x) dx = 1/2
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return math.Exp(-x)
			},
		},
		// Test with more complex functions - relax tolerance
		{
			name:         "sin(x)",
			expectedArea: 0.5, // ∫₀^∞ sin(x)*e^(-x) dx = 1/2
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			name:         "cos(x)",
			expectedArea: 0.5, // ∫₀^∞ cos(x)*e^(-x) dx = 1/2
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		for _, strategy := range strategies {
			t.Run(fmt.Sprintf("%s - Order %d", testCase.name, strategy.Order()), func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()
				result, err := strategy.Integrate(ctx, testCase.expr, 0.0, math.Inf(1))

				assert.NoError(t, err, "Should integrate without error")
				assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
					"Expected area should match calculated area within tolerance")
			})
		}
	}
}

func TestGaussLaguerreErrorCases(t *testing.T) {
	// Arrange
	t.Parallel()

	strategy, err := NewGaussLaguerre(2)
	assert.NoError(t, err, "Should create Gauss-Laguerre strategy without error")

	testCases := []struct {
		name          string
		leftInterval  float64
		rightInterval float64
		expectedError error
	}{
		{
			name:          "Invalid interval: [1, +∞)",
			leftInterval:  1.0,
			rightInterval: math.Inf(1),
			expectedError: ErrLaguerreIntervalsMustBePositiveInfinite,
		},
		{
			name:          "Invalid interval: [0, 1]",
			leftInterval:  0.0,
			rightInterval: 1.0,
			expectedError: ErrLaguerreIntervalsMustBePositiveInfinite,
		},
		{
			name:          "Invalid interval: [-1, +∞)",
			leftInterval:  -1.0,
			rightInterval: math.Inf(1),
			expectedError: ErrLaguerreIntervalsMustBePositiveInfinite,
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			expr := func(x float64) float64 { return 1.0 }

			_, err := strategy.Integrate(ctx, expr, testCase.leftInterval, testCase.rightInterval)

			assert.Error(t, err, "Should return error for invalid interval")
			assert.Equal(t, testCase.expectedError, err, "Should return correct error type")
		})
	}
}

func TestGaussLaguerreInvalidOrder(t *testing.T) {
	// Arrange
	t.Parallel()

	invalidOrders := []int{0, 1, 5, 10, -1}

	// Act & Assert
	for _, order := range invalidOrders {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			t.Parallel()

			_, err := NewGaussLaguerre(order)

			assert.Error(t, err, "Should return error for invalid order")
			assert.Equal(t, ErrInvalidOrder, err)
		})
	}
}

func TestGaussLaguerreValidOrders(t *testing.T) {
	// Arrange
	t.Parallel()

	validOrders := []int{2, 3, 4}

	// Act & Assert
	for _, order := range validOrders {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			t.Parallel()

			strategy, err := NewGaussLaguerre(order)

			assert.NoError(t, err, "Should create strategy without error")
			assert.NotNil(t, strategy, "Strategy should not be nil")
			assert.Equal(t, order, strategy.Order(), "Order should match")
			assert.Equal(t, "Gauss-Laguerre Quadrature", strategy.Describe())
		})
	}
}

func TestGaussLaguerreOrder2HighTolerance(t *testing.T) {
	// Arrange
	t.Parallel()

	// Create strategy for order 2 only
	strategy, err := NewGaussLaguerre(2)
	assert.NoError(t, err, "Should create Gauss-Laguerre strategy without error")

	testCases := []gaussLaguerreTestCase{
		// Polynomials multiplied by weight function - Gauss-Laguerre integrates f(x)*e^(-x) from 0 to +∞
		{
			name:         "1 (constant)",
			expectedArea: 1.0, // ∫₀^∞ e^(-x) dx = 1
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return 1.0
			},
		},
		{
			name:         "x (linear)",
			expectedArea: 1.0, // ∫₀^∞ x*e^(-x) dx = 1
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:         "x² (quadratic)",
			expectedArea: 2.0, // ∫₀^∞ x²*e^(-x) dx = 2
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:         "x³ (cubic)",
			expectedArea: 6.0, // ∫₀^∞ x³*e^(-x) dx = 6
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:         "x⁴ (quartic)",
			expectedArea: 24.0, // ∫₀^∞ x⁴*e^(-x) dx = 24
			tolerance:    1e1,
			expr: func(x float64) float64 {
				return x * x * x * x
			},
		},
		{
			name:         "x⁵ (quintic)",
			expectedArea: 120.0, // ∫₀^∞ x⁵*e^(-x) dx = 120
			tolerance:    6e1,
			expr: func(x float64) float64 {
				return x * x * x * x * x
			},
		},
		// Test with exponential functions
		{
			name:         "e^(-x) (exponential)",
			expectedArea: 0.5, // ∫₀^∞ e^(-x)*e^(-x) dx = 1/2
			tolerance:    2e-1,
			expr: func(x float64) float64 {
				return math.Exp(-x)
			},
		},
		// Test with more complex functions - very high tolerance for order 2
		{
			name:         "sin(x)",
			expectedArea: 0.5, // ∫₀^∞ sin(x)*e^(-x) dx = 1/2
			tolerance:    2e-1,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			name:         "cos(x)",
			expectedArea: 0.5, // ∫₀^∞ cos(x)*e^(-x) dx = 1/2
			tolerance:    2e-1,
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Order 2 - %s", testCase.name), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			result, err := strategy.Integrate(ctx, testCase.expr, 0.0, math.Inf(1))

			assert.NoError(t, err, "Should integrate without error")
			assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
				"Expected area should match calculated area within tolerance")
		})
	}
}

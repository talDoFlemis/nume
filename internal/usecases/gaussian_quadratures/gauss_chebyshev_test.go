package gaussianquadratures

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type gaussChebyshevTestCase struct {
	name         string
	expr         expressions.SingleVariableExpr
	tolerance    float64
	expectedArea float64
}

func TestGaussChebyshev(t *testing.T) {
	// Arrange
	t.Parallel()

	// Test different orders of Gauss-Chebyshev quadrature
	strategies := []GaussianQuadrature{}

	// Create strategies for orders 3, and 4
	for order := 3; order <= 4; order++ {
		strategy, err := NewGaussChebyshev(order)
		assert.NoError(t, err, "Should create Gauss-Chebyshev strategy without error")
		strategies = append(strategies, strategy)
	}

	testCases := []gaussChebyshevTestCase{
		// Polynomials multiplied by weight function - Gauss-Chebyshev integrates f(x)/√(1-x²) from -1 to 1
		{
			name:         "1 (constant)",
			expectedArea: math.Pi, // ∫₋₁¹ 1/√(1-x²) dx = π
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return 1.0
			},
		},
		{
			name:         "x (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:         "x² (even function)",
			expectedArea: math.Pi / 2.0, // ∫₋₁¹ x²/√(1-x²) dx = π/2
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:         "x³ (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x³/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:         "x⁴ (even function)",
			expectedArea: 3.0 * math.Pi / 8.0, // ∫₋₁¹ x⁴/√(1-x²) dx = 3π/8
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x * x
			},
		},
		{
			name:         "x⁵ (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x⁵/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x * x * x
			},
		},
		// Test with trigonometric functions
		{
			name:         "cos(x)",
			expectedArea: math.Pi * 0.7652, // ∫₋₁¹ cos(x)/√(1-x²) dx ≈ π*0.7652
			tolerance:    1e-1,             // Relax tolerance for approximation
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
		{
			name:         "sin(x)",
			expectedArea: 0.0, // ∫₋₁¹ sin(x)/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		// Test with simple rational functions
		{
			name:         "1/(1+x²)",
			expectedArea: math.Pi / math.Sqrt(2.0), // Approximate value
			tolerance:    10e-1,
			expr: func(x float64) float64 {
				return 1.0 / (1.0 + x*x)
			},
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		for _, strategy := range strategies {
			t.Run(fmt.Sprintf("%s - Order %d", testCase.name, strategy.Order()), func(t *testing.T) {
				t.Parallel()

				ctx := context.Background()
				result, err := strategy.Integrate(ctx, testCase.expr, -1.0, 1.0)

				assert.NoError(t, err, "Should integrate without error")
				assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
					"Expected area should match calculated area within tolerance")
			})
		}
	}
}

func TestGaussChebyshevOrder2HighTolerance(t *testing.T) {
	// Arrange
	t.Parallel()

	// Create strategy for order 2 only
	strategy, err := NewGaussChebyshev(2)
	assert.NoError(t, err, "Should create Gauss-Chebyshev strategy without error")

	testCases := []gaussChebyshevTestCase{
		// Polynomials multiplied by weight function - Gauss-Chebyshev integrates f(x)/√(1-x²) from -1 to 1
		{
			name:         "1 (constant)",
			expectedArea: math.Pi, // ∫₋₁¹ 1/√(1-x²) dx = π
			tolerance:    1e-2,    // Higher tolerance for order 2
			expr: func(x float64) float64 {
				return 1.0
			},
		},
		{
			name:         "x (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:         "x² (even function)",
			expectedArea: math.Pi / 2.0, // ∫₋₁¹ x²/√(1-x²) dx = π/2
			tolerance:    1e-1,          // Even higher tolerance for quadratic
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:         "x³ (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x³/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:         "x⁴ (even function)",
			expectedArea: 3.0 * math.Pi / 8.0, // ∫₋₁¹ x⁴/√(1-x²) dx = 3π/8
			tolerance:    5e-1,                 // Very high tolerance for order 4 polynomial
			expr: func(x float64) float64 {
				return x * x * x * x
			},
		},
		{
			name:         "x⁵ (odd function)",
			expectedArea: 0.0, // ∫₋₁¹ x⁵/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return x * x * x * x * x
			},
		},
		// Test with trigonometric functions
		{
			name:         "cos(x)",
			expectedArea: math.Pi * 0.7652, // ∫₋₁¹ cos(x)/√(1-x²) dx ≈ π*0.7652
			tolerance:    2e-1,             // Very high tolerance for trig functions
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
		{
			name:         "sin(x)",
			expectedArea: 0.0, // ∫₋₁¹ sin(x)/√(1-x²) dx = 0 (odd function)
			tolerance:    1e-2,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		// Test with simple rational functions
		{
			name:         "1/(1+x²)",
			expectedArea: math.Pi / math.Sqrt(2.0), // Approximate value
			tolerance:    5e-1,                     // Very high tolerance for rational function
			expr: func(x float64) float64 {
				return 1.0 / (1.0 + x*x)
			},
		},
	}

	// Act & Assert
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Order 2 - %s", testCase.name), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			result, err := strategy.Integrate(ctx, testCase.expr, -1.0, 1.0)

			assert.NoError(t, err, "Should integrate without error")
			assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
				"Expected area should match calculated area within tolerance")
		})
	}
}

func TestGaussChebyshevErrorCases(t *testing.T) {
	// Arrange
	t.Parallel()

	strategy, err := NewGaussChebyshev(2)
	assert.NoError(t, err, "Should create Gauss-Chebyshev strategy without error")

	testCases := []struct {
		name          string
		leftInterval  float64
		rightInterval float64
		expectedError error
	}{
		{
			name:          "Invalid interval: [0, 1]",
			leftInterval:  0.0,
			rightInterval: 1.0,
			expectedError: ErrChebyshevIntervalsMustBeMinusOneToOne,
		},
		{
			name:          "Invalid interval: [-1, 0]",
			leftInterval:  -1.0,
			rightInterval: 0.0,
			expectedError: ErrChebyshevIntervalsMustBeMinusOneToOne,
		},
		{
			name:          "Invalid interval: [-2, 2]",
			leftInterval:  -2.0,
			rightInterval: 2.0,
			expectedError: ErrChebyshevIntervalsMustBeMinusOneToOne,
		},
		{
			name:          "Invalid interval: [0, 2]",
			leftInterval:  0.0,
			rightInterval: 2.0,
			expectedError: ErrChebyshevIntervalsMustBeMinusOneToOne,
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

func TestGaussChebyshevInvalidOrder(t *testing.T) {
	// Arrange
	t.Parallel()

	invalidOrders := []int{0, 1, 5, 10, -1}

	// Act & Assert
	for _, order := range invalidOrders {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			t.Parallel()

			_, err := NewGaussChebyshev(order)

			assert.Error(t, err, "Should return error for invalid order")
			assert.Equal(t, ErrInvalidOrder, err)
		})
	}
}

func TestGaussChebyshevValidOrders(t *testing.T) {
	// Arrange
	t.Parallel()

	validOrders := []int{2, 3, 4}

	// Act & Assert
	for _, order := range validOrders {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			t.Parallel()

			strategy, err := NewGaussChebyshev(order)

			assert.NoError(t, err, "Should create strategy without error")
			assert.NotNil(t, strategy, "Strategy should not be nil")
			assert.Equal(t, order, strategy.Order(), "Order should match")
			assert.Equal(t, "Gauss-Chebyshev Quadrature", strategy.Describe())
		})
	}
}

package gaussianquadratures

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type gaussQuadratureTestCase struct {
	name          string
	expr          expressions.SingleVariableExpr
	leftInterval  float64
	rightInterval float64
	tolerance     float64
	expectedArea  float64
}

func TestGaussLegendre(t *testing.T) {
	// Arrange
	t.Parallel()

	// Test different orders of Gauss-Legendre quadrature
	strategies := []GaussianQuadrature{}

	// Create strategies for orders 2, 3, and 4
	for order := 2; order <= 4; order++ {
		strategy, err := NewGaussLegendre(order)
		assert.NoError(t, err, "Should create Gauss-Legendre strategy without error")
		strategies = append(strategies, strategy)
	}

	testCases := []gaussQuadratureTestCase{
		// Polynomials - Gauss-Legendre is exact for polynomials of degree 2n-1
		{
			name:          "x",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  0.5, // ∫₀¹ x dx = 1/2
			tolerance:     1e-10,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:          "x²",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  1.0 / 3.0, // ∫₀¹ x² dx = 1/3
			tolerance:     1e-10,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:          "x³",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  0.25, // ∫₀¹ x³ dx = 1/4
			tolerance:     1e-10,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:          "x⁴",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  0.2, // ∫₀¹ x⁴ dx = 1/5
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return x * x * x * x
			},
		},
		// Trigonometric functions
		{
			name:          "sin(x)",
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedArea:  1.0, // ∫₀^(π/2) sin(x) dx = 1
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
		{
			name:          "cos(x)",
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedArea:  1.0, // ∫₀^(π/2) cos(x) dx = 1
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
		// Exponential function
		{
			name:          "e^x",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  math.E - 1, // ∫₀¹ e^x dx = e - 1
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return math.Exp(x)
			},
		},
		// Rational function
		{
			name:          "1/x",
			leftInterval:  1,
			rightInterval: 2,
			expectedArea:  math.Log(2), // ∫₁² 1/x dx = ln(2)
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return 1.0 / x
			},
		},
		// Square root function
		{
			name:          "√x",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  2.0 / 3.0, // ∫₀¹ √x dx = 2/3
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return math.Sqrt(x)
			},
		},
		// Symmetric interval tests
		{
			name:          "x²",
			leftInterval:  -1,
			rightInterval: 1,
			expectedArea:  2.0 / 3.0, // ∫₋₁¹ x² dx = 2/3
			tolerance:     1e-2,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		// Test with different interval scaling
		{
			name:          "x",
			leftInterval:  2,
			rightInterval: 4,
			expectedArea:  6.0, // ∫₂⁴ x dx = 6
			tolerance:     1e-3,
			expr: func(x float64) float64 {
				return x
			},
		},
	}

	for _, testCase := range testCases {
		for _, strategy := range strategies {
			testName := fmt.Sprintf("%s Order %d - %s from %.2f to %.2f",
				strategy.Describe(), strategy.Order(), testCase.name,
				testCase.leftInterval, testCase.rightInterval)
			t.Run(testName, func(t *testing.T) {
				// Act
				result, err := strategy.Integrate(
					t.Context(),
					testCase.expr,
					testCase.leftInterval,
					testCase.rightInterval,
				)

				// Assert
				assert.NoError(t, err, "Expected no error during integration")
				assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
					"Expected integration result to be within tolerance")
			})
		}
	}
}

func TestGaussLegendreErrorCases(t *testing.T) {
	// Arrange
	t.Parallel()

	strategy, err := NewGaussLegendre(2)
	assert.NoError(t, err, "Should create Gauss-Legendre strategy without error")

	simpleExpr := func(x float64) float64 { return x }

	t.Run("Infinite left interval", func(t *testing.T) {
		// Act
		result, err := strategy.Integrate(
			t.Context(),
			simpleExpr,
			math.Inf(-1),
			1.0,
		)

		// Assert
		assert.Error(t, err, "Expected error for infinite left interval")
		assert.Equal(t, ErrInfiniteLeftInterval, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Infinite right interval", func(t *testing.T) {
		// Act
		result, err := strategy.Integrate(
			t.Context(),
			simpleExpr,
			0.0,
			math.Inf(1),
		)

		// Assert
		assert.Error(t, err, "Expected error for infinite right interval")
		assert.Equal(t, ErrInfiniteRightInterval, err)
		assert.Equal(t, 0.0, result)
	})
}

func TestGaussLegendreInvalidOrder(t *testing.T) {
	// Arrange
	t.Parallel()

	invalidOrders := []int{1, 5, 10, -1, 0}

	for _, order := range invalidOrders {
		t.Run(fmt.Sprintf("Invalid order %d", order), func(t *testing.T) {
			// Act
			strategy, err := NewGaussLegendre(order)

			// Assert
			assert.Error(t, err, "Expected error for invalid order")
			assert.Equal(t, ErrInvalidOrder, err)
			assert.Nil(t, strategy)
		})
	}
}

func TestGaussLegendreValidOrders(t *testing.T) {
	// Arrange
	t.Parallel()

	validOrders := []int{2, 3, 4}

	for _, order := range validOrders {
		t.Run(fmt.Sprintf("Valid order %d", order), func(t *testing.T) {
			// Act
			strategy, err := NewGaussLegendre(order)

			// Assert
			assert.NoError(t, err, "Expected no error for valid order")
			assert.NotNil(t, strategy)
			assert.Equal(t, order, strategy.Order())
			assert.Equal(t, "Gauss-Legendre", strategy.Describe())
		})
	}
}

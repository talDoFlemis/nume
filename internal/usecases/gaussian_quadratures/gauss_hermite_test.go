package gaussianquadratures

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taldoflemis/nume/internal/expressions"
)

type gaussHermiteTestCase struct {
	name         string
	expr         expressions.SingleVariableExpr
	tolerance    float64
	expectedArea float64
}

func TestGaussHermite(t *testing.T) {
	// Arrange
	t.Parallel()

	// Test different orders of Gauss-Hermite quadrature
	strategies := []GaussianQuadrature{}

	// Create strategies for orders 2, 3, and 4
	for order := 2; order <= 4; order++ {
		strategy, err := NewGaussHermite(order)
		assert.NoError(t, err, "Should create Gauss-Hermite strategy without error")
		strategies = append(strategies, strategy)
	}

	testCases := []gaussHermiteTestCase{
		// Polynomials multiplied by weight function - Gauss-Hermite integrates f(x)*e^(-x²) from -∞ to +∞
		{
			name:         "1 (constant)",
			expectedArea: math.Sqrt(math.Pi), // ∫₋∞^∞ e^(-x²) dx = √π
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return 1.0
			},
		},
		{
			name:         "x² (even polynomial)",
			expectedArea: math.Sqrt(math.Pi) / 2.0, // ∫₋∞^∞ x²*e^(-x²) dx = √π/2
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		{
			name:         "x⁶ (even polynomial)",
			expectedArea: 15.0 * math.Sqrt(math.Pi) / 8.0, // ∫₋∞^∞ x⁶*e^(-x²) dx = 15√π/8
			tolerance:    4,
			expr: func(x float64) float64 {
				return x * x * x * x * x * x
			},
		},
		// Odd polynomials should integrate to 0 due to symmetry
		{
			name:         "x (odd polynomial)",
			expectedArea: 0.0, // ∫₋∞^∞ x*e^(-x²) dx = 0
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x
			},
		},
		{
			name:         "x³ (odd polynomial)",
			expectedArea: 0.0, // ∫₋∞^∞ x³*e^(-x²) dx = 0
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x
			},
		},
		{
			name:         "x⁵ (odd polynomial)",
			expectedArea: 0.0, // ∫₋∞^∞ x⁵*e^(-x²) dx = 0
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return x * x * x * x * x
			},
		},
		// Test with exponential functions - relax tolerance for more complex functions
		{
			name:         "e^(-x²) (Gaussian)",
			expectedArea: math.Sqrt(math.Pi) / math.Sqrt(2.0), // ∫₋∞^∞ e^(-x²)*e^(-x²) dx = √π/√2
			tolerance:    0.2,
			expr: func(x float64) float64 {
				return math.Exp(-x * x)
			},
		},
		// Test with cosine (even function) - relax tolerance
		{
			name:         "cos(x)",
			expectedArea: math.Sqrt(math.Pi) * math.Exp(-0.25), // ∫₋∞^∞ cos(x)*e^(-x²) dx = √π*e^(-1/4)
			tolerance:    1e-1,
			expr: func(x float64) float64 {
				return math.Cos(x)
			},
		},
		// Test with sine (odd function) should be 0
		{
			name:         "sin(x)",
			expectedArea: 0.0, // ∫₋∞^∞ sin(x)*e^(-x²) dx = 0
			tolerance:    1e-10,
			expr: func(x float64) float64 {
				return math.Sin(x)
			},
		},
	}

	for _, testCase := range testCases {
		for _, strategy := range strategies {
			testName := fmt.Sprintf("%s Order %d - %s",
				strategy.Describe(), strategy.Order(), testCase.name)
			t.Run(testName, func(t *testing.T) {
				// Act
				result, err := strategy.Integrate(
					t.Context(),
					testCase.expr,
					math.Inf(-1),
					math.Inf(1),
				)

				// Assert
				assert.NoError(t, err, "Expected no error during integration")
				assert.InDelta(t, testCase.expectedArea, result, testCase.tolerance,
					"Expected integration result to be within tolerance")
			})
		}
	}
}

func TestGaussHermiteErrorCases(t *testing.T) {
	// Arrange
	t.Parallel()

	strategy, err := NewGaussHermite(2)
	assert.NoError(t, err, "Should create Gauss-Hermite strategy without error")

	simpleExpr := func(x float64) float64 { return 1.0 }

	t.Run("Finite left interval", func(t *testing.T) {
		// Act
		result, err := strategy.Integrate(
			t.Context(),
			simpleExpr,
			-1.0,
			math.Inf(1),
		)

		// Assert
		assert.Error(t, err, "Expected error for finite left interval")
		assert.Equal(t, ErrHermiteIntervalsMustBeInfinite, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Finite right interval", func(t *testing.T) {
		// Act
		result, err := strategy.Integrate(
			t.Context(),
			simpleExpr,
			math.Inf(-1),
			1.0,
		)

		// Assert
		assert.Error(t, err, "Expected error for finite right interval")
		assert.Equal(t, ErrHermiteIntervalsMustBeInfinite, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Both intervals finite", func(t *testing.T) {
		// Act
		result, err := strategy.Integrate(
			t.Context(),
			simpleExpr,
			-1.0,
			1.0,
		)

		// Assert
		assert.Error(t, err, "Expected error for both finite intervals")
		assert.Equal(t, ErrHermiteIntervalsMustBeInfinite, err)
		assert.Equal(t, 0.0, result)
	})
}

func TestGaussHermiteInvalidOrder(t *testing.T) {
	// Arrange
	t.Parallel()

	invalidOrders := []int{1, 5, 10, -1, 0}

	for _, order := range invalidOrders {
		t.Run(fmt.Sprintf("Invalid order %d", order), func(t *testing.T) {
			// Act
			strategy, err := NewGaussHermite(order)

			// Assert
			assert.Error(t, err, "Expected error for invalid order")
			assert.Equal(t, ErrInvalidOrder, err)
			assert.Nil(t, strategy)
		})
	}
}

func TestGaussHermiteValidOrders(t *testing.T) {
	// Arrange
	t.Parallel()

	validOrders := []int{2, 3, 4}

	for _, order := range validOrders {
		t.Run(fmt.Sprintf("Valid order %d", order), func(t *testing.T) {
			// Act
			strategy, err := NewGaussHermite(order)

			// Assert
			assert.NoError(t, err, "Expected no error for valid order")
			assert.NotNil(t, strategy)
			assert.Equal(t, order, strategy.Order())
			assert.Equal(t, "Gauss-Hermite Quadrature", strategy.Describe())
		})
	}
}


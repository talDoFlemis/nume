package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/taldoflemis/nume/internal/expressions"
)

type GaussChebyshev struct {
	order   int
	nodes   map[int][]float64
	weights map[int][]float64
}

const (
	chebyshevMaximumOrder = 4
	chebyshevMinimumOrder = 2
)

var ErrChebyshevIntervalsMustBeMinusOneToOne = errors.New("chebyshev quadrature requires interval [-1, 1]")

var _ GaussianQuadrature = (*GaussChebyshev)(nil)

func NewGaussChebyshev(order int) (*GaussChebyshev, error) {
	if order < chebyshevMinimumOrder || order > chebyshevMaximumOrder {
		slog.Error("Invalid order for Gauss-Chebyshev quadrature", slog.Int("order", order))
		return nil, ErrInvalidOrder
	}

	nodes := make(map[int][]float64)
	weights := make(map[int][]float64)

	// Gauss-Chebyshev quadrature nodes and weights using mathematical constants
	// These are based on Chebyshev polynomials of the first kind
	// Order 2
	nodes[2] = []float64{
		-math.Cos(math.Pi / 4.0),
		math.Cos(math.Pi / 4.0),
	}
	weights[2] = []float64{
		math.Pi / 2.0,
		math.Pi / 2.0,
	}

	// Order 3
	nodes[3] = []float64{
		-math.Cos(math.Pi / 6.0),
		0.0,
		math.Cos(math.Pi / 6.0),
	}
	weights[3] = []float64{
		math.Pi / 3.0,
		math.Pi / 3.0,
		math.Pi / 3.0,
	}

	// Order 4
	nodes[4] = []float64{
		-math.Cos(math.Pi / 8.0),
		-math.Cos(3.0 * math.Pi / 8.0),
		math.Cos(3.0 * math.Pi / 8.0),
		math.Cos(math.Pi / 8.0),
	}
	weights[4] = []float64{
		math.Pi / 4.0,
		math.Pi / 4.0,
		math.Pi / 4.0,
		math.Pi / 4.0,
	}

	return &GaussChebyshev{
		order:   order,
		nodes:   nodes,
		weights: weights,
	}, nil
}

// Describe implements GaussianQuadrature.
func (g *GaussChebyshev) Describe() string {
	return "Gauss-Chebyshev Quadrature"
}

// Integrate implements GaussianQuadrature.
func (g *GaussChebyshev) Integrate(
	ctx context.Context,
	expr expressions.SingleVariableExpr,
	leftInterval,
	rightInterval float64,
) (float64, error) {
	return calculatePartition(ctx, g, expr, leftInterval, rightInterval)
}

// Order implements GaussianQuadrature.
func (g *GaussChebyshev) Order() int {
	return g.order
}

// Validate implements GaussianQuadrature.
func (g *GaussChebyshev) Validate(ctx context.Context, leftInterval, rightInterval float64) error {
	if leftInterval != -1.0 || rightInterval != 1.0 {
		slog.ErrorContext(ctx, "Left interval must be -1 and right interval must be 1, "+
			"cannot perform Gauss-Chebyshev quadrature. Use another quadrature method.")
		return ErrChebyshevIntervalsMustBeMinusOneToOne
	}
	return nil
}

// GetNodes implements GaussianQuadrature.
func (g *GaussChebyshev) GetNodes() []float64 {
	return g.nodes[g.order]
}

// GetWeights implements GaussianQuadrature.
func (g *GaussChebyshev) GetWeights() []float64 {
	return g.weights[g.order]
}

// GetOffset implements GaussianQuadrature.
func (g *GaussChebyshev) GetOffset(leftInterval, rightInterval float64) float64 {
	// Gauss-Chebyshev quadrature doesn't need offset transformation
	return 0.0
}

// GetScalingFactor implements GaussianQuadrature.
func (g *GaussChebyshev) GetScalingFactor(leftInterval, rightInterval float64) float64 {
	// Gauss-Chebyshev quadrature doesn't need scaling transformation
	return 1.0
}

// AllowPartitioning implements GaussianQuadrature.
func (g *GaussChebyshev) AllowPartitioning() bool {
	// Gauss-Chebyshev quadrature is for [-1, 1] interval and doesn't support partitioning
	return false
}

package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/taldoflemis/nume/internal/expressions"
)

type GaussHermite struct {
	order   int
	nodes   map[int][]float64
	weights map[int][]float64
}

const (
	hermiteMaximumOrder = 4
	hermiteMinimumOrder = 2
)

var ErrHermiteIntervalsMustBeInfinite = errors.New("hermite quadrature requires infinite intervals")

var _ GaussianQuadrature = (*GaussHermite)(nil)

func NewGaussHermite(order int) (*GaussHermite, error) {
	if order < hermiteMinimumOrder || order > hermiteMaximumOrder {
		slog.Error("Invalid order for Gauss-Hermite quadrature", slog.Int("order", order))
		return nil, ErrInvalidOrder
	}

	nodes := make(map[int][]float64)
	weights := make(map[int][]float64)

	// Gauss-Hermite quadrature nodes and weights using mathematical constants
	// Order 2
	nodes[2] = []float64{
		-math.Sqrt(2.0) / 2.0,
		math.Sqrt(2.0) / 2.0,
	}
	weights[2] = []float64{
		math.Sqrt(math.Pi) / 2.0,
		math.Sqrt(math.Pi) / 2.0,
	}

	// Order 3
	nodes[3] = []float64{
		-math.Sqrt(6.0) / 2.0,
		0.0,
		math.Sqrt(6.0) / 2.0,
	}
	weights[3] = []float64{
		math.Sqrt(math.Pi) / 6.0,
		2.0 * math.Sqrt(math.Pi) / 3.0,
		math.Sqrt(math.Pi) / 6.0,
	}

	// Order 4
	nodes[4] = []float64{
		// Calculated using numpy
		-1.650680123885784, -0.52464762327529, 0.52464762327529, 1.650680123885784,
	}
	weights[4] = []float64{
		0.081312835447245, 0.804914090005513, 0.804914090005513, 0.081312835447245,
	}

	return &GaussHermite{
		order:   order,
		nodes:   nodes,
		weights: weights,
	}, nil
}

// Describe implements GaussianQuadrature.
func (g *GaussHermite) Describe() string {
	return "Gauss-Hermite Quadrature"
}

// Integrate implements GaussianQuadrature.
func (g *GaussHermite) Integrate(
	ctx context.Context,
	expr expressions.SingleVariableExpr,
	leftInterval,
	rightInterval float64,
) (float64, error) {
	return calculatePartition(ctx, g, expr, leftInterval, rightInterval)
}

// Order implements GaussianQuadrature.
func (g *GaussHermite) Order() int {
	return g.order
}

// Validate implements GaussianQuadrature.
func (g *GaussHermite) Validate(ctx context.Context, leftInterval, rightInterval float64) error {
	if leftInterval != math.Inf(-1) || rightInterval != math.Inf(1) {
		slog.ErrorContext(ctx, "Left and right intervals must be infinite, cannot perform Gauss-Hermite quadrature. Use another quadrature method.")
		return ErrHermiteIntervalsMustBeInfinite
	}
	return nil
}

// GetNodes implements GaussianQuadrature.
func (g *GaussHermite) GetNodes() []float64 {
	return g.nodes[g.order]
}

// GetWeights implements GaussianQuadrature.
func (g *GaussHermite) GetWeights() []float64 {
	return g.weights[g.order]
}

// GetOffset implements GaussianQuadrature.
func (g *GaussHermite) GetOffset(leftInterval, rightInterval float64) float64 {
	// Gauss-Hermite quadrature doesn't need offset transformation
	return 0.0
}

// GetScalingFactor implements GaussianQuadrature.
func (g *GaussHermite) GetScalingFactor(leftInterval, rightInterval float64) float64 {
	// Gauss-Hermite quadrature doesn't need scaling transformation
	return 1.0
}

// AllowPartitioning implements GaussianQuadrature.
func (g *GaussHermite) AllowPartitioning() bool {
	// Gauss-Hermite quadrature is for (-∞, +∞) interval and doesn't support partitioning
	return false
}

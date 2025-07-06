package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/taldoflemis/nume/internal/expressions"
)

type GaussLegendre struct {
	order   int
	nodes   map[int][]float64
	weights map[int][]float64
}

const (
	maximumOrder = 4
	minimumOrder = 2
)

var ErrInvalidOrder = errors.New("invalid order for quadrature, must be between 2 and 4")

var _ GaussianQuadrature = (*GaussLegendre)(nil)

func NewGaussLegendre(order int) (*GaussLegendre, error) {
	if order < minimumOrder || order > maximumOrder {
		slog.Error("Invalid order for Gauss-Legendre quadrature", slog.Int("order", order))
		return nil, ErrInvalidOrder
	}
	nodes := make(map[int][]float64)
	weights := make(map[int][]float64)

	// 2 Points
	nodes[2] = []float64{-1.0 / math.Sqrt(3.0), 1.0 / math.Sqrt(3.0)}
	weights[2] = []float64{1.0, 1.0}

	// 3 Points
	nodes[3] = []float64{-math.Sqrt(3.0 / 5.0), 0.0, math.Sqrt(3.0 / 5.0)}
	weights[3] = []float64{5.0 / 9.0, 8.0 / 9.0, 5.0 / 9.0}

	// 4 Points
	nodes[4] = []float64{
		-math.Sqrt((3.0 + 2.0*math.Sqrt(6.0/5.0)) / 7.0),
		-math.Sqrt((3.0 - 2.0*math.Sqrt(6.0/5.0)) / 7.0),
		math.Sqrt((3.0 - 2.0*math.Sqrt(6.0/5.0)) / 7.0),
		math.Sqrt((3.0 + 2.0*math.Sqrt(6.0/5.0)) / 7.0),
	}
	weights[4] = []float64{
		((18.0 - math.Sqrt(30.0)) / 36.0),
		((18.0 + math.Sqrt(30.0)) / 36.0),
		((18.0 + math.Sqrt(30.0)) / 36.0),
		((18.0 - math.Sqrt(30.0)) / 36.0),
	}

	return &GaussLegendre{
		order:   order,
		nodes:   nodes,
		weights: weights,
	}, nil
}

var (
	ErrInfiniteLeftInterval  = errors.New("left interval is infinite")
	ErrInfiniteRightInterval = errors.New("right interval is infinite")
)

// Integrate implements GaussianQuadrature.
func (g *GaussLegendre) Integrate(
	ctx context.Context,
	expr expressions.SingleVariableExpr,
	leftInterval,
	rightInterval float64,
) (float64, error) {
	return calculatePartition(ctx, g, expr, leftInterval, rightInterval)
}

// Describe implements GaussianQuadrature.
func (g *GaussLegendre) Describe() string {
	return "Gauss-Legendre"
}

// Order implements GaussianQuadrature.
func (g *GaussLegendre) Order() int {
	return g.order
}

// Validate implements GaussianQuadrature.
func (g *GaussLegendre) Validate(ctx context.Context, leftInterval, rightInterval float64) error {
	if leftInterval == math.Inf(-1) {
		slog.ErrorContext(ctx, "Left interval is infinite, cannot perform Gauss-Legendre quadrature. Use another quadrature method.")
		return ErrInfiniteLeftInterval
	}

	if rightInterval == math.Inf(1) {
		slog.ErrorContext(ctx, "Right interval is infinite, cannot perform Gauss-Legendre quadrature. Use another quadrature method.")
		return ErrInfiniteRightInterval
	}

	if leftInterval == rightInterval {
		return ErrZeroWidthInterval
	}

	return nil
}

// GetNodes implements GaussianQuadrature.
func (g *GaussLegendre) GetNodes() []float64 {
	return g.nodes[g.order]
}

// GetWeights implements GaussianQuadrature.
func (g *GaussLegendre) GetWeights() []float64 {
	return g.weights[g.order]
}

// GetOffset implements GaussianQuadrature.
func (g *GaussLegendre) GetOffset(leftInterval, rightInterval float64) float64 {
	// Gauss-Legendre quadrature uses dynamic offset calculation
	return (rightInterval + leftInterval) / 2.0
}

// GetScalingFactor implements GaussianQuadrature.
func (g *GaussLegendre) GetScalingFactor(leftInterval, rightInterval float64) float64 {
	// Gauss-Legendre quadrature uses dynamic scaling factor calculation
	return (rightInterval - leftInterval) / 2.0
}

// AllowPartitioning implements GaussianQuadrature.
func (g *GaussLegendre) AllowPartitioning() bool {
	// Gauss-Legendre quadrature supports partitioning for arbitrary intervals
	return true
}

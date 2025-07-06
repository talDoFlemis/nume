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
	nodes := g.nodes[g.order]
	weights := g.weights[g.order]

	slog.DebugContext(ctx, "Calculating quadrature",
		slog.String("method", g.Describe()),
		slog.Any("expression", expr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Int("order", g.order),
		slog.Any("nodes", nodes),
		slog.Any("weights", weights),
	)

	if leftInterval != math.Inf(-1) || rightInterval != math.Inf(1) {
		slog.ErrorContext(ctx, "Left and right intervals must be infinite, cannot perform Gauss-Hermite quadrature. Use another quadrature method.")
		return 0, ErrHermiteIntervalsMustBeInfinite
	}

	accumulatedArea := 0.0

	for i := range nodes {
		slog.DebugContext(ctx, "Processing node",
			slog.Float64("node", nodes[i]),
			slog.Float64("weight", weights[i]),
			slog.Float64("accumulatedArea", accumulatedArea),
		)
		accumulatedArea += weights[i] * expr(nodes[i])
	}

	slog.InfoContext(ctx, "Final accumulated area",
		slog.Float64("accumulatedArea", accumulatedArea),
	)

	return accumulatedArea, nil
}

// Order implements GaussianQuadrature.
func (g *GaussHermite) Order() int {
	return g.order
}

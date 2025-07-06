package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/taldoflemis/nume/internal/expressions"
)

type GaussLaguerre struct {
	order   int
	nodes   map[int][]float64
	weights map[int][]float64
}

const (
	laguerreMaximumOrder = 4
	laguerreMinimumOrder = 2
)

var ErrLaguerreIntervalsMustBePositiveInfinite = errors.New(
	"laguerre quadrature requires interval [0, +∞)",
)

var _ GaussianQuadrature = (*GaussLaguerre)(nil)

func NewGaussLaguerre(order int) (*GaussLaguerre, error) {
	if order < laguerreMinimumOrder || order > laguerreMaximumOrder {
		slog.Error("Invalid order for Gauss-Laguerre quadrature", slog.Int("order", order))
		return nil, ErrInvalidOrder
	}

	nodes := make(map[int][]float64)
	weights := make(map[int][]float64)

	// Gauss-Laguerre quadrature nodes and weights using mathematical constants
	// These are the roots of Laguerre polynomials and their corresponding weights
	// Order 2 - roots of L₂(x) = x² - 4x + 2
	nodes[2] = []float64{
		0.585786437626905,
		3.414213562373095,
	}
	weights[2] = []float64{
		0.853553390593274, 0.146446609406726,
	}

	// Order 3 - roots of L₃(x) = -x³ + 9x² - 18x + 6
	nodes[3] = []float64{
		0.415774556783479, 2.294280360279042, 6.289945082937479,
	}
	weights[3] = []float64{
		0.711093009929173, 0.278517733569241, 0.010389256501586,
	}

	// Order 4 - using correct Laguerre polynomial roots
	nodes[4] = []float64{
		// Calculated using numpy
		0.322547689619392,
		1.745761101158346,
		4.536620296921128,
		9.395070912301133,
	}
	weights[4] = []float64{
		6.031541043416337e-01,
		3.574186924377996e-01,
		3.888790851500541e-02,
		5.392947055613296e-04,
	}

	return &GaussLaguerre{
		order:   order,
		nodes:   nodes,
		weights: weights,
	}, nil
}

// Describe implements GaussianQuadrature.
func (g *GaussLaguerre) Describe() string {
	return "Gauss-Laguerre Quadrature"
}

// Integrate implements GaussianQuadrature.
func (g *GaussLaguerre) Integrate(
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

	if leftInterval != 0.0 || rightInterval != math.Inf(1) {
		slog.ErrorContext(ctx, "Left interval must be 0 and right interval must be +∞, "+
			"cannot perform Gauss-Laguerre quadrature. Use another quadrature method.")
		return 0, ErrLaguerreIntervalsMustBePositiveInfinite
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
func (g *GaussLaguerre) Order() int {
	return g.order
}

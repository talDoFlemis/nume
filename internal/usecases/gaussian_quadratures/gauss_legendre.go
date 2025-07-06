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

var ErrInvalidOrder = errors.New("invalid order for Gauss-Legendre quadrature, must be between 2 and 4")

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
	nodes := g.nodes[g.order]
	weights := g.weights[g.order]

	slog.DebugContext(ctx, "Calculating Gauss-Legendre quadrature",
		slog.Any("expression", expr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Int("order", g.order),
		slog.Any("nodes", nodes),
		slog.Any("weights", weights),
	)

	if leftInterval == math.Inf(-1) {
		slog.ErrorContext(ctx, "Left interval is infinite, cannot perform Gauss-Legendre quadrature. Use another quadrature method.")
		return 0, ErrInfiniteLeftInterval
	}

	if rightInterval == math.Inf(1) {
		slog.ErrorContext(ctx, "Right interval is infinite, cannot perform Gauss-Legendre quadrature. Use another quadrature method.")
		return 0, ErrInfiniteRightInterval
	}

	if leftInterval == rightInterval {
		panic("Left and right intervals are equal, cannot perform Gauss-Legendre quadrature")
	}

	scaleFactor := (rightInterval - leftInterval) / 2.0
	offset := (rightInterval + leftInterval) / 2.0

	slog.DebugContext(ctx, "Scale factor and offset calculated",
		slog.Float64("scaleFactor", scaleFactor),
		slog.Float64("offset", offset),
	)

	accumulatedArea := 0.0

	for i := range nodes {
		slog.DebugContext(ctx, "Processing node",
			slog.Float64("node", nodes[i]),
			slog.Float64("weight", weights[i]),
			slog.Float64("accumulatedArea", accumulatedArea),
		)

		transformedX := scaleFactor*nodes[i] + offset
		accumulatedArea += weights[i] * expr(transformedX)
	}

	accumulatedArea = accumulatedArea * scaleFactor

	slog.InfoContext(ctx, "Final accumulated area",
		slog.Float64("accumulatedArea", accumulatedArea),
	)

	return accumulatedArea, nil
}

// Describe implements GaussianQuadrature.
func (g *GaussLegendre) Describe() string {
	return "Gauss-Legendre"
}

// Order implements GaussianQuadrature.
func (g *GaussLegendre) Order() int {
	return g.order
}

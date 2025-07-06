package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

var ErrZeroWidthInterval = errors.New("interval width is zero")

type GaussianQuadrature interface {
	Integrate(ctx context.Context, expr expressions.SingleVariableExpr, leftInterval, rightInterval float64) (float64, error)
	Describe() string
	Order() int
}

type GaussCalculatorUseCase struct {
	strategy GaussianQuadrature
}

func NewGaussCalculatorUseCase(strategy GaussianQuadrature) *GaussCalculatorUseCase {
	return &GaussCalculatorUseCase{
		strategy: strategy,
	}
}

func (u *GaussCalculatorUseCase) Calculate(
	ctx context.Context,
	expr expressions.SingleVariableExpr,
	leftInterval,
	rightInterval float64,
	maxNumberOfPartitions uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Calculating Gauss quadrature",
		slog.Any("expression", expr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Uint64("maxNumberOfPartitions", maxNumberOfPartitions),
		slog.String("strategy", u.strategy.Describe()),
		slog.Int("order", u.strategy.Order()),
	)

	if leftInterval == rightInterval {
		slog.ErrorContext(ctx, "Left and right intervals are equal")
		return 0, ErrZeroWidthInterval
	}
	return 0.0, nil
}

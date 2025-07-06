package gaussianquadratures

import (
	"context"
	"errors"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

var ErrZeroWidthInterval = errors.New("interval width is zero")

type GaussianQuadrature interface {
	Integrate(
		ctx context.Context,
		expr expressions.SingleVariableExpr,
		leftInterval, rightInterval float64,
	) (float64, error)
	Validate(ctx context.Context, leftInterval, rightInterval float64) error
	GetNodes() []float64
	GetWeights() []float64
	GetOffset(leftInterval, rightInterval float64) float64
	GetScalingFactor(leftInterval, rightInterval float64) float64
	AllowPartitioning() bool
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
	numberOfPartitions uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Calculating Gauss quadrature",
		slog.Any("expression", expr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Uint64("numberOfPartitions", numberOfPartitions),
		slog.String("strategy", u.strategy.Describe()),
		slog.Int("order", u.strategy.Order()),
	)

	if leftInterval == rightInterval {
		slog.ErrorContext(ctx, "Left and right intervals are equal")
		return 0, ErrZeroWidthInterval
	}

	if !u.strategy.AllowPartitioning() {
		slog.DebugContext(ctx, "Strategy does not allow partitioning, calculating directly")
		return u.strategy.Integrate(ctx, expr, leftInterval, rightInterval)
	}

	if numberOfPartitions == 0 {
		slog.ErrorContext(ctx, "Max number of partitions is zero")
		return 0.0, errors.New("max number of partitions must be greater than zero")
	}

	delta := (rightInterval - leftInterval) / float64(numberOfPartitions)

	accumulatedArea := 0.0

	for i := leftInterval; i <= rightInterval; i += delta {
		slog.DebugContext(ctx, "Calculating area for partition",
			slog.Float64("left", i),
			slog.Float64("right", i+delta),
			slog.Uint64("partition", uint64(i/delta)),
		)
		partitionArea, err := u.strategy.Integrate(ctx, expr, i, i+delta)
		if err != nil {
			slog.ErrorContext(ctx, "Error integrating partition", slog.Any("error", err))
			return 0.0, errors.New("error integrating partition: " + err.Error())
		}

		slog.DebugContext(ctx, "Calculated area for partition",
			slog.Float64("partitionArea", partitionArea),
		)

		accumulatedArea += partitionArea
	}

	slog.InfoContext(ctx, "Gauss quadrature integration completed",
		slog.Float64("totalArea", accumulatedArea),
	)

	return accumulatedArea, nil
}

func calculatePartition(
	ctx context.Context,
	strategy GaussianQuadrature,
	expr expressions.SingleVariableExpr,
	leftInterval,
	rightInterval float64,
) (float64, error) {
	nodes := strategy.GetNodes()
	weights := strategy.GetWeights()

	slog.DebugContext(ctx, "Calculating quadrature",
		slog.String("method", strategy.Describe()),
		slog.Any("expression", expr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Int("order", strategy.Order()),
		slog.Any("nodes", nodes),
		slog.Any("weights", weights),
	)

	slog.DebugContext(ctx, "Validating intervals")
	err := strategy.Validate(ctx, leftInterval, rightInterval)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid intervals", slog.Any("error", err))
		return 0.0, err
	}

	slog.DebugContext(ctx, "Valid intervals")

	accumulatedArea := 0.0

	scaleFactor := strategy.GetScalingFactor(leftInterval, rightInterval)
	offset := strategy.GetOffset(leftInterval, rightInterval)

	slog.DebugContext(ctx, "Transformation parameters",
		slog.Float64("scaleFactor", scaleFactor),
		slog.Float64("offset", offset),
	)

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

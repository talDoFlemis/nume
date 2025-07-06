package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

type DoubleIntegralUseCase struct {
}

func NewDoubleIntegralUseCase() *DoubleIntegralUseCase {
	return &DoubleIntegralUseCase{}
}

var ErrZeroWidthInterval = errors.New(
	"left and right intervals are equal, cannot perform double integral",
)

func (d *DoubleIntegralUseCase) CalculateArea(
	ctx context.Context,
	expr expressions.DualVariableExpr,
	leftIntervalX, rightIntervalX,
	leftIntervalY, rightIntervalY float64,
	numberOfPartitions uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Calculating double integral area",
		slog.Any("expression", expr),
		slog.Float64("leftIntervalX", leftIntervalX),
		slog.Float64("rightIntervalX", rightIntervalX),
		slog.Float64("leftIntervalY", leftIntervalY),
		slog.Float64("rightIntervalY", rightIntervalY),
		slog.Uint64("numberOfPartitions", numberOfPartitions),
	)

	if leftIntervalX == rightIntervalX || leftIntervalY == rightIntervalY {
		return 0, ErrZeroWidthInterval
	}



	if numberOfPartitions == 0 {
		slog.WarnContext(ctx, "Number of partitions is zero, using default value of 1")
		numberOfPartitions = 1
	}

	// Calculate step sizes for both dimensions
	deltaX := (rightIntervalX - leftIntervalX) / float64(numberOfPartitions)
	deltaY := (rightIntervalY - leftIntervalY) / float64(numberOfPartitions)

	accumulatedArea := 0.0

	// Double Riemann sum using midpoint rule
	for i := uint64(0); i < numberOfPartitions; i++ {
		for j := uint64(0); j < numberOfPartitions; j++ {
			// Calculate midpoint coordinates
			midX := leftIntervalX + (float64(i)+0.5)*deltaX
			midY := leftIntervalY + (float64(j)+0.5)*deltaY

			// Evaluate function at midpoint and add to accumulated area
			functionValue := expr(midX, midY)
			accumulatedArea += functionValue * deltaX * deltaY
		}
	}

	return accumulatedArea, nil
}

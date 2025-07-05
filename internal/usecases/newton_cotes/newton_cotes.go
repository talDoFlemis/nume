package newtoncotes

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

type FormulaType string

const (
	OpenFormulaType   FormulaType = "open"
	ClosedFormulaType FormulaType = "closed"
)

type NewtonCotesOrder int

const (
	FirstOrder NewtonCotesOrder = iota + 1
	SecondOrder
	ThirdOrder
)

type NewtonCotesStrategy interface {
	Integrate(
		ctx context.Context,
		simpleExpr expressions.SingleVariableExpr,
		leftInterval float64,
		rightInterval float64,
	) (float64, error) // Integrates the expression using the Newton-Cotes formula
	Description() string     // Returns a description of the strategy (e.g., "Trapezoidal Rule")
	Order() NewtonCotesOrder // Returns the polynomial order of the strategy
	Type() FormulaType       // Returns the type of formula ("closed" or "open")
}

type NewtonCotesUseCase struct {
	strategy NewtonCotesStrategy
}

func NewNewtonCotesUseCase(strategy NewtonCotesStrategy) *NewtonCotesUseCase {
	return &NewtonCotesUseCase{
		strategy: strategy,
	}
}

func (u *NewtonCotesUseCase) Calculate(
	ctx context.Context,
	simpleExpr expressions.SingleVariableExpr,
	leftInterval float64,
	rightInterval float64,
	numberOfPartitions uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Starting Newton-Cotes integration",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
		slog.Uint64("numberOfPartitions", numberOfPartitions),
		slog.String("strategy", u.strategy.Description()),
		slog.Int("order", int(u.strategy.Order())),
		slog.String("type", string(u.strategy.Type())),
	)

	acumulatedArea := 0.0
	delta := (rightInterval - leftInterval) / float64(numberOfPartitions)

	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	for i := leftInterval; i <= rightInterval; i += delta {
		slog.DebugContext(ctx, "Calculating area for partition",
			slog.Float64("left", i),
			slog.Float64("right", i+delta),
			slog.Uint64("partition", uint64(i/delta)),
			slog.Float64("currentArea", acumulatedArea),
		)

		partitionArea, err := u.strategy.Integrate(ctx, simpleExpr, i, i+delta)
		if err != nil {
			slog.ErrorContext(ctx, "Error integrating partition", "err", err)
			return 0, fmt.Errorf("error integrating partition [%f, %f]: %w", i, i+delta, err)
		}

		slog.DebugContext(ctx, "Calculated area for partition",
			slog.Float64("partitionArea", partitionArea),
		)

		acumulatedArea += partitionArea
	}

	slog.InfoContext(ctx, "Newton-Cotes integration completed",
		slog.Float64("totalArea", acumulatedArea),
	)

	return acumulatedArea, nil
}

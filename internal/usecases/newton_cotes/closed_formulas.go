package newtoncotes

import (
	"context"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

type TrapezoidalRule struct{}

var _ NewtonCotesStrategy = (*TrapezoidalRule)(nil)

// Description implements NewtonCotesStrategy.
func (t *TrapezoidalRule) Description() string {
	return "Trapezoidal Rule"
}

// Integrate implements NewtonCotesStrategy.
func (t *TrapezoidalRule) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Trapezoidal Rule",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)

	delta := (rightInterval - leftInterval)

	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (delta / 2.0) * (simpleExpr(leftInterval) + simpleExpr(rightInterval)), nil
}

// Order implements NewtonCotesStrategy.
func (t *TrapezoidalRule) Order() NewtonCotesOrder {
	return FirstOrder
}

// Type implements NewtonCotesStrategy.
func (t *TrapezoidalRule) Type() FormulaType {
	return ClosedFormulaType
}

type SimpsonsOneThirdRule struct{}

var _ NewtonCotesStrategy = (*SimpsonsOneThirdRule)(nil)

// Description implements NewtonCotesStrategy.
func (s *SimpsonsOneThirdRule) Description() string {
	return "Simpson's One-Third Rule"
}

// Integrate implements NewtonCotesStrategy.
func (s *SimpsonsOneThirdRule) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Simpson's One-Third Rule",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)

	delta := (rightInterval - leftInterval) / 2.0

	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (delta / 3.0) * (simpleExpr(leftInterval) + 4*simpleExpr(leftInterval+delta) + simpleExpr(rightInterval)), nil
}

// Order implements NewtonCotesStrategy.
func (s *SimpsonsOneThirdRule) Order() NewtonCotesOrder {
	return SecondOrder
}

// Type implements NewtonCotesStrategy.
func (s *SimpsonsOneThirdRule) Type() FormulaType {
	return ClosedFormulaType
}

type SimpsonsThreeEighthsRule struct{}

var _ NewtonCotesStrategy = (*SimpsonsThreeEighthsRule)(nil)

// Description implements NewtonCotesStrategy.
func (s *SimpsonsThreeEighthsRule) Description() string {
	return "Simpson's Three-Eighths Rule"
}

// Integrate implements NewtonCotesStrategy.
func (s *SimpsonsThreeEighthsRule) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Simpson's Three-Eighths Rule",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)

	delta := (rightInterval - leftInterval) / 3.0

	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (3.0 * delta / 8.0) * (simpleExpr(leftInterval) + 3*simpleExpr(leftInterval+delta) + 3*simpleExpr(leftInterval+2*delta) + simpleExpr(rightInterval)), nil
}

// Order implements NewtonCotesStrategy.
func (s *SimpsonsThreeEighthsRule) Order() NewtonCotesOrder {
	return ThirdOrder
}

// Type implements NewtonCotesStrategy.
func (s *SimpsonsThreeEighthsRule) Type() FormulaType {
	return ClosedFormulaType
}

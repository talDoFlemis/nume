package newtoncotes

import (
	"context"
	"log/slog"

	"github.com/taldoflemis/nume/internal/expressions"
)

var (
	_ NewtonCotesStrategy = (*OpenTrapezoidalRule)(nil)
	_ NewtonCotesStrategy = (*MilneRule)(nil)
	_ NewtonCotesStrategy = (*ThirdDegreeOpenNewtonCotesStrategy)(nil)
)

type OpenTrapezoidalRule struct{}

// Description implements NewtonCotesStrategy.
func (o *OpenTrapezoidalRule) Description() string {
	return "Open Trapezoidal Rule"
}

// Integrate implements NewtonCotesStrategy.
func (o *OpenTrapezoidalRule) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Open Trapezoidal Rule",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)

	delta := (rightInterval - leftInterval) / 3.0

	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (3 * delta / 2.0) * (simpleExpr(leftInterval+delta) + simpleExpr(leftInterval+2*delta)), nil
}

// Order implements NewtonCotesStrategy.
func (o *OpenTrapezoidalRule) Order() NewtonCotesOrder {
	return FirstOrder
}

// Type implements NewtonCotesStrategy.
func (o *OpenTrapezoidalRule) Type() FormulaType {
	return OpenFormulaType
}

type MilneRule struct{}

// Description implements NewtonCotesStrategy.
func (m *MilneRule) Description() string {
	return "Milne's Rule"
}

// Integrate implements NewtonCotesStrategy.
func (m *MilneRule) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Milne's Rule",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)
	delta := (rightInterval - leftInterval) / 4.0
	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (4 * delta / 3.0) * (2*simpleExpr(leftInterval+delta) - simpleExpr(leftInterval+2*delta) + 2*simpleExpr(leftInterval+3*delta)), nil
}

// Order implements NewtonCotesStrategy.
func (m *MilneRule) Order() NewtonCotesOrder {
	return SecondOrder
}

// Type implements NewtonCotesStrategy.
func (m *MilneRule) Type() FormulaType {
	return OpenFormulaType
}

type ThirdDegreeOpenNewtonCotesStrategy struct{}

// Description implements NewtonCotesStrategy.
func (t *ThirdDegreeOpenNewtonCotesStrategy) Description() string {
	return "Third Degree Open Newton-Cotes Formula that I'm calling marcelinho"
}

// Integrate implements NewtonCotesStrategy.
func (t *ThirdDegreeOpenNewtonCotesStrategy) Integrate(ctx context.Context, simpleExpr expressions.SingleVariableExpr, leftInterval float64, rightInterval float64) (float64, error) {
	slog.DebugContext(ctx, "Integrating using Third Degree Open Newton-Cotes Formula",
		slog.Any("simpleExpr", simpleExpr),
		slog.Float64("leftInterval", leftInterval),
		slog.Float64("rightInterval", rightInterval),
	)

	delta := (rightInterval - leftInterval) / 5.0
	slog.DebugContext(ctx, "Calculated delta for integration", slog.Float64("delta", delta))

	return (5 * delta / 24.0) * (11*simpleExpr(leftInterval+delta) + simpleExpr(leftInterval+2*delta) + simpleExpr(leftInterval+3*delta) + 11*simpleExpr(leftInterval+4*delta)), nil
}

// Order implements NewtonCotesStrategy.
func (t *ThirdDegreeOpenNewtonCotesStrategy) Order() NewtonCotesOrder {
	return ThirdOrder
}

// Type implements NewtonCotesStrategy.
func (t *ThirdDegreeOpenNewtonCotesStrategy) Type() FormulaType {
	return OpenFormulaType
}

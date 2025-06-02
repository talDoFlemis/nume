package usecases

import (
	"context"
	"errors"

	"github.com/taldoflemis/nume/internal/expressions"
)

var (
	ErrDeltaIsZero = errors.New("delta is zero")
)

type ErrorOrder uint8

const (
	LinearErrorOrder    ErrorOrder = 0
	QuadraticErrorOrder ErrorOrder = 1
	CubicErrorOrder     ErrorOrder = 2
	QuarticErrorOrder   ErrorOrder = 3
)

type DifferenceStrategy interface {
	Derivative(
		ctx context.Context,
		simpleExpr expressions.SingleVariableExpr,
		delta float64,
	) (expressions.SingleVariableExpr, error)
	DoubleDerivative(
		ctx context.Context,
		simpleExpr expressions.SingleVariableExpr,
		delta float64,
	) (expressions.SingleVariableExpr, error)
	TripleDerivative(
		ctx context.Context,
		simpleExpr expressions.SingleVariableExpr,
		delta float64,
		errorOrder ErrorOrder,
	) (expressions.SingleVariableExpr, error)
}

var (
	_ DifferenceStrategy = (*ForwardDifferenceStrategy)(nil)
	_ DifferenceStrategy = (*BackwardDifferenceStrategy)(nil)
	_ DifferenceStrategy = (*CentralDifferenceStrategy)(nil)
)

type ForwardDifferenceStrategy struct {
}

func (*ForwardDifferenceStrategy) Derivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(variable+delta) - simpleExpr(variable)
		denominator := delta

		return numerator / denominator
	}, nil
}

func (*ForwardDifferenceStrategy) DoubleDerivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(
			variable+2*delta,
		) - 2*simpleExpr(
			variable+delta,
		) + simpleExpr(
			variable,
		)

		denominator := delta * delta

		return numerator / denominator
	}, nil
}

// TripleDerivative implements DifferenceStrategy.
func (f *ForwardDifferenceStrategy) TripleDerivative(
	ctx context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
	errorOrder ErrorOrder,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}

	fn := simpleExpr

	switch errorOrder {
	case LinearErrorOrder:
		fn = func(variable float64) float64 {
			numerator := simpleExpr(variable+3*delta) - 3*simpleExpr(variable+2*delta) +
				+3*simpleExpr(variable+delta) - simpleExpr(variable)
			denominator := delta * delta * delta
			return numerator / denominator
		}
	default:
		return nil, errors.New("unsupported error order for triple derivative in forward difference strategy")
	}

	return fn, nil
}

type BackwardDifferenceStrategy struct {
}

func (*BackwardDifferenceStrategy) Derivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(variable) - simpleExpr(variable-delta)
		denominator := delta
		return numerator / denominator
	}, nil
}

func (*BackwardDifferenceStrategy) DoubleDerivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(
			variable,
		) - 2*simpleExpr(
			variable-delta,
		) + simpleExpr(
			variable-2*delta,
		)
		denominator := delta * delta
		return numerator / denominator
	}, nil
}

// TripleDerivative implements DifferenceStrategy.
func (b *BackwardDifferenceStrategy) TripleDerivative(ctx context.Context, simpleExpr expressions.SingleVariableExpr, delta float64, errorOrder ErrorOrder) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	fn := simpleExpr

	switch errorOrder {
	case LinearErrorOrder:
		fn = func(variable float64) float64 {
			numerator := -simpleExpr(variable-3*delta) + 3*simpleExpr(variable-2*delta) +
				-3*simpleExpr(variable-delta) + simpleExpr(variable)
			denominator := delta * delta * delta
			return numerator / denominator
		}
	default:
		return nil, errors.New("unsupported error order for triple derivative in backward difference strategy")
	}

	return fn, nil
}

type CentralDifferenceStrategy struct {
}

func (*CentralDifferenceStrategy) Derivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(variable+delta) - simpleExpr(variable-delta)
		//nolint:mnd
		denominator := 2 * delta
		return numerator / denominator
	}, nil
}

func (*CentralDifferenceStrategy) DoubleDerivative(
	_ context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(
			variable+delta,
		) - 2*simpleExpr(
			variable,
		) + simpleExpr(
			variable-delta,
		)
		denominator := delta * delta
		return numerator / denominator
	}, nil
}

// TripleDerivative implements DifferenceStrategy.
func (c *CentralDifferenceStrategy) TripleDerivative(ctx context.Context,
	simpleExpr expressions.SingleVariableExpr,
	delta float64,
	errorOrder ErrorOrder,
) (expressions.SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}

	fn := simpleExpr

	switch errorOrder {
	case QuadraticErrorOrder:
		fn = func(variable float64) float64 {
			numerator := simpleExpr(variable+2*delta) - 2*simpleExpr(variable+delta) +
				+2*simpleExpr(variable-delta) - simpleExpr(variable-2*delta)
			denominator := delta * delta * delta * 2
			return numerator / denominator
		}
	default:
		return nil, errors.New("unsupported error order for triple derivative in central difference strategy")
	}

	return fn, nil
}

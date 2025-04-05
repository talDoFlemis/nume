package usecases

import (
	"context"
	"errors"
)

var (
	ErrDeltaIsZero = errors.New("delta is zero")
)

type DifferenceStrategy interface {
	Derivative(
		ctx context.Context,
		simpleExpr SingleVariableExpr,
		delta float64,
	) (SingleVariableExpr, error)
	DoubleDerivative(
		ctx context.Context,
		simpleExpr SingleVariableExpr,
		delta float64,
	) (SingleVariableExpr, error)
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
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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

type BackwardDifferenceStrategy struct {
}

func (*BackwardDifferenceStrategy) Derivative(
	_ context.Context,
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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

type CentralDifferenceStrategy struct {
}

func (*CentralDifferenceStrategy) Derivative(
	_ context.Context,
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
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

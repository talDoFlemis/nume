package usecases

import (
	"context"
	"errors"
	"math"
)

var (
	ErrDeltaIsZero = errors.New("delta is zero")
)

type DifferenceStrategy interface {
	Derivate(
		ctx context.Context,
		simpleExpr SingleVariableExpr,
		delta float64,
	) (SingleVariableExpr, error)
	DoubleDerivate(
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

func (f *ForwardDifferenceStrategy) Derivate(
	ctx context.Context,
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

func (f *ForwardDifferenceStrategy) DoubleDerivate(
	ctx context.Context,
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

		denominator := math.Pow(delta, 2)

		return numerator / denominator
	}, nil
}

type BackwardDifferenceStrategy struct {
}

func (b *BackwardDifferenceStrategy) Derivate(
	ctx context.Context,
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

func (b *BackwardDifferenceStrategy) DoubleDerivate(
	ctx context.Context,
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
		denominator := math.Pow(delta, 2)
		return numerator / denominator
	}, nil
}

type CentralDifferenceStrategy struct {
}

func (b *CentralDifferenceStrategy) Derivate(
	ctx context.Context,
	simpleExpr SingleVariableExpr,
	delta float64,
) (SingleVariableExpr, error) {
	if delta == 0 {
		return nil, ErrDeltaIsZero
	}
	return func(variable float64) float64 {
		numerator := simpleExpr(variable+delta) - simpleExpr(variable-delta)
		denominator := 2 * delta
		return numerator / denominator
	}, nil
}

func (b *CentralDifferenceStrategy) DoubleDerivate(
	ctx context.Context,
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
		denominator := math.Pow(delta, 2)
		return numerator / denominator
	}, nil
}

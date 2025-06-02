package usecases

import (
	"context"
	"log/slog"
	"math"

	"github.com/taldoflemis/nume/internal/expressions"
)

type DerivativeUseCase struct {
	philosophyStrategy DifferenceStrategy
}

func NewDerivativeUseCase(philosophyStrategy DifferenceStrategy) *DerivativeUseCase {
	return &DerivativeUseCase{
		philosophyStrategy: philosophyStrategy,
	}
}

func (d *DerivativeUseCase) Derivative(
	ctx context.Context,
	value float64,
	simpleExpr expressions.SingleVariableExpr,
	initialDelta float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Starting first derivative calculation",
		"simplified_expression", simpleExpr, "value", value, "epsilon", epsilon, "max_iterations", maxNumberOfIterations,
	)

	result, err := d.ImproveDerivative(
		ctx,
		value,
		simpleExpr,
		d.philosophyStrategy.Derivative,
		initialDelta,
		epsilon,
		maxNumberOfIterations,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error calculating first derivative", "error", err)
		return 0, err
	}

	slog.InfoContext(ctx, "First derivative calculation completed", "result", result)
	return result, nil
}

func (d *DerivativeUseCase) SecondDerivative(
	ctx context.Context,
	value float64,
	simpleExpr expressions.SingleVariableExpr,
	initialDelta float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Starting second derivative calculation",
		"simplified_expression", simpleExpr, "value", value, "epsilon", epsilon, "max_iterations", maxNumberOfIterations,
	)

	result, err := d.ImproveDerivative(
		ctx,
		value,
		simpleExpr,
		d.philosophyStrategy.DoubleDerivative,
		initialDelta,
		epsilon,
		maxNumberOfIterations,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error calculating second derivative", "error", err)
		return 0, err
	}

	slog.InfoContext(ctx, "Second derivative calculation completed", "result", result)
	return result, nil
}

func (d *DerivativeUseCase) TripleDerivative(
	ctx context.Context,
	value float64,
	simpleExpr expressions.SingleVariableExpr,
	initialDelta float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (float64, error) {
	panic("not implemented yet")
}

func (d *DerivativeUseCase) ImproveDerivative(
	ctx context.Context,
	value float64,
	simpleExpr expressions.SingleVariableExpr,
	derivativeFn func(ctx context.Context, simpleExpresion expressions.SingleVariableExpr, delta float64) (expressions.SingleVariableExpr, error),
	initialDelta float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (float64, error) {
	slog.DebugContext(ctx, "Starting to improve derivative calculation",
		"simplified_expression", simpleExpr, "value", value, "epsilon", epsilon, "max_iterations", maxNumberOfIterations,
		"initial_delta", initialDelta,
		"derivative_function", derivativeFn,
	)

	currentDelta := initialDelta
	currentError := math.Inf(1)
	bestResult := 0.0

	for i := 0; i < int(maxNumberOfIterations); i++ {
		slog.DebugContext(ctx, "Current iteration", "iteration", i, "delta", currentDelta)

		derivative, err := derivativeFn(ctx, simpleExpr, currentDelta)
		if err != nil {
			slog.ErrorContext(ctx, "Error calculating derivative", "error", err, "iteration", i, "delta", currentDelta)
			return 0, err
		}

		result := derivative(value)

		slog.DebugContext(ctx, "Current iteration result", "iteration", i, "result", result, "delta", currentDelta)

		absDifference := math.Abs(result - bestResult)
		denominator := max(math.Abs(result), math.Abs(bestResult), 1e-15)
		relativeError := absDifference / denominator

		if relativeError < epsilon {
			slog.InfoContext(ctx, "Converged to result", "result", result, "delta", currentDelta)
			return result, nil
		}

		if relativeError > currentError {
			slog.InfoContext(ctx, "Error increased, taking the current result as best", "result", result, "current_error", currentError, "relative_error", relativeError)
			return result, nil
		}

		slog.DebugContext(ctx, "Result not converged and error is decreasing, adjusting delta", "result", result, "delta", currentDelta, "relative_error", relativeError)

		currentDelta /= 2.0
		bestResult = result
		currentError = relativeError
	}

	slog.InfoContext(ctx, "Max iterations reached without convergence", "max_iterations", maxNumberOfIterations, "last_result", bestResult)
	return bestResult, nil
}

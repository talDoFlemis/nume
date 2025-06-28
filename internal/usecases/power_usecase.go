package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"

	"gonum.org/v1/gonum/mat"
)

type PowerUseCase struct{}

func NewPowerUseCase() *PowerUseCase {
	return &PowerUseCase{}
}

type PowerResult struct {
	Eigenvalue    float64
	Eigenvector   []float64
	NumIterations uint64
}

func (u *PowerUseCase) RegularPower(
	ctx context.Context,
	matrix [][]float64,
	initialGuess []float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (*PowerResult, error) {
	slog.DebugContext(ctx, "Starting the regular power method",
		slog.Any("matrix", matrix),
		slog.Any("initialGuess", initialGuess),
		slog.Float64("epsilon", epsilon),
		slog.Uint64("maxNumberOfIterations", maxNumberOfIterations),
	)

	if all(initialGuess, func(value float64) bool { return value == 0 }) {
		slog.ErrorContext(ctx, "Initial guess cannot be zero")
		return nil, errors.New("zero initial guess")
	}

	if len(matrix) == 0 || len(matrix[0]) == 0 {
		slog.ErrorContext(ctx, "Matrix cannot be empty")
		return nil, errors.New("empty matrix")
	}

	if len(matrix[0]) != len(initialGuess) {
		slog.ErrorContext(ctx, "Matrix and initial guess dimensions do not match",
			slog.Int("matrixRows", len(matrix)),
			slog.Int("matrixCols", len(matrix[0])),
		)
		return nil, errors.New("matrix and initial guess dimensions do not match")
	}

	A := constructMatrix(matrix)
	initialGuessVector := constructVector(initialGuess)

	slog.DebugContext(ctx, "Normalizing the initial guess vector")

	bestEigenvector := mat.NewVecDense(initialGuessVector.Len(), nil)
	// Normalize the initialGuess vector
	bestEigenvector.ScaleVec(1/initialGuessVector.Norm(2), initialGuessVector)

	currentError := math.Inf(1)
	currentIteration := uint64(0)
	Y := mat.NewVecDense(initialGuessVector.Len(), nil)

	var bestEigenvalue float64

	for currentIteration < maxNumberOfIterations {
		currentIteration++

		slog.DebugContext(ctx, "Iteration",
			slog.Uint64("iteration", currentIteration),
			slog.Float64("currentError", currentError),
			slog.String("bestEigenvector", fmt.Sprintf("%v", bestEigenvector.RawVector().Data)),
			slog.Float64("bestEigenvalue", bestEigenvalue),
		)

		Y.MulVec(A, bestEigenvector)

		slog.DebugContext(ctx, "Multiplying matrix A with the current Y eigenvector",
			slog.String("Y", fmt.Sprintf("%v", Y.RawVector().Data)),
		)

		normY := Y.Norm(2)
		if normY == 0 {
			slog.WarnContext(ctx, "Norm is 0, cannot continue iterating",
				slog.Any("Y", mat.Formatted(Y)),
			)
			break
		}

		// Takes the largest element in absolute value from Y
		possibleBestEigenvalue := Y.AtVec(0)
		for _, element := range Y.RawVector().Data {
			if math.Abs(element) > possibleBestEigenvalue {
				possibleBestEigenvalue = element
			}
		}

		slog.DebugContext(ctx, "Largest absolute element in Y",
			slog.Float64("largestElement", possibleBestEigenvalue),
		)

		// Calculate the iteration error with relative error
		iterationError := math.Abs((possibleBestEigenvalue - bestEigenvalue) / possibleBestEigenvalue)
		slog.DebugContext(ctx, "Calculated iteration error",
			slog.Float64("iterationError", iterationError),
		)

		currentError = iterationError
		bestEigenvalue = possibleBestEigenvalue

		// Scale down the the best eigenvector to be a fraction of the eigenvalue from the multiplication from the previous iteration with A
		bestEigenvector.ScaleVec(1/bestEigenvalue, Y)

		if iterationError < epsilon {
			slog.DebugContext(ctx, "The current error is less than epsilon, stopping the iterations",
				slog.Float64("iterationError", iterationError),
				slog.Float64("epsilon", epsilon),
			)
			break
		}
	}

	slog.InfoContext(ctx, "Finished the regular power method",
		slog.Float64("bestEigenvalue", bestEigenvalue),
		slog.String("bestEigenvector", fmt.Sprintf("%v", bestEigenvector.RawVector().Data)),
		slog.Uint64("numIterations", currentIteration),
		slog.Float64("finalError", currentError),
		slog.Float64("epsilon", epsilon),
	)

	return &PowerResult{
		Eigenvalue:    bestEigenvalue,
		Eigenvector:   bestEigenvector.RawVector().Data,
		NumIterations: currentIteration,
	}, nil
}

func constructMatrix(matrix [][]float64) *mat.Dense {
	rows := len(matrix)
	cols := len(matrix[0])
	data := make([]float64, rows*cols)

	for i := range rows {
		for j := range cols {
			data[i*cols+j] = matrix[i][j]
		}
	}

	return mat.NewDense(rows, cols, data)
}

func constructVector(vector []float64) *mat.VecDense {
	return mat.NewVecDense(len(vector), vector)
}

func all(values []float64, condition func(float64) bool) bool {
	for _, item := range values {
		if !condition(item) {
			return false
		}
	}

	return true
}

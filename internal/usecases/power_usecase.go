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

	result, err := u.innerRegularPower(ctx, A, initialGuessVector, epsilon, maxNumberOfIterations)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to compute the regular power method", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute the regular power method: %w", err)
	}

	slog.InfoContext(ctx, "Finished the regular power method",
		slog.Float64("bestEigenvalue", result.Eigenvalue),
		slog.String("bestEigenvector", fmt.Sprintf("%v", result.Eigenvector)),
		slog.Uint64("numIterations", result.NumIterations),
		slog.Float64("epsilon", epsilon),
	)

	return result, nil
}

func (u *PowerUseCase) InversePower(
	ctx context.Context,
	matrix [][]float64,
	initialGuess []float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (*PowerResult, error) {
	slog.DebugContext(ctx, "Starting the inverse power method",
		slog.Any("matrix", matrix),
		slog.Any("initialGuess", initialGuess),
		slog.Float64("epsilon", epsilon),
		slog.Uint64("maxNumberOfIterations", maxNumberOfIterations),
	)

	originalMatrix := constructMatrix(matrix)

	var inverseMatrix mat.Dense

	slog.DebugContext(ctx, "Computing the inverse of the matrix")
	err := inverseMatrix.Inverse(originalMatrix)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to compute the inverse of the matrix", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute the inverse of the matrix: %w", err)
	}

	slog.DebugContext(ctx, "Inverse matrix computed successfully",
		slog.Any("inverseMatrix", inverseMatrix.RawMatrix().Data),
	)

	result, err := u.innerRegularPower(ctx, &inverseMatrix, constructVector(initialGuess), epsilon, maxNumberOfIterations)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to compute the inverse power method", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute the inverse power method: %w", err)
	}

	// Normalize the eigenvector with the real eigenvalue
	eigenvalue := 1.0 / result.Eigenvalue

	slog.InfoContext(ctx, "Finished the inverse power method",
		slog.Float64("bestEigenvalue", eigenvalue),
		slog.String("bestEigenvector", fmt.Sprintf("%v", result.Eigenvector)),
		slog.Uint64("numIterations", result.NumIterations),
		slog.Float64("epsilon", epsilon),
	)

	return &PowerResult{
		Eigenvector:   result.Eigenvector,
		Eigenvalue:    eigenvalue,
		NumIterations: result.NumIterations,
	}, nil
}

func (u *PowerUseCase) FarthestEigenvaluePower(
	ctx context.Context,
	matrix [][]float64,
	initialGuess []float64,
	scalarToGoFarthest float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (*PowerResult, error) {
	slog.DebugContext(ctx, "Starting the Farthest power method",
		slog.Any("matrix", matrix),
		slog.Any("initialGuess", initialGuess),
		slog.Float64("epsilon", epsilon),
		slog.Uint64("maxNumberOfIterations", maxNumberOfIterations),
		slog.Float64("scalarToGoFarthest", scalarToGoFarthest),
	)

	slog.DebugContext(ctx, "Creating matrix and scalar farthest matrix")

	A := constructMatrix(matrix)
	scalarFarthestMatrix := mat.NewDense(len(matrix[0]), len(matrix[0]), nil)
	for i := 0; i < len(matrix[0]); i++ {
		scalarFarthestMatrix.Set(i, i, -1.0*scalarToGoFarthest)
	}

	slog.DebugContext(ctx, "Scalar farthest matrix created",
		slog.Any("scalarFarthestMatrix", scalarFarthestMatrix.RawMatrix().Data),
	)

	var matrixToFindLargestPowerResult mat.Dense
	matrixToFindLargestPowerResult.Add(A, scalarFarthestMatrix)

	slog.DebugContext(ctx, "Matrix to find largest power result",
		slog.Any("matrixToFindLargestPowerResult", matrixToFindLargestPowerResult.RawMatrix().Data),
	)

	initialGuessVector := constructVector(initialGuess)

	result, err := u.innerRegularPower(ctx, &matrixToFindLargestPowerResult, initialGuessVector, epsilon, maxNumberOfIterations)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to compute the farthest power method", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute the farthest power method: %w", err)
	}

	farthestEigenvalue := result.Eigenvalue + scalarToGoFarthest

	// Extract the correct eigenvector from the original matrix using eigenvalue decomposition
	eigenvector, err := u.extractEigenvectorFromMatrix(ctx, A, farthestEigenvalue)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract eigenvector from original matrix", slog.Any("error", err))
		return nil, fmt.Errorf("failed to extract eigenvector from original matrix: %w", err)
	}

	slog.InfoContext(ctx, "Finished the Farthest power method",
		slog.Float64("farthestEigenvalue", farthestEigenvalue),
		slog.String("farthestEigenvector", fmt.Sprintf("%v", eigenvector)),
		slog.Uint64("numIterations", result.NumIterations),
	)

	return &PowerResult{
		Eigenvalue:    farthestEigenvalue,
		Eigenvector:   eigenvector,
		NumIterations: result.NumIterations,
	}, nil
}

func (u *PowerUseCase) NearestEigenvaluePower(
	ctx context.Context,
	matrix [][]float64,
	initialGuess []float64,
	scalarToGoNearest float64,
	epsilon float64,
	maxNumberOfIterations uint64,
) (*PowerResult, error) {
	slog.DebugContext(ctx, "Starting the NearestEigenvaluePower method",
		slog.Any("matrix", matrix),
		slog.Any("initialGuess", initialGuess),
		slog.Float64("epsilon", epsilon),
		slog.Uint64("maxNumberOfIterations", maxNumberOfIterations),
		slog.Float64("scalarToGoNearest", scalarToGoNearest),
	)

	slog.DebugContext(ctx, "Creating matrix and scalar nearest matrix")

	A := constructMatrix(matrix)
	scalarNearestMatrix := mat.NewDense(len(matrix[0]), len(matrix[0]), nil)
	for i := 0; i < len(matrix[0]); i++ {
		scalarNearestMatrix.Set(i, i, -1.0*scalarToGoNearest)
	}

	slog.DebugContext(ctx, "Scalar nearest matrix created",
		slog.Any("scalarNearestMatrix", scalarNearestMatrix.RawMatrix().Data),
	)

	var matrixToFindSmallestPowerResult mat.Dense
	matrixToFindSmallestPowerResult.Add(A, scalarNearestMatrix)

	slog.DebugContext(ctx, "Matrix to find smallest power result",
		slog.Any("matrixToFindSmallestPowerResult", matrixToFindSmallestPowerResult.RawMatrix().Data),
	)

	matrixAsSlice := denseToSliceOfSlices(&matrixToFindSmallestPowerResult)

	result, err := u.InversePower(ctx, matrixAsSlice, initialGuess, epsilon, maxNumberOfIterations)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to compute the nearest eigenvalue power method", slog.Any("error", err))
		return nil, fmt.Errorf("failed to compute the nearest eigenvalue power method: %w", err)
	}

	nearestEigenvalue := result.Eigenvalue + scalarToGoNearest

	// Extract the correct eigenvector from the original matrix using eigenvalue decomposition
	eigenvector, err := u.extractEigenvectorFromMatrix(ctx, A, nearestEigenvalue)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract eigenvector from original matrix", slog.Any("error", err))
		return nil, fmt.Errorf("failed to extract eigenvector from original matrix: %w", err)
	}

	slog.InfoContext(ctx, "Finished the NearestEigenvaluePower method",
		slog.Float64("nearestEigenvalue", nearestEigenvalue),
		slog.String("nearestEigenvector", fmt.Sprintf("%v", eigenvector)),
		slog.Uint64("numIterations", result.NumIterations),
	)

	return &PowerResult{
		Eigenvalue:    nearestEigenvalue,
		Eigenvector:   eigenvector,
		NumIterations: result.NumIterations,
	}, nil
}

func (u *PowerUseCase) innerRegularPower(ctx context.Context,
	matrix *mat.Dense,
	initialGuess *mat.VecDense,
	epsilon float64,
	maxNumberOfIterations uint64,
) (*PowerResult, error) {
	slog.DebugContext(ctx, "Starting the inner regular power method",
		slog.Any("matrix", matrix.RawMatrix().Data),
		slog.Any("initialGuess", initialGuess.RawVector().Data),
		slog.Float64("epsilon", epsilon),
		slog.Uint64("maxNumberOfIterations", maxNumberOfIterations),
	)

	slog.DebugContext(ctx, "Normalizing the initial guess vector")

	bestEigenvector := mat.NewVecDense(initialGuess.Len(), nil)
	// Normalize the initialGuess vector
	const l2Norm = 2
	bestEigenvector.ScaleVec(1/initialGuess.Norm(l2Norm), initialGuess)

	currentError := math.Inf(1)
	currentIteration := uint64(0)
	Y := mat.NewVecDense(initialGuess.Len(), nil)

	var bestEigenvalue float64

	for currentIteration < maxNumberOfIterations {
		currentIteration++

		slog.DebugContext(ctx, "Iteration",
			slog.Uint64("iteration", currentIteration),
			slog.Float64("currentError", currentError),
			slog.String("bestEigenvector", fmt.Sprintf("%v", bestEigenvector.RawVector().Data)),
			slog.Float64("bestEigenvalue", bestEigenvalue),
		)

		Y.MulVec(matrix, bestEigenvector)

		slog.DebugContext(ctx, "Multiplying matrix A with the calculated Y eigenvector",
			slog.String("Y", fmt.Sprintf("%v", Y.RawVector().Data)),
		)

		normY := Y.Norm(l2Norm)
		if normY == 0 {
			slog.WarnContext(ctx, "Norm is 0, cannot continue iterating",
				slog.Any("Y", mat.Formatted(Y)),
			)
			break
		}

		// Takes the largest element in absolute value from Y
		possibleBestEigenvalue := mat.Dot(Y, bestEigenvector)

		bestEigenvector.ScaleVec(1/normY, Y)

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

		if iterationError < epsilon {
			slog.DebugContext(ctx, "The current error is less than epsilon, stopping the iterations",
				slog.Float64("iterationError", iterationError),
				slog.Float64("epsilon", epsilon),
			)
			break
		}
	}

	slog.InfoContext(ctx, "Finished the inner regular power method",
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

func denseToSliceOfSlices(m *mat.Dense) [][]float64 {
	r, c := m.Dims()

	result := make([][]float64, r)
	for i := range result {
		result[i] = make([]float64, c)
	}

	for i := range r {
		for j := range c {
			result[i][j] = m.At(i, j)
		}
	}
	return result
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

// extractEigenvectorFromMatrix uses Gonum's eigenvalue decomposition to find
// the eigenvector corresponding to the given eigenvalue from the original matrix
func (u *PowerUseCase) extractEigenvectorFromMatrix(ctx context.Context, matrix *mat.Dense, targetEigenvalue float64) ([]float64, error) {
	slog.DebugContext(ctx, "Extracting eigenvector from matrix using eigenvalue decomposition",
		slog.Float64("targetEigenvalue", targetEigenvalue),
	)

	var eig mat.Eigen
	ok := eig.Factorize(matrix, mat.EigenRight)
	if !ok {
		return nil, errors.New("eigenvalue decomposition failed")
	}

	eigenvalues := eig.Values(nil)
	n := len(eigenvalues)

	// Find the eigenvalue closest to our target
	bestIndex := 0
	bestDiff := math.Abs(real(eigenvalues[0]) - targetEigenvalue)

	for i := 1; i < n; i++ {
		diff := math.Abs(real(eigenvalues[i]) - targetEigenvalue)
		if diff < bestDiff {
			bestDiff = diff
			bestIndex = i
		}
	}

	slog.DebugContext(ctx, "Found closest eigenvalue",
		slog.Float64("foundEigenvalue", real(eigenvalues[bestIndex])),
		slog.Float64("targetEigenvalue", targetEigenvalue),
		slog.Float64("difference", bestDiff),
		slog.Int("eigenvalueIndex", bestIndex),
	)

	// Extract the corresponding eigenvector
	var eigenvectors mat.CDense
	eig.VectorsTo(&eigenvectors)

	r, _ := eigenvectors.Dims()
	eigenvector := make([]float64, r)

	for i := 0; i < r; i++ {
		eigenvector[i] = real(eigenvectors.At(i, bestIndex))
	}

	slog.DebugContext(ctx, "Extracted eigenvector",
		slog.Any("eigenvector", eigenvector),
	)

	return eigenvector, nil
}

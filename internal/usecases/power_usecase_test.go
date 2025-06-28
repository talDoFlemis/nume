package usecases

import (
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

type powerTestCase struct {
	name                string
	matrix              [][]float64
	initialGuess        []float64
	epsilon             float64
	expectedEigenvalue  float64
	expectedEigenvector []float64
}

func TestRegularPowerMethod(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Arrange
	t.Parallel()

	tests := []powerTestCase{
		{
			name: "Simple 2x2 Matrix",
			matrix: [][]float64{
				{2, 3},
				{5, 4},
			},
			initialGuess:        []float64{1, 1},
			epsilon:             1e-5,
			expectedEigenvalue:  7,
			expectedEigenvector: []float64{3.0 / 5, 1},
		},
		{
			name: "Simple 3x3 Matrix",
			matrix: [][]float64{
				{0, 2, 4},
				{1, 1, -2},
				{-2, 0, 5},
			},
			initialGuess:        []float64{1, 1, 1},
			epsilon:             1e-5,
			expectedEigenvalue:  3,
			expectedEigenvector: []float64{1, -0.5, 1},
		},
		{
			name: "Another 3x3 Matrix",
			matrix: [][]float64{
				{10, 6, 7},
				{1, 7, -2},
				{2, 2, 2},
			},
			initialGuess:        []float64{1, 1, 1},
			epsilon:             1e-5,
			expectedEigenvalue:  (math.Sqrt(129) + 13.0) / 2.0,
			expectedEigenvector: []float64{(math.Sqrt(129) + 7.0) / 4.0, 0.5, 1},
		},
		{
			name: "Eigenvector 3x3 Matrix",
			matrix: [][]float64{
				{1, -1, 0},
				{-1, 2, -1},
				{0, -1, 1},
			},
			initialGuess:        []float64{1, -1, 1},
			epsilon:             1e-5,
			expectedEigenvalue:  3,
			expectedEigenvector: []float64{1, -2, 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewPowerUseCase()

			// Act
			result, err := useCase.RegularPower(t.Context(), tc.matrix, tc.initialGuess, tc.epsilon, 100)

			// Assert
			assert.NoError(t, err, "Expected no error for test case: %s", tc.name)
			assert.InDelta(t, tc.expectedEigenvalue, result.Eigenvalue, tc.epsilon*10)
			matchVectorsWithTolerance(t, tc.expectedEigenvector, result.Eigenvector, tc.epsilon*10)
		})
	}
}

func matchVectorsWithTolerance(t *testing.T, expected, actual []float64, tolerance float64) {
	actualVec := constructVector(actual)
	normalizedActualVec := mat.NewVecDense(actualVec.Len(), nil)
	normalizedActualVec.ScaleVec(1/actualVec.Norm(2), actualVec)

	expectedVec := constructVector(expected)
	normalizedExpectedVec := mat.NewVecDense(expectedVec.Len(), nil)
	normalizedExpectedVec.ScaleVec(1/expectedVec.Norm(2), expectedVec)

	for i := range expected {
		actualValue := normalizedActualVec.AtVec(i)
		expectedValue := normalizedExpectedVec.AtVec(i)
		assert.InDelta(t, expectedValue, actualValue, tolerance,
			"Expected normalized value %v but got %v at index %d", expectedValue, actualValue, i)
	}
}

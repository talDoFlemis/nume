package usecases

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

type householderMethodTest struct {
	inputMatrix               [][]float64
	expectedTridiagonalMatrix [][]float64
	epsilon                   float64
}

func TestHouseholderMethod(t *testing.T) {
	// Arrange
	t.Parallel()
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	tests := []householderMethodTest{
		{
			inputMatrix: [][]float64{
				{3, 1, -1},
				{1, 3, -1},
				{-1, -1, 5},
			},
			expectedTridiagonalMatrix: [][]float64{
				{3.0, -1.4142, 0.0},
				{-1.4142, 5.0, 1.0},
				{0.0, 1, 3.0},
			},
			epsilon: 1e-3,
		},
	}

	for _, tc := range tests {
		testName := fmt.Sprintf("Householder method with input matrix: %v", tc.inputMatrix)

		t.Run(testName, func(t *testing.T) {
			useCase := NewSimilarityTransformationUseCase()

			// Act
			result, err := useCase.householderMethod(t.Context(), tc.inputMatrix)

			// Assert
			assert.NoError(t, err)
			compareMatricesWithTolerance(t, tc.expectedTridiagonalMatrix, result.T, tc.epsilon)
		})
	}
}

func compareMatricesWithTolerance(t *testing.T, expectedMatrix [][]float64, actualMatrix *mat.Dense, epsilon float64) {
	for i := range len(expectedMatrix) {
		for j := range len(expectedMatrix[i]) {
			assert.InDelta(t, expectedMatrix[i][j], actualMatrix.At(i, j), epsilon)
		}
	}
}

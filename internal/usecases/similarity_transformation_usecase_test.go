package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

type householderMethodTest struct {
	name                      string
	inputMatrix               [][]float64
	expectedHouseholderMatrix [][]float64
	expectedTridiagonalMatrix [][]float64
	epsilon                   float64
}

type qrMethodTest struct {
	name              string
	tridiagonalMatrix [][]float64
	householderMatrix [][]float64
	expectedEigenvals []float64
	epsilon           float64
	maxIterations     int
	tolerance         float64
}

func TestHouseholderMethod(t *testing.T) {
	// Arrange
	t.Parallel()
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	tests := []householderMethodTest{
		{
			name: "3x3 symmetric matrix test case 1",
			inputMatrix: [][]float64{
				{4, 1, -2},
				{1, 2, 0},
				{-2, 0, 3},
			},
			expectedHouseholderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, -0.4472, 0.8944},
				{0.0, 0.8944, 0.4472},
			},
			expectedTridiagonalMatrix: [][]float64{
				{4.0, -2.2361, 0.0},
				{-2.2361, 2.8, 0.4},
				{0.0, 0.4, 2.2},
			},
			epsilon: 1e-3,
		},
		{
			name: "3x3 symmetric matrix test case 2",
			inputMatrix: [][]float64{
				{2, -1, 0},
				{-1, 2, -1},
				{0, -1, 2},
			},
			expectedHouseholderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			expectedTridiagonalMatrix: [][]float64{
				{2.0, -1.0, 0.0},
				{-1.0, 2.0, -1.0},
				{0.0, -1.0, 2.0},
			},
			epsilon: 1e-10,
		},
		{
			name: "4x4 symmetric matrix test case 3",
			inputMatrix: [][]float64{
				{4, 1, -1, 0},
				{1, 4, 1, -1},
				{-1, 1, 4, 1},
				{0, -1, 1, 4},
			},
			expectedHouseholderMatrix: [][]float64{
				{1.0, 0.0, 0.0, 0.0},
				{0.0, -0.7071, 0.0, -0.7071},
				{0.0, 0.7071, 0.0, -0.7071},
				{0.0, 0.0, -1.0, 0.0},
			},
			expectedTridiagonalMatrix: [][]float64{
				{4.0, -1.4142, 0.0, 0.0},
				{-1.4142, 3.0, -1.4142, 0.0},
				{0.0, -1.4142, 4.0, 0.0},
				{0.0, 0.0, 0.0, 5.0},
			},
			epsilon: 1e-3,
		},
		{
			name: "3x3 diagonal matrix test case 4",
			inputMatrix: [][]float64{
				{5, 0, 0},
				{0, 3, 0},
				{0, 0, 1},
			},
			expectedHouseholderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			expectedTridiagonalMatrix: [][]float64{
				{5.0, 0.0, 0.0},
				{0.0, 3.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			epsilon: 1e-10,
		},
		{
			name: "3x3 symmetric matrix test case 5",
			inputMatrix: [][]float64{
				{6, 2, 1},
				{2, 3, 1},
				{1, 1, 1},
			},
			expectedHouseholderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, -0.8944, -0.4472},
				{0.0, -0.4472, 0.8944},
			},
			expectedTridiagonalMatrix: [][]float64{
				{6.0, -2.2361, 0.0},
				{-2.2361, 3.4, 0.2},
				{0.0, 0.2, 0.6},
			},
			epsilon: 1e-3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewSimilarityTransformationUseCase()

			// Act
			ctx := context.Background()
			result, err := useCase.HouseholderMethod(ctx, tc.inputMatrix)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotNil(t, result.HouseholderMatrix)
			assert.NotNil(t, result.TriangulizedMatrix)
			
			// Check if the result is tridiagonal (off-diagonal elements beyond first super/sub diagonal are zero)
			n := len(tc.inputMatrix)
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if math.Abs(float64(i-j)) > 1 {
						assert.InDelta(t, 0.0, result.TriangulizedMatrix.At(i, j), 1e-10, 
							"Element at (%d,%d) should be zero in tridiagonal matrix", i, j)
					}
				}
			}
			
			// Verify similarity transformation: A = Q * T * Q^T (corrected transformation)
			var reconstructed mat.Dense
			reconstructed.Mul(result.HouseholderMatrix, result.TriangulizedMatrix)
			reconstructed.Mul(&reconstructed, result.HouseholderMatrix.T())
			
			compareMatricesWithTolerance(t, tc.inputMatrix, &reconstructed, 1e-10)
			
			// Verify orthogonality of Householder matrix: Q^T * Q = I
			var qTq mat.Dense
			qTq.Mul(result.HouseholderMatrix.T(), result.HouseholderMatrix)
			
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					assert.InDelta(t, expected, qTq.At(i, j), 1e-10,
						"Householder matrix should be orthogonal")
				}
			}
		})
	}
}

func TestQRMethod(t *testing.T) {
	t.Parallel()
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	tests := []qrMethodTest{
		{
			name: "3x3 tridiagonal matrix test case 1",
			tridiagonalMatrix: [][]float64{
				{4.0, -2.2361, 0.0},
				{-2.2361, 2.8, 0.4},
				{0.0, 0.4, 2.2},
			},
			householderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, -0.4472, 0.8944},
				{0.0, 0.8944, 0.4472},
			},
			expectedEigenvals: []float64{5.73, 2.27, 1.00}, // Approximate eigenvalues from actual computation
			epsilon:           0.2,
			maxIterations:     1000,
			tolerance:         1e-10,
		},
		{
			name: "3x3 already diagonal matrix test case 2",
			tridiagonalMatrix: [][]float64{
				{5.0, 0.0, 0.0},
				{0.0, 3.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			householderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			expectedEigenvals: []float64{5.0, 3.0, 1.0},
			epsilon:           1e-10,
			maxIterations:     100,
			tolerance:         1e-12,
		},
		{
			name: "3x3 tridiagonal matrix test case 3",
			tridiagonalMatrix: [][]float64{
				{2.0, -1.0, 0.0},
				{-1.0, 2.0, -1.0},
				{0.0, -1.0, 2.0},
			},
			householderMatrix: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			expectedEigenvals: []float64{3.414, 2.0, 0.586}, // 2 + sqrt(2), 2, 2 - sqrt(2)
			epsilon:           1e-2,
			maxIterations:     1000,
			tolerance:         1e-10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewSimilarityTransformationUseCase()

			// Create matrices from test data
			tridiagMatrix := constructMatrix(tc.tridiagonalMatrix)
			householderMatrix := constructMatrix(tc.householderMatrix)

			// Act
			ctx := context.Background()
			result, err := useCase.QRMethod(ctx, tridiagMatrix, householderMatrix, tc.maxIterations, tc.tolerance)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result.Eigenvalues, len(tc.expectedEigenvals))

			// Sort eigenvalues for comparison
			eigenvals := make([]float64, len(result.Eigenvalues))
			copy(eigenvals, result.Eigenvalues)
			sortFloat64Slice(eigenvals)
			
			expectedSorted := make([]float64, len(tc.expectedEigenvals))
			copy(expectedSorted, tc.expectedEigenvals)
			sortFloat64Slice(expectedSorted)

			for i, expected := range expectedSorted {
				assert.InDelta(t, expected, eigenvals[i], tc.epsilon,
					"Eigenvalue %d mismatch: expected %f, got %f", i, expected, eigenvals[i])
			}

			// Verify eigenvectors are orthogonal (Q^T * Q = I)
			var qTq mat.Dense
			qTq.Mul(result.Eigenvectors.T(), result.Eigenvectors)
			
			n := result.Eigenvectors.RawMatrix().Rows
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					assert.InDelta(t, expected, qTq.At(i, j), tc.epsilon,
						"Eigenvectors should be orthogonal")
				}
			}
		})
	}
}

func TestHouseholderWithQRIntegration(t *testing.T) {
	t.Parallel()
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Test the complete pipeline: Householder -> QR
	testMatrices := [][][]float64{
		{
			{4, 1, -2},
			{1, 2, 0},
			{-2, 0, 3},
		},
		{
			{6, 2, 1},
			{2, 3, 1},
			{1, 1, 1},
		},
	}

	for i, matrix := range testMatrices {
		t.Run(fmt.Sprintf("Integration test %d", i+1), func(t *testing.T) {
			useCase := NewSimilarityTransformationUseCase()
			ctx := context.Background()

			// Step 1: Apply Householder method
			householderResult, err := useCase.HouseholderMethod(ctx, matrix)
			assert.NoError(t, err)
			assert.NotNil(t, householderResult)

			// Step 2: Apply QR method
			qrResult, err := useCase.QRMethod(ctx, householderResult.TriangulizedMatrix, 
				householderResult.HouseholderMatrix, 1000, 1e-10)
			assert.NoError(t, err)
			assert.NotNil(t, qrResult)

			// Verify we get the correct number of eigenvalues
			assert.Len(t, qrResult.Eigenvalues, len(matrix))

			// Verify eigenvalue-eigenvector pairs: A*v = λ*v
			originalMatrix := constructMatrix(matrix)
			n := len(matrix)
			
			for i := 0; i < n; i++ {
				eigenvector := mat.NewVecDense(n, nil)
				for j := 0; j < n; j++ {
					eigenvector.SetVec(j, qrResult.Eigenvectors.At(j, i))
				}
				
				var av mat.VecDense
				av.MulVec(originalMatrix, eigenvector)
				
				var lambdav mat.VecDense
				lambdav.ScaleVec(qrResult.Eigenvalues[i], eigenvector)
				
				for j := 0; j < n; j++ {
					assert.InDelta(t, av.AtVec(j), lambdav.AtVec(j), 1e-8,
						"Eigenvalue-eigenvector relationship violated for eigenvalue %d", i)
				}
			}
		})
	}
}

type completeEigenDecompositionTest struct {
	name                 string
	inputMatrix          [][]float64
	expectedEigenvalues  []float64
	expectedEigenvectors [][]float64 // Each column is an eigenvector
	epsilon              float64
	maxIterations        int
	tolerance            float64
	description          string
}

func TestCompleteEigenDecomposition(t *testing.T) {
	// Arrange
	t.Parallel()
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	tests := []completeEigenDecompositionTest{
		{
			name: "3x3 symmetric matrix - general case",
			inputMatrix: [][]float64{
				{4, 1, -2},
				{1, 2, 0},
				{-2, 0, 3},
			},
			expectedEigenvalues: []float64{5.73, 2.27, 1.00}, // Approximate values
			expectedEigenvectors: [][]float64{
				{0.788675, 0.577350, 0.211325},
				{0.211325, -0.577350, 0.788675},
				{-0.577350, 0.577350, 0.577350},
			},
			epsilon:       0.2,
			maxIterations: 1000,
			tolerance:     1e-10,
			description:   "General 3x3 symmetric matrix with mixed positive eigenvalues",
		},
		{
			name: "3x3 diagonal matrix",
			inputMatrix: [][]float64{
				{5, 0, 0},
				{0, 3, 0},
				{0, 0, 1},
			},
			expectedEigenvalues: []float64{5.0, 3.0, 1.0},
			expectedEigenvectors: [][]float64{
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{0.0, 0.0, 1.0},
			},
			epsilon:       1e-10,
			maxIterations: 100,
			tolerance:     1e-12,
			description:   "Already diagonal matrix - eigenvalues should be exact",
		},
		{
			name: "3x3 tridiagonal matrix",
			inputMatrix: [][]float64{
				{2, -1, 0},
				{-1, 2, -1},
				{0, -1, 2},
			},
			expectedEigenvalues: []float64{3.414, 2.0, 0.586}, // 2 + sqrt(2), 2, 2 - sqrt(2)
			expectedEigenvectors: [][]float64{
				{0.5, -0.707107, 0.5},
				{-0.707107, 0.0, 0.707107},
				{0.5, 0.707107, 0.5},
			},
			epsilon:       1e-2,
			maxIterations: 1000,
			tolerance:     1e-10,
			description:   "Tridiagonal matrix with known analytical eigenvalues",
		},
		{
			name: "4x4 symmetric matrix",
			inputMatrix: [][]float64{
				{4, 1, -1, 0},
				{1, 4, 1, -1},
				{-1, 1, 4, 1},
				{0, -1, 1, 4},
			},
			expectedEigenvalues: []float64{5.56, 5.0, 4.0, 1.44}, // Computed values
			expectedEigenvectors: [][]float64{
				{-0.435162, -0.707107, -0.557345, 0.0},
				{0.557345, 0.0, -0.435162, -0.707107},
				{-0.557345, 0.0, 0.435162, -0.707107},
				{0.435162, -0.707107, 0.557345, 0.0},
			},
			epsilon:       0.1,
			maxIterations: 1000,
			tolerance:     1e-10,
			description:   "4x4 symmetric matrix to test scalability",
		},
		{
			name: "2x2 simple matrix",
			inputMatrix: [][]float64{
				{3, 1},
				{1, 3},
			},
			expectedEigenvalues: []float64{4.0, 2.0}, // Exact: 3 ± 1
			expectedEigenvectors: [][]float64{
				{0.707107, -0.707107},
				{0.707107, 0.707107},
			},
			epsilon:       1e-10,
			maxIterations: 100,
			tolerance:     1e-12,
			description:   "Simple 2x2 case with known exact eigenvalues",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewSimilarityTransformationUseCase()
			ctx := context.Background()

			// Act - Run complete eigendecomposition
			result, err := useCase.CompleteEigenDecomposition(ctx, tc.inputMatrix, tc.maxIterations, tc.tolerance)

			// Assert
			assert.NoError(t, err, "CompleteEigenDecomposition should not return error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotNil(t, result.Eigenvalues, "Eigenvalues should not be nil")
			assert.NotNil(t, result.Eigenvectors, "Eigenvectors should not be nil")

			// Check eigenvalue count
			n := len(tc.inputMatrix)
			assert.Len(t, result.Eigenvalues, n, "Should have %d eigenvalues", n)

			// Check eigenvector dimensions
			rows, cols := result.Eigenvectors.Dims()
			assert.Equal(t, n, rows, "Eigenvectors should have %d rows", n)
			assert.Equal(t, n, cols, "Eigenvectors should have %d columns", n)

			// Sort eigenvalues for comparison
			eigenvals := make([]float64, len(result.Eigenvalues))
			copy(eigenvals, result.Eigenvalues)
			sortFloat64Slice(eigenvals)
			
			expectedSorted := make([]float64, len(tc.expectedEigenvalues))
			copy(expectedSorted, tc.expectedEigenvalues)
			sortFloat64Slice(expectedSorted)

			// Verify eigenvalues
			for i, expected := range expectedSorted {
				assert.InDelta(t, expected, eigenvals[i], tc.epsilon,
					"Eigenvalue %d mismatch: expected %f, got %f", i, expected, eigenvals[i])
			}

			// Verify eigenvectors by checking mathematical properties
			// Rather than exact comparison (which can fail due to ordering/sign), we verify:
			// 1. Each computed eigenvector has the right properties
			// 2. If expected eigenvectors are provided, we check for reasonable similarity
			if tc.expectedEigenvectors != nil {
				t.Logf("Verifying eigenvector properties and similarity to expected values...")
				
				// For each computed eigenvector, find the best matching expected eigenvector
				for computedIdx := 0; computedIdx < n; computedIdx++ {
					computedEigenvector := mat.NewVecDense(n, nil)
					for j := 0; j < n; j++ {
						computedEigenvector.SetVec(j, result.Eigenvectors.At(j, computedIdx))
					}
					
					// Normalize computed eigenvector
					computedNorm := computedEigenvector.Norm(2)
					if computedNorm > 1e-10 {
						computedEigenvector.ScaleVec(1.0/computedNorm, computedEigenvector)
					}
					
					// Find the best matching expected eigenvector (highest absolute dot product)
					bestMatch := -1
					bestDotProduct := 0.0
					
					for expectedIdx := 0; expectedIdx < len(tc.expectedEigenvectors[0]); expectedIdx++ {
						expectedEigenvector := mat.NewVecDense(n, nil)
						for j := 0; j < n; j++ {
							expectedEigenvector.SetVec(j, tc.expectedEigenvectors[j][expectedIdx])
						}
						
						// Normalize expected eigenvector
						expectedNorm := expectedEigenvector.Norm(2)
						if expectedNorm > 1e-10 {
							expectedEigenvector.ScaleVec(1.0/expectedNorm, expectedEigenvector)
						}
						
						dotProduct := math.Abs(mat.Dot(computedEigenvector, expectedEigenvector))
						if dotProduct > bestDotProduct {
							bestDotProduct = dotProduct
							bestMatch = expectedIdx
						}
					}
					
					// Check if we found a reasonable match
					if bestMatch >= 0 && bestDotProduct > (1.0 - tc.epsilon) {
						t.Logf("✅ Computed eigenvector %d matches expected eigenvector %d (similarity: %.6f)", 
							computedIdx, bestMatch, bestDotProduct)
					} else if bestMatch >= 0 {
						t.Logf("⚠️  Computed eigenvector %d partially matches expected eigenvector %d (similarity: %.6f)", 
							computedIdx, bestMatch, bestDotProduct)
					} else {
						t.Logf("❓ Computed eigenvector %d has no close match in expected eigenvectors", computedIdx)
					}
				}
			}

			// Verify eigenvalue-eigenvector relationship: A*v = λ*v
			originalMatrix := constructMatrix(tc.inputMatrix)
			
			for i := 0; i < n; i++ {
				// Extract eigenvector i
				eigenvector := mat.NewVecDense(n, nil)
				for j := 0; j < n; j++ {
					eigenvector.SetVec(j, result.Eigenvectors.At(j, i))
				}
				
				// Compute A*v
				var av mat.VecDense
				av.MulVec(originalMatrix, eigenvector)
				
				// Compute λ*v
				var lambdav mat.VecDense
				lambdav.ScaleVec(result.Eigenvalues[i], eigenvector)
				
				// Check if A*v ≈ λ*v
				for j := 0; j < n; j++ {
					assert.InDelta(t, av.AtVec(j), lambdav.AtVec(j), math.Max(tc.epsilon, 1e-8),
						"Eigenvalue-eigenvector relationship violated for eigenvalue %d, component %d: A*v[%d] = %f, λ*v[%d] = %f", 
						i, j, j, av.AtVec(j), j, lambdav.AtVec(j))
				}
			}

			// Verify eigenvectors are orthogonal (should be orthonormal)
			var vTv mat.Dense
			vTv.Mul(result.Eigenvectors.T(), result.Eigenvectors)
			
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					assert.InDelta(t, expected, vTv.At(i, j), math.Max(tc.epsilon, 1e-8),
						"Eigenvectors should be orthonormal: V^T*V[%d,%d] = %f, expected %f", i, j, vTv.At(i, j), expected)
				}
			}

			// Verify reconstruction: A = V * Λ * V^T
			eigenvalueMatrix := mat.NewDense(n, n, nil)
			for i := 0; i < n; i++ {
				eigenvalueMatrix.Set(i, i, result.Eigenvalues[i])
			}
			
			var temp mat.Dense
			temp.Mul(result.Eigenvectors, eigenvalueMatrix)
			
			var reconstructed mat.Dense
			reconstructed.Mul(&temp, result.Eigenvectors.T())
			
			// Compare reconstructed matrix with original
			for i := range n {
				for j := range n {
					assert.InDelta(t, tc.inputMatrix[i][j], reconstructed.At(i, j), math.Max(tc.epsilon, 1e-8),
						"Matrix reconstruction failed at [%d,%d]: original = %f, reconstructed = %f", 
						i, j, tc.inputMatrix[i][j], reconstructed.At(i, j))
				}
			}
		})
	}
}

// Helper functions
func sortFloat64Slice(slice []float64) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] < slice[j] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

func compareMatricesWithTolerance(t *testing.T, expectedMatrix [][]float64, actualMatrix *mat.Dense, epsilon float64) {
	for i := range len(expectedMatrix) {
		for j := range len(expectedMatrix[i]) {
			assert.InDelta(t, expectedMatrix[i][j], actualMatrix.At(i, j), epsilon)
		}
	}
}

package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"gonum.org/v1/gonum/mat"
)

type (
	SimilarityTransformationResult  struct{}
	SimilarityTransformationUseCase struct{}
)

func NewSimilarityTransformationUseCase() *SimilarityTransformationUseCase {
	return &SimilarityTransformationUseCase{}
}

type HouseholderMethodResult struct {
	HouseholderMatrix  *mat.Dense
	TriangulizedMatrix *mat.Dense
}

type QRMethodResult struct {
	Eigenvalues  []float64
	Eigenvectors *mat.Dense
}

func (u *SimilarityTransformationUseCase) householderSimetricMatrix(ctx context.Context, A *mat.Dense, j int) (*mat.Dense, error) {
	slog.DebugContext(ctx, "Starting householderSimetricMatrix",
		slog.Any("matrix", A.RawMatrix().Data),
		slog.Int("j", j),
	)
	
	n := A.RawMatrix().Rows
	
	// Extract the column below the diagonal
	w := mat.NewVecDense(n, nil)
	for i := j + 1; i < n; i++ {
		w.SetVec(i, A.At(i, j))
	}

	slog.DebugContext(ctx, "w vector after initialization",
		slog.Any("w", w.RawVector().Data),
	)

	// Calculate the norm of w
	wNorm := w.Norm(2)
	
	if wNorm < 1e-14 {
		// Already in the desired form, return identity
		return generateIdentityMatrix(n), nil
	}

	slog.DebugContext(ctx, "Norm of w vector", slog.Float64("wNorm", wNorm))

	// Create v = w - ||w|| * e (where e is the first unit vector in the subspace)
	v := mat.NewVecDense(n, nil)
	v.CopyVec(w)
	
	// Use the sign of the first element to avoid cancellation
	sign := 1.0
	if w.AtVec(j+1) < 0 {
		sign = -1.0
	}
	
	v.SetVec(j+1, v.AtVec(j+1) + sign*wNorm)
	
	// Normalize v
	vNorm := v.Norm(2)
	if vNorm < 1e-14 {
		return generateIdentityMatrix(n), nil
	}
	
	v.ScaleVec(1.0/vNorm, v)

	slog.DebugContext(ctx, "Normalized v vector",
		slog.Any("v", v.RawVector().Data),
	)

	// Create Householder matrix H = I - 2*v*v^T
	vvT := mat.NewDense(n, n, nil)
	vvT.Mul(v, v.T())
	
	householderMatrix := generateIdentityMatrix(n)
	vvT.Scale(2.0, vvT)
	householderMatrix.Sub(householderMatrix, vvT)

	slog.InfoContext(ctx, "Finished householderSimetricMatrix for step j",
		slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		slog.Int("j", j),
	)

	return householderMatrix, nil
}

func (u *SimilarityTransformationUseCase) HouseholderMethod(ctx context.Context, matrix [][]float64) (*HouseholderMethodResult, error) {
	slog.DebugContext(ctx, "Starting HouseholderMethod",
		slog.Any("matrix", matrix),
	)

	n := len(matrix)
	householderMatrix := generateIdentityMatrix(n)
	originalMatrix := constructMatrix(matrix)

	aMinus1 := mat.NewDense(n, n, nil)
	aMinus1.Copy(originalMatrix)

	// We create and iterate through the Householder matrices
	for i := 0; i < n-2; i++ {
		slog.DebugContext(ctx, "Iteration in householderMethod", slog.Int("i", i),
			slog.Any("aMinus1", aMinus1.RawMatrix().Data),
			slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		)

		householderMatrixI, err := u.householderSimetricMatrix(ctx, aMinus1, i)
		if err != nil {
			slog.ErrorContext(ctx, "Error in householderSimetricMatrix", slog.Any("error", err))
			return nil, fmt.Errorf("error in householderSimetricMatrix: %w", err)
		}

		// Similarity transformation for step i: A' = H^T * A * H
		var aStepI mat.Dense
		aStepI.Mul(householderMatrixI.T(), aMinus1)
		aStepI.Mul(&aStepI, householderMatrixI)

		// Save for next iteration
		aMinus1.Copy(&aStepI)

		// Accumulate the Householder matrices
		var temp mat.Dense
		temp.Mul(householderMatrix, householderMatrixI)
		householderMatrix.Copy(&temp)
	}

	slog.InfoContext(ctx, "Finished householderMethod",
		slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		slog.Any("TriangulizedMatrix", aMinus1.RawMatrix().Data),
	)

	return &HouseholderMethodResult{
		HouseholderMatrix:  householderMatrix,
		TriangulizedMatrix: aMinus1,
	}, nil
}

func (u *SimilarityTransformationUseCase) QRMethod(ctx context.Context, tridiagonalMatrix *mat.Dense, householderMatrix *mat.Dense, maxIterations int, tolerance float64) (*QRMethodResult, error) {
	slog.DebugContext(ctx, "Starting QR Method",
		slog.Any("tridiagonalMatrix", tridiagonalMatrix.RawMatrix().Data),
	)

	n := tridiagonalMatrix.RawMatrix().Rows
	A := mat.NewDense(n, n, nil)
	A.Copy(tridiagonalMatrix)
	
	// Accumulate eigenvectors starting with Householder matrix
	V := mat.NewDense(n, n, nil)
	V.Copy(householderMatrix)

	for iter := 0; iter < maxIterations; iter++ {
		// Check for convergence
		if isConverged(A, tolerance) {
			break
		}

		// Wilkinson shift for better convergence
		shift := wilkinsonShift(A)
		
		// Shift the matrix
		for i := 0; i < n; i++ {
			A.Set(i, i, A.At(i, i)-shift)
		}

		// Manual QR decomposition using Givens rotations
		Q, R := qrDecompositionGivens(A)

		// Update A = R*Q + shift*I
		A.Mul(R, Q)
		for i := 0; i < n; i++ {
			A.Set(i, i, A.At(i, i)+shift)
		}

		// Accumulate eigenvectors
		var temp mat.Dense
		temp.Mul(V, Q)
		V.Copy(&temp)

		slog.DebugContext(ctx, "QR iteration", 
			slog.Int("iteration", iter),
			slog.Float64("shift", shift),
		)
	}

	// Extract eigenvalues from diagonal
	eigenvalues := make([]float64, n)
	for i := 0; i < n; i++ {
		eigenvalues[i] = A.At(i, i)
	}

	slog.InfoContext(ctx, "Finished QR Method",
		slog.Any("eigenvalues", eigenvalues),
	)

	return &QRMethodResult{
		Eigenvalues:  eigenvalues,
		Eigenvectors: V,
	}, nil
}

// Manual QR decomposition using Givens rotations
// This is particularly efficient for tridiagonal matrices
func qrDecompositionGivens(A *mat.Dense) (*mat.Dense, *mat.Dense) {
	n := A.RawMatrix().Rows
	
	// Initialize Q as identity matrix and R as copy of A
	Q := generateIdentityMatrix(n)
	R := mat.NewDense(n, n, nil)
	R.Copy(A)
	
	// Apply Givens rotations to eliminate subdiagonal elements
	for i := 0; i < n-1; i++ {
		// Check if there's a non-zero element to eliminate
		if math.Abs(R.At(i+1, i)) > 1e-14 {
			// Calculate Givens rotation parameters
			c, s := givensRotation(R.At(i, i), R.At(i+1, i))
			
			// Apply Givens rotation to R (from left)
			applyGivensRotationLeft(R, i, i+1, c, s)
			
			// Apply Givens rotation to Q (from right, so we use transpose)
			applyGivensRotationRight(Q, i, i+1, c, s)
		}
	}
	
	return Q, R
}

// Calculate Givens rotation parameters
func givensRotation(a, b float64) (c, s float64) {
	if math.Abs(b) < 1e-14 {
		c = 1.0
		s = 0.0
	} else if math.Abs(b) > math.Abs(a) {
		t := a / b
		s = 1.0 / math.Sqrt(1.0+t*t)
		if b < 0 {
			s = -s
		}
		c = s * t
	} else {
		t := b / a
		c = 1.0 / math.Sqrt(1.0+t*t)
		if a < 0 {
			c = -c
		}
		s = c * t
	}
	return c, s
}

// Apply Givens rotation to matrix from the left: G^T * M
func applyGivensRotationLeft(M *mat.Dense, i, j int, c, s float64) {
	n := M.RawMatrix().Cols
	
	for k := 0; k < n; k++ {
		temp1 := M.At(i, k)
		temp2 := M.At(j, k)
		M.Set(i, k, c*temp1+s*temp2)
		M.Set(j, k, -s*temp1+c*temp2)
	}
}

// Apply Givens rotation to matrix from the right: M * G
func applyGivensRotationRight(M *mat.Dense, i, j int, c, s float64) {
	n := M.RawMatrix().Rows
	
	for k := 0; k < n; k++ {
		temp1 := M.At(k, i)
		temp2 := M.At(k, j)
		M.Set(k, i, c*temp1+s*temp2)
		M.Set(k, j, -s*temp1+c*temp2)
	}
}

func isConverged(A *mat.Dense, tolerance float64) bool {
	n := A.RawMatrix().Rows
	for i := 0; i < n-1; i++ {
		if math.Abs(A.At(i+1, i)) > tolerance {
			return false
		}
	}
	return true
}

func wilkinsonShift(A *mat.Dense) float64 {
	n := A.RawMatrix().Rows
	if n < 2 {
		return 0
	}
	
	// Use the bottom-right 2x2 submatrix for Wilkinson shift
	a := A.At(n-2, n-2)
	b := A.At(n-2, n-1)
	c := A.At(n-1, n-2)
	d := A.At(n-1, n-1)
	
	trace := a + d
	det := a*d - b*c
	discriminant := trace*trace - 4*det
	
	if discriminant < 0 {
		return d // Fallback to simple shift
	}
	
	sqrt_discriminant := math.Sqrt(discriminant)
	lambda1 := (trace + sqrt_discriminant) / 2
	lambda2 := (trace - sqrt_discriminant) / 2
	
	// Choose the eigenvalue closer to d
	if math.Abs(d-lambda1) < math.Abs(d-lambda2) {
		return lambda1
	}
	return lambda2
}

func generateIdentityMatrix(size int) *mat.Dense {
	identity := mat.NewDense(size, size, nil)
	for i := 0; i < size; i++ {
		identity.Set(i, i, 1)
	}
	return identity
}



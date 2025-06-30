package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"gonum.org/v1/gonum/mat"
)

type (
	SimilarityTransformationResult  struct{}
	SimilarityTransformationUseCase struct{}
)

func NewSimilarityTransformationUseCase() *SimilarityTransformationUseCase {
	return &SimilarityTransformationUseCase{}
}

type householderMethodResult struct {
	H *mat.Dense
	T *mat.Dense
}

func (u *SimilarityTransformationUseCase) householderSimetricMatrix(ctx context.Context, A *mat.Dense, j int) (*mat.Dense, error) {
	slog.DebugContext(ctx, "Starting householderSimetricMatrix",
		slog.Any("matrix", A.RawMatrix().Data),
		slog.Int("j", j),
	)
	w := mat.NewVecDense(A.RawMatrix().Rows, nil)
	wLine := mat.NewVecDense(A.RawMatrix().Rows, nil)

	for i := j + 1; i < A.RawMatrix().Rows; i++ {
		w.SetVec(i, A.At(i, j))
	}

	slog.DebugContext(ctx, "w vector after initialization",
		slog.Any("w", w.RawVector().Data),
	)

	sizeOfW := w.Norm(2)

	slog.DebugContext(ctx, "Size of w vector", slog.Float64("sizeOfW", sizeOfW))

	wLine.SetVec(j+1, sizeOfW)
	var N mat.VecDense

	N.AddScaledVec(w, -1, wLine)

	var normalizedN mat.VecDense
	normalizedN.ScaleVec(1/N.Norm(2), &N)

	slog.DebugContext(ctx, "Normalized N vector",
		slog.Any("normalizedN", normalizedN.RawVector().Data),
	)

	var rightSideMatrix mat.Dense
	rightSideMatrix.Mul(&normalizedN, normalizedN.T())
	rightSideMatrix.Scale(-2, &rightSideMatrix)

	slog.DebugContext(ctx, "Right side matrix after multiplication",
		slog.Any("rightSideMatrix", rightSideMatrix.RawMatrix().Data),
	)

	identityMatrix := generateIdentityMatrix(A.RawMatrix().Rows)

	householderMatrix := mat.NewDense(A.RawMatrix().Rows, A.RawMatrix().Rows, nil)
	householderMatrix.Add(identityMatrix, &rightSideMatrix)

	slog.InfoContext(ctx, "Finished householderSimetricMatrix for step j",
		slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		slog.Int("j", j),
	)

	return householderMatrix, nil
}

func (u *SimilarityTransformationUseCase) householderMethod(ctx context.Context, matrix [][]float64) (*householderMethodResult, error) {
	slog.DebugContext(ctx, "Starting householderMethod",
		slog.Any("matrix", matrix),
	)

	householderMatrix := generateIdentityMatrix(len(matrix))
	originalMatrix := constructMatrix(matrix)

	aMinus1 := mat.NewDense(len(matrix), len(matrix), nil)
	aMinus1.Copy(originalMatrix)

	// We create and iterate through the Householder matrices
	for i := range len(matrix) - 2 {
		slog.DebugContext(ctx, "Iteration in householderMethod", slog.Int("i", i),
			slog.Any("aMinus1", aMinus1.RawMatrix().Data),
			slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		)

		householderMatrixI, err := u.householderSimetricMatrix(ctx, aMinus1, i)
		if err != nil {
			slog.ErrorContext(ctx, "Error in householderSimetricMatrix", slog.Any("error", err))
			return nil, fmt.Errorf("error in householderSimetricMatrix: %w", err)
		}
		// Similarity transformation for step i
		var aStepI mat.Dense
		aStepI.Mul(householderMatrixI.T(), aMinus1)
		aStepI.Mul(&aStepI, householderMatrixI)

		// Save for next iteration
		aMinus1.Copy(&aStepI)

		// Accumulate the Householder matrices
		householderMatrix.Mul(householderMatrix, householderMatrixI)
	}

	slog.InfoContext(ctx, "Finished householderMethod",
		slog.Any("householderMatrix", householderMatrix.RawMatrix().Data),
		slog.Any("aMinus1", aMinus1.RawMatrix().Data),
	)

	return &householderMethodResult{
		H: householderMatrix,
		T: aMinus1,
	}, nil
}

func generateIdentityMatrix(size int) *mat.Dense {
	identity := mat.NewDense(size, size, nil)
	for i := range size {
		identity.Set(i, i, 1)
	}
	return identity
}

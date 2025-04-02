package usecases

type SingleVariableExpr func(float64) float64

type DerivativeUseCase struct {
}

func NewDerivativeUseCase() *DerivativeUseCase {
	return &DerivativeUseCase{}
}

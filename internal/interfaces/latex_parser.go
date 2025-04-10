package interfaces

import (
	"context"

	"github.com/taldoflemis/nume/internal/latex"
)

type LatexParser interface {
	ParseExpression(ctx context.Context, input string) (*latex.ExpressionNode, error)
}

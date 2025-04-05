package interfaces

import (
	"context"

	"github.com/taldoflemis/nume/internal/ast"
	"github.com/taldoflemis/nume/internal/expressions"
)

type EvaluableExpressionGenerator interface {
	GenerateSingleVariableExpression(
		ctx context.Context,
		node *ast.SingleVariableExpressionNode,
	) (expressions.SingleVariableExpr, error)
}

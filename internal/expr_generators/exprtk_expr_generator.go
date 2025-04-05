package exprgenerators

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/Pramod-Devireddy/go-exprtk"

	"github.com/taldoflemis/nume/internal/ast"
	"github.com/taldoflemis/nume/internal/expressions"
	"github.com/taldoflemis/nume/internal/interfaces"
)

type ExprTKExpressionGenerator struct {
}

var (
	_ (interfaces.EvaluableExpressionGenerator) = (*ExprTKExpressionGenerator)(nil)
)

func (e *ExprTKExpressionGenerator) GenerateSingleVariableExpression(
	ctx context.Context,
	node *ast.SingleVariableExpressionNode,
) (expressions.SingleVariableExpr, error) {
	exprtkObj := exprtk.NewExprtk()
	// Delete this object because we are using a CGO wrapper
	defer exprtkObj.Delete()

	exprtkObj.SetExpression(node.Expression)
	exprtkObj.AddStringVariable(node.VariableIdentifier)

	err := exprtkObj.CompileExpression()
	if err != nil {
		slog.ErrorContext(ctx, "failed to compile expression", slog.Any("err", err))
		return nil, err
	}

	return func(f float64) float64 {
		exprtkObj.SetStringVariableValue(
			node.VariableIdentifier,
			strconv.FormatFloat(f, 'E', -1, 64),
		)
		return exprtkObj.GetEvaluatedValue()
	}, nil
}

package parsers

import (
	"context"
	"log/slog"
	"math"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/taldoflemis/nume/internal/interfaces"
	"github.com/taldoflemis/nume/internal/latex"
)

func ptr[T any](v T) *T {
	return &v
}

type participleExpr interface {
	toLatexNode() latex.ExpressionNode
}

type participleExpression struct {
	Expression additionExpressionNode `@@`
}

// toLatexNode implements participleExpr.
func (p *participleExpression) toLatexNode() latex.ExpressionNode {
	return p.Expression.toLatexNode()
}

var (
	_ participleExpr = (*participleExpression)(nil)
	_ participleExpr = (*additionExpressionNode)(nil)
	_ participleExpr = (*multiplicationExpressionNode)(nil)
	_ participleExpr = (*powerExpressionNode)(nil)
	_ participleExpr = (*unaryExpressionNode)(nil)
)

var (
	_ primaryExpressionNode = (*participleVariableExpressionNode)(nil)
	_ primaryExpressionNode = (*participleNumberExpressionNode)(nil)
	_ primaryExpressionNode = (*participleConstantExpressionNode)(nil)
	_ primaryExpressionNode = (*parenthesesExpressionNode)(nil)
	_ primaryExpressionNode = (*squirlyExpressionNode)(nil)
	_ primaryExpressionNode = (*participleSquareRootExpressionNode)(nil)
)

type additionExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Multiplication multiplicationExpression `@@`
	Operator       string                   `( @("+" | "-")`
	Next           *additionExpressionNode  ` @@ )*`
}

// toLatexNode implements participleExpr.
func (a *additionExpressionNode) toLatexNode() latex.ExpressionNode {
	if a.Operator == "" {
		return a.Multiplication.toLatexNode()
	}

	var operator string
	switch a.Operator {
	case "+":
		operator = string(latex.PlusOperator)
	case "-":
		operator = string(latex.MinusOperator)
	default:
		panic("unknown operator: " + a.Operator)
	}

	return &latex.BinaryExpressionNode{
		LHS:      a.Multiplication.toLatexNode(),
		Operator: operator,
		RHS:      a.Next.toLatexNode(),
	}
}

type multiplicationExpression interface {
	participleExpr
	mult()
}

var (
	_ multiplicationExpression = (*multiplicationExpressionNode)(nil)
	_ multiplicationExpression = (*participleFractionExpressionNode)(nil)
)

type multiplicationExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Power    *powerExpressionNode     `@@`
	Operator string                   `( @("*" | "/")`
	Next     multiplicationExpression ` @@ )*`
}

// mult implements multiplicationExpression.
func (m *multiplicationExpressionNode) mult() {
}

// toLatexNode implements participleExpr.
func (m *multiplicationExpressionNode) toLatexNode() latex.ExpressionNode {
	if m.Operator == "" {
		return m.Power.toLatexNode()
	}

	var operator string
	switch m.Operator {
	case "*":
		operator = string(latex.MulOperator)
	case "/":
		operator = string(latex.DivOperator)
	default:
		panic("unknown operator for multiplication: " + m.Operator)
	}

	return &latex.BinaryExpressionNode{
		LHS:      m.Power.toLatexNode(),
		Operator: operator,
		RHS:      m.Next.toLatexNode(),
	}
}

type powerExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Unary    *unaryExpressionNode `@@`
	Operator string               `( @("^")`
	Next     *powerExpressionNode ` @@ )*`
}

// toLatexNode implements participleExpr.
func (p *powerExpressionNode) toLatexNode() latex.ExpressionNode {
	if p.Operator == "" {
		return p.Unary.toLatexNode()
	}

	var operator string
	switch p.Operator {
	case "^":
		operator = string(latex.PowerOperator)
	default:
		panic("unknown operator for power: " + p.Operator)
	}

	return &latex.BinaryExpressionNode{
		LHS:      p.Unary.toLatexNode(),
		Operator: operator,
		RHS:      p.Next.toLatexNode(),
	}
}

type unaryExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Operator string                `( @("+" | "-")`
	Unary    *unaryExpressionNode  ` @@ )`
	Primary  primaryExpressionNode `| @@`
}

// toLatexNode implements participleExpr.
func (u *unaryExpressionNode) toLatexNode() latex.ExpressionNode {
	if u.Operator == "" {
		res := u.Primary.toLatexNode()
		return res
	}

	var operator string
	switch u.Operator {
	case "+":
		operator = string(latex.PlusOperator)
	case "-":
		operator = string(latex.MinusOperator)
	default:
		panic("unknown operator for unary: " + u.Operator)
	}

	return &latex.UnaryExpressionNode{
		Operator:      operator,
		SubExpression: u.Primary.toLatexNode(),
	}
}

type primaryExpressionNode interface {
	participleExpr
	primary()
}

type participleVariableExpressionNode struct {
	Pos        lexer.Position
	EndPos     lexer.Position
	Tokens     []lexer.Token
	Identifier *string `@Ident`
}

// primary implements primaryExpressionNode.
func (p *participleVariableExpressionNode) primary() {
}

// toLatexNode implements ParticipleExpr.
func (p *participleVariableExpressionNode) toLatexNode() latex.ExpressionNode {
	return &latex.VariableExpressionNode{
		Identifier: *p.Identifier,
	}
}

type participleNumberExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token
	Value  *float64 `@(Float|Int)`
}

// primary implements primaryExpressionNode.
func (p *participleNumberExpressionNode) primary() {
}

// toLatexNode implements ParticipleExpr.
func (p *participleNumberExpressionNode) toLatexNode() latex.ExpressionNode {
	return &latex.NumberExpression{
		Value: *p.Value,
	}
}

type participleConstantExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token
	Value  *string `"\\" @("epsilon" | "pi")`
}

// primary implements primaryExpressionNode.
func (p *participleConstantExpressionNode) primary() {
}

// toLatexNode implements ParticipleExpr.
func (p *participleConstantExpressionNode) toLatexNode() latex.ExpressionNode {
	value := 0.0

	switch *p.Value {
	case `epsilon`:
		value = math.E
	case `pi`:
		value = math.Pi
	}

	return &latex.NumberExpression{
		Value: value,
	}
}

type participleSquareRootExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Index    *participleExpression `"\\" "sqrt" ("[" @@ "]")?`
	Radicand squirlyExpressionNode `@@`
}

// primary implements primaryExpressionNode.
func (p *participleSquareRootExpressionNode) primary() {
}

// toLatexNode implements ParticipleExpr.
func (p *participleSquareRootExpressionNode) toLatexNode() latex.ExpressionNode {
	var index latex.ExpressionNode

	index = &latex.NumberExpression{
		Value: 2.0,
	}

	if p.Index != nil {
		index = p.Index.toLatexNode()
	}

	return &latex.SquareRootExpressionNode{
		Index:    index,
		Radicand: p.Radicand.toLatexNode(),
	}
}

type participleFractionExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Numerator   squirlyExpressionNode `"\\" "frac" @@`
	Denominator squirlyExpressionNode `@@`
}

// mult implements multiplicationExpression.
func (p *participleFractionExpressionNode) mult() {
}

// toLatexNode implements ParticipleExpr.
func (p *participleFractionExpressionNode) toLatexNode() latex.ExpressionNode {
	return &latex.BinaryExpressionNode{
		LHS:      p.Numerator.toLatexNode(),
		Operator: string(latex.DivOperator),
		RHS:      p.Denominator.toLatexNode(),
	}
}

type parenthesesExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Expr participleExpression `"(" @@ ")"`
}

// primary implements primaryExpressionNode.
func (p *parenthesesExpressionNode) primary() {
}

// toLatexNode implements primaryExpressionNode.
func (p *parenthesesExpressionNode) toLatexNode() latex.ExpressionNode {
	return p.Expr.toLatexNode()
}

type squirlyExpressionNode struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Expr participleExpression `"{" @@ "}"`
}

// primary implements primaryExpressionNode.
func (s *squirlyExpressionNode) primary() {
}

// toLatexNode implements primaryExpressionNode.
func (s *squirlyExpressionNode) toLatexNode() latex.ExpressionNode {
	return s.Expr.toLatexNode()
}

type ParticipalMathJaxParser struct {
	parser *participle.Parser[participleExpression]
}

var (
	_ interfaces.LatexParser = (*ParticipalMathJaxParser)(nil)
)

func NewParticipalLatexParser() (*ParticipalMathJaxParser, error) {
	parser, err := participle.Build[participleExpression](
		participle.UseLookahead(99999),
		participle.Union[multiplicationExpression](
			&participleFractionExpressionNode{},
			&multiplicationExpressionNode{},
		),
		participle.Union[primaryExpressionNode](
			&participleVariableExpressionNode{},
			&participleConstantExpressionNode{},
			&participleNumberExpressionNode{},
			&parenthesesExpressionNode{},
			&squirlyExpressionNode{},
			&participleSquareRootExpressionNode{},
		),
	)
	if err != nil {
		slog.Error("failed to build participle parser", slog.Any("error", err))
		return nil, err
	}

	return &ParticipalMathJaxParser{
		parser: parser,
	}, nil
}

func (p *ParticipalMathJaxParser) ParseExpression(
	ctx context.Context,
	input string,
) (*latex.ExpressionNode, error) {
	panic("unimplemented")
}

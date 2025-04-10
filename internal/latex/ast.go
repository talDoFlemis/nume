package latex

import (
	"bytes"
	"fmt"
)

type ExpressionNode interface {
	String() string
	expression()
}

var (
	_ ExpressionNode = (*BinaryExpressionNode)(nil)
	_ ExpressionNode = (*UnaryExpressionNode)(nil)
	_ ExpressionNode = (*SquareRootExpressionNode)(nil)
	_ ExpressionNode = (*NumberExpression)(nil)
	_ ExpressionNode = (*VariableExpressionNode)(nil)
	_ ExpressionNode = (*VariableExpressionNode)(nil)
)

const (
	escapedBackslash = "\\"
)

type Operator string

const (
	PlusOperator  Operator = "+"
	MinusOperator Operator = "-"
	MulOperator   Operator = "*"
	DivOperator   Operator = "/"
	PowerOperator Operator = "^"
)

type BinaryExpressionNode struct {
	LHS      ExpressionNode
	Operator string
	RHS      ExpressionNode
}

// String implements ExpressionNode.
func (b *BinaryExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(b.LHS.String())
	out.WriteString(" " + b.Operator + " ")
	out.WriteString(b.RHS.String())
	out.WriteString(")")

	return out.String()
}

// expression implements ExpressionNode.
func (b *BinaryExpressionNode) expression() {
}

type UnaryExpressionNode struct {
	Operator      string
	SubExpression ExpressionNode
}

// String implements ExpressionNode.
func (u *UnaryExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(u.Operator)
	out.WriteString(u.SubExpression.String())
	out.WriteString(")")

	return out.String()
}

// expression implements ExpressionNode.
func (u *UnaryExpressionNode) expression() {
}

type SquareRootExpressionNode struct {
	Index    ExpressionNode
	Radicand ExpressionNode
}

// String implements ExpressionNode.
func (s *SquareRootExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(escapedBackslash + "sqrt[" + s.Index.String() + "]")
	out.WriteString("{" + s.Radicand.String() + "}")
	out.WriteString(")")

	return out.String()
}

// expression implements ExpressionNode.
func (s *SquareRootExpressionNode) expression() {}

type NumberExpression struct {
	Value float64
}

// String implements ExpressionNode.
func (n *NumberExpression) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// expression implements ExpressionNode.
func (n *NumberExpression) expression() {}

type VariableExpressionNode struct {
	Identifier string
}

// String implements ExpressionNode.
func (v *VariableExpressionNode) String() string {
	return v.Identifier
}

// expression implements ExpressionNode.
func (v *VariableExpressionNode) expression() {}

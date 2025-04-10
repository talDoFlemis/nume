package parsers

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taldoflemis/nume/internal/latex"
)

func TestVariableExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		expectedExpression *latex.VariableExpressionNode
	}{
		{
			name:  "Parse Single Letter",
			input: "x",
			expectedExpression: &latex.VariableExpressionNode{
				Identifier: "x",
			},
		},
		{
			name:  "Parse Word",
			input: "gabrigas",
			expectedExpression: &latex.VariableExpressionNode{
				Identifier: "gabrigas",
			},
		},
		{
			name:  "Parse Word with _",
			input: "my_variable",
			expectedExpression: &latex.VariableExpressionNode{
				Identifier: "my_variable",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

func TestParseNumberExpression(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name               string
		input              string
		expectedExpression *latex.NumberExpression
	}{
		{
			name:  "Parse integer",
			input: "42",
			expectedExpression: &latex.NumberExpression{
				Value: 42,
			},
		},
		{
			name:  "Parse float",
			input: "42.69",
			expectedExpression: &latex.NumberExpression{
				Value: 42.69,
			},
		},
		{
			name:  "Parse exponential",
			input: "42e-1",
			expectedExpression: &latex.NumberExpression{
				Value: 4.2,
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

func TestParseConstantExpression(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name               string
		input              string
		expectedExpression *latex.NumberExpression
	}{
		{
			name:  "Parse Euler number",
			input: `\epsilon`,
			expectedExpression: &latex.NumberExpression{
				Value: math.E,
			},
		},
		{
			name:  "Parse PI number",
			input: `\pi`,
			expectedExpression: &latex.NumberExpression{
				Value: math.Pi,
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

func TestParseSquareRoot(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name               string
		input              string
		expectedExpression *latex.SquareRootExpressionNode
	}{
		{
			name:  "Parse simple square root",
			input: `\sqrt{1}`,
			expectedExpression: &latex.SquareRootExpressionNode{
				Index: &latex.NumberExpression{
					Value: 2.0,
				},
				Radicand: &latex.NumberExpression{
					Value: 1.0,
				},
			},
		},
		{
			name:  "Parse square root with index",
			input: `\sqrt[4]{1}`,
			expectedExpression: &latex.SquareRootExpressionNode{
				Index: &latex.NumberExpression{
					Value: 4.0,
				},
				Radicand: &latex.NumberExpression{
					Value: 1.0,
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

func TestParseFrac(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name               string
		input              string
		expectedExpression *latex.BinaryExpressionNode
	}{
		{
			name:  "Parse fraq with group",
			input: `\frac{1}{2}`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS: &latex.NumberExpression{
					Value: 1.0,
				},
				Operator: string(latex.DivOperator),
				RHS: &latex.NumberExpression{
					Value: 2.0,
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

func TestBinaryExpression(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name               string
		input              string
		expectedExpression *latex.BinaryExpressionNode
	}{
		{
			name:  "1 + 2",
			input: `1 + 2`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS: &latex.NumberExpression{
					Value: 1.0,
				},
				Operator: string(latex.PlusOperator),
				RHS: &latex.NumberExpression{
					Value: 2.0,
				},
			},
		},
		{
			name:  "1 + 2 * 3",
			input: `1 + 2 * 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.PlusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.MulOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "1 - 2 * 3",
			input: `1 - 2 * 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.MinusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.MulOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "1 - 2 / 3",
			input: `1 - 2 / 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.MinusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.DivOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  `1 - \frac{2}{3}`,
			input: `1 - \frac{2}{3}`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.MinusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.DivOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "1 * 2 + 3",
			input: `1 * 2 + 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 1.0},
					Operator: string(latex.MulOperator),
					RHS:      &latex.NumberExpression{Value: 2.0},
				},
				Operator: string(latex.PlusOperator),
				RHS:      &latex.NumberExpression{Value: 3.0},
			},
		},
		{
			name:  "1 + 2 + 3",
			input: `1 + 2 + 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS: &latex.NumberExpression{
					Value: 1.0,
				},
				Operator: string(latex.PlusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.PlusOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "1 + 2 ^ 3",
			input: `1 + 2 ^ 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.PlusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.PowerOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "{1} + {2} ^ 3",
			input: `{1} + {2} ^ 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS:      &latex.NumberExpression{Value: 1.0},
				Operator: string(latex.PlusOperator),
				RHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 2.0},
					Operator: string(latex.PowerOperator),
					RHS:      &latex.NumberExpression{Value: 3.0},
				},
			},
		},
		{
			name:  "(1 + 2) ^ 3",
			input: `(1 + 2) ^ 3`,
			expectedExpression: &latex.BinaryExpressionNode{
				LHS: &latex.BinaryExpressionNode{
					LHS:      &latex.NumberExpression{Value: 1.0},
					Operator: string(latex.PlusOperator),
					RHS:      &latex.NumberExpression{Value: 2.0},
				},
				Operator: string(latex.PowerOperator),
				RHS: &latex.NumberExpression{
					Value: 3.0,
				},
			},
		},
		// {
		// 	name:  "a + b * c + d / e - f",
		// 	input: `a + b * c + d / e - f`,
		// 	expectedExpression: &latex.BinaryExpressionNode{
		// 		LHS: &latex.BinaryExpressionNode{
		// 			LHS: &latex.BinaryExpressionNode{
		// 				LHS: &latex.VariableExpressionNode{
		// 					Identifier: "a",
		// 				},
		// 				RHS: &latex.BinaryExpressionNode{
		// 					LHS: &latex.VariableExpressionNode{
		// 						Identifier: "b",
		// 					},
		// 					Operator: string(latex.MulOperator),
		// 					RHS: &latex.VariableExpressionNode{
		// 						Identifier: "c",
		// 					},
		// 				},
		// 				Operator: string(latex.PlusOperator),
		// 			},
		// 			RHS: &latex.BinaryExpressionNode{
		// 				LHS: &latex.VariableExpressionNode{
		// 					Identifier: "d",
		// 				},
		// 				RHS: &latex.VariableExpressionNode{
		// 					Identifier: "e",
		// 				},
		// 				Operator: string(latex.DivOperator),
		// 			},
		// 			Operator: string(latex.PlusOperator),
		// 		},
		// 		Operator: string(latex.MinusOperator),
		// 		RHS: &latex.VariableExpressionNode{
		// 			Identifier: "f",
		// 		},
		// 	},
		// },
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parser, err := NewParticipalLatexParser()
			require.NoError(t, err)

			result, err := parser.parser.ParseString("", test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expectedExpression, result.Expression.toLatexNode())
		})
	}
}

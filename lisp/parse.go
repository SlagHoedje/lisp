package lisp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	Dump(depth int)
	Evaluate() (int, error)
}

type ExpressionNode struct {
	Operation string
	Nodes     []Node
}

func (e ExpressionNode) Dump(depth int) {
	fmt.Println(strings.Repeat("  ", depth), "Expression", e.Operation)
	for _, n := range e.Nodes {
		n.Dump(depth + 1)
	}
}

func (e ExpressionNode) repeatOperation(op func(int, int) int) (int, error) {
	ret, err := e.Nodes[0].Evaluate()
	if err != nil {
		return 0, err
	}

	for _, node := range e.Nodes[1:] {
		val, err := node.Evaluate()
		if err != nil {
			return 0, err
		}

		ret = op(ret, val)
	}

	return ret, nil
}

func (e ExpressionNode) Evaluate() (int, error) {
	switch e.Operation {
	case "+":
		if len(e.Nodes) < 2 {
			return 0, errors.New("not enough arguments for addition operation")
		}

		return e.repeatOperation(func(x, y int) int {
			return x + y
		})
	case "-":
		if len(e.Nodes) < 2 {
			return 0, errors.New("not enough arguments for subtraction operation")
		}

		return e.repeatOperation(func(x, y int) int {
			return x - y
		})
	case "*":
		if len(e.Nodes) < 2 {
			return 0, errors.New("not enough arguments for multiplication operation")
		}

		return e.repeatOperation(func(x, y int) int {
			return x * y
		})
	case "/":
		if len(e.Nodes) != 2 {
			return 0, errors.New("incorrect number of arguments for division operation")
		}

		x, err := e.Nodes[0].Evaluate()
		if err != nil {
			return 0, nil
		}

		y, err := e.Nodes[1].Evaluate()
		if err != nil {
			return 0, nil
		}

		if y == 0 {
			return 0, errors.New("division by zero")
		}

		return x / y, nil
	default:
		return 0, errors.New("unknown operation " + e.Operation)
	}
}

type ValueNode struct {
	Value int
}

func (v ValueNode) Dump(depth int) {
	fmt.Println(strings.Repeat("  ", depth), "Value", v.Value)
}

func (v ValueNode) Evaluate() (int, error) {
	return v.Value, nil
}

type UnexpectedToken Token

func (t UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token in input: %v", t.Value)
}

type UnexpectedEOI struct{}

func (_ UnexpectedEOI) Error() string {
	return "unexpected end of input"
}

func trimWhitespace(input *[]Token) bool {
	if len(*input) == 0 {
		return false
	}

	if (*input)[0].Type == Whitespace {
		*input = (*input)[1:]
		return true
	} else {
		return false
	}
}

func findMatchingClose(input []Token) int {
	depth := 0

	for i, token := range input {
		switch token.Type {
		case Open:
			depth++
		case Close:
			depth--

			if depth == 0 {
				return i
			}
		}
	}

	return -1
}

func ParseExpression(input []Token) (ExpressionNode, error) {
	trimWhitespace(&input)

	if len(input) == 0 {
		return ExpressionNode{}, UnexpectedEOI{}
	}

	if input[0].Type != Identifier {
		return ExpressionNode{}, UnexpectedToken(input[0])
	}

	ret := ExpressionNode{}
	ret.Operation, input = input[0].Value, input[1:]

	for trimWhitespace(&input) {
		if len(input) == 0 {
			break
		}

		switch input[0].Type {
		case Number:
			value, _ := strconv.Atoi(input[0].Value)
			ret.Nodes = append(ret.Nodes, ValueNode{value})

			input = input[1:]
		case Open:
			closeIndex := findMatchingClose(input)
			if closeIndex == -1 {
				return ExpressionNode{}, UnexpectedEOI{}
			}

			nestedExpression, err := ParseExpression(input[1:closeIndex])
			if err != nil {
				return ExpressionNode{}, err
			}

			ret.Nodes = append(ret.Nodes, nestedExpression)

			input = input[closeIndex+1:]
		default:
			return ExpressionNode{}, UnexpectedToken(input[0])
		}
	}

	return ret, nil
}

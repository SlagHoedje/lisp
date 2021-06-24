package lisp

import (
	"errors"
	"fmt"
)

type Builtin func(*Environment, []Node) Node

type IncorrectType struct {
	Expected string
	Actual   string
}

func (i IncorrectType) Error() string {
	return fmt.Sprintf("expected %v, got %v", i.Expected, i.Actual)
}

func Add(_ *Environment, args []Node) Node {
	sum := 0
	for _, node := range args {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum += int(n)
	}
	return NumberNode(sum)
}

func Sub(_ *Environment, args []Node) Node {
	if len(args) == 1 {
		n, ok := args[0].(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
		}
		return -n
	}

	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum -= n
	}

	return sum
}

func Mul(_ *Environment, args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum *= n
	}

	return sum
}

func Div(_ *Environment, args []Node) Node {
	sum, ok := args[0].(NumberNode)
	if !ok {
		return ErrorNode{IncorrectType{"Number", args[0].TypeString()}}
	}

	for _, node := range args[1:] {
		n, ok := node.(NumberNode)
		if !ok {
			return ErrorNode{IncorrectType{"Number", node.TypeString()}}
		}
		sum /= n
	}

	return sum
}

func Head(_ *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{fmt.Errorf("expected 1 argument, got %v", len(args))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take head of empty list")}
	}

	return expr.Nodes[0]
}

func Tail(_ *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	if len(expr.Nodes) == 0 {
		return ErrorNode{errors.New("cannot take tail of empty list")}
	}

	return ExpressionNode{QExpression, expr.Nodes[1:]}
}

func List(_ *Environment, args []Node) Node {
	return ExpressionNode{QExpression, args}
}

func Eval(env *Environment, args []Node) Node {
	if len(args) > 1 {
		return ErrorNode{errors.New(fmt.Sprintf("expected 1 argument, got %v", len(args)))}
	}

	expr, ok := args[0].(ExpressionNode)
	if ok {
		return expr.EvalAsSExpr(env)
	} else {
		return args[0].Evaluate(env)
	}
}

func Join(_ *Environment, args []Node) Node {
	nodes := make([]Node, 0)
	for _, n := range args {
		expr, ok := n.(ExpressionNode)
		if !ok || expr.Type != QExpression {
			return ErrorNode{IncorrectType{"Q-Expression", n.TypeString()}}
		}
		nodes = append(nodes, expr.Nodes...)
	}
	return ExpressionNode{QExpression, nodes}
}

func Def(env *Environment, args []Node) Node {
	expr, ok := args[0].(ExpressionNode)
	if !ok || expr.Type != QExpression {
		return ErrorNode{IncorrectType{"Q-Expression", args[0].TypeString()}}
	}

	for _, node := range expr.Nodes {
		_, ok := node.(IdentifierNode)
		if !ok {
			return ErrorNode{IncorrectType{"Identifier", node.TypeString()}}
		}
	}

	args = args[1:]
	if len(args) != len(expr.Nodes) {
		return ErrorNode{fmt.Errorf("expected %v arguments, got %v", len(expr.Nodes)+1, len(args)+1)}
	}

	for i := range args {
		(*env)[expr.Nodes[i].(IdentifierNode)] = args[i]
	}

	return ExpressionNode{Type: SExpression}
}

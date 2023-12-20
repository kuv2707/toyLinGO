package interpreter

import (
	"fmt"
	"math"
	"toylingo/parser"
	"toylingo/utils"
)

func evaluateOperator(node parser.TreeNode, ctx *scopeContext) Variable {
	if node.Label == "literal" {
		return evaluatePrimary(node, ctx)
	}
	if utils.IsOneOf(node.Description, []string{"+", "-", "*", "/"}) {
		if len(node.Children) == 1 {
			return evaluateUnary(node, ctx, node.Description)
		} else {

			return evaluateDMAS(ctx, node, node.Description)
		}
	} else if utils.IsOneOf(node.Description, []string{"<", ">", "<=", ">=", "==", "!="}) {
		return evaluateComparison(ctx, node, node.Description)
	} else if node.Description == "=" {
		return evaluateAssignment(ctx, node)
	} else if node.Description == "#" {
		return evaluatePrint(ctx, node)
	} else if utils.IsOneOf(node.Description, []string{"&&", "||"}) {
		return evaluateLogical(ctx, node, node.Description)
	}
	interrupt("invalid operator " + node.Description)
	return Variable{}
}

func evaluateExpression(node *parser.TreeNode, ctx *scopeContext) Variable {
	switch node.Label {
	case "operator":
		return evaluateOperator(*node, ctx)
	case "literal":
		fallthrough
	case "primary":
		return evaluatePrimary(*node, ctx)
	case "call":
		ret := evaluateFuncCall(*node, ctx)
		if ret == nil {
			interrupt("function " + node.Description + " does not return a value but is expected to")
		}
		return *ret
	default:
		node.PrintTree("->")
		panic("invalid expression " + node.Label)
	}
}

func evaluateAssignment(ctx *scopeContext, node parser.TreeNode) Variable {
	variableName := node.Children[0].Description
	variableValue := evaluateExpression(node.Children[1], ctx)
	val, alreadyExists := ctx.variables[variableName]
	if alreadyExists {
		//todo: make a function to copy value from one pointer to another: memcpy
		writeBits(*val.pointer, int64(math.Float64bits(getNumber(variableValue))), 8)
		return val
	}
	ctx.variables[variableName] = variableValue
	variableValue.pointer.temp = false

	return variableValue
}

func evaluateLogical(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "bool" && right.vartype == "bool" {
		value := false
		switch operator {
		case "&&":
			value = getBool(left) && getBool(right)
		case "||":
			value = getBool(left) || getBool(right)
		}
		memaddr := malloc(1, ctx.scopeType, true)
		val := int64(0)
		if value {
			val = 1
		}
		writeBits(*memaddr, val, 1)
		return Variable{memaddr, TYPE_BOOLEAN}
	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}

func evaluateDMAS(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "number" && right.vartype == "number" {
		leftVal := getValue(left).(float64)
		rightVal := getValue(right).(float64)
		value := 0.0
		switch operator {
		case "+":
			value = leftVal + rightVal
		case "-":
			value = leftVal - rightVal
		case "*":
			value = leftVal * rightVal
		case "/":
			value = leftVal / rightVal
		}
		memaddr := malloc(8, ctx.scopeType, true)
		writeBits(*memaddr, int64(math.Float64bits(value)), 8)
		return Variable{memaddr, "number"}
	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}

func evaluateComparison(ctx *scopeContext, node parser.TreeNode, operator string) Variable {
	left := evaluateExpression(node.Children[0], ctx)
	right := evaluateExpression(node.Children[1], ctx)
	if left.vartype == "number" && right.vartype == "number" {
		value := false
		leftVal := getValue(left).(float64)
		rightVal := getValue(right).(float64)
		switch operator {
		case "<":
			value = leftVal < rightVal
		case ">":
			value = leftVal > rightVal
		case "<=":
			value = leftVal <= rightVal
		case ">=":
			value = leftVal >= rightVal
		case "==":
			value = leftVal == rightVal
		case "!=":
			value = leftVal != rightVal
		}
		memaddr := malloc(1, ctx.scopeType, true)
		val := int64(0)
		if value {
			val = 1
		}
		writeBits(*memaddr, val, 8)
		return Variable{memaddr, TYPE_BOOLEAN}
	} else {
		interrupt("invalid operands to binary operator " + operator)
	}
	return Variable{}
}

func evaluateUnary(node parser.TreeNode, ctx *scopeContext, operator string) Variable {
	val := evaluateExpression(node.Children[0], ctx)
	if val.vartype == "number" {
		switch operator {
		case "+":
			return val
		case "-":
			memaddr := malloc(8, ctx.scopeType, true)
			writeBits(*memaddr, int64(math.Float64bits(-getNumber(val))), 8)
			return Variable{memaddr, TYPE_NUMBER}
		default:
			interrupt("invalid unary operator " + operator)
		}
	}
	interrupt("invalid operand to unary operator " + operator)
	return Variable{}
}

func evaluatePrimary(node parser.TreeNode, ctx *scopeContext) Variable {
	val := node.Description
	if len(node.Children) > 0 {
		//func call or array or object
		return evaluateCompositeDS(node, ctx)
	}
	if utils.IsNumber(val) {
		memaddr := malloc(8, ctx.scopeType, true)
		writeBits(*memaddr, int64(math.Float64bits(utils.StringToNumber(val))), 8)
		return Variable{memaddr, TYPE_NUMBER}
	} else if utils.IsBoolean(val) {
		memaddr := malloc(1, ctx.scopeType, true)
		boolnum := 0
		if utils.StringToBoolean(val) {
			boolnum = 1
		}
		writeBits(*memaddr, int64(boolnum), 1)
		return Variable{memaddr, TYPE_BOOLEAN}
	} else {
		if v, exists := ctx.variables[val]; exists {
			copy := copyVariable(v, ctx.scopeType)
			return copy
		} else {
			interrupt("variable " + val + " does not exist in current scope")
		}
	}
	return Variable{}
}

func evaluateFuncCall(node parser.TreeNode, ctx *scopeContext) *Variable {
	funcNode := ctx.functions[node.Description]
	newCtx := pushScopeContext(TYPE_FUNCTION)
	lastValidLine := LineNo
	for i := 0; i < len(funcNode.Properties["args"].Children); i++ {
		argName := funcNode.Properties["args"].Children[i].Description
		argNode := node.Properties["args"+fmt.Sprint(i)]
		if argNode == nil {
			LineNo = lastValidLine
			interrupt("missing argument " + argName + " in function call " + funcNode.Description)
		}
		lastValidLine = argNode.LineNo
		argValue := evaluateExpression(argNode, newCtx)
		argValue.pointer.temp = false
		newCtx.variables[argName] = argValue
	}
	debug_info("calling", funcNode.Description)
	executeScope(funcNode.Properties["body"], newCtx)
	return newCtx.returnValue
}

func evaluateCompositeDS(node parser.TreeNode, ctx *scopeContext) Variable {
	// switch node.Description{
	// case "array":
	// 	return evaluateArray(node, ctx)
	// }
	panic("invalid composite data structure")

}

//returns a Variable with pointer to an array
// func evaluateArray(node parser.TreeNode, ctx *scopeContext) Variable {
// 	len:=len(node.Children)

// }

func evaluatePrint(ctx *scopeContext, node parser.TreeNode) Variable {
	value := evaluateExpression(node.Children[0], ctx)
	fmt.Print(utils.Colors["WHITE"], getNumber(value), utils.Colors["RESET"])
	return value
}

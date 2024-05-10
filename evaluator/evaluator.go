package evaluator

import (
	"fmt"

	ast "github.com/Artypuppet/monkey/ast"
	object "github.com/Artypuppet/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// function that creates new error structs
// It takes in the same arguments that would have been passed to sprintf.
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// helper function to check if object is of type ERROR_OBJ
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// The top level function to evaluate nodes in the ast
// It calls itself recursively whenever it encounters an expression
// Otherwise it delegates it to other functions
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)

		if isError(left) {
			return left
		}

		right := Eval(node.Right)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return nil
}

// hlper function to get the reference to boolean object
func nativeBoolToBooleanObject(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

// This function evaulates all the statements in a Statement slice
// by calling Eval on each statement.
func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

// This function evalulates prefix expressions i.e. expressions
// that have ! and - as their prefix.
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorRight(right)
	case "-":
		return evalMinusPrefixOperatorRight(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// This function evaluates the right for the ! operator i.e. it
// inverses the value of the right. We treat null to be false hence
// its negation results in true.
func evalBangOperatorRight(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// This function evaluates the right for the - operator when it is encountered
// as a prefix.
func evalMinusPrefixOperatorRight(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

// This function calls other functions to evaluate infix expression based on the operator
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "==":
		// we only need to check the ptr value since they will always be the same for the objects defined at the top.
		// pitfall is that the !(585 > 9) == 71 return false since we comparing ptrs the address will ofcourse
		// be different for lhs which will be the FALSE struct and whatever address 71 has.
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// This function evaluates an integer infix operation where both left and right are integers.
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// This function evalutes an if expression
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

// This function simply checks if the boolean expression should considered truthy
// TRUE is ofcourse considered true while FALSE and NULL are considered false
// The default has been kept true because if e.g. the result of the condition is an integer
// literal then in the language we consider it to be true.
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// function that evaluates block statements.
// It's different as it makes sure to terminate the execution of the block
// when it encounters a return statement.
func evalBlockStatement(node *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

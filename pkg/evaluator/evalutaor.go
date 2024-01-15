package evaluator

import (
	"fmt"

	"github.com/hculpan/kabkey/pkg/ast"

	"github.com/hculpan/kabkey/pkg/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}

		return applyFunction(node, function, args)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	}

	return nil
}

func applyFunction(node ast.Node, fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError(node, "not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)

	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError(node, "identifier not found: %s", node.Value)
	}

	return val
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if IsError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if IsError(left) {
		return left
	}

	right := Eval(node.Right, env)
	if IsError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case left.Type() != right.Type():
		return newError(node, "type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	}

	switch {
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch node.Operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
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
		return newError(node, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if IsError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return evalMinusPrefixOperatorExpression(node.Right, right)
	case "!":
		return evalBangOperatorExpression(node.Right, right)
	default:
		return newError(node, "unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(node ast.Expression, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError(node, "unknown operator: -%s", right.Type())
	}

	return &object.Integer{Value: -right.(*object.Integer).Value}
}

func evalBangOperatorExpression(node ast.Expression, right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		intLit, ok := right.(*object.Integer)
		if ok {
			return nativeBoolToBooleanObject(intLit.Value == 0)
		}

		return FALSE
	}
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	case TRUE:
		return true
	default:
		if intObj, ok := obj.(*object.Integer); ok {
			return intObj.Value != 0
		}

		return true
	}
}

func newError(node ast.Node, format string, a ...interface{}) *object.Error {
	return object.NewError(fmt.Sprintf(format, a...), node.NodeToken().LineNo, node.NodeToken().Position)
}

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

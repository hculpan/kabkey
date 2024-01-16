package evaluator

import (
	"bytes"
	"fmt"

	"github.com/hculpan/kabkey/pkg/object"
)

var builtins = map[string]object.BuiltinFunction{
	"print":   print,
	"println": println,
	"len":     length,
	"printf":  printf,
	"inspect": inspect,
	"type":    typeout,
}

func typeout(env *object.Environment, args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("incorrect number of paramets to 'typeout': expected 1, got %d", len(args))}
	}

	return &object.String{Value: string(args[0].Type())}
}

func inspect(env *object.Environment, args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("incorrect number of paramets to 'inspect': expected 1, got %d", len(args))}
	}

	return &object.String{Value: args[0].Inspect()}
}

func printf(env *object.Environment, args []object.Object) object.Object {
	if len(args) == 0 {
		return &object.Error{Message: "missing parameter in call to 'printf'"}
	} else if args[0].Type() != object.STRING_OBJ {
		return &object.Error{Message: fmt.Sprintf("first parameter to 'printf' must be %s, got %s", string(object.STRING_OBJ), args[0].Type())}
	}

	params := []interface{}{}
	for i, a := range args[1:] {
		switch t := a.(type) {
		case *object.Boolean:
			params = append(params, t.Value)
		case *object.Integer:
			params = append(params, t.Value)
		case *object.String:
			params = append(params, t.Value)
		default:
			return &object.Error{Message: fmt.Sprintf("invalid parameter to 'printf': type %s not supported", args[i+1].Type())}
		}
	}

	fmt.Printf(replaceEscapedChars(args[0].(*object.String).Value), params...)
	return nil
}

// ReplaceEscapedChars replaces escaped characters with their ASCII values.
func replaceEscapedChars(s string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			if i+1 < len(s) {
				next := s[i+1]
				switch next {
				case 'n':
					buffer.WriteByte('\n')
				case 't':
					buffer.WriteByte('\t')
				case '\\':
					buffer.WriteByte('\\')
				// Add other escape sequences here as needed
				default:
					buffer.WriteByte('\\')
					buffer.WriteByte(next)
				}
				i++ // Skip the next character as it is part of the escape sequence
			} else {
				buffer.WriteByte('\\')
			}
		} else {
			buffer.WriteByte(s[i])
		}
	}

	return buffer.String()
}

func print(env *object.Environment, args []object.Object) object.Object {
	result := &object.String{Value: ""}

	for _, o := range args {
		result.Value += o.Inspect()
	}

	fmt.Print(result.Value)

	return nil
}

func println(env *object.Environment, args []object.Object) object.Object {
	result := print(env, args)
	fmt.Println()
	return result
}

func length(env *object.Environment, args []object.Object) object.Object {
	if len(args) == 0 {
		return &object.Error{Message: "missing parameter in call to 'len'"}
	} else if len(args) > 1 {
		return &object.Error{Message: "too many parameters in call to 'len'"}
	}

	switch t := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(t.Value))}
	default:
		return &object.Error{Message: fmt.Sprintf("type %s not support for 'length'", args[0].Type())}
	}
}

func LoadBuiltins(env *object.Environment) {
	for k, v := range builtins {
		env.Set(k, &object.Function{Env: env, NativeImpl: v})
	}
}

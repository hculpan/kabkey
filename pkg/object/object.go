package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hculpan/kabkey/pkg/ast"
)

type ObjectType string

type BuiltinFunction func(env *Environment, arg []Object) Object

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ   = "STRING"
)

var extendedErrorOutput bool = true

func SetExtendedErrorOutput(v bool) {
	extendedErrorOutput = v
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct {
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message  string
	LineNo   int
	Position int
	Filename string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	if extendedErrorOutput {
		return fmt.Sprintf("[%d:%d] ERROR: %s", e.LineNo, e.Position, e.Message)
	} else {
		return fmt.Sprintf("ERROR: %s", e.Message)
	}
}

func NewError(msg string, lineNo, position int) *Error {
	return &Error{
		Message:  msg,
		LineNo:   lineNo,
		Position: position,
	}
}

type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	NativeImpl BuiltinFunction
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

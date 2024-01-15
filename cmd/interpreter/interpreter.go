package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hculpan/kabkey/pkg/evaluator"
	"github.com/hculpan/kabkey/pkg/lexer"
	"github.com/hculpan/kabkey/pkg/object"
	"github.com/hculpan/kabkey/pkg/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing file parameter")
		os.Exit(1)
	}

	input, err := LoadFileToString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	if len(l.Errors()) > 0 {
		printErrors(os.Stdout, l.Errors())
		os.Exit(1)
	} else if len(p.Errors()) > 0 {
		printErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluator.LoadBuiltins(env)
	o := evaluator.Eval(program, env)

	if o != nil {
		io.WriteString(os.Stdout, o.Inspect())
		io.WriteString(os.Stdout, "\n")
	}

}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func LoadFileToString(filename string) (string, error) {
	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", err
	}

	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Convert content to string and return
	return string(content), nil
}

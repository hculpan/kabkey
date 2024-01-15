package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/hculpan/kabkey/pkg/evaluator"
	"github.com/hculpan/kabkey/pkg/lexer"
	"github.com/hculpan/kabkey/pkg/object"
	"github.com/hculpan/kabkey/pkg/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	object.SetExtendedErrorOutput(false)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.NewLexer(line)

		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(l.Errors()) != 0 {
			printErrors(out, l.Errors())
			continue
		} else if len(p.Errors()) != 0 {
			printErrors(out, p.Errors())
			continue
		}

		o := evaluator.Eval(program, env)

		if o != nil {
			io.WriteString(out, o.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

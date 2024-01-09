package main

import (
	"fmt"
	"os"

	"github.com/hculpan/kabkey/pkg/repl"
)

func main() {
	fmt.Printf("Type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

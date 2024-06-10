package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf("%s", PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if _, err := io.WriteString(out, program.String()); err != nil {
			return
		}
		if _, err := io.WriteString(out, "\n"); err != nil {
			return
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			return
		}
	}
}

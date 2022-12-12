package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pspiagicw/uranus/pkg/evaluator"
	"github.com/pspiagicw/uranus/pkg/lexer"
	"github.com/pspiagicw/uranus/pkg/object"
	"github.com/pspiagicw/uranus/pkg/parser"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
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

		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")

        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out , evaluated.Inspect())
            io.WriteString(out, "\n")
        }
	}
}

const MONKEY_FACE = `
              (
               )
              (
        /\  .-"""-.  /\
       //\\/  ,,,  \//\\
       |/\| ,;;;;;, |/\|
       //\\\;-"""-;///\\
      //  \/   .   \/  \\
     (| ,-_| \ | / |_-, |)
       //'__\.-.-./__'\\
      // /.-(() ())-.\ \\
     (\ |)   '---'   (| /)
      ' (|           |) '
        \)           (/
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

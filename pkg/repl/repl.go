package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pspiagicw/uranus/pkg/lexer"
	"github.com/pspiagicw/uranus/pkg/parser"
)


const PROMPT = ">>> "


func Start(in io.Reader , out io.Writer) {
	scanner := bufio.NewScanner(in)

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
            printParserErrors(out , p.Errors())
        }

		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	fmt.Printf("%+v\n" , tok)
		// }

        io.WriteString(out , program.String())
        io.WriteString(out , "\n")
	}
}
func printParserErrors(out io.Writer , errors []string) {
    for _ , msg := range errors {
        io.WriteString(out , "\t" + msg + "\t")
    }
}

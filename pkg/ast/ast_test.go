package ast

import (
	"testing"

	"github.com/pspiagicw/uranus/pkg/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "var"},
					Value: "var",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let var = anotherVar;" {
		t.Errorf("program.String() was wrong: got='%s'", program.String())

	}

}

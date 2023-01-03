package parser

import (
	"fmt"
	"testing"

	"github.com/pspiagicw/uranus/pkg/ast"
	"github.com/pspiagicw/uranus/pkg/lexer"
)

func TestIndexExpressions(t *testing.T) {
	input := `myArray[1 + 1]`

	program := newProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression, got=%T", stmt.Expression)
	}

	assertIdentifier(t, indexExp.Left, "myArray")
	assertInfixExpression(t, indexExp.Index, 1, "+", 1)

}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	program := newProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("length of Array.Elements not 3, got=%d", len(array.Elements))

	}

	assertIntegerLiteral(t, array.Elements[0], 1)
	assertInfixExpression(t, array.Elements[1], 2, "*", 2)
	assertInfixExpression(t, array.Elements[2], 3, "+", 3)

}
func TestStringLitearlExpressions(t *testing.T) {
	input := `"hello world"`
	program := newProgram(t, input)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. Got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q, got=%q", "hello world", literal.Value)
	}
}
func TestCallExpressions(t *testing.T) {
	input := `add(1 , 2 * 3, 4 + 5);`
	program := newProgram(t, input)

	assertStatementCount(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}
	assertIdentifier(t, exp.Function, "add")
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	assertLiteralExpression(t, exp.Arguments[0], 1)
	assertInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	assertInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{`fn() {}`, []string{}},
		{`fn(x) {}`, []string{"x"}},
		{`fn(x, y, z) {}`, []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := newProgram(t, tt.input)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ExressionStatement , got=%T",
				program.Statements[0],
			)
		}
		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("ExpressionStatement is not FunctionLiteral , got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf(
				"length parameters wrong , want=%d , got=%d",
				len(tt.expectedParams),
				len(function.Parameters),
			)
		}

		for i, ident := range tt.expectedParams {
			assertLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestFuncLiteral(t *testing.T) {
	input := `fn(x ,y) { x + y; }`

	program := newProgram(t, input)

	assertStatementCount(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Program.Statements[0] not a ExpressionStatement , got=%T",
			program.Statements[0],
		)
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression Statement does not contain function literal , got=%T", stmt)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong, want 2, got=%d\n", len(function.Parameters))
	}

	assertLiteralExpression(t, function.Parameters[0], "x")
	assertLiteralExpression(t, function.Parameters[1], "y")

	assertStatementCount(t, function.Body.Statements, 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"function body is not ast.ExpressionStatement , got=%T",
			function.Body.Statements[0],
		)
	}

	assertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestIfExpressions(t *testing.T) {
	inputs := []struct {
		input       string
		alternative bool
	}{
		{`if ( x < y ) { x }`, false},
		{`if ( x < y ) { x } else { y }`, true},
	}

	for _, input := range inputs {

		program := newProgram(t, input.input)

		assertStatementCount(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Program.Statements[0] not a ExpressionStatement , got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)

		if !ok {
			t.Fatalf("Expression Statement does not contain if expression , got=%T", stmt)
		}

		assertInfixExpression(t, exp.Condition, "x", "<", "y")

		assertStatementCount(t, exp.Consequence.Statements, 1)

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"Statements[0] is not *ast.ExpressionStatement , got=%T",
				exp.Consequence.Statements[0],
			)
		}

		assertIdentifier(t, consequence.Expression, "x")

		if input.alternative {
			alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

			if !ok {
				t.Fatalf(
					"Alternative Statements[0] is not *ast.ExpressionStatement , got=%T",
					exp.Consequence.Statements[0],
				)
			}

			assertIdentifier(t, alternative.Expression, "y")

		} else {
			if exp.Alternative != nil {
				t.Errorf("exp.Alternative was not nil, got=%+v", exp.Alternative)
			}
		}

	}
}

func TestBooleanParsing(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{
			"true;",
			true,
		},
		{
			"false;",
			false,
		},
	}

	for _, tt := range tests {

		program := newProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("No of statements in program not 1 , got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Program.Statements[0] not a ExpressionStatement , got=%T",
				program.Statements[0],
			)
		}

		bol, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("ExpressionStatement not a Boolean , got =%T", stmt.Expression)
		}

		if bol.Value != tt.value {
			t.Errorf("Ident Value not '%v' , got=%v", tt.value, bol.Value)
		}
		if bol.TokenLiteral() != fmt.Sprintf("%v", tt.value) {
			t.Errorf("Identifier Token not '%v' , got=%s", tt.value, bol.TokenLiteral())
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == true",
			"((3 > 5) == true)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {

		program := newProgram(t, tt.input)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q , got =%q", tt.expected, actual)
		}
	}
}

// / === INFIX Expressions ===
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		// {"true == true" , true, "==" , true},
		// {"false == false" , false, "==" , false},
	}

	for _, tt := range infixTests {

		program := newProgram(t, tt.input)
		assertStatementCount(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement not a expression Statement , got=%T", program.Statements[0])
		}

		assertInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)

	}
}

// / === PREFIX Expressinos ===
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {

		program := newProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("No of statements in program not 1 , got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement not a expression Statement , got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("Expression not Prefix Expression , got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf(
				"expression operation not the same , got =%s , wanted= %s",
				exp.Operator,
				tt.operator,
			)
		}

		assertIntegerLiteral(t, exp.Right, tt.integerValue)

	}
}

/// === IDENTIFIER Tests ===

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	program := newProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("No of statements in program not 1 , got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Program.Statements[0] not a ExpressionStatement , got=%T", program.Statements[0])
	}

	assertIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	program := newProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("No of statements in program not 1 , got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Program.Statements[0] not a ExpressionStatement , got=%T", program.Statements[0])
	}

	assertIntegerLiteral(t, stmt.Expression, 5)
}

// === RETURN Statements Tests ===

func TestReturnStatement(t *testing.T) {
	input := `
    return 5;
    return 10;
    return 99818;
    `

	program := newProgram(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("No of statements in program not 3 , got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("Statement not Return Statement , got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("Statement Token Literal not 'return' , got=%s", returnStmt.TokenLiteral())
		}
	}
}

// === LET Statement Tests ===

func TestLetStatement(t *testing.T) {
	input := `
    let x = 5;
    let y = 10;
    let foobar = 8383838;
    `

	program := newProgram(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("No of statements in program not 3 , got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		assertLetStatement(t, stmt, tt.expectedIdentifier)

	}
}

func TestInvalidLetStatements(t *testing.T) {
	input := `
    let x  5;
    let = 10;
    let 8383838;
    `

	l := lexer.New(input)

	p := New(l)

	_ = p.ParseProgram()

	assertParserErrors(t, p)
}

// / === Testing Utilities ===
func assertParserErrors(t testing.TB, p *Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		t.Errorf("Errors expected , not errors present")
	}
}

func checkParserErrors(t testing.TB, p *Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) != 0 {
		t.Errorf("Parser encountered errors! , no of errors %d", len(errors))
	}

	for _, msg := range errors {
		t.Errorf("Parser Error: %s", msg)
	}
}

func assertLetStatement(t testing.TB, stmt ast.Statement, expectedIdentifier string) {
	t.Helper()

	if stmt.TokenLiteral() != "let" {
		t.Fatalf("Token Literal not 'let' , got=%s", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Fatalf("Statement not let Statement , got=%T", stmt)
	}

	if letStmt.Name.Value != expectedIdentifier {
		t.Errorf(
			"LetStatement variable name not correct, got=%s , wanted=%s",
			letStmt.Name.Value,
			expectedIdentifier,
		)
	}

	if letStmt.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf(
			"LetStatement variable token not correct , got=%s , wanted=%s",
			letStmt.Name.TokenLiteral(),
			expectedIdentifier,
		)
	}
}

func assertIdentifier(t testing.TB, exp ast.Expression, value string) {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier , got=%T", exp)
	}

	if ident.Value != value {
		t.Errorf("Identifier value not equal , got %s , wanted %s", ident.Value, value)
	}

	if ident.TokenLiteral() != value {
		t.Errorf(
			"Identifier tokenliteral not equal , got %s , wanted %s",
			ident.TokenLiteral(),
			value,
		)
	}
}

func assertStatementCount(t *testing.T, statements []ast.Statement, expected int) {
	t.Helper()
	if len(statements) != expected {
		t.Fatalf("No of statements in program not %d , got=%d", expected, len(statements))
	}
}

func assertLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		assertIntegerLiteral(t, exp, int64(v))
	case int64:
		assertIntegerLiteral(t, exp, v)
	case string:
		assertIdentifier(t, exp, v)
	default:
		t.Errorf("type of exp not handled , got=%T", exp)
	}
}

func assertInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	t.Helper()

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp not *ast.InfixExpression , got=%T", exp)
	}

	assertLiteralExpression(t, opExp.Left, left)

	if opExp.Operator != operator {
		t.Errorf("exp.Operator not '%s' , got '%s'", operator, opExp.Operator)
	}

	assertLiteralExpression(t, opExp.Right, right)
}

func newProgram(t *testing.T, input string) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("Program should not be nil")
	}

	return program
}

func assertIntegerLiteral(t testing.TB, exp ast.Expression, value int64) {
	t.Helper()

	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integ not IntegerLiteral , got=%T", exp)
	}

	if integ.Value != value {
		t.Errorf("exp Value not %d , got %d", value, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Token Literal of Integer Literal , got %d ,got %s", value, exp.TokenLiteral())
	}
}

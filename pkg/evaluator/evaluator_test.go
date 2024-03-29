package evaluator

import (
	"testing"

	"github.com/pspiagicw/uranus/pkg/lexer"
	"github.com/pspiagicw/uranus/pkg/object"
	"github.com/pspiagicw/uranus/pkg/parser"
)

func TestConcatenation(t *testing.T) {
	input := `"hello " + "world"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not string, got=%T , (%+v)", evaluated, evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("value of string wrong, got=%q", str.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T , (+%v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}
func TestClosures(t *testing.T) {
	input := `
    let newAdder = fn(x) {
        fn(y) { x + y };
    };

    let addTwo = newAdder(2);

    addTwo(2);
    `

	testIntegerObject(t, testEval(input), 4)
}
func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x;}; identity(5)", 5},
		{"let double = fn(x) { x * 2;}; double(2)", 4},
		{"let add = fn(x,y) { x + y;}; add(1,4)", 5},
		{"let add = fn(x,y) { x + y;}; add(5+5,add(5,5))", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}

}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Function has wrong number of paramaeter, got=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	expectedBody := `(x + 2)`

	if fn.Body.String() != expectedBody {
		t.Fatalf("Body is not equal %q, got=%q", expectedBody, fn.Body.String())
	}

}
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = b + a + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
		{"foobar",
			"identifier not found: foobar",
		},

		{"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{"-true",
			"unknown operator: -BOOLEAN",
		},
		{"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("not error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
        if (10 > 1) {
            if (10 > 1) {
                return 10;
            }
            return 1;
        }

        `, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func testNullObject(t testing.TB, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func testBooleanObject(t testing.TB, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean.got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%T, want=%t", result.Value, expected)

		return false
	}
	return true
}
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"15", 15},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	e := object.NewEnvironment()

	program := p.ParseProgram()

	return Eval(program, e)
}
func testIntegerObject(t testing.TB, obj object.Object, expected int64) bool {

	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("Object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Object has wrong value, got=%d , want=%d", result.Value, expected)
		return false
	}
	return true
}

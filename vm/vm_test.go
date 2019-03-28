package vm

import (
	"fmt"
	"testing"

	"github.com/skatsuta/monkey-compiler/ast"
	"github.com/skatsuta/monkey-compiler/compiler"
	"github.com/skatsuta/monkey-compiler/lexer"
	"github.com/skatsuta/monkey-compiler/object"
	"github.com/skatsuta/monkey-compiler/parser"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
	}

	runVMTests(t, tests)
}

type vmTestCase struct {
	input string
	want  interface{}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
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

	runVMTests(t, tests)
}

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		complr := compiler.New()
		if err := complr.Compile(program); err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(complr.Bytecode())
		if err := vm.Run(); err != nil {
			t.Fatalf("vm error: %s", err)
		}

		got := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.want, got)
	}
}

func parse(input string) *ast.Program {
	return parser.New(lexer.New(input)).ParseProgram()
}

func testExpectedObject(t *testing.T, want interface{}, got object.Object) {
	t.Helper()

	switch want := want.(type) {
	case bool:
		if err := testBooleanObject(bool(want), got); err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case int:
		if err := testIntegerObject(int64(want), got); err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	default:
		t.Errorf("testExpectedObject failed: unknown type %T (%#v)", got, got)
	}
}

func testBooleanObject(want bool, got object.Object) error {
	result, ok := got.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%#v)", got, got)
	}

	if result.Value != want {
		return fmt.Errorf("object has wrong value. want=%t, got=%t", want, result.Value)
	}

	return nil
}

func testIntegerObject(want int64, got object.Object) error {
	result, ok := got.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%#v)", got, got)
	}

	if result.Value != want {
		return fmt.Errorf("object has wrong value. want=%d, got=%d", want, result.Value)
	}

	return nil
}

package compiler

import (
	"fmt"
	"testing"

	"github.com/skatsuta/monkey-compiler/ast"
	"github.com/skatsuta/monkey-compiler/code"
	"github.com/skatsuta/monkey-compiler/lexer"
	"github.com/skatsuta/monkey-compiler/object"
	"github.com/skatsuta/monkey-compiler/parser"
)

type compilerTestCase struct {
	input      string
	wantConsts []interface{}
	wantInsns  []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      "1; 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 + 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 - 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 * 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "2 / 1",
			wantConsts: []interface{}{2, 1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "-1",
			wantConsts: []interface{}{1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      "true",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "false",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 > 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 < 2",
			wantConsts: []interface{}{2, 1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 == 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "1 != 2",
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "true == false",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "true != false",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "!true",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		cmplr := New()
		if err := cmplr.Compile(program); err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := cmplr.Bytecode()

		if err := testInstructions(tt.wantInsns, bytecode.Instructions); err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		if err := testConstants(tt.wantConsts, bytecode.Constants); err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	return parser.New(lexer.New(input)).ParseProgram()
}

func testInstructions(want []code.Instructions, got code.Instructions) error {
	concat := concatInstructions(want)

	if len(got) != len(concat) {
		return fmt.Errorf("wrong instructions length.\nwant:\n%s\ngot:\n%s", concat, got)
	}

	for i, insn := range concat {
		if got[i] != insn {
			return fmt.Errorf("wrong instruction at pos %d.\nwant=%q\ngot=%q", i, concat, got)
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := make(code.Instructions, 0, len(s))
	for _, insns := range s {
		out = append(out, insns...)
	}
	return out
}

func testConstants(want []interface{}, got []object.Object) error {
	if len(got) != len(want) {
		return fmt.Errorf("wrong number of constants. want=%d, got=%d", len(want), len(got))
	}

	for i, c := range want {
		switch c := c.(type) {
		case int:
			if e := testIntegerObject(int64(c), got[i]); e != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, e)
			}
		}
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

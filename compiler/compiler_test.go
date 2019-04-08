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

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      `if (true) { 10 }; 3333;`,
			wantConsts: []interface{}{10, 3333},
			wantInsns: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNil),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
		{
			input:      `if (true) { 10 } else { 20 }; 3333;`,
			wantConsts: []interface{}{10, 20, 3333},
			wantInsns: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 0010
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
			`,
			wantConsts: []interface{}{1, 2},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
			let one = 1;
			one;
			`,
			wantConsts: []interface{}{1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			wantConsts: []interface{}{1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      `"monkey"`,
			wantConsts: []interface{}{"monkey"},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:      `"mon" + "key"`,
			wantConsts: []interface{}{"mon", "key"},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      "[]",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "[1, 2, 3]",
			wantConsts: []interface{}{1, 2, 3},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "[1 + 2, 3 - 4, 5 * 6]",
			wantConsts: []interface{}{1, 2, 3, 4, 5, 6},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      "{}",
			wantConsts: []interface{}{},
			wantInsns: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "{1: 2, 3: 4, 5: 6}",
			wantConsts: []interface{}{1, 2, 3, 4, 5, 6},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "{1: 2 + 3, 4: 5 * 6}",
			wantConsts: []interface{}{1, 2, 3, 4, 5, 6},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:      "[1, 2, 3][1 + 1]",
			wantConsts: []interface{}{1, 2, 3, 1, 1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "{1: 2}[2 - 1]",
			wantConsts: []interface{}{1, 2, 2, 1},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "fn() { return 5 + 10 }",
			wantConsts: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: "fn() { 5 + 10 }",
			wantConsts: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: "fn() { 1; 2 }",
			wantConsts: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: "fn() { }",
			wantConsts: []interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	c := New()
	if c.scopeIdx != 0 {
		t.Errorf("scopeIdx wrong. want=%d, got=%d", 0, c.scopeIdx)
	}
	globalSymTab := c.symTab

	c.emit(code.OpMul)

	c.enterScope()
	if c.scopeIdx != 1 {
		t.Errorf("scopeIdx wrong. want=%d, got=%d", 1, c.scopeIdx)
	}

	c.emit(code.OpSub)

	scope := c.currentScope()
	insnsLen := len(scope.insns)
	if insnsLen != 1 {
		t.Errorf("instructions length wrong. want=%d, got=%d", 1, insnsLen)
	}
	if last := scope.lastInsn; last.Opcode != code.OpSub {
		t.Errorf("lastInsn.Opcode wrong. want=%d, got=%d", code.OpSub, last.Opcode)
	}

	if c.symTab.outer != globalSymTab {
		t.Errorf("compiler did not enclose global symbol table")
	}

	c.leaveScope()
	if c.scopeIdx != 0 {
		t.Errorf("scopeIdx wrong. want=%d, got=%d", 0, c.scopeIdx)
	}

	if c.symTab != globalSymTab {
		t.Errorf("compiler did not restore global symbol table")
	}
	if c.symTab.hasOuter() {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	c.emit(code.OpAdd)

	scope = c.currentScope()
	insnsLen = len(scope.insns)
	if insnsLen != 2 {
		t.Errorf("instructions length wrong. want=%d, got=%d", 2, insnsLen)
	}
	if last := scope.lastInsn; last.Opcode != code.OpAdd {
		t.Errorf("lastInsn.Opcode wrong. want=%d, got=%d", code.OpAdd, last.Opcode)
	}
	if prev := scope.prevInsn; prev.Opcode != code.OpMul {
		t.Errorf("prevInsn.Opcode wrong. want=%d, got=%d", code.OpMul, prev.Opcode)
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "fn() { 24 }()",
			wantConsts: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 1), // The compiled function
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let noArg = fn() { 24 };
			noArg();
			`,
			wantConsts: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 1), // The compiled function
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let oneArg = fn(a) { a };
			oneArg(24);
			`,
			wantConsts: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				24,
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0), // The compiled function
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let manyArg = fn(a, b, c) { a; b; c; };
			manyArg(24, 25, 26);
			`,
			wantConsts: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 2),
					code.Make(code.OpReturnValue),
				},
				24,
				25,
				26,
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0), // The compiled function
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let num = 55;
			fn() { num }
			`,
			wantConsts: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			fn() {
				let num = 55;
				num
			}
			`,
			wantConsts: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			fn() {
				let a = 55;
				let b = 77;
				a + b
			}
			`,
			wantConsts: []interface{}{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			len([]);
			push([], 1);
			`,
			wantConsts: []interface{}{1},
			wantInsns: []code.Instructions{
				code.Make(code.OpGetBuiltin, 0),
				code.Make(code.OpArray, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpGetBuiltin, 5),
				code.Make(code.OpArray, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: "fn() { len([]) }",
			wantConsts: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetBuiltin, 0),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			wantInsns: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:      "first(rest(push([1, 2, 3], 4)))",
			wantConsts: []interface{}{1, 2, 3, 4},
			wantInsns: []code.Instructions{
				code.Make(code.OpGetBuiltin, 2),
				code.Make(code.OpGetBuiltin, 4),
				code.Make(code.OpGetBuiltin, 5),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 2),
				code.Make(code.OpCall, 1),
				code.Make(code.OpCall, 1),
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
			return fmt.Errorf("wrong instruction at pos %d.\nwant:\n%s\ngot:\n%s", i, concat, got)
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

		case string:
			if err := testStringObject(c, got[i]); err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}

		case []code.Instructions:
			fn, ok := got[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T", i, got[i])
			}

			if err := testInstructions(c, fn.Instructions); err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}

		default:
			return fmt.Errorf("constant %d - unsupported constant type: %T", i, c)
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

func testStringObject(want string, got object.Object) error {
	result, ok := got.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%#v)", got, got)
	}

	if result.Value != want {
		return fmt.Errorf("object has wrong value. want=%q, got=%q", want, result.Value)
	}

	return nil
}

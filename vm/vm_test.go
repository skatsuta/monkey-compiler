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

type vmTestCase struct {
	input string
	want  interface{}
}

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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVMTests(t, tests)
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
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5 })", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVMTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Nil},
		{"if (false) { 10 }", Nil},
	}

	runVMTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVMTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVMTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 - 4, 5 * 6]", []int{3, -1, 30}},
	}

	runVMTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			input: "{}",
			want:  map[object.HashKey]int64{},
		},
		{
			input: "{1: 2, 2: 3}",
			want: map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			input: "{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			want: map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVMTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Nil},
		{"[1, 2, 3][99]", Nil},
		{"[1][-1]", Nil},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Nil},
		{"{}[0]", Nil},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fivePlusTen = fn() { 5 + 10; };
			fivePlusTen();
			`,
			want: 15,
		},
		{
			input: `
			let one = fn() { 1; };
			let two = fn() { 2; };
			one() + two();
			`,
			want: 3,
		},
		{
			input: `
			let a = fn() { 1; };
			let b = fn() { a() + 1; };
			let c = fn() { b() + 1; };
			c();
			`,
			want: 3,
		},
	}

	runVMTests(t, tests)
}

func TestFunctionsWithReturnStatements(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let earlyExit = fn() { return 99; 100; };
			earlyExit();
			`,
			want: 99,
		},
		{
			input: `
			let earlyExit = fn() { return 99; return 100; };
			earlyExit();
			`,
			want: 99,
		},
	}

	runVMTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() { };
			noReturn();
			`,
			want: Nil,
		},
		{
			input: `
			let noReturn = fn() { };
			let noReturnTwo = fn() { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
			want: Nil,
		},
	}

	runVMTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() { 1; };
			let returnsOneReturner = fn() { returnsOne; };
			returnsOneReturner()();
			`,
			want: 1,
		},
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

	case string:
		if err := testStringObject(want, got); err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}

	case []int:
		arr, ok := got.(*object.Array)
		if !ok {
			t.Errorf("object is not Array. got=%T (%#v)", got, got)
			return
		}

		if len(arr.Elements) != len(want) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(want), len(arr.Elements))
			return
		}

		for i, el := range want {
			if err := testIntegerObject(int64(el), arr.Elements[i]); err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case map[object.HashKey]int64:
		hash, ok := got.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%#v)", got, got)
		}

		if len(hash.Pairs) != len(want) {
			t.Errorf(
				"hash has wrong number of pairs. want=%d (%#v), got=%d (%#v)",
				len(want), want, len(hash.Pairs), hash.Pairs,
			)
		}

		for wantKey, wantVal := range want {
			pair, ok := hash.Pairs[wantKey]
			if !ok {
				t.Errorf("no pair for given key %v in pairs", wantKey)
			}

			if err := testIntegerObject(wantVal, pair.Value); err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case *object.Nil:
		if got != Nil {
			t.Errorf("object is not Nil: %T (%#v)", got, got)
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
		return fmt.Errorf("object is not Integer. got=%T (%+v)", got, got)
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

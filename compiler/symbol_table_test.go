package compiler

import "testing"

func TestDefine(t *testing.T) {
	want := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
		"c": Symbol{Name: "c", Scope: LocalScope, Index: 0},
		"d": Symbol{Name: "d", Scope: LocalScope, Index: 1},
		"e": Symbol{Name: "e", Scope: LocalScope, Index: 0},
		"f": Symbol{Name: "f", Scope: LocalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != want["a"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "a", want["a"], a)
	}

	b := global.Define("b")
	if b != want["b"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "b", want["b"], b)
	}

	firstLocal := NewEnclosedSymbolTable(global)

	c := firstLocal.Define("c")
	if c != want["c"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "c", want["c"], c)
	}

	d := firstLocal.Define("d")
	if d != want["d"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "d", want["d"], d)
	}

	secondLocal := NewEnclosedSymbolTable(firstLocal)

	e := secondLocal.Define("e")
	if e != want["e"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "e", want["e"], e)
	}

	f := secondLocal.Define("f")
	if f != want["f"] {
		t.Errorf("symbol %q: want=%#v, got=%#v", "f", want["f"], f)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	wantSymbols := []Symbol{
		Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, want := range wantSymbols {
		got, ok := global.Resolve(want.Name)
		if !ok {
			t.Errorf("name %q not resolvable", got.Name)
			continue
		}

		if got != want {
			t.Errorf("expected %q to resolve to %#v, but got %#v", want.Name, want, got)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expectedSymbols := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, want := range expectedSymbols {
		got, ok := local.Resolve(want.Name)
		if !ok {
			t.Errorf("name %q not resolvable", want.Name)
			continue
		}

		if got != want {
			t.Errorf("expected %q to resolve to %+v, but got %+v", want.Name, want, got)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table       *SymbolTable
		wantSymbols []Symbol
	}{
		{
			table: firstLocal,
			wantSymbols: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			table: secondLocal,
			wantSymbols: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, want := range tt.wantSymbols {
			got, ok := tt.table.Resolve(want.Name)
			if !ok {
				t.Errorf("name %q not resolvable", want.Name)
				continue
			}

			if got != want {
				t.Errorf("expected %q to resolve to %+v, but got %+v", want.Name, want, got)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(global)

	wantSymbols := []Symbol{
		{Name: "a", Scope: BuiltinScope, Index: 0},
		{Name: "c", Scope: BuiltinScope, Index: 1},
		{Name: "e", Scope: BuiltinScope, Index: 2},
		{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, v := range wantSymbols {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, want := range wantSymbols {
			got, ok := table.Resolve(want.Name)
			if !ok {
				t.Errorf("name %q not resolvable", want.Name)
				continue
			}

			if got != want {
				t.Errorf("expected %q to resolve to %+v, but got %+v", want.Name, want, got)
			}
		}
	}
}

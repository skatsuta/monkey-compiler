package compiler

import "testing"

func TestDefine(t *testing.T) {
	want := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
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

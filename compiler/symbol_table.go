package compiler

// SymbolScope represents a scope of symbols.
type SymbolScope string

const (
	// GlobalScope represents a global scope, i.e. top level context of a program.
	GlobalScope SymbolScope = "GLOBAL"
	// LocalScope represents a local scope, i.e. a function level context.
	LocalScope SymbolScope = "LOCAL"
	// BuiltinScope represents a scope for built-in functions.
	BuiltinScope SymbolScope = "BUILTIN"
	// FreeScope represents a scope for closures referencing free variables.
	FreeScope SymbolScope = "FREE"
)

// Symbol is a symbol defined in a scope with an identifier (name).
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable is a mapping table of identifiers (names) and defined symbols.
type SymbolTable struct {
	FreeSymbols []Symbol

	outer *SymbolTable

	store   map[string]Symbol
	numDefs int
}

// NewSymbolTable creates a new symbol table.
func NewSymbolTable() *SymbolTable {
	return NewEnclosedSymbolTable(nil)
}

// NewEnclosedSymbolTable creates a new symbol table with an outer one.
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		FreeSymbols: make([]Symbol, 0),
		outer:       outer,
		store:       make(map[string]Symbol),
	}
}

// Define defines an identifier as a symbol in a scope.
func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Scope: GlobalScope, Index: s.numDefs}
	if s.hasOuter() {
		sym.Scope = LocalScope
	}

	s.store[name] = sym
	s.numDefs++
	return sym
}

// Resolve resolves an identifier and returns a defined symbol and `true` if any.
// If the identifier is not found anywhere within a chain of symbol tables, it returns an empty
// symbol and `false`.
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	if sym, exists := s.store[name]; exists || !s.hasOuter() {
		return sym, exists
	}

	sym, exists := s.outer.Resolve(name)
	if exists && (sym.Scope == LocalScope || sym.Scope == FreeScope) {
		// Define an outer local or free variable as a free variable in the current scope
		sym = s.defineFree(sym)
	}
	return sym, exists
}

// hasOuter returns true if `s` has an outer symbol table, otherwise false.
func (s *SymbolTable) hasOuter() bool {
	return s.outer != nil
}

// DefineBuiltin defines
func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	sym := Symbol{Name: name, Scope: BuiltinScope, Index: index}
	s.store[name] = sym
	return sym
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	return s.define(original.Name, FreeScope, len(s.FreeSymbols)-1)
}

func (s *SymbolTable) define(name string, scope SymbolScope, index int) Symbol {
	sym := Symbol{Name: name, Scope: scope, Index: index}
	s.store[name] = sym
	return sym
}

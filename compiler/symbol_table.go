package compiler

// SymbolScope represents a scope of symbols.
type SymbolScope string

const (
	// GlobalScope represents a global scope, i.e. top level context of programs.
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol is a symbol defined in a scope with an identifier (name).
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable is a mapping table of identifiers (names) and defined symbols.
type SymbolTable struct {
	store   map[string]Symbol
	numDefs int
}

// NewSymbolTable creates a new symbol table.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

// Define defines an identifier as a symbol in a scope.
func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Scope: GlobalScope, Index: s.numDefs}
	s.store[name] = sym
	s.numDefs++
	return sym
}

// Resolve resolves an identifier and returns a defined symbol and `true` if any.
// If the identifier is not found in a symbol table, it returns an empty symbol and `false`.
func (s *SymbolTable) Resolve(name string) (sym Symbol, exists bool) {
	sym, exists = s.store[name]
	return sym, exists
}

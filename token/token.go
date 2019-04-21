package token

// Type is a token type.
type Type string

const (
	// Illegal is a token type for illegal tokens.
	Illegal Type = "Illegal"
	// EOF is a token type that represents end of file.
	EOF = "EOF"

	// Ident is a token type for identifiers.
	Ident = "Ident" // add, foobar, x, y, ...
	// Int is a token type for integers.
	Int = "Int"
	// Float is a token type for floating point numbers.
	Float = "Float"
	// String is a token type for strings.
	String = "String"

	// Bang is a token type for NOT operator.
	Bang = "!"
	// Plus is a token type for addition.
	Plus = "+"
	// Minus is a token type for subtraction.
	Minus = "-"
	// Astarisk is a token type for multiplication.
	Astarisk = "*"
	// Slash is a token type for division.
	Slash = "/"
	// LT is a token ype for 'less than' operator.
	LT = "<"
	// GT is a token ype for 'greater than' operator.
	GT = ">"
	// LE is a token type for 'less than or equal to' operator.
	LE = "<="
	// GE is a token type for 'greater than or equal to' operator.
	GE = ">="
	// Eq is a token type for equality operator.
	Eq = "=="
	// NEq is a token type for not equality operator.
	NEq = "!="
	// And is a token type for binary And logical operator.
	And = "&&"
	// Or is a token type for binary Or logical operator.
	Or = "||"
	// Assign is a token type for assignment operators.
	Assign = "="

	// Comma is a token type for commas.
	Comma = ","
	// Semicolon is a token type for semicolons.
	Semicolon = ";"
	// Colon is a token type for colons.
	Colon = ":"

	// LParen is a token type for left parentheses.
	LParen = "("
	// RParen is a token type for right parentheses.
	RParen = ")"
	// LBrace is a token type for left braces.
	LBrace = "{"
	// RBrace is a token type for right braces.
	RBrace = "}"
	// LBracket is a token type for left brackets.
	LBracket = "["
	// RBracket is a token type for right brackets.
	RBracket = "]"

	// Function is a token type for functions.
	Function = "Function"
	// Let is a token type for lets.
	Let = "Let"
	// True is a token type for true.
	True = "True"
	// False is a token type for false.
	False = "False"
	// Nil is a token type for nil.
	Nil = "Nil"
	// If is a token type for if.
	If = "If"
	// Else is a token type for else.
	Else = "Else"
	// Return is a token type for return.
	Return = "Return"
	// Macro is a token type for macros.
	Macro = "Macro"
)

// Token represents a token which has a token type and literal.
type Token struct {
	Type    Type
	Literal string
}

// Language keywords
var keywords = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"nil":    Nil,
	"if":     If,
	"else":   Else,
	"return": Return,
	"macro":  Macro,
}

// LookupIdent checks the language keywords to see whether the given identifier is a keyword.
// If it is, it returns the keyword's Type constant. If it isn't, it just gets back IDENT.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}

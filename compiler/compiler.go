package compiler

import (
	"fmt"
	"sort"

	"github.com/skatsuta/monkey-compiler/ast"
	"github.com/skatsuta/monkey-compiler/code"
	"github.com/skatsuta/monkey-compiler/object"
)

// EmittedInstruction represents an instruction emitted at a position.
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// CompilationScope represents a function scope at compilation.
type CompilationScope struct {
	insns              code.Instructions
	lastInsn, prevInsn EmittedInstruction
}

// Compiler is a bytecode compiler.
type Compiler struct {
	// consts is a slice that serves as a constant pool.
	consts []object.Object

	symTab *SymbolTable

	scopes   []CompilationScope
	scopeIdx int
}

// New creates a new Compiler.
func New() *Compiler {
	return NewWithState(NewSymbolTable(), make([]object.Object, 0))
}

// NewWithState creates a new Compiler with a given symbol table and constant pool.
func NewWithState(symTab *SymbolTable, consts []object.Object) *Compiler {
	mainScope := CompilationScope{
		insns: make(code.Instructions, 0),
	}
	return &Compiler{
		consts: consts,
		symTab: symTab,
		scopes: []CompilationScope{mainScope},
	}
}

// Compile compiles an AST node to a bytecode.
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}

		c.emit(code.OpPop)

	case *ast.LetStatement:
		// Compile the right-hand side expression
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		// Define an identifier as a symbol in a proper scope
		sym := c.symTab.Define(node.Name.Value)
		if sym.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, sym.Index)
		} else {
			c.emit(code.OpSetLocal, sym.Index)
		}

	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return nil
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown unary operator: %s", node.Operator)
		}

	case *ast.InfixExpression:
		// Reverse the two operands if the operator is "<" (less than)
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}

			if err := c.Compile(node.Left); err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(node.Consequence); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastInstruction()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInsns())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNil)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastInstruction()
			}
		}

		afterAlternativePos := len(c.currentInsns())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}

		c.emit(code.OpCall)

	case *ast.Ident:
		sym, ok := c.symTab.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %q", node.Value)
		}

		if sym.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, sym.Index)
		} else {
			c.emit(code.OpGetLocal, sym.Index)
		}

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IntegerLiteral:
		i := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(i))

	case *ast.StringLiteral:
		s := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(s))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			if err := c.Compile(el); err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		l := len(node.Pairs)
		keys := make([]ast.Expression, 0, l)
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			if err := c.Compile(k); err != nil {
				return err
			}
			if err := c.Compile(node.Pairs[k]); err != nil {
				return err
			}
		}

		c.emit(code.OpHash, l*2)

	case *ast.FunctionLiteral:
		c.enterScope()

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastInsnWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		// Take the number of local bindings defined in the current scope from the symbol table
		// before leaving the scope, in order to pass the number to the function later on
		numLocals := c.symTab.numDefs
		insns := c.leaveScope()

		compiledFn := &object.CompiledFunction{
			Instructions: insns,
			NumLocals:    numLocals,
		}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	}

	return nil
}

// addConstant adds a constant object to the compiler's constant pool and returns an identifier
// for the constant.
func (c *Compiler) addConstant(obj object.Object) (id int) {
	c.consts = append(c.consts, obj)
	return len(c.consts) - 1
}

// emit generates a bytecode corresponding to `op` and `operands`, adds it to the compiler's
// internal bytecode instruction sequence and returns the starting position of the instruction.
func (c *Compiler) emit(op code.Opcode, operands ...int) (pos int) {
	insn := code.Make(op, operands...)
	pos = c.addInstruction(insn)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentScope() CompilationScope {
	return c.scopes[c.scopeIdx]
}

func (c *Compiler) currentInsns() code.Instructions {
	return c.scopes[c.scopeIdx].insns
}

func (c *Compiler) addInstruction(insn []byte) (pos int) {
	insns := c.currentInsns()
	pos = len(insns)
	c.scopes[c.scopeIdx].insns = append(insns, insn...)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.scopes[c.scopeIdx].prevInsn = c.currentScope().lastInsn
	c.scopes[c.scopeIdx].lastInsn = EmittedInstruction{Opcode: op, Position: pos}
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	scope := c.currentScope()
	return len(scope.insns) > 0 && scope.lastInsn.Opcode == op
}

func (c *Compiler) removeLastInstruction() {
	scope := c.currentScope()
	c.scopes[c.scopeIdx].insns = scope.insns[:scope.lastInsn.Position]
	c.scopes[c.scopeIdx].lastInsn = scope.prevInsn
}

func (c *Compiler) replaceInstruction(pos int, newInsn []byte) {
	// The underlying assumption here is that we only replace instructions of the same type,
	// with the same non-variable length.
	insns := c.currentInsns()
	copy(insns[pos:pos+len(newInsn)], newInsn)
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.currentInsns()[opPos])
	c.replaceInstruction(opPos, code.Make(op, operand))
}

func (c *Compiler) replaceLastInsnWithReturn() {
	lastPos := c.currentScope().lastInsn.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIdx].lastInsn.Opcode = code.OpReturnValue
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		insns: make(code.Instructions, 0),
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIdx++

	// Create a new nested symbol table
	c.symTab = NewEnclosedSymbolTable(c.symTab)
}
func (c *Compiler) leaveScope() code.Instructions {
	insns := c.currentInsns()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIdx--

	// Restore the outer symbol table
	c.symTab = c.symTab.outer

	return insns
}

// Bytecode returns a bytecode generated by the compiler.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInsns(),
		Constants:    c.consts,
	}
}

// Bytecode represents a bytecode.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

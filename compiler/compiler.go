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

// Compiler is a bytecode compiler.
type Compiler struct {
	// insns holds the generated bytecode.
	insns code.Instructions
	// consts is a slice that serves as a constant pool.
	consts []object.Object

	lastInsn, prevInsn EmittedInstruction

	symTab *SymbolTable
}

// New creates a new Compiler.
func New() *Compiler {
	return NewWithState(NewSymbolTable(), make([]object.Object, 0))
}

// NewWithState creates a new Compiler with a given symbol table and constant pool.
func NewWithState(s *SymbolTable, consts []object.Object) *Compiler {
	return &Compiler{
		insns:  make(code.Instructions, 0),
		consts: consts,
		symTab: s,
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

	case *ast.LetStatement:
		// Compile the right-hand side expression
		if err := c.Compile(node.Value); err != nil {
			return err
		}

		sym := c.symTab.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, sym.Index)

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}

		c.emit(code.OpPop)

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

		afterConsequencePos := len(c.insns)
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

		afterAlternativePos := len(c.insns)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.Ident:
		sym, ok := c.symTab.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %q", node.Value)
		}

		c.emit(code.OpGetGlobal, sym.Index)

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

func (c *Compiler) addInstruction(insn []byte) (pos int) {
	pos = len(c.insns)
	c.insns = append(c.insns, insn...)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.prevInsn = c.lastInsn
	c.lastInsn = EmittedInstruction{Opcode: op, Position: pos}
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	return c.lastInsn.Opcode == op
}

func (c *Compiler) removeLastInstruction() {
	c.insns = c.insns[:c.lastInsn.Position]
	c.lastInsn = c.prevInsn
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.insns[opPos])
	newInsn := code.Make(op, operand)
	c.replaceInstruction(opPos, newInsn)
}

func (c *Compiler) replaceInstruction(pos int, newInsn []byte) {
	// The underlying assumption here is that we only replace instructions of the same type,
	// with the same non-variable length.
	copy(c.insns[pos:pos+len(newInsn)], newInsn)
}

// Bytecode returns a bytecode generated by the compiler.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.insns,
		Constants:    c.consts,
	}
}

// Bytecode represents a bytecode.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

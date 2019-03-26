package code

import "fmt"

// Instructions represents a sequence of instructions.
type Instructions []byte

// Opcode represents an opcode.
type Opcode byte

// Byte returns the corresponding byte of the opcode `op`.
func (op Opcode) Byte() byte {
	return byte(op)
}

const (
	// OpConstant represents an opcode which pushes a constant value on to the stack.
	OpConstant Opcode = iota
)

// Definition represents the definition of an opcode.
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},
}

// Lookup performs a lookup for `op` in the definitions of opcodes.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

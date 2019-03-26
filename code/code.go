package code

import (
	"encoding/binary"
	"fmt"
)

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

// Make makes a bytecode instruction sequence from an opcode and operands.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return nil
	}

	insnLen := 1
	for _, w := range def.OperandWidths {
		insnLen += w
	}

	insn := make([]byte, insnLen)
	insn[0] = op.Byte()

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2: // 2 bytes
			binary.BigEndian.PutUint16(insn[offset:], uint16(o))
		}
		offset += width
	}

	return insn
}

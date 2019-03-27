package code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Opcode represents an opcode.
type Opcode byte

const (
	// OpConstant represents an opcode which pushes a constant value on to the stack.
	OpConstant Opcode = iota
	// OpAdd represents an opcode for integer addition.
	OpAdd
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

// Instructions represents a sequence of instructions.
type Instructions []byte

func (insns Instructions) String() string {
	var out strings.Builder

	i := 0
	for i < len(insns) {
		def, err := Lookup(insns[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, insns[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, insns.formatInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (insns Instructions) formatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand length %d does not match defined %d",
			len(operands), operandCount)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s 0x%X", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operand width for %s: %d\n", def.Name, operandCount)
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
	insn[0] = byte(op)

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

// ReadOperands reads operands from bytecode instructions based on the definition of an opcode.
// It returns operands read and the offset describing the starting position of next opcode.
func ReadOperands(def *Definition, insns Instructions) (operands []int, offset int) {
	operands = make([]int, len(def.OperandWidths))

	for i, width := range def.OperandWidths {
		switch width {
		case 2: // 2 bytes
			operands[i] = int(ReadUint16(insns[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint16 reads a single uint16 value from bytecode instruction sequence.
func ReadUint16(insns Instructions) uint16 {
	return binary.BigEndian.Uint16(insns)
}

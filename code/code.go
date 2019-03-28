package code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Opcode represents an opcode.
type Opcode byte

const (
	// OpConstant is an opcode to push a constant value on to the stack.
	OpConstant Opcode = iota
	// OpPop is an opcode to pop the topmost element off the stack.
	OpPop
	// OpAdd is an opcode for addition (+).
	OpAdd
	// OpSub is an opcode for subtraction (-).
	OpSub
	// OpMul is an opcode for multiplication (*).
	OpMul
	// OpDiv is an opcode for division (/).
	OpDiv
	// OpTrue is an opcode to push `true` value on to the stack.
	OpTrue
	// OpFalse is an opcode to push `false` value on to the stack.
	OpFalse
	// OpEqual is an opcode to check the equality of the two topmost elements on the stack.
	OpEqual
	// OpNotEqual is an opcode to check the inequality of the two topmost elements on the stack.
	OpNotEqual
	// OpGreaterThan is an opcode to check the second topmost element is greater than the first.
	OpGreaterThan
	// OpMinus is an opcode to negate integers.
	OpMinus
	// OpBang is an opcode to negate booleans.
	OpBang
)

// Definition represents the definition of an opcode.
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:    {Name: "OpConstant", OperandWidths: []int{2}},
	OpPop:         {Name: "OpPop", OperandWidths: nil},
	OpAdd:         {Name: "OpAdd", OperandWidths: nil},
	OpSub:         {Name: "OpSub", OperandWidths: nil},
	OpMul:         {Name: "OpMul", OperandWidths: nil},
	OpDiv:         {Name: "OpDiv", OperandWidths: nil},
	OpTrue:        {Name: "OpTrue", OperandWidths: nil},
	OpFalse:       {Name: "OpFalse", OperandWidths: nil},
	OpEqual:       {Name: "OpEqual", OperandWidths: nil},
	OpNotEqual:    {Name: "OpNotEqual", OperandWidths: nil},
	OpGreaterThan: {Name: "OpGreaterThan", OperandWidths: nil},
	OpMinus:       {Name: "OpMinus", OperandWidths: nil},
	OpBang:        {Name: "OpBang", OperandWidths: nil},
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
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s 0x%X", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operand width for %s: %d", def.Name, operandCount)
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

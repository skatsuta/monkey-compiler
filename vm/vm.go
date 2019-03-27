package vm

import (
	"errors"

	"github.com/skatsuta/monkey-compiler/code"
	"github.com/skatsuta/monkey-compiler/compiler"
	"github.com/skatsuta/monkey-compiler/object"
)

// StackSize is an initial stack size.
const StackSize = 2048

// VM is a virtual machine which interprets and executes bytecode instructions.
type VM struct {
	consts []object.Object
	insns  code.Instructions

	stack []object.Object
	// Stackpointer always points to the *next* value. Top of stack is `stack[sp-1]`.
	sp int
}

// New creates a new VM instance which executes the given bytecode.
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		insns:  bytecode.Instructions,
		consts: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// StackTop returns an object at the top of stack.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// Run executes bytecode instructions.
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.insns); ip++ {
		op := code.Opcode(vm.insns[ip])

		switch op {
		case code.OpConstant:
			// Read a 2-byte operand from the next position
			constIdx := code.ReadUint16(vm.insns[ip+1:])
			// Because the operand is 2-byte width and we already read it,
			// increment the pointer by 2 (bytes)
			ip += 2

			if err := vm.push(vm.consts[constIdx]); err != nil {
				return err
			}

		case code.OpAdd:
			right := vm.pop().(*object.Integer).Value
			left := vm.pop().(*object.Integer).Value
			vm.push(&object.Integer{Value: left + right})
		}
	}

	return nil
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return errors.New("stack overflow")
	}

	// Push the object to the top of stack
	vm.stack[vm.sp] = obj
	// Increment the stack pointer
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	// Pop an object at the top of stack
	obj := vm.stack[vm.sp-1]
	// Decrement the stack pointer
	vm.sp--

	return obj
}
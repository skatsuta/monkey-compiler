package vm

import (
	"github.com/skatsuta/monkey-compiler/code"
	"github.com/skatsuta/monkey-compiler/object"
)

// Frame represents a stack frame.
type Frame struct {
	fn *object.CompiledFunction
	// Instruction pointer.
	ip int
	// Base pointer points to the bottom of the stack of the current stack frame.
	// It's also called "frame pointer".
	bp int
}

// NewFrame creates a new stack frame for a given compiled function.
func NewFrame(fn *object.CompiledFunction, bp int) *Frame {
	return &Frame{fn: fn, ip: -1, bp: bp}
}

// Instructions returns bytecode instructions of a function the stack frame is created for.
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}

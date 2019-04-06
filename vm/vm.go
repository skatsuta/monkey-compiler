package vm

import (
	"errors"
	"fmt"

	"github.com/skatsuta/monkey-compiler/code"
	"github.com/skatsuta/monkey-compiler/compiler"
	"github.com/skatsuta/monkey-compiler/object"
)

const (
	// StackSize is an initial stack size.
	StackSize = 2048

	// GlobalSize is an upper limit of the number of global bindings the VM can support.
	GlobalSize = 1 << 16 // 16 bits
)

var (
	// True is the boolean `true` value.
	True = &object.Boolean{Value: true}
	// False is the boolean `false` value.
	False = &object.Boolean{Value: false}
	// Nil represents the zero value.
	Nil = &object.Nil{}
)

// VM is a virtual machine which interprets and executes bytecode instructions.
type VM struct {
	consts []object.Object
	insns  code.Instructions

	stack []object.Object
	// Stackpointer always points to the *next* value. Top of stack is `stack[sp-1]`.
	sp int

	// globals store
	globals []object.Object
}

// New creates a new VM instance which executes the given bytecode.
func New(bytecode *compiler.Bytecode) *VM {
	return NewWithGlobalStore(bytecode, make([]object.Object, GlobalSize))
}

// NewWithGlobalStore creates a new VM instance which executes the given bytecode with the
// given globals store.
func NewWithGlobalStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	return &VM{
		insns:  bytecode.Instructions,
		consts: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: globals,
	}
}

// StackTop returns an object on top of the stack.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// LastPoppedStackElem returns an object which was popped off the stack most recently.
func (vm *VM) LastPoppedStackElem() object.Object {
	// vm.sp always points to the *next free* slot in vm.stack
	return vm.stack[vm.sp]
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

		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}

		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}

		case code.OpNil:
			if err := vm.push(Nil); err != nil {
				return err
			}

		case code.OpArray:
			numElems := int(code.ReadUint16(vm.insns[ip+1:]))
			ip += 2

			startIdx := vm.sp - numElems
			arr := vm.buildArray(startIdx, vm.sp)
			vm.sp = startIdx

			if err := vm.push(arr); err != nil {
				return err
			}

		case code.OpHash:
			numElems := int(code.ReadUint16(vm.insns[ip+1:]))
			ip += 2

			startIdx := vm.sp - numElems
			hash, err := vm.buildHash(startIdx, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = startIdx

			if err := vm.push(hash); err != nil {
				return err
			}

		case code.OpIndex:
			idx := vm.pop()
			left := vm.pop()

			if err := vm.execIndexExpr(left, idx); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpBang:
			if err := vm.execBangOp(); err != nil {
				return err
			}

		case code.OpMinus:
			if err := vm.execMinusOp(); err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.execBinaryOp(op); err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := vm.execComparison(op); err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(vm.insns[ip+1:]))
			// Since we're in a loop that increments `ip` with each iteration, we need to set `ip`
			// to the offset *right before the one* we want.
			ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.insns[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpSetGlobal:
			globalIdx := code.ReadUint16(vm.insns[ip+1:])
			ip += 2

			vm.globals[globalIdx] = vm.pop()

		case code.OpGetGlobal:
			globalIdx := code.ReadUint16(vm.insns[ip+1:])
			ip += 2

			if err := vm.push(vm.globals[globalIdx]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return errors.New("stack overflow")
	}

	// Push the object on to the stack
	vm.stack[vm.sp] = obj
	// Increment the stack pointer
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	// Pop an object off the stack
	obj := vm.stack[vm.sp-1]
	// Decrement the stack pointer
	vm.sp--

	return obj
}

func (vm *VM) buildArray(startIdx, endIdx int) object.Object {
	elems := make([]object.Object, endIdx-startIdx)

	for i := startIdx; i < endIdx; i++ {
		elems[i-startIdx] = vm.stack[i]
	}

	return &object.Array{Elements: elems}
}

func (vm *VM) buildHash(startIdx, endIdx int) (object.Object, error) {
	m := make(map[object.HashKey]object.HashPair)

	for i := startIdx; i < endIdx; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: val}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		m[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: m}, nil
}

func (vm *VM) execBangOp() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False, Nil:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) execMinusOp() error {
	operand := vm.pop()

	typ := operand.Type()
	if typ != object.IntegerType {
		return fmt.Errorf("unsupported type for negation: %s", typ)
	}

	val := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -val})
}

func (vm *VM) execBinaryOp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.IntegerType && rightType == object.IntegerType:
		return vm.execBinaryIntOp(op, left, right)
	case leftType == object.StringType && rightType == object.StringType:
		return vm.execBinaryStrOp(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s and %s", leftType, rightType)
	}

}

func (vm *VM) execBinaryIntOp(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) execBinaryStrOp(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return vm.push(&object.String{Value: leftVal + rightVal})
}

func (vm *VM) execIndexExpr(left, idx object.Object) error {
	leftType := left.Type()
	switch {
	case leftType == object.ArrayType && idx.Type() == object.IntegerType:
		return vm.execArrayIndex(left, idx)
	case leftType == object.HashType:
		return vm.execHashIndex(left, idx)
	default:
		return fmt.Errorf("index operator not supported: %s", leftType)
	}
}

func (vm *VM) execArrayIndex(array, idx object.Object) error {
	arr := array.(*object.Array)
	i := idx.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(Nil)
	}

	return vm.push(arr.Elements[i])
}

func (vm *VM) execHashIndex(hash, idx object.Object) error {
	h := hash.(*object.Hash)

	key, ok := idx.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", idx.Type())
	}

	pair, ok := h.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Nil)
	}

	return vm.push(pair.Value)
}

func (vm *VM) execComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.IntegerType && rightType == object.IntegerType {
		return vm.execIntComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator %d: %s and %s", op, leftType, rightType)
	}
}

func (vm *VM) execIntComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal == rightVal))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal != rightVal))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator %d for integers", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Nil:
		return false
	default:
		return true
	}
}

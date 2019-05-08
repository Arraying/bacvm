package bacvm

import (
	"fmt"
	"strconv"
)

type (
	// Function represents a native function.
	Function func(*VM, []Variable) Variable
	// Instruction represents an instruction
	Instruction struct {
		Argument  string
		Operation string
	}
	// Native represents a native function function matter.
	Native struct {
		Name string
		Task Function
	}
	// Scope represents a scope.
	Scope struct {
		previous  *Scope
		variables map[string]Variable
	}
	// VM represents the BacVM virtual machine.
	VM struct {
		Scope        *Scope
		Verbose      bool
		buffer       *buffer
		err          error
		exit         bool
		feeding      *feeding
		index        int
		instructions []*Instruction
		natives      map[string]Function
		reading      *reading
	}
	buffer struct {
		value    string
		previous *buffer
	}
	operation func(*VM, string)
	reading   struct {
		current  bool
		previous *reading
	}
)

var operations = map[string]operation{
	"bd": bd,
	"ex": ex,
	"ff": ff,
	"fi": fi,
	"gc": gc,
	"gt": gt,
	"pb": pb,
	"pu": pu,
	"rn": rn,
	"rt": rt,
	"sf": sf,
	"si": si,
	"vl": vl,
}

// String turns the instruction into a string so it can be pretty printed.
func (ins *Instruction) String() string {
	return ins.Operation + " " + ins.Argument
}

// Natives specifies an array of native functions.
func (vm *VM) Natives(natives ...*Native) {
	vm.natives = make(map[string]Function)
	for _, native := range natives {
		vm.natives[native.Name] = native.Task
	}
}

// Run runs a set of operations in the virtual machine.
func (vm *VM) Run(get func() []*Instruction) error {
	if vm.Verbose {
		fmt.Println("BacVM version", Version())
	}
	vm.instructions = get()
	vm.Scope = &Scope{}
	vm.reading = &reading{
		current: true,
	}
	for vm.index < len(vm.instructions) && !vm.exit {
		ins := vm.instructions[vm.index]
		if vm.Verbose {
			fmt.Println("Executing instruction", ins)
		}
		op, contains := operations[ins.Operation]
		if !contains {
			return ErrorOperationUnknown
		}
		if vm.reading != nil && !vm.reading.current {
			if ins.Operation != "rn" && ins.Operation != "rt" {
				vm.index++
				continue
			}
		}
		op(vm, ins.Argument)
		if vm.err != nil {
			return vm.err
		}
		vm.index++
	}
	return nil
}

// bufferPush pushes a new value onto the buffer, the previous value will be stored.
func (vm *VM) bufferPush(data string) {
	newBuffer := &buffer{
		previous: vm.buffer,
		value:    data,
	}
	vm.buffer = newBuffer
}

// bufferPop pops the latest value from the buffer, and the previously buffered value will be the new buffered value.
// If the buffer is empty then an empty string will return.
func (vm *VM) bufferPop() (string, error) {
	currentBuffer := vm.buffer
	if currentBuffer == nil {
		return "", ErrorBufferEmpty
	}
	vm.buffer = currentBuffer.previous
	return currentBuffer.value, nil
}

// bufferVariable is a method specifically for buffering variables.
// This will take the string value of the variable and buffer it.
func (vm *VM) bufferVariable(variable Variable) {
	vm.bufferPush(variable.Value())
}

// scopeResolve determines the scope in which a variable is registered, and if found executes the callback.
func (vm *VM) scopeResolve(variable string, result func(*Scope)) {
	working := vm.Scope
	for working != nil {
		if _, here := working.variables[variable]; here {
			result(working)
			return
		}
		working = working.previous
	}
}

// Get gets a variable by identifier.
// This will check the current scope, as well as all previous scopes.
func (scope *Scope) Get(identifier string) Variable {
	scope.initialize()
	variable, contains := scope.variables[identifier]
	if !contains && scope.previous != nil {
		return scope.previous.Get(identifier)
	}
	return variable
}

// Put puts a variable into the current scope, thus "overriding" any previous values in previous scopes.
// The previous values are not deleted, and will be restored once the current scope is destroyed.
func (scope *Scope) Put(identifier string, variable Variable) {
	scope.initialize()
	scope.variables[identifier] = variable
}

// initialize initializes the scope.
func (scope *Scope) initialize() {
	if scope.variables == nil {
		scope.variables = make(map[string]Variable, 0)
	}
}

// bd pops the buffer, and destroys the outcome.
// This is used when something buffers a value of no interest (i.e. result of function call not needed).
func bd(vm *VM, argument string) {
	_, err := vm.bufferPop()
	if err != nil {
		vm.err = err
	}
}

// ex exits the virtual machine.
func ex(vm *VM, argument string) {
	vm.exit = true
}

// ff finalizes the current feed, attempting to perform the actual operation.
func ff(vm *VM, argument string) {
	if vm.feeding == nil {
		vm.err = ErrorFeedSize
		return
	}
	current := vm.feeding
	if err := current.current.finalize(vm); err != nil {
		vm.err = err
		return
	}
	vm.feeding = current.previous
}

// fi initializes a new feed.
func fi(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	feeding := &feeding{
		previous: vm.feeding,
	}
	switch argument {
	case feederComparisonTypeEg, feederComparisonTypeEq, feederComparisonTypeEs, feederComparisonTypeG, feederComparisonTypeS:
		feeding.current = &feederComparison{
			kind: argument,
		}
	case feederFunc:
		feeding.current = &feederFunction{
			arguments: make([]Variable, 0),
		}
	case feederVar:
		feeding.current = &feederVariable{}
	default:
		vm.err = ErrorFeedType
		return
	}
	vm.feeding = feeding
}

// gc garbage cleans a variable.
// This will only delete the variable in the first found scope, other scopes will still contain the value.
func gc(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	vm.scopeResolve(argument, func(scope *Scope) {
		delete(scope.variables, argument)
	})
}

// gt jumps to an instruction.
func gt(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	inst, err := strconv.Atoi(argument)
	if err != nil {
		vm.err = ErrorOperationArgument
		return
	}
	if inst < 0 || inst >= len(vm.instructions) {
		vm.err = ErrorOperationArgument
		return
	}
	vm.index = inst - 1
}

// pb pops the buffer and performs a push with it.
func pb(vm *VM, argument string) {
	arg, err := vm.bufferPop()
	if err != nil {
		vm.err = err
		return
	}
	pu(vm, arg)
}

// pu pushes the given argument into the current feed operation
func pu(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	if vm.feeding == nil {
		vm.err = ErrorFeedSize
		return
	}
	vm.err = vm.feeding.current.feed(vm, argument)
}

// // TODO -- run currently errors, if the current section is not meant to be flipped, it will still get flipped
// // possible solution: current reading and all previous readings need to be true

// // rn inverts the current reading, if applicable
// func rn(vm *VM, argument string) {
// 	if vm.reading == nil {
// 		vm.err = ErrorReadingClose
// 		return
// 	}
// 	previously := vm.reading.truthy()
// 	if !previously {
// 		return
// 	}
// 	vm.reading.current = !(vm.reading.previouslyTrue() && vm.reading.current)
// }

// rn inverts the current reading.
func rn(vm *VM, argument string) {
	if vm.reading == nil {
		vm.err = ErrorReadingClose
		return
	}
	if !vm.reading.previously() {
		return
	}
	vm.reading.current = !vm.reading.current
}

// rt stops the current reading, if applicable
func rt(vm *VM, argument string) {
	if vm.reading == nil {
		vm.err = ErrorReadingClose
		return
	}
	current := vm.reading
	vm.reading = current.previous
}

// sf destroys the current scope, alongside all scope-specific variables in it.
// If a variable with the same name was declared in previous scopes then this variable will obtain the value of the variable in the current scope.
func sf(vm *VM, argument string) {
	sco := vm.Scope
	if sco == nil || sco.previous == nil {
		vm.err = ErrorScopeMin
		return
	}
	vm.Scope = sco.previous
	for name, variable := range sco.variables {
		vm.scopeResolve(name, func(scope *Scope) {
			scope.variables[name] = variable
		})
	}
}

// si creates a new scope.
func si(vm *VM, argument string) {
	sco := &Scope{
		previous: vm.Scope,
	}
	vm.Scope = sco
}

// vl loads a variable into the buffer, by the name.
func vl(vm *VM, argument string) {
	variable := vm.Scope.Get(argument)
	vm.bufferVariable(variable)
}

// truthy checks if ALL the previous values are true.
func (reading *reading) truthy() bool {
	if reading.previous == nil {
		return true
	}
	return reading.previous.current && reading.previous.truthy()
}

func (reading *reading) previously() bool {
	if reading.previous == nil {
		return true
	}
	return reading.previous.current
}

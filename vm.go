package bacvm

import (
	"fmt"
	"strconv"
)

type (
	// Function represents a native function.
	Function func(*VM, []*Variable) Variable
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
		variables map[string]*Variable
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
	}
	buffer struct {
		value    string
		previous *buffer
	}
	operation func(*VM, string)
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
	"sf": sf,
	"si": si,
	"vl": vl,
}

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
	for vm.index < len(vm.instructions) && !vm.exit {
		ins := vm.instructions[vm.index]
		if vm.Verbose {
			fmt.Println("Executing instruction", ins)
		}
		op, contains := operations[ins.Operation]
		if !contains {
			return ErrorOperationUnknown
		}
		op(vm, ins.Argument)
		if vm.err != nil {
			return vm.err
		}
		vm.index++
	}
	return nil
}

func (vm *VM) bufferPush(data string) {
	newBuffer := &buffer{
		previous: vm.buffer,
		value:    data,
	}
	vm.buffer = newBuffer
}

func (vm *VM) bufferPop() (string, error) {
	currentBuffer := vm.buffer
	if currentBuffer == nil {
		return "", ErrorBufferEmpty
	}
	vm.buffer = currentBuffer.previous
	return currentBuffer.value, nil
}

func (vm *VM) bufferVariable(variable *Variable) {
	vm.bufferPush(variable.ValueString())
	vm.bufferPush(string(variable.Type))
}

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
func (scope *Scope) Get(identifier string) *Variable {
	scope.initialize()
	variable, contains := scope.variables[identifier]
	if !contains && scope.previous != nil {
		return scope.previous.Get(identifier)
	}
	return variable
}

// Put puts a variable.
func (scope *Scope) Put(variable *Variable) {
	scope.initialize()
	scope.variables[variable.Name] = variable
}

func (scope *Scope) initialize() {
	if scope.variables == nil {
		scope.variables = make(map[string]*Variable, 0)
	}
}

func bd(vm *VM, argument string) {
	_, err := vm.bufferPop()
	if err != nil {
		vm.err = err
	}
}

func ex(vm *VM, argument string) {
	vm.exit = true
}

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

func fi(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	feeding := &feeding{
		previous: vm.feeding,
	}
	switch argument {
	case feederComparisonTypeA, feederComparisonTypeEg, feederComparisonTypeEq, feederComparisonTypeEs, feederComparisonTypeG, feederComparisonTypeO, feederComparisonTypeS:
		feeding.current = &feederComparison{
			kind: argument,
		}
	case feederVar:
		feeding.current = &feederVariable{
			variable: &Variable{},
		}
	default:
		vm.err = ErrorFeedType
		return
	}
	vm.feeding = feeding
}

func gc(vm *VM, argument string) {
	if argument == "" {
		vm.err = ErrorOperationArgument
		return
	}
	vm.scopeResolve(argument, func(scope *Scope) {
		delete(scope.variables, argument)
	})
}

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

func pb(vm *VM, argument string) {
	arg, err := vm.bufferPop()
	if err != nil {
		vm.err = err
		return
	}
	pu(vm, arg)
}

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

func pv(vm *VM, argument string) {

}

func sf(vm *VM, argument string) {
	sco := vm.Scope
	if sco == nil || sco.previous == nil {
		vm.err = ErrorScopeMin
		return
	}
	vm.Scope = sco.previous
	for _, variable := range sco.variables {
		vm.scopeResolve(variable.Name, func(scope *Scope) {
			scope.variables[variable.Name] = variable
		})
	}
	// 	working := vm.Scope
	// outer:
	// 	for working != nil {
	// 		for _, vari := range sco.variables {
	// 			if _, here := working.variables[vari.Name]; here {
	// 				working.variables[vari.Name] = vari
	// 				break outer
	// 			}
	// 		}
	// 		working = working.previous
	// 	}
}

func si(vm *VM, argument string) {
	sco := &Scope{
		previous: vm.Scope,
	}
	vm.Scope = sco
}

func vl(vm *VM, argument string) {
	variable := vm.Scope.Get(argument)
	if variable == nil {
		vm.err = ErrorVariableExistance
		return
	}
	vm.bufferVariable(variable)
}

package native

import (
	"fmt"

	"github.com/arraying/bacvm"
)

// StdOut prints the specified variables, with no new line.
func StdOut(vm *bacvm.VM, variables []bacvm.Variable) bacvm.Variable {
	return stdout(vm, variables, func(in ...interface{}) (int, error) {
		return fmt.Print(in...)
	})
}

// StdOutLn prints the specified variables, with a new line.
func StdOutLn(vm *bacvm.VM, variables []bacvm.Variable) bacvm.Variable {
	return stdout(vm, variables, func(in ...interface{}) (int, error) {
		return fmt.Println(in...)
	})
}

// stdout gets the input ready to be printed, the handler being the consumer dealing with prints.
// It returns the number of bytes written. If this number is bigger than 53 bits it will be reduced to 0.
// Similarly, 0 will return if there is an error.
func stdout(vm *bacvm.VM, variables []bacvm.Variable, handle func(...interface{}) (int, error)) bacvm.Variable {
	args := make([]interface{}, len(variables), len(variables))
	for i, variable := range variables {
		args[i] = variable
	}
	var out int
	if out, err := handle(args...); err != nil || !bacvm.VariableBounded(out) {
		out = 0
	}
	return bacvm.Variable(out)
}
